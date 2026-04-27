package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"runtime"
	"runtime/debug"
	"time"
	"zk-htlc/actors"
	"zk-htlc/circuit"
	"zk-htlc/zkp"

	"github.com/consensys/gnark/frontend"
)

// BatchTestMetrics 批量测试结果
type BatchTestMetrics struct {
	BatchSize       int     `json:"batch_size"`
	MerkleDepth     int     `json:"merkle_depth"`
	SetupTimeMs     float64 `json:"setup_time_ms"`
	ProveTimeMs     float64 `json:"prove_time_ms"`
	VerifyTimeMs    float64 `json:"verify_time_ms"`
	ConstraintCount int     `json:"constraint_count"`
	ProofSizeBytes  int     `json:"proof_size_bytes"`
}

func main() {
	// 解析命令行参数
	batchMode := flag.Bool("batch", false, "批量测试模式（16/64/128/256）")
	singleBatch := flag.Int("size", 16, "单次测试的批量大小")
	output := flag.String("output", "results/batch_metrics.json", "输出文件路径")
	flag.Parse()

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║   批量锁定/解锁性能测试 (Batch Lock/Unlock Benchmark)   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝\n")

	if *batchMode {
		runBatchTests(*output)
	} else {
		metrics := runSingleTest(*singleBatch)
		printMetrics(metrics)
		exportMetrics([]*BatchTestMetrics{metrics}, *output)
	}

	fmt.Println("\n🎉 测试完成！")
}

// runBatchTests 批量测试模式（16/64/128/256）
func runBatchTests(outputPath string) {
	batchSizes := []int{16, 64, 128, 256}
	allMetrics := make([]*BatchTestMetrics, 0)

	for idx, size := range batchSizes {
		fmt.Printf("\n═══════════════════════════════════════════════════════════\n")
		fmt.Printf("    测试 [%d/%d]: 批量大小 = %d\n", idx+1, len(batchSizes), size)
		fmt.Printf("═══════════════════════════════════════════════════════════\n\n")

		// 预热 + 清理
		if idx == 0 {
			fmt.Println("🔥 Warm-up run...")
			_ = runSingleTest(8)
			runtime.GC()
			debug.FreeOSMemory()
			time.Sleep(1 * time.Second)
		}

		// 7 次测试取中位数
		const runs = 7
		var setupTimes, proveTimes, verifyTimes []time.Duration

		for run := 1; run <= runs; run++ {
			fmt.Printf("\n[Run %d/%d] Testing batch size %d...\n", run, runs, size)

			runtime.GC()
			debug.FreeOSMemory()
			time.Sleep(200 * time.Millisecond)

			metrics := runSingleTest(size)

			setupTimes = append(setupTimes, time.Duration(metrics.SetupTimeMs*1e6))
			proveTimes = append(proveTimes, time.Duration(metrics.ProveTimeMs*1e6))
			verifyTimes = append(verifyTimes, time.Duration(metrics.VerifyTimeMs*1e6))

			if run == 1 {
				// 保存第一次运行的固定参数
				allMetrics = append(allMetrics, metrics)
			}
		}

		// 更新为中位数
		medianMetrics := allMetrics[len(allMetrics)-1]
		medianMetrics.SetupTimeMs = medianDuration(setupTimes).Seconds() * 1000
		medianMetrics.ProveTimeMs = medianDuration(proveTimes).Seconds() * 1000
		medianMetrics.VerifyTimeMs = medianDuration(verifyTimes).Seconds() * 1000

		printMetrics(medianMetrics)

		if idx < len(batchSizes)-1 {
			fmt.Println("\n⏸️  Cooling down 2s...")
			time.Sleep(2 * time.Second)
		}
	}

	// 导出 JSON
	exportMetrics(allMetrics, outputPath)
	fmt.Printf("\n✅ 结果已保存到: %s\n", outputPath)
}

