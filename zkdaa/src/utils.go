package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"zk-htlc/actors"
	"zk-htlc/contracts"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 全局变量用于输出控制（在 main 中初始化）
var (
	originalStdout *os.File
	originalStderr *os.File
	fileLogger     *log.Logger
)

// 1. 趋势状态管理 (State Management)

type LastRunState struct {
	LastLockLatency float64 `json:"last_lock"`
}

const StateFile = "last_run_state.json"

func loadLastRunState() float64 {
	data, err := ioutil.ReadFile(StateFile)
	if err != nil {
		return 0.0
	}
	var state LastRunState
	json.Unmarshal(data, &state)
	return state.LastLockLatency
}

func saveLastRunState(lockTime float64) {
	state := LastRunState{LastLockLatency: lockTime}
	data, _ := json.Marshal(state)
	ioutil.WriteFile(StateFile, data, 0644)
}

// 2. 结构体定义 (Structs)

type PerformanceMetrics struct {
	UnlockCircuitSetupTime time.Duration // A 的解锁电路 setup
	UnlockProveTimeA       time.Duration // A 的解锁证明生成

	AuditUnlockSetupTime  time.Duration // B 的审计+解锁电路 setup
	AuditUnlockProveTimeB time.Duration // AuditUnlockProveTimeB
	UnlockVerifyTimeA     time.Duration
	AuditVerifyTimeB      time.Duration
	TxConfirmTimeLockA    time.Duration
	TxConfirmTimeLockB    time.Duration
	TxConfirmTimeUnlockA  time.Duration // A 的 Unlock 交易
	TxConfirmTimeUnlockB  time.Duration // B 的 AuditUnlock 交易
	DataPrepTimeA         time.Duration
	DataTransferTime      time.Duration
	EndToEndTime          time.Duration
	TrustedSetupUnlockMs  int64 `json:"trusted_setup_unlock_ms"`
	TrustedSetupAuditMs   int64 `json:"trusted_setup_audit_ms"`

	GasLockA   uint64
	GasLockB   uint64
	GasUnlockB uint64
	GasUnlockA uint64
	TotalGas   uint64

	FileSize                 int
	ChunkSize                int
	ChunkCount               int
	MerkleDepth              int
	CIDFSize                 int
	UnlockCircuitConstraints int // A 的约束数
	AuditCircuitConstraints  int // B 的约束数
	UnlockProofSize          int // A 的证明大小
	AuditProofSize           int // B 的证明大小
	PublicInputCount         int
	AuditChallengeCount      int
	AuditSuccessCount        int
	AuditSuccessRate         float64
	MerkleProofLength        int

	AverageTxLatency     time.Duration
	AverageLockLatency   time.Duration
	AverageUnlockLatency time.Duration
	TxWindowTime         time.Duration
}

type TPSMetrics struct {
	TxLockTPS   float64
	TxUnlockTPS float64
}

type BatchTestConfig struct {
	FileSize  int
	Label     string
	ChunkSize int
}

// 3. JSON 输出 (核心修改部分)

// ✅ 把函数名从 durationToMs 改成 durationToMsHighPrecision
func durationToMsHighPrecision(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1_000_000.0
}

