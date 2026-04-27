package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	mathrand "math/rand"
	"os"
	"runtime"
	"sort"
	"time"

	"zk-htlc/circuit"
	"zk-htlc/contracts"

	// "github.com/consensys/gnark-crypto/ecc"
	// "github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 🔥 Benchmark 模式主函数

func runBenchmarkMode(fileSize, chunkSize int, addrA, addrB string, rounds int, output string) {
	// 区块链配置
	realChainRPC := "http://127.0.0.1:8545"
	myPrivateKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	fmt.Printf("\n🚀 Benchmark Mode Started\n")
	fmt.Printf("   File Size: %d bytes\n", fileSize)
	fmt.Printf("   Chunk Size: %d bytes\n", chunkSize)
	fmt.Printf("   Rounds: %d\n", rounds)
	fmt.Printf("   Contracts: A=%s, B=%s\n\n", addrA, addrB)

	// 初始化区块链连接
	client, auth, _, err := initBatchClient(realChainRPC, myPrivateKey)
	if err != nil {
		log.Fatalf("❌ Blockchain setup failed: %v", err)
	}
	defer client.Close()

	instA := setupContract(addrA, client)
	instB := setupContract(addrB, client)

	// 🔥 一次性加载Handler（避免重复Setup）
	chunkCount := fileSize / chunkSize
	actualDepth := calcMerkleDepth(chunkCount)

	fmt.Printf("   [Init] Loading keys (depth=%d)...\n", actualDepth)
	unlockHandler, unlockSetupTime, err := NewUnlockHandler()
	if err != nil {
		log.Fatalf("❌ Failed to load Unlock Handler: %v", err)
	}
	fmt.Printf("   ✅ Unlock Handler loaded in %v\n", unlockSetupTime)

	auditHandler, auditSetupTime, err := NewAuditUnlockHandler(actualDepth)
	if err != nil {
		log.Fatalf("❌ Failed to load Audit Handler: %v", err)
	}
	fmt.Printf("   ✅ Audit Handler loaded in %v\n", auditSetupTime)

	// 🔥 Warm-up阶段
	fmt.Println("\n🔥 Warming up (20 rounds)...")
	for i := 0; i < 20; i++ {
		runSingleIteration(unlockHandler, auditHandler, fileSize, chunkSize, instA, instB, client, auth, true)
		if i%5 == 0 {
			fmt.Printf("   Warmup progress: %d/20\n", i)
		}
	}
	runtime.GC()
	time.Sleep(1 * time.Second)
	fmt.Println("   ✅ Warmup completed\n")

	// 🔥 正式采样
	var proveTimesA, proveTimesB []float64
	var verifyTimesA, verifyTimesB []float64
	var txLockTimesA, txLockTimesB []float64
	var txUnlockTimesA, txUnlockTimesB []float64
	var dataPrepTimes []float64
	var e2eTimes []float64

	var totalGasLockA, totalGasLockB, totalGasUnlockA, totalGasUnlockB uint64

	fmt.Printf("📊 Running %d benchmark rounds...\n", rounds)
	successCount := 0

	for i := 0; i < rounds; i++ {
		m := runSingleIteration(unlockHandler, auditHandler, fileSize, chunkSize, instA, instB, client, auth, false)

		if m != nil {
			proveTimesA = append(proveTimesA, durationToMsHighPrecision(m.UnlockProveTimeA))
			proveTimesB = append(proveTimesB, durationToMsHighPrecision(m.AuditUnlockProveTimeB))
			verifyTimesA = append(verifyTimesA, durationToMsHighPrecision(m.UnlockVerifyTimeA))
			verifyTimesB = append(verifyTimesB, durationToMsHighPrecision(m.AuditVerifyTimeB))

			// 链上指标单独收集
			txLockTimesA = append(txLockTimesA, durationToMsHighPrecision(m.TxConfirmTimeLockA))
			txLockTimesB = append(txLockTimesB, durationToMsHighPrecision(m.TxConfirmTimeLockB))
			txUnlockTimesA = append(txUnlockTimesA, durationToMsHighPrecision(m.TxConfirmTimeUnlockA))
			txUnlockTimesB = append(txUnlockTimesB, durationToMsHighPrecision(m.TxConfirmTimeUnlockB))

			dataPrepTimes = append(dataPrepTimes, durationToMsHighPrecision(m.DataPrepTimeA))
			e2eTimes = append(e2eTimes, durationToMsHighPrecision(m.EndToEndTime))

			totalGasLockA += m.GasLockA
			totalGasLockB += m.GasLockB
			totalGasUnlockA += m.GasUnlockA
			totalGasUnlockB += m.GasUnlockB

			successCount++
		}

		if (i+1)%10 == 0 {
			fmt.Printf("   Progress: %d/%d (Success: %d)\n", i+1, rounds, successCount)
		}
	}

	if successCount == 0 {
		log.Fatal("❌ No successful iterations!")
	}

	fmt.Printf("\n✅ Benchmark completed: %d/%d successful\n\n", successCount, rounds)

	// 🔥 统计学处理 - 计算中位数
	finalMetrics := &PerformanceMetrics{
		// ZKP Performance
		UnlockProveTimeA:      time.Duration(calculateMedian(proveTimesA) * 1e6),
		AuditUnlockProveTimeB: time.Duration(calculateMedian(proveTimesB) * 1e6),
		UnlockVerifyTimeA:     time.Duration(calculateMedian(verifyTimesA) * 1e6),
		AuditVerifyTimeB:      time.Duration(calculateMedian(verifyTimesB) * 1e6),

		// Blockchain Performance
		TxConfirmTimeLockA:   time.Duration(calculateMedian(txLockTimesA) * 1e6),
		TxConfirmTimeLockB:   time.Duration(calculateMedian(txLockTimesB) * 1e6),
		TxConfirmTimeUnlockA: time.Duration(calculateMedian(txUnlockTimesA) * 1e6),
		TxConfirmTimeUnlockB: time.Duration(calculateMedian(txUnlockTimesB) * 1e6),

		// Data Preparation
		DataPrepTimeA: time.Duration(calculateMedian(dataPrepTimes) * 1e6),

		// E2E
		EndToEndTime: time.Duration(calculateMedian(e2eTimes) * 1e6),

		// Gas (平均值)
		GasLockA:   totalGasLockA / uint64(successCount),
		GasLockB:   totalGasLockB / uint64(successCount),
		GasUnlockA: totalGasUnlockA / uint64(successCount),
		GasUnlockB: totalGasUnlockB / uint64(successCount),
		TotalGas:   (totalGasLockA + totalGasLockB + totalGasUnlockA + totalGasUnlockB) / uint64(successCount),

		// Setup Times (从文件读取)
		UnlockCircuitSetupTime: unlockSetupTime,
		AuditUnlockSetupTime:   auditSetupTime,

		// Metadata
		FileSize:                 fileSize,
		ChunkSize:                chunkSize,
		ChunkCount:               chunkCount,
		MerkleDepth:              actualDepth,
		UnlockCircuitConstraints: unlockHandler.GetConstraints(),
		AuditCircuitConstraints:  auditHandler.GetConstraints(),
	}

	// 读取 Trusted Setup 时间
	if content, err := os.ReadFile("build/setup_time_unlock.txt"); err == nil {
		fmt.Sscanf(string(content), "%d", &finalMetrics.TrustedSetupUnlockMs)
	}
	if content, err := os.ReadFile("build/setup_time_audit.txt"); err == nil {
		fmt.Sscanf(string(content), "%d", &finalMetrics.TrustedSetupAuditMs)
	}

	// 输出结果
	fmt.Println("📈 Final Median Results:")
	fmt.Printf("   Prove A:    %.3f ms\n", calculateMedian(proveTimesA))
	fmt.Printf("   Prove B:    %.3f ms\n", calculateMedian(proveTimesB))
	fmt.Printf("   Verify A:   %.3f ms\n", calculateMedian(verifyTimesA))
	fmt.Printf("   Verify B:   %.3f ms\n", calculateMedian(verifyTimesB))
	fmt.Printf("   E2E Time:   %.3f ms\n", calculateMedian(e2eTimes))

	// 导出 JSON
	exportMetricsToJSON(finalMetrics, output)
	fmt.Printf("\n✅ Results exported to: %s\n", output)
}