// runSingleTest 运行单次测试
func runSingleTest(batchSize int) *BatchTestMetrics {
	metrics := &BatchTestMetrics{
		BatchSize: batchSize,
	}

	// Step 1: Setup（初始化电路）
	fmt.Printf("\n[1/4] Setup for batch size %d...\n", batchSize)
	setupStart := time.Now()
	handler, err := zkp.NewBatchZKPHandler(batchSize)
	if err != nil {
		log.Fatalf("❌ Setup 失败: %v", err)
	}
	metrics.SetupTimeMs = time.Since(setupStart).Seconds() * 1000
	metrics.ConstraintCount = handler.GetConstraintCount()
	metrics.MerkleDepth = handler.GetMerkleDepth()

	fmt.Printf("   ✅ Setup: %.2f ms (约束: %d, 深度: %d)\n",
		metrics.SetupTimeMs, metrics.ConstraintCount, metrics.MerkleDepth)

	// Step 2: Operator 构建批量 Merkle 树
	fmt.Printf("\n[2/4] Building batch Merkle tree...\n")
	operator := actors.NewOperator()
	if err := operator.GenerateMockBatch(batchSize); err != nil {
		log.Fatalf("❌ 构建 Merkle 树失败: %v", err)
	}
	fmt.Printf("   ✅ Merkle tree built (root: %x...)\n", operator.MerkleRoot.Bytes()[:8])

	// Step 3: Prove（生成证明）
	fmt.Printf("\n[3/4] Generating proof...\n")

	// 获取第 0 个交易的 Merkle 路径
	tx := operator.TxLocks[0]
	proofElements, err := operator.GetProofPath(0)
	if err != nil {
		log.Fatalf("❌ 获取 Merkle 路径失败: %v", err)
	}

	// 🔍 调试：验证链下 Merkle 证明
	fmt.Printf("\n[DEBUG] 验证链下 Merkle 证明...\n")
	fmt.Printf("  - 叶子哈希 (tx.H): %x\n", tx.H.Bytes()[:8])
	fmt.Printf("  - Merkle 根: %x\n", operator.MerkleRoot.Bytes()[:8])
	fmt.Printf("  - 证明长度: %d\n", len(proofElements))
	fmt.Printf("  - 预期深度: %d\n", metrics.MerkleDepth)

	isValid := operator.MerkleTree.VerifyProof(tx.H, 0, proofElements)
	if !isValid {
		log.Fatal("❌ 链下 Merkle 证明验证失败！")
	}
	fmt.Println("  ✅ 链下验证通过")

	// 🔧 转换证明路径为电路格式
	proofPath := make([]frontend.Variable, len(proofElements))
	helpers := make([]frontend.Variable, len(proofElements))
	leafCounts := make([]frontend.Variable, len(proofElements))
	leafNumBytes := make([]frontend.Variable, len(proofElements))

	fmt.Println("\n[DEBUG] 证明路径详情:")
	for i, elem := range proofElements {
		// Hash 已经是 *big.Int 类型
		proofPath[i] = elem.Hash

		// 🔧 修复：方向标记转换
		// elem.IsLeft 表示兄弟节点是否在左侧
		// 如果兄弟在左，说明当前节点在右 → helper=1
		// 如果兄弟在右，说明当前节点在左 → helper=0
		if elem.IsLeft {
			helpers[i] = 1 // 当前节点在右侧
		} else {
			helpers[i] = 0 // 当前节点在左侧
		}

		// LeafCount 转换
		if elem.LeafCount != nil {
			leafCounts[i] = elem.LeafCount
		} else {
			leafCounts[i] = big.NewInt(1)
		}

		// 🔧 修复：直接使用 ProofElement 的 LeafNumBytes 字段（就是 leafCount）
		if elem.LeafNumBytes != nil {
			leafNumBytes[i] = elem.LeafNumBytes
		} else {
			// 如果为空，使用 LeafCount
			if elem.LeafCount != nil {
				leafNumBytes[i] = elem.LeafCount
			} else {
				leafNumBytes[i] = big.NewInt(32)
			}
		}

		fmt.Printf("  [%d] Hash: %x..., IsLeft: %v, LeafCount: %v, LeafNumBytes: %v\n",
			i, elem.Hash.Bytes()[:min(4, len(elem.Hash.Bytes()))], elem.IsLeft,
			elem.LeafCount, leafNumBytes[i])
	}

	// 🔧 如果证明路径长度小于电路深度，需要填充
	circuitDepth := handler.GetMerkleDepth()
	if len(proofElements) < circuitDepth {
		fmt.Printf("\n⚠️  证明长度(%d) < 电路深度(%d), 需要填充\n",
			len(proofElements), circuitDepth)

		// 使用零值填充
		for i := len(proofElements); i < circuitDepth; i++ {
			proofPath = append(proofPath, big.NewInt(0))
			helpers = append(helpers, big.NewInt(0))
			leafCounts = append(leafCounts, big.NewInt(1))
			leafNumBytes = append(leafNumBytes, big.NewInt(1))
		}
	}

	// 构造 witness
	assignment := &circuit.BatchUnlockCircuit{
		Preimage:           tx.Preimage,
		SerialNumber:       tx.SerialNumber,
		TxIndex:            big.NewInt(int64(tx.Index)),
		ProofPath:          proofPath,
		Helpers:            helpers,
		LeafCounts:         leafCounts,
		LeafNumBytes:       leafNumBytes,
		MerkleRoot:         operator.MerkleRoot,
		SerialNumberPublic: tx.SerialNumber,
	}

	// 🔍 调试：打印关键 witness 数据
	fmt.Println("\n[DEBUG] Witness 数据:")
	fmt.Printf("  - Preimage: %x\n", tx.Preimage.Bytes()[:8])
	fmt.Printf("  - SerialNumber: %x\n", tx.SerialNumber.Bytes()[:8])
	fmt.Printf("  - TxIndex: %d\n", tx.Index)
	fmt.Printf("  - MerkleRoot: %x\n", operator.MerkleRoot.Bytes()[:8])
	fmt.Printf("  - ProofPath length: %d\n", len(proofPath))
	fmt.Printf("  - Helpers length: %d\n", len(helpers))

	proveStart := time.Now()
	proof, err := handler.Prove(assignment)
	if err != nil {
		log.Fatalf("❌ Prove 失败: %v", err)
	}
	metrics.ProveTimeMs = time.Since(proveStart).Seconds() * 1000

	// 计算证明大小
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	metrics.ProofSizeBytes = buf.Len()

	fmt.Printf("   ✅ Prove: %.2f ms (大小: %d bytes)\n",
		metrics.ProveTimeMs, metrics.ProofSizeBytes)

	// Step 4: Verify（验证证明）
	fmt.Printf("\n[4/4] Verifying proof...\n")
	verifyStart := time.Now()
	err = handler.Verify(proof, []*big.Int{operator.MerkleRoot, tx.SerialNumber})
	if err != nil {
		log.Fatalf("❌ Verify 失败: %v", err)
	}
	metrics.VerifyTimeMs = time.Since(verifyStart).Seconds() * 1000

	fmt.Printf("   ✅ Verify: %.2f ms\n", metrics.VerifyTimeMs)

	return metrics
}