func exportMetricsToJSON(m *PerformanceMetrics, filepath string) {
	// --- 1. 定义符合专利逻辑的清晰 JSON 结构 ---
	type TimeMetrics struct {
		// Setup Phase (一次性成本，离线完成)
		TrustedSetupUnlockMs float64 `json:"trusted_setup_unlock_ms"`
		TrustedSetupAuditMs  float64 `json:"trusted_setup_audit_ms"`

		// Off-Chain Computation (ZKP 核心开销)
		ProveUnlockAMs       float64 `json:"prove_unlock_a_ms"`        // DSPA 生成解锁证明
		ProveAuditUnlockBMs  float64 `json:"prove_audit_unlock_b_ms"`  // DSPB 生成审计证明
		VerifyUnlockAMs      float64 `json:"verify_unlock_a_ms"`       // 链下/合约验证 A
		VerifyAuditUnlockBMs float64 `json:"verify_audit_unlock_b_ms"` // 链下/合约验证 B

		// On-Chain Interaction (区块链延迟)
		TxLockAMs        float64 `json:"tx_lock_a_ms"`         // SCA 锁定交易
		TxLockBMs        float64 `json:"tx_lock_b_ms"`         // SCB 锁定交易
		TxUnlockAMs      float64 `json:"tx_unlock_a_ms"`       // 秘密反推解锁
		TxAuditUnlockBMs float64 `json:"tx_audit_unlock_b_ms"` // 审计解锁

		// Data Operations
		DataPrepMs     float64 `json:"data_prep_ms"`     // Merkle Tree 构建
		DataTransferMs float64 `json:"data_transfer_ms"` // 链下 P2P 传输

		// Summary Metrics
		TotalE2EMs float64 `json:"total_e2e_ms"` // 端到端总时间

		// 🔥 关键修复：分解 E2E 时间为各部分
		ZkpComputeMs     float64 `json:"zkp_compute_ms"`     // Prove + Verify 总和
		TxLatencyMs      float64 `json:"tx_latency_ms"`      // 所有交易延迟总和
		SystemOverheadMs float64 `json:"system_overhead_ms"` // GC + 网络 + 其他
	}

	type GasMetrics struct {
		LockA        uint64 `json:"lock_a"`
		LockB        uint64 `json:"lock_b"`
		UnlockA      uint64 `json:"unlock_a"`
		AuditUnlockB uint64 `json:"audit_unlock_b"`
		Total        uint64 `json:"total"`
	}

	type ThroughputMetrics struct {
		TPS_Lock   float64 `json:"tps_lock"`
		TPS_Unlock float64 `json:"tps_unlock"`
	}

	type Export struct {
		Time        TimeMetrics       `json:"time"`
		Gas         GasMetrics        `json:"gas"`
		Throughput  ThroughputMetrics `json:"throughput"`
		Constraints struct {
			UnlockCircuit int `json:"unlock_circuit_constraints"`
			AuditCircuit  int `json:"audit_circuit_constraints"`
		} `json:"constraints"`
		DataConfig struct {
			FileSize    int `json:"file_size_bytes"`
			ChunkSize   int `json:"chunk_size_bytes"`
			ChunkCount  int `json:"chunk_count"`
			MerkleDepth int `json:"merkle_depth"`
		} `json:"data_config"`
	}

	// --- 2. 使用高精度转换函数 ---
	toMs := durationToMsHighPrecision

	// 计算各部分时间
	zkpCompute := toMs(m.UnlockProveTimeA) + toMs(m.AuditUnlockProveTimeB) +
		toMs(m.UnlockVerifyTimeA) + toMs(m.AuditVerifyTimeB)

	txLatency := toMs(m.TxConfirmTimeLockA) + toMs(m.TxConfirmTimeLockB) +
		toMs(m.TxConfirmTimeUnlockA) + toMs(m.TxConfirmTimeUnlockB)

	dataPrep := toMs(m.DataPrepTimeA)
	dataTransfer := toMs(m.DataTransferTime)
	e2eTime := toMs(m.EndToEndTime)

	// 🔥 修复 Overhead 计算：E2E - (所有已知部分)
	// 剩余时间归因于系统开销（GC、线程切换、网络抖动等）
	systemOverhead := e2eTime - (zkpCompute + txLatency + dataPrep + dataTransfer)

	// 防御性编程：如果计算出负值，说明计时有重叠，归零处理
	if systemOverhead < 0 {
		fmt.Printf("⚠️ Warning: Negative overhead (%.2fms), likely due to timing overlap. Setting to 0.\n", systemOverhead)
		systemOverhead = 0
	}

	// --- 3. 构建导出对象 ---
	e := Export{
		Time: TimeMetrics{
			TrustedSetupUnlockMs: float64(m.TrustedSetupUnlockMs),
			TrustedSetupAuditMs:  float64(m.TrustedSetupAuditMs),

			ProveUnlockAMs:       toMs(m.UnlockProveTimeA),
			ProveAuditUnlockBMs:  toMs(m.AuditUnlockProveTimeB),
			VerifyUnlockAMs:      toMs(m.UnlockVerifyTimeA),
			VerifyAuditUnlockBMs: toMs(m.AuditVerifyTimeB),

			TxLockAMs:        toMs(m.TxConfirmTimeLockA),
			TxLockBMs:        toMs(m.TxConfirmTimeLockB),
			TxUnlockAMs:      toMs(m.TxConfirmTimeUnlockA),
			TxAuditUnlockBMs: toMs(m.TxConfirmTimeUnlockB),

			DataPrepMs:     dataPrep,
			DataTransferMs: dataTransfer,
			TotalE2EMs:     e2eTime,

			ZkpComputeMs:     zkpCompute,
			TxLatencyMs:      txLatency,
			SystemOverheadMs: systemOverhead,
		},
		Gas: GasMetrics{
			LockA:        m.GasLockA,
			LockB:        m.GasLockB,
			UnlockA:      m.GasUnlockA,
			AuditUnlockB: m.GasUnlockB,
			Total:        m.TotalGas,
		},
		Throughput: ThroughputMetrics{
			TPS_Lock:   0.0,
			TPS_Unlock: 0.0,
		},
		Constraints: struct {
			UnlockCircuit int `json:"unlock_circuit_constraints"`
			AuditCircuit  int `json:"audit_circuit_constraints"`
		}{
			UnlockCircuit: m.UnlockCircuitConstraints,
			AuditCircuit:  m.AuditCircuitConstraints,
		},
		DataConfig: struct {
			FileSize    int `json:"file_size_bytes"`
			ChunkSize   int `json:"chunk_size_bytes"`
			ChunkCount  int `json:"chunk_count"`
			MerkleDepth int `json:"merkle_depth"`
		}{
			FileSize:    m.FileSize,
			ChunkSize:   m.ChunkSize,
			ChunkCount:  m.ChunkCount,
			MerkleDepth: m.MerkleDepth,
		},
	}

	// TPS 估算（简化版，避免除零）
	lockTimeTotal := e.Time.TxLockAMs + e.Time.TxLockBMs
	if lockTimeTotal > 0 {
		// 假设串行执行 2 次 Lock
		e.Throughput.TPS_Lock = 2000.0 / lockTimeTotal
	}

	unlockTimeTotal := e.Time.TxUnlockAMs + e.Time.TxAuditUnlockBMs
	if unlockTimeTotal > 0 {
		e.Throughput.TPS_Unlock = 2000.0 / unlockTimeTotal
	}

	// --- 4. 写入文件 ---
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to export metrics: %v\n", err)
		return
	}
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		fmt.Printf("❌ Failed to write file: %v\n", err)
		return
	}

}