// 单次迭代（不含Setup）

func runSingleIteration(unlockH *UnlockHandler, auditH *AuditUnlockHandler,
	fileSize, chunkSize int, instA, instB *contracts.DataMigration,
	client *ethclient.Client, auth *bind.TransactOpts,
	isWarmup bool) *PerformanceMetrics {

	e2eStart := time.Now()

	// 数据准备
	prepStart := time.Now()
	data := generateDummyData(fileSize)
	chunkHashes, cidf, proofPaths, helpers, leafCounts, leafNumBytes := buildMerkleTreeForCircuit(data, chunkSize)
	dataPrepTime := time.Since(prepStart)

	// 生成密钥和哈希
	zMask := new(big.Int)
	zMask.SetString("AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55", 16)
	snII := randomFieldElement()
	snI := new(big.Int).Xor(snII, zMask)
	snI.Mod(snI, bn254FieldModulus)
	preI := randomFieldElement()
	h1 := mimcHashBig(preI, snI)
	h2 := mimcHashBig(cidf, snII)

	// Warm-up 模式：只执行 Prove，不上链
	if isWarmup {
		// 执行 Prove 让 JIT 编译
		unlockH.Prove(preI, snI, h1)

		challengeIdx := mathrand.Intn(len(chunkHashes))
		assignB := &circuit.AuditUnlockCircuit{
			ProofPath:    make([]frontend.Variable, auditH.depth),
			Helpers:      make([]frontend.Variable, auditH.depth),
			LeafCounts:   make([]frontend.Variable, auditH.depth),
			LeafNumBytes: make([]frontend.Variable, auditH.depth),
			Sn:           snII,
			ChunkIndex:   challengeIdx,
			ChunkHash:    chunkHashes[challengeIdx],
			H:            h2,
		}
		fillWitnessSlice(assignB, proofPaths[challengeIdx], helpers[challengeIdx], leafCounts, leafNumBytes)
		auditH.Prove(assignB)

		return nil // Warm-up 不记录结果
	}

	// 正式测试：完整流程
	m := &PerformanceMetrics{
		FileSize:      fileSize,
		ChunkSize:     chunkSize,
		DataPrepTimeA: dataPrepTime,
	}

	// Lock 阶段
	dataIdA := strToDataID(fmt.Sprintf("flow_A_%d", time.Now().UnixNano()))
	dataIdB := strToDataID(fmt.Sprintf("flow_B_%d", time.Now().UnixNano()))
	timeout := big.NewInt(3600)

	nonce, _ := client.PendingNonceAt(context.Background(), auth.From)

	authA := cloneAuth(auth, nonce)
	m.TxConfirmTimeLockA, m.GasLockA = sendAndWait(client, "LockA", func() (*types.Transaction, error) {
		return instA.Lock(authA, batchTo32Bytes(h1), dataIdA, timeout)
	})

	authB := cloneAuth(auth, nonce+1)
	m.TxConfirmTimeLockB, m.GasLockB = sendAndWait(client, "LockB", func() (*types.Transaction, error) {
		return instB.Lock(authB, batchTo32Bytes(h2), dataIdB, timeout)
	})

	// DSPB Prove
	challengeIdx := mathrand.Intn(len(chunkHashes))
	assignB := &circuit.AuditUnlockCircuit{
		ProofPath:    make([]frontend.Variable, auditH.depth),
		Helpers:      make([]frontend.Variable, auditH.depth),
		LeafCounts:   make([]frontend.Variable, auditH.depth),
		LeafNumBytes: make([]frontend.Variable, auditH.depth),
		Sn:           snII,
		ChunkIndex:   challengeIdx,
		ChunkHash:    chunkHashes[challengeIdx],
		H:            h2,
	}
	fillWitnessSlice(assignB, proofPaths[challengeIdx], helpers[challengeIdx], leafCounts, leafNumBytes)

	proofB, proveTimeB, err := auditH.Prove(assignB)
	if err != nil {
		fmt.Printf("❌ ProveB failed: %v\n", err)
		return nil
	}
	m.AuditUnlockProveTimeB = proveTimeB

	// Verify B
	verifyTimeB, err := auditH.Verify(proofB, challengeIdx, chunkHashes[challengeIdx], h2)
	if err != nil {
		fmt.Printf("❌ VerifyB failed: %v\n", err)
		return nil
	}
	m.AuditVerifyTimeB = verifyTimeB

	// SCB Unlock
	solProofB, pubB := formatProofForSolidityAuditUnlock(proofB, challengeIdx, chunkHashes[challengeIdx], h2)
	nonce, _ = client.PendingNonceAt(context.Background(), auth.From)
	authB = cloneAuth(auth, nonce)
	m.TxConfirmTimeUnlockB, m.GasUnlockB = sendAndWait(client, "AuditUnlockB", func() (*types.Transaction, error) {
		return instB.AuditUnlock(authB, solProofB, pubB)
	})

	// User Unlock Prove
	snI_calculated := new(big.Int).Xor(snII, zMask)
	snI_calculated.Mod(snI_calculated, bn254FieldModulus)
	proofA, proveTimeA, err := unlockH.Prove(preI, snI_calculated, h1)
	if err != nil {
		fmt.Printf("❌ ProveA failed: %v\n", err)
		return nil
	}
	m.UnlockProveTimeA = proveTimeA

	// Verify A
	verifyTimeA, err := unlockH.Verify(proofA, snI_calculated, h1)
	if err != nil {
		fmt.Printf("❌ VerifyA failed: %v\n", err)
		return nil
	}
	m.UnlockVerifyTimeA = verifyTimeA

	// SCA Unlock
	solProofA, pubA := formatProofForSolidityUnlock(proofA, h1, snI_calculated)
	nonce, _ = client.PendingNonceAt(context.Background(), auth.From)
	authA = cloneAuth(auth, nonce)
	m.TxConfirmTimeUnlockA, m.GasUnlockA = sendAndWait(client, "UnlockA", func() (*types.Transaction, error) {
		return instA.Unlock(authA, solProofA, pubA)
	})

	m.TotalGas = m.GasLockA + m.GasLockB + m.GasUnlockA + m.GasUnlockB
	m.EndToEndTime = time.Since(e2eStart)

	return m
}

// 统计学工具函数

func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	mid := len(values) / 2
	if len(values)%2 == 0 {
		return (values[mid-1] + values[mid]) / 2
	}
	return values[mid]
}
