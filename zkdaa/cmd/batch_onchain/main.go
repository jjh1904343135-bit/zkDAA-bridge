package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
	"zk-htlc/actors"
	"zk-htlc/circuit"
	"zk-htlc/contracts"
	"zk-htlc/zkp"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SetupMetricsCache 从 setup_metrics.json 加载的 Setup 数据
type SetupMetricsCache struct {
	BatchSize       int     `json:"batch_size"`
	MerkleDepth     int     `json:"merkle_depth"`
	ConstraintCount int     `json:"constraint_count"`
	CompileTimeMs   float64 `json:"compile_time_ms"`
	SetupTimeMs     float64 `json:"setup_time_ms"`
	TotalTimeMs     float64 `json:"total_time_ms"`
	PKSizeBytes     int64   `json:"pk_size_bytes"`
	VKSizeBytes     int64   `json:"vk_size_bytes"`
	Timestamp       string  `json:"timestamp"`
}

// BatchTestMetrics 批量测试性能指标
type BatchTestMetrics struct {
	BatchSize       int `json:"batch_size"`
	MerkleDepth     int `json:"merkle_depth"`
	ConstraintCount int `json:"constraint_count"`

	// Setup 阶段时间（从缓存加载）
	CompileTimeMs  float64 `json:"compile_time_ms"`
	SetupTimeMs    float64 `json:"setup_time_ms"`
	LoadKeysTimeMs float64 `json:"load_keys_time_ms"`

	// 执行阶段时间（运行时测量）
	BuildTreeTimeMs float64 `json:"build_tree_time_ms"`
	ProveTimeMs     float64 `json:"prove_time_ms"`
	VerifyTimeMs    float64 `json:"verify_time_ms"`

	SubmitRootGas  uint64 `json:"submit_root_gas"`
	UnlockGas      uint64 `json:"unlock_gas"`
	ProofSizeBytes int    `json:"proof_size_bytes"`
	Success        bool   `json:"success"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

// 全局变量：Setup 指标缓存
var setupMetricsMap = make(map[int]*SetupMetricsCache)

func main() {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║   批量锁定/解锁 - 真实链上验证实验 (完整版)              ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝\n")

	// 加载 Setup 指标缓存
	if err := loadSetupMetrics(); err != nil {
		log.Printf("⚠️  加载 Setup 指标失败: %v\n", err)
		fmt.Println("   💡 提示: 请先运行 setup_generator.go 生成 Setup 数据")
	} else {
		fmt.Printf("✅ 已加载 %d 个批量大小的 Setup 指标\n", len(setupMetricsMap))
		for batchSize := range setupMetricsMap {
			fmt.Printf("   - 批量 %d: Setup 时间 %.2f 秒\n",
				batchSize, setupMetricsMap[batchSize].SetupTimeMs/1000.0)
		}
		fmt.Println()
	}

	// 解析命令行参数
	batchMode := false
	if len(os.Args) > 1 && os.Args[1] == "-batch" {
		batchMode = true
	}

	// 连接区块链
	fmt.Println("🔗 连接到区块链节点...")
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("❌ 连接区块链失败: %v\n", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 ChainID 失败: %v", err)
	}
	fmt.Printf("✅ 已连接到区块链 (ChainID: %s)\n", chainID.String())

	// 设置账户
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatalf("❌ 解析私钥失败: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("❌ 创建交易器失败: %v", err)
	}

	// 连接合约
	contractAddr := common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9")
	fmt.Printf("📄 连接到合约: %s\n", contractAddr.Hex())

	instance, err := contracts.NewBatchDataMigration(contractAddr, client)
	if err != nil {
		log.Fatalf("❌ 连接合约失败: %v\n", err)
	}
	fmt.Println("✅ 合约连接成功\n")

	if batchMode {
		runBatchTests(auth, instance, client)
	} else {
		fmt.Println("🔍 单次测试模式 (批量大小: 16)")
		fmt.Println("💡 提示: 使用 -batch 参数运行完整性能测试\n")
		metrics := testBatchSize(16, auth, instance, client, false)
		printMetrics(metrics)
	}
}

func loadSetupMetrics() error {
	data, err := os.ReadFile("zkp/setup_metrics.json")
	if err != nil {
		return err
	}

	var metrics []SetupMetricsCache
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	for i := range metrics {
		setupMetricsMap[metrics[i].BatchSize] = &metrics[i]
	}

	return nil
}

func runBatchTests(auth *bind.TransactOpts, instance *contracts.BatchDataMigration, client *ethclient.Client) {
	batchSizes := []int{16, 64, 128, 256}
	runsPerBatch := 7
	allMetrics := make([]*BatchTestMetrics, 0)

	fmt.Println("🚀 批量性能测试模式")
	fmt.Printf("测试配置: %v\n", batchSizes)
	fmt.Printf("每个批量运行次数: %d (取中位数)\n\n", runsPerBatch)

	for i, batchSize := range batchSizes {
		fmt.Printf("\n╔═══════════════════════════════════════════════════════════╗\n")
		fmt.Printf("║ [%d/%d] 批量大小: %-4d                                     ║\n", i+1, len(batchSizes), batchSize)
		fmt.Printf("╚═══════════════════════════════════════════════════════════╝\n\n")

		runs := make([]*BatchTestMetrics, 0)

		for run := 1; run <= runsPerBatch; run++ {
			fmt.Printf("▶ 运行 %d/%d...\n", run, runsPerBatch)
			metrics := testBatchSize(batchSize, auth, instance, client, true)
			runs = append(runs, metrics)

			if metrics.Success {
				fmt.Printf("  ✅ Prove: %.2fms, Verify: %.2fms, Gas: %d\n",
					metrics.ProveTimeMs, metrics.VerifyTimeMs, metrics.UnlockGas)
			} else {
				fmt.Printf("  ❌ 失败: %s\n", metrics.ErrorMessage)
			}

			if run < runsPerBatch {
				time.Sleep(1 * time.Second)
			}
		}

		medianMetrics := calculateMedian(runs)
		allMetrics = append(allMetrics, medianMetrics)

		fmt.Println("\n📊 中位数结果:")
		printMetrics(medianMetrics)

		if i < len(batchSizes)-1 {
			fmt.Println("\n⏳ 等待 3 秒后继续...")
			time.Sleep(3 * time.Second)
		}
	}

	printSummary(allMetrics)
	saveResults(allMetrics)
}

func testBatchSize(
	batchSize int,
	auth *bind.TransactOpts,
	instance *contracts.BatchDataMigration,
	client *ethclient.Client,
	quietMode bool,
) *BatchTestMetrics {
	metrics := &BatchTestMetrics{
		BatchSize: batchSize,
	}

	// 从缓存加载 Setup 指标
	if cached, ok := setupMetricsMap[batchSize]; ok {
		metrics.MerkleDepth = cached.MerkleDepth
		metrics.ConstraintCount = cached.ConstraintCount
		metrics.CompileTimeMs = cached.CompileTimeMs
		metrics.SetupTimeMs = cached.SetupTimeMs
	}

	if !quietMode {
		fmt.Println("\n[1/5] Loading ZKP keys...")
	}

	loadKeysStart := time.Now()
	handler, err := zkp.NewBatchZKPHandler(batchSize)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("加载密钥失败: %v", err)
		return metrics
	}
	metrics.LoadKeysTimeMs = float64(time.Since(loadKeysStart).Microseconds()) / 1000.0

	if metrics.ConstraintCount == 0 {
		metrics.ConstraintCount = handler.GetConstraintCount()
		metrics.MerkleDepth = handler.GetMerkleDepth()
	}

	if !quietMode {
		fmt.Printf("   ✅ 约束数量: %d\n", metrics.ConstraintCount)
		fmt.Printf("   ✅ Merkle 深度: %d\n", metrics.MerkleDepth)
		if metrics.SetupTimeMs > 0 {
			fmt.Printf("   📋 Setup 时间 (缓存): %.2f 秒\n", metrics.SetupTimeMs/1000.0)
		}
		fmt.Printf("   ⏱️  密钥加载: %.2f ms\n", metrics.LoadKeysTimeMs)
	}

	if !quietMode {
		fmt.Println("\n[2/5] Building Merkle tree...")
	}

	buildTreeStart := time.Now()
	operator := actors.NewOperator()
	if err := operator.GenerateMockBatch(batchSize); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("构建 Merkle 树失败: %v", err)
		return metrics
	}
	metrics.BuildTreeTimeMs = float64(time.Since(buildTreeStart).Microseconds()) / 1000.0

	tx0 := operator.TxLocks[0]
	if !quietMode {
		fmt.Printf("   ✅ 生成 %d 笔交易\n", batchSize)
		fmt.Printf("   ⏱️  构建 Merkle 树: %.2f ms\n", metrics.BuildTreeTimeMs)
		fmt.Printf("   📍 Merkle Root: %s...\n", operator.MerkleRoot.String()[:20])
	}

	if !quietMode {
		fmt.Println("\n[3/5] Submitting batch root...")
	}
	rootBytes := operator.MerkleRoot.Bytes()
	var root32 [32]byte
	copy(root32[32-len(rootBytes):], rootBytes)

	txSubmit, err := instance.SubmitBatchRoot(auth, big.NewInt(int64(batchSize)), root32)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("提交 Root 失败: %v", err)
		return metrics
	}

	receiptSubmit, err := bind.WaitMined(context.Background(), client, txSubmit)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("等待 Root 确认失败: %v", err)
		return metrics
	}
	metrics.SubmitRootGas = receiptSubmit.GasUsed
	if !quietMode {
		fmt.Printf("   ✅ Gas: %d\n", metrics.SubmitRootGas)
	}

	if !quietMode {
		fmt.Println("\n[4/5] Generating ZKP proof...")
	}

	proofElements, err := operator.GetProofPath(0)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("获取 Merkle 路径失败: %v", err)
		return metrics
	}

	proofPath := make([]frontend.Variable, len(proofElements))
	helpers := make([]frontend.Variable, len(proofElements))
	leafCounts := make([]frontend.Variable, len(proofElements))
	leafNumBytes := make([]frontend.Variable, len(proofElements))

	for i, elem := range proofElements {
		proofPath[i] = elem.Hash
		if elem.IsLeft {
			helpers[i] = 1
		} else {
			helpers[i] = 0
		}
		leafCounts[i] = elem.LeafCount
		leafNumBytes[i] = elem.LeafNumBytes
	}

	circuitDepth := handler.GetMerkleDepth()
	for i := len(proofElements); i < circuitDepth; i++ {
		proofPath = append(proofPath, big.NewInt(0))
		helpers = append(helpers, big.NewInt(0))
		leafCounts = append(leafCounts, big.NewInt(1))
		leafNumBytes = append(leafNumBytes, big.NewInt(1))
	}

	assignment := &circuit.BatchUnlockCircuit{
		Preimage:           tx0.Preimage,
		SerialNumber:       tx0.SerialNumber,
		TxIndex:            big.NewInt(int64(tx0.Index)),
		ProofPath:          proofPath,
		Helpers:            helpers,
		LeafCounts:         leafCounts,
		LeafNumBytes:       leafNumBytes,
		MerkleRoot:         operator.MerkleRoot,
		SerialNumberPublic: tx0.SerialNumber,
	}

	proveStart := time.Now()
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("创建 witness 失败: %v", err)
		return metrics
	}

	publicWitness, err := witness.Public()
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("提取公开 witness 失败: %v", err)
		return metrics
	}

	proof, err := handler.Prove(assignment)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("生成证明失败: %v", err)
		return metrics
	}
	metrics.ProveTimeMs = float64(time.Since(proveStart).Microseconds()) / 1000.0

	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	metrics.ProofSizeBytes = buf.Len()
	if !quietMode {
		fmt.Printf("   ✅ 证明大小: %d bytes\n", metrics.ProofSizeBytes)
		fmt.Printf("   ⏱️  Prove 时间: %.2f ms\n", metrics.ProveTimeMs)
	}

	verifyStart := time.Now()
	publicInputs := []*big.Int{operator.MerkleRoot, tx0.SerialNumber}
	err = handler.Verify(proof, publicInputs)
	metrics.VerifyTimeMs = float64(time.Since(verifyStart).Microseconds()) / 1000.0
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("链下验证失败: %v", err)
		return metrics
	}
	if !quietMode {
		fmt.Printf("   ✅ 链下验证通过 (%.2f ms)\n", metrics.VerifyTimeMs)
	}

	if !quietMode {
		fmt.Println("\n[5/5] Submitting to blockchain...")
	}

	publicVector := publicWitness.Vector().(fr.Vector)
	solProof, solPubInputs := formatProofForSolidity(proof, publicVector)

	txUnlock, err := instance.Unlock(auth, big.NewInt(int64(batchSize)), solProof, solPubInputs)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("Unlock 失败: %v", err)
		return metrics
	}

	receiptUnlock, err := bind.WaitMined(context.Background(), client, txUnlock)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = fmt.Sprintf("等待 Unlock 确认失败: %v", err)
		return metrics
	}

	metrics.UnlockGas = receiptUnlock.GasUsed
	metrics.Success = true
	if !quietMode {
		fmt.Printf("   ✅ 链上验证通过! Gas: %d\n", metrics.UnlockGas)
	}

	return metrics
}

func formatProofForSolidity(proof groth16.Proof, publicVector fr.Vector) ([8]*big.Int, [2]*big.Int) {
	proofBN254, ok := proof.(*groth16_bn254.Proof)
	if !ok {
		log.Fatal("❌ 证明类型错误")
	}

	var solProof [8]*big.Int

	solProof[0] = new(big.Int)
	solProof[1] = new(big.Int)
	proofBN254.Ar.X.BigInt(solProof[0])
	proofBN254.Ar.Y.BigInt(solProof[1])

	solProof[2] = new(big.Int)
	solProof[3] = new(big.Int)
	solProof[4] = new(big.Int)
	solProof[5] = new(big.Int)
	proofBN254.Bs.X.A1.BigInt(solProof[2])
	proofBN254.Bs.X.A0.BigInt(solProof[3])
	proofBN254.Bs.Y.A1.BigInt(solProof[4])
	proofBN254.Bs.Y.A0.BigInt(solProof[5])

	solProof[6] = new(big.Int)
	solProof[7] = new(big.Int)
	proofBN254.Krs.X.BigInt(solProof[6])
	proofBN254.Krs.Y.BigInt(solProof[7])

	var solPublicInputs [2]*big.Int

	if len(publicVector) != 2 {
		log.Fatalf("❌ 公开输入向量长度错误: 期望 2, 实际 %d", len(publicVector))
	}

	modulus := fr.Modulus()

	solPublicInputs[0] = new(big.Int)
	solPublicInputs[1] = new(big.Int)
	publicVector[0].BigInt(solPublicInputs[0])
	publicVector[1].BigInt(solPublicInputs[1])

	solPublicInputs[0].Mod(solPublicInputs[0], modulus)
	solPublicInputs[1].Mod(solPublicInputs[1], modulus)

	return solProof, solPublicInputs
}

func printMetrics(m *BatchTestMetrics) {
	fmt.Println("\n┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│              性能指标 (Performance Metrics)             │")
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Printf("│  批量大小 (Batch Size):        %20d      │\n", m.BatchSize)
	fmt.Printf("│  Merkle 深度 (Depth):          %20d      │\n", m.MerkleDepth)
	fmt.Printf("│  约束数量 (Constraints):       %20d      │\n", m.ConstraintCount)
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Println("│  🔧 Setup 阶段 (一次性成本)                            │")
	if m.SetupTimeMs > 0 {
		fmt.Printf("│    Setup (可信设置):           %20.2f 秒  │\n", m.SetupTimeMs/1000.0)
	}
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Println("│  ⚡ 执行阶段 (每次运行)                                 │")
	fmt.Printf("│    密钥加载:                   %20.2f ms   │\n", m.LoadKeysTimeMs)
	fmt.Printf("│    构建 Merkle 树:             %20.2f ms   │\n", m.BuildTreeTimeMs)
	fmt.Printf("│    证明生成 (Prove):           %20.2f ms   │\n", m.ProveTimeMs)
	fmt.Printf("│    链下验证 (Verify):          %20.2f ms   │\n", m.VerifyTimeMs)
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Println("│  ⛽ Gas 消耗                                            │")
	fmt.Printf("│    提交 Root:                  %20d      │\n", m.SubmitRootGas)
	fmt.Printf("│    Unlock:                     %20d      │\n", m.UnlockGas)
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Printf("│  📦 证明大小:                   %20d B    │\n", m.ProofSizeBytes)
	if m.Success {
		fmt.Println("│  ✅ 状态:                                     成功     │")
	} else {
		fmt.Println("│  ❌ 状态:                                     失败     │")
		fmt.Printf("│  错误: %-50s │\n", truncate(m.ErrorMessage, 50))
	}
	fmt.Println("└─────────────────────────────────────────────────────────┘")
}

func printSummary(allMetrics []*BatchTestMetrics) {
	fmt.Println("\n\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                   测试汇总 (Summary)                      ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝\n")

	fmt.Println("┌──────────┬─────────┬────────────┬────────────┬───────────┬──────────┬────────────┐")
	fmt.Println("│ 批量大小 │ 约束数  │ Setup(秒)  │ Prove(ms)  │ Verify(ms)│ Gas      │ 状态       │")
	fmt.Println("├──────────┼─────────┼────────────┼────────────┼───────────┼──────────┼────────────┤")

	for _, m := range allMetrics {
		status := "✅"
		if !m.Success {
			status = "❌"
		}
		setupSec := m.SetupTimeMs / 1000.0
		fmt.Printf("│ %-8d │ %-7d │ %10.2f │ %10.2f │ %9.2f │ %-8d │ %-10s │\n",
			m.BatchSize, m.ConstraintCount, setupSec, m.ProveTimeMs, m.VerifyTimeMs, m.UnlockGas, status)
	}

	fmt.Println("└──────────┴─────────┴────────────┴────────────┴───────────┴──────────┴────────────┘")

	successCount := 0
	for _, m := range allMetrics {
		if m.Success {
			successCount++
		}
	}
	fmt.Printf("\n📊 成功率: %d/%d (%.1f%%)\n", successCount, len(allMetrics),
		float64(successCount)/float64(len(allMetrics))*100)
}

func saveResults(allMetrics []*BatchTestMetrics) {
	data, err := json.MarshalIndent(allMetrics, "", "  ")
	if err != nil {
		fmt.Printf("⚠️  保存结果失败: %v\n", err)
		return
	}

	filename := fmt.Sprintf("results\batch_test_results_%d.json", time.Now().Unix())
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("⚠️  写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("\n💾 测试结果已保存到: %s\n", filename)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func calculateMedian(runs []*BatchTestMetrics) *BatchTestMetrics {
	if len(runs) == 0 {
		return nil
	}

	successRuns := make([]*BatchTestMetrics, 0)
	for _, r := range runs {
		if r.Success {
			successRuns = append(successRuns, r)
		}
	}

	if len(successRuns) == 0 {
		return runs[0]
	}

	result := &BatchTestMetrics{
		BatchSize:       successRuns[0].BatchSize,
		MerkleDepth:     successRuns[0].MerkleDepth,
		ConstraintCount: successRuns[0].ConstraintCount,
		ProofSizeBytes:  successRuns[0].ProofSizeBytes,
		Success:         true,
		CompileTimeMs:   successRuns[0].CompileTimeMs,
		SetupTimeMs:     successRuns[0].SetupTimeMs,
	}

	loadKeysTimes := make([]float64, len(successRuns))
	buildTreeTimes := make([]float64, len(successRuns))
	proveTimes := make([]float64, len(successRuns))
	verifyTimes := make([]float64, len(successRuns))
	submitRootGas := make([]uint64, len(successRuns))
	unlockGas := make([]uint64, len(successRuns))

	for i, r := range successRuns {
		loadKeysTimes[i] = r.LoadKeysTimeMs
		buildTreeTimes[i] = r.BuildTreeTimeMs
		proveTimes[i] = r.ProveTimeMs
		verifyTimes[i] = r.VerifyTimeMs
		submitRootGas[i] = r.SubmitRootGas
		unlockGas[i] = r.UnlockGas
	}

	result.LoadKeysTimeMs = medianFloat64(loadKeysTimes)
	result.BuildTreeTimeMs = medianFloat64(buildTreeTimes)
	result.ProveTimeMs = medianFloat64(proveTimes)
	result.VerifyTimeMs = medianFloat64(verifyTimes)
	result.SubmitRootGas = medianUint64(submitRootGas)
	result.UnlockGas = medianUint64(unlockGas)

	return result
}

func medianFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2.0
	}
	return sorted[mid]
}

func medianUint64(values []uint64) uint64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]uint64, len(values))
	copy(sorted, values)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}