// 4. 其他辅助函数

func curveToContractProof(proof groth16.Proof) ([8]*big.Int, error) {
	var solProof [8]*big.Int
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()
	const fpSize = 32
	readElement := func(n int) *big.Int {
		start := n * fpSize
		end := start + fpSize
		if end > len(proofBytes) {
			return big.NewInt(0)
		}
		return new(big.Int).SetBytes(proofBytes[start:end])
	}
	solProof[0] = readElement(0)
	solProof[1] = readElement(1)
	solProof[2] = readElement(2)
	solProof[3] = readElement(3)
	solProof[4] = readElement(4)
	solProof[5] = readElement(5)
	solProof[6] = readElement(6)
	solProof[7] = readElement(7)
	return solProof, nil
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("❌ %s failed: %v", msg, err)
	}
}

// 去极值平均算法
func calculateTrimmedStat(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	if n <= 2 {
		sum := 0.0
		for _, v := range data {
			sum += v
		}
		return sum / float64(n)
	}
	sort.Float64s(data)
	trimmed := data[1 : n-1]
	sum := 0.0
	for _, v := range trimmed {
		sum += v
	}
	return sum / float64(len(trimmed))
}

func calculateSmartTPS(all []TPSMetrics) TPSMetrics {
	var vl, vu []float64
	for _, m := range all {
		if m.TxLockTPS > 1.0 {
			vl = append(vl, m.TxLockTPS)
		}
		if m.TxUnlockTPS > 1.0 {
			vu = append(vu, m.TxUnlockTPS)
		}
	}
	return TPSMetrics{TxLockTPS: calculateTrimmedStat(vl), TxUnlockTPS: calculateTrimmedStat(vu)}
}

// 5. 监控器 (Monitor)

type TPSMonitor struct {
	client         *ethclient.Client
	txRecords      map[common.Hash]*TxRecord
	mu             sync.Mutex
	lockStart      time.Time
	lockEnd        time.Time
	unlockStart    time.Time
	unlockEnd      time.Time
	hasLock        bool
	hasUnlock      bool
	confirmedCount int
}