// printMetrics 打印性能指标
func printMetrics(m *BatchTestMetrics) {
	fmt.Println("\n┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│              性能指标 (Performance Metrics)             │")
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Printf("│  批量大小 (Batch Size):        %20d      │\n", m.BatchSize)
	fmt.Printf("│  Merkle 深度 (Depth):          %20d      │\n", m.MerkleDepth)
	fmt.Printf("│  约束数量 (Constraints):       %20d      │\n", m.ConstraintCount)
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Printf("│  ⏱️  Setup 时间:                %20.2f ms   │\n", m.SetupTimeMs)
	fmt.Printf("│  ⏱️  Prove 时间:                %20.2f ms   │\n", m.ProveTimeMs)
	fmt.Printf("│  ⏱️  Verify 时间:               %20.2f ms   │\n", m.VerifyTimeMs)
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	fmt.Printf("│  📦 证明大小:                   %20d B    │\n", m.ProofSizeBytes)
	fmt.Println("└─────────────────────────────────────────────────────────┘")
}

// exportMetrics 导出 JSON
func exportMetrics(allMetrics []*BatchTestMetrics, filepath string) error {
	data, err := json.MarshalIndent(allMetrics, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, data, 0644)
}

// medianDuration 计算中位数
func medianDuration(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(values))
	copy(sorted, values)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted[len(sorted)/2]
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