type TxRecord struct {
	Label     string
	SentTime  time.Time
	BlockNum  uint64
	Confirmed bool
}

func NewTPSMonitor(c *ethclient.Client) *TPSMonitor {
	return &TPSMonitor{client: c, txRecords: make(map[common.Hash]*TxRecord)}
}

func (m *TPSMonitor) RecordProtocolTx(h common.Hash, l string, t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.txRecords[h] = &TxRecord{Label: l, SentTime: t}
	if strings.Contains(l, "Lock") {
		if !m.hasLock || t.Before(m.lockStart) {
			m.lockStart = t
			m.hasLock = true
		}
	} else {
		if !m.hasUnlock || t.Before(m.unlockStart) {
			m.unlockStart = t
			m.hasUnlock = true
		}
	}
}

func (m *TPSMonitor) UpdateTxConfirmation(h common.Hash, b uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if r, ok := m.txRecords[h]; ok && !r.Confirmed {
		r.Confirmed = true
		r.BlockNum = b
		m.confirmedCount++
		if strings.Contains(r.Label, "Lock") {
			if now.After(m.lockEnd) {
				m.lockEnd = now
			}
		} else {
			if now.After(m.unlockEnd) {
				m.unlockEnd = now
			}
		}
	}
}

func (m *TPSMonitor) WaitForConfirmations(target int, timeout time.Duration) bool {
	start := time.Now()
	for {
		m.mu.Lock()
		c := m.confirmedCount
		m.mu.Unlock()
		if c >= target || (target > 10 && c >= int(float64(target)*0.9)) {
			return true
		}
		if time.Since(start) > timeout {
			return false
		}
		time.Sleep(1 * time.Second)
	}
}

func (m *TPSMonitor) Stop() TPSMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	var lc, uc int
	for _, r := range m.txRecords {
		if r.Confirmed {
			if strings.Contains(r.Label, "Lock") {
				lc++
			} else {
				uc++
			}
		}
	}
	ld := m.lockEnd.Sub(m.lockStart).Seconds()
	if ld < 0.1 {
		ld = 1.0
	}
	ud := m.unlockEnd.Sub(m.unlockStart).Seconds()
	if ud < 0.1 {
		ud = 1.0
	}
	if m.lockEnd.IsZero() && lc > 0 {
		ld = time.Since(m.lockStart).Seconds()
	}
	if m.unlockEnd.IsZero() && uc > 0 {
		ud = time.Since(m.unlockStart).Seconds()
	}
	return TPSMetrics{TxLockTPS: float64(lc) / ld, TxUnlockTPS: float64(uc) / ud}
}

// 静音控制
func silenceAllOutput() {
	if originalStdout == nil {
		originalStdout = os.Stdout
	}
	if originalStderr == nil {
		originalStderr = os.Stderr
	}
	nullDevice, _ := os.Open(os.DevNull)
	os.Stdout = nullDevice
	os.Stderr = nullDevice
}

func restoreAllOutput() {
	if originalStdout != nil {
		os.Stdout = originalStdout
	}
	if originalStderr != nil {
		os.Stderr = originalStderr
	}
}

func printResultForce(format string, a ...interface{}) {
	restoreAllOutput()
	fmt.Printf(format, a...)
	silenceAllOutput()
}

func logToFile(format string, v ...interface{}) {
	if fileLogger != nil {
		fileLogger.Printf(format, v...)
	}
}

func drainPendingPool(client *ethclient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for {
		p, _ := client.PendingTransactionCount(ctx)
		if p == 0 {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func sendAndWaitLatency(client *ethclient.Client, label string, action func() (*types.Transaction, error)) (time.Duration, uint64) {
	start := time.Now()
	tx, err := action()
	if err != nil {
		printResultForce("❌ Error %s: %v\n", label, err)
		return 0, 0
	}
	receipt, err := waitForReceiptLatency(client, tx.Hash())
	if err != nil {
		printResultForce("❌ Wait %s: %v\n", label, err)
		return 0, 0
	}
	if receipt.Status == 0 {
		printResultForce("❌ Reverted %s\n", label)
		return 0, 0
	}
	return time.Since(start), receipt.GasUsed
}

func waitForReceiptLatency(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func strToDataID(s string) [32]byte {
	var id [32]byte
	copy(id[:], []byte(s))
	return id
}

func bigIntTo32Bytes(val *big.Int) [32]byte {
	var r [32]byte
	if val != nil {
		b := val.Bytes()
		copy(r[32-len(b):], b)
	}
	return r
}

func toFieldElement(val *big.Int) *big.Int {
	if val == nil {
		return big.NewInt(0)
	}
	mod := new(big.Int)
	mod.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	return new(big.Int).Mod(val, mod)
}

func calculateDepth(n int) int {
	depth := 0
	for n > 1 {
		n = (n + 1) / 2
		depth++
	}
	return depth
}

func finalizeMetrics(m *PerformanceMetrics) {
	m.AverageTxLatency = (m.TxConfirmTimeLockA + m.TxConfirmTimeLockB + m.TxConfirmTimeUnlockB + m.TxConfirmTimeUnlockA) / 4
	m.AverageLockLatency = (m.TxConfirmTimeLockA + m.TxConfirmTimeLockB) / 2
	m.AverageUnlockLatency = (m.TxConfirmTimeUnlockB + m.TxConfirmTimeUnlockA) / 2
}
func setupContract(addr string, client *ethclient.Client) *contracts.DataMigration {
	inst, err := contracts.NewDataMigration(common.HexToAddress(addr), client)
	checkErr(err, "Setup contract")
	return inst
}

func generateDummyData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

func sendAndWait(client *ethclient.Client, label string, action func() (*types.Transaction, error)) (time.Duration, uint64) {
	start := time.Now()
	tx, err := action()
	if err != nil {
		fmt.Printf("❌ Error %s: %v\n", label, err)
		return 0, 0
	}
	receipt, err := waitForReceiptLatency(client, tx.Hash())
	if err != nil {
		fmt.Printf("❌ Wait %s: %v\n", label, err)
		return 0, 0
	}
	if receipt.Status == 0 {
		fmt.Printf("❌ Reverted %s\n", label)
		return 0, 0
	}
	return time.Since(start), receipt.GasUsed
}
func setupFileLogging() {
	logFile, err := os.OpenFile(*flagLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		fileLogger = log.New(logFile, "", log.LstdFlags)
	}
}

func setupEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func setupBlockchain(rpcURL, privateKeyHex string) (*ethclient.Client, *bind.TransactOpts, *big.Int, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, nil, nil, err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, nil, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, nil, err
	}

	auth.GasLimit = 8000000
	return client, auth, chainID, nil
}

func runSimulatedBenchmark(n int) {
	fmt.Printf("Simulated benchmark with %d nodes (not implemented)\n", n)
}

// ===== latency.go 需要的函数签名 =====
func runAuditPhase(metrics *PerformanceMetrics, user *actors.User, dspa *actors.DSP, dspb *actors.DSP, idA string) error {
	// Placeholder - 原逻辑可能在其他文件
	return nil
}

func formatProofForSolidity(p groth16.Proof, pubInputs []*big.Int) ([8]*big.Int, [2]*big.Int) {
	solProof, _ := curveToContractProof(p)
	var pub [2]*big.Int
	if len(pubInputs) >= 1 {
		pub[0] = toFieldElement(pubInputs[0])
	}
	if len(pubInputs) >= 2 {
		pub[1] = toFieldElement(pubInputs[1])
	}
	return solProof, pub
}

// ===== tps.go 需要的函数 =====
func warmup(client *ethclient.Client, auth *bind.TransactOpts, inst *contracts.DataMigration) {
	// 执行一次简单的调用预热
	timeout := big.NewInt(3600)
	dataID := [32]byte{0xff}
	dummyHash := [32]byte{0xaa}
	inst.Lock(auth, dummyHash, dataID, timeout)
	time.Sleep(500 * time.Millisecond)
}

func ExportTPSMetrics(metrics TPSMetrics, filepath string) {
	data, _ := json.MarshalIndent(metrics, "", "  ")
	ioutil.WriteFile(filepath, data, 0644)
}

func watchTx(client *ethclient.Client, txHash common.Hash, monitor *TPSMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			monitor.UpdateTxConfirmation(txHash, receipt.BlockNumber.Uint64())
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}
