package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type SetupMetrics struct {
	BatchSize       int     `json:"batch_size"`
	MerkleDepth     int     `json:"merkle_depth"`
	ConstraintCount int     `json:"constraint_count"`
	CompileTimeMs   float64 `json:"compile_time_ms"`
	SetupTimeMs     float64 `json:"setup_time_ms"`
	SetupMinMs      float64 `json:"setup_min_ms"`
	SetupMaxMs      float64 `json:"setup_max_ms"`
	SetupStdDevMs   float64 `json:"setup_std_dev_ms"`
	TotalTimeMs     float64 `json:"total_time_ms"`
	PKSizeBytes     int64   `json:"pk_size_bytes"`
	VKSizeBytes     int64   `json:"vk_size_bytes"`
	Timestamp       string  `json:"timestamp"`
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║     为所有批量大小生成 Setup 和验证器合约                ║")
	fmt.Println("║     单个 HTLC 验证 + 不同深度的 Merkle 树                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝\n")

	// 优化运行时设置
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

	// 预热 GC，减少性能波动
	runtime.GC()
	time.Sleep(time.Second)

	batchSizes := []int{16, 64, 128, 256}
	allMetrics := make([]*SetupMetrics, 0)

	for i, batchSize := range batchSizes {
		fmt.Printf("\n[%d/%d] ========== 批量大小: %d ==========\n", i+1, len(batchSizes), batchSize)

		// 每次 Setup 前强制 GC，减少波动
		runtime.GC()
		time.Sleep(500 * time.Millisecond) // 增加到500ms

		metrics, err := generateSetupForBatchSize(batchSize)
		if err != nil {
			log.Printf("❌ 批量 %d 失败: %v\n", batchSize, err)
			continue
		}

		allMetrics = append(allMetrics, metrics)

		fmt.Printf("✅ 批量 %d 完成 (总耗时: %.2f 秒)\n", batchSize, metrics.TotalTimeMs/1000.0)
	}

	saveSetupMetrics(allMetrics)

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    全部完成！                             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝\n")

	printSummary(allMetrics)
}

func generateSetupForBatchSize(batchSize int) (*SetupMetrics, error) {
	totalStart := time.Now()

	depth := calculateDepth(batchSize)
	fmt.Printf("批量大小: %d, Merkle 深度: %d\n", batchSize, depth)

	// 你的原电路设计是正确的！只验证单个 HTLC
	dummyCircuit := circuit.BatchUnlockCircuit{
		ProofPath:    make([]frontend.Variable, depth),
		Helpers:      make([]frontend.Variable, depth),
		LeafCounts:   make([]frontend.Variable, depth),
		LeafNumBytes: make([]frontend.Variable, depth),
	}

	// 编译电路
	fmt.Println("  📝 编译电路...")
	compileStart := time.Now()
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &dummyCircuit)
	if err != nil {
		return nil, fmt.Errorf("编译失败: %w", err)
	}
	compileTime := time.Since(compileStart)
	constraintCount := ccs.GetNbConstraints()
	fmt.Printf("  ✅ 约束数量: %d (耗时: %.2f ms)\n", constraintCount, float64(compileTime.Milliseconds()))

	// 强制 GC 后再 Setup
	runtime.GC()
	time.Sleep(500 * time.Millisecond) // 增加到500ms

	// Setup - 预热 + 多次运行 + 去除异常值
	fmt.Println("  🔑 执行 Setup (预热 + 多次测试 + 去除异常值)...")

	// 预热运行 2 次（不计入统计）
	fmt.Println("     [预热阶段]")
	for warmup := 0; warmup < 2; warmup++ {
		runtime.GC()
		time.Sleep(200 * time.Millisecond)
		_, _, _ = groth16.Setup(ccs)
		fmt.Printf("     预热 %d/2 完成\n", warmup+1)
	}

	// 正式测量 7 次
	fmt.Println("     [正式测量]")
	const numRuns = 7 // 增加到7次，便于去除异常值
	setupTimes := make([]time.Duration, numRuns)

	var pk groth16.ProvingKey
	var vk groth16.VerifyingKey

	for run := 0; run < numRuns; run++ {
		runtime.GC()
		time.Sleep(200 * time.Millisecond)

		setupStart := time.Now()
		pk, vk, err = groth16.Setup(ccs)
		if err != nil {
			return nil, fmt.Errorf("Setup 失败: %w", err)
		}
		setupTimes[run] = time.Since(setupStart)
		fmt.Printf("     第 %d/%d 次: %.2f 秒\n", run+1, numRuns, setupTimes[run].Seconds())
	}

	// 计算统计数据（去除最大最小值）
	trimmedMean := calculateTrimmedMean(setupTimes)
	minTime, maxTime := getMinMax(setupTimes)
	stdDev := calculateStdDev(setupTimes)

	fmt.Printf("\n  ✅ Setup 性能统计:\n")
	fmt.Printf("     去除异常值后平均: %.2f 秒\n", trimmedMean.Seconds())
	fmt.Printf("     范围: %.2f - %.2f 秒 (±%.1f%%)\n",
		minTime.Seconds(), maxTime.Seconds(),
		(maxTime.Seconds()-minTime.Seconds())/trimmedMean.Seconds()*100)
	fmt.Printf("     标准差: %.2f 秒\n", stdDev.Seconds())

	// 保存密钥
	pkPath, vkPath, err := saveKeys(batchSize, pk, vk)
	if err != nil {
		return nil, err
	}

	pkInfo, _ := os.Stat(pkPath)
	vkInfo, _ := os.Stat(vkPath)

	// 导出验证器合约
	if err := exportVerifierContract(batchSize, vk); err != nil {
		return nil, err
	}

	totalTime := time.Since(totalStart)

	return &SetupMetrics{
		BatchSize:       batchSize,
		MerkleDepth:     depth,
		ConstraintCount: constraintCount,
		CompileTimeMs:   float64(compileTime.Milliseconds()),
		SetupTimeMs:     float64(trimmedMean.Milliseconds()), // 使用去除异常值后的平均值
		SetupMinMs:      float64(minTime.Milliseconds()),
		SetupMaxMs:      float64(maxTime.Milliseconds()),
		SetupStdDevMs:   float64(stdDev.Milliseconds()),
		TotalTimeMs:     float64(totalTime.Milliseconds()),
		PKSizeBytes:     pkInfo.Size(),
		VKSizeBytes:     vkInfo.Size(),
		Timestamp:       time.Now().Format(time.RFC3339),
	}, nil
}

func saveKeys(batchSize int, pk groth16.ProvingKey, vk groth16.VerifyingKey) (string, string, error) {
	fmt.Println("  💾 保存 Setup 密钥...")

	pkPath := fmt.Sprintf("zkp/batch_pk_%d.bin", batchSize)
	pkFile, err := os.Create(pkPath)
	if err != nil {
		return "", "", fmt.Errorf("创建 PK 文件失败: %w", err)
	}
	defer pkFile.Close()

	if _, err = pk.WriteTo(pkFile); err != nil {
		return "", "", fmt.Errorf("写入 PK 失败: %w", err)
	}
	pkInfo, _ := os.Stat(pkPath)
	fmt.Printf("     ✅ Proving Key: %s (%.2f MB)\n", pkPath, float64(pkInfo.Size())/1024/1024)

	vkPath := fmt.Sprintf("zkp/batch_vk_%d.bin", batchSize)
	vkFile, err := os.Create(vkPath)
	if err != nil {
		return "", "", fmt.Errorf("创建 VK 文件失败: %w", err)
	}
	defer vkFile.Close()

	if _, err = vk.WriteTo(vkFile); err != nil {
		return "", "", fmt.Errorf("写入 VK 失败: %w", err)
	}
	vkInfo, _ := os.Stat(vkPath)
	fmt.Printf("     ✅ Verifying Key: %s (%.2f KB)\n", vkPath, float64(vkInfo.Size())/1024)

	return pkPath, vkPath, nil
}

func exportVerifierContract(batchSize int, vk groth16.VerifyingKey) error {
	fmt.Println("  📤 导出验证器合约...")

	tempFile, err := os.CreateTemp("", "verifier_*.sol")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if err := vk.ExportSolidity(tempFile); err != nil {
		tempFile.Close()
		return fmt.Errorf("导出 Solidity 失败: %w", err)
	}
	tempFile.Close()

	content, err := os.ReadFile(tempPath)
	if err != nil {
		return fmt.Errorf("读取临时文件失败: %w", err)
	}

	modifiedContent := strings.ReplaceAll(string(content),
		"contract Groth16Verifier",
		fmt.Sprintf("contract BatchUnlockVerifier%d", batchSize))

	solPath := fmt.Sprintf("blockchain-contracts/contracts/BatchUnlockVerifier_%d.sol", batchSize)
	if err := os.WriteFile(solPath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("写入最终文件失败: %w", err)
	}

	solInfo, _ := os.Stat(solPath)
	fmt.Printf("     ✅ 验证器合约: %s (%.2f KB)\n", solPath, float64(solInfo.Size())/1024)
	fmt.Printf("     📝 合约名称: BatchUnlockVerifier%d\n", batchSize)

	return nil
}

func saveSetupMetrics(metrics []*SetupMetrics) {
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		fmt.Printf("⚠️  保存 Setup 指标失败: %v\n", err)
		return
	}

	filename := "zkp/setup_metrics.json"
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("⚠️  写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("\n💾 Setup 指标已保存到: %s\n", filename)
}

func printSummary(metrics []*SetupMetrics) {
	fmt.Println("📁 生成的文件汇总:")
	fmt.Println("┌──────────┬─────────┬────────────┬────────────┬──────────┬─────────────┬─────────────┐")
	fmt.Println("│ 批量大小 │ 约束数  │ 编译(ms)   │ Setup(秒)  │ 标准差(s)│ PK(MB)      │ VK(KB)      │")
	fmt.Println("├──────────┼─────────┼────────────┼────────────┼──────────┼─────────────┼─────────────┤")

	for _, m := range metrics {
		fmt.Printf("│ %-8d │ %-7d │ %10.2f │ %10.2f │ %8.2f │ %11.2f │ %11.2f │\n",
			m.BatchSize,
			m.ConstraintCount,
			m.CompileTimeMs,
			m.SetupTimeMs/1000.0,
			m.SetupStdDevMs/1000.0,
			float64(m.PKSizeBytes)/1024/1024,
			float64(m.VKSizeBytes)/1024)
	}

	fmt.Println("└──────────┴─────────┴────────────┴────────────┴──────────┴─────────────┴─────────────┘")

	// 分析约束增长（基于 Merkle 深度）
	fmt.Println("\n🔍 Merkle 深度与约束关系分析:")
	for i := 1; i < len(metrics); i++ {
		prev := metrics[i-1]
		curr := metrics[i]

		depthIncrease := curr.MerkleDepth - prev.MerkleDepth
		constraintIncrease := curr.ConstraintCount - prev.ConstraintCount
		constraintsPerLevel := float64(constraintIncrease) / float64(depthIncrease)

		setupRatio := curr.SetupTimeMs / prev.SetupTimeMs

		fmt.Printf("   批量 %d→%d: 深度 +%d 层, 约束 +%d (≈%.0f 约束/层), Setup %.2fx\n",
			prev.BatchSize, curr.BatchSize,
			depthIncrease, constraintIncrease,
			constraintsPerLevel, setupRatio)
	}

	// 计算平均每层约束数
	if len(metrics) > 1 {
		totalConstraintIncrease := metrics[len(metrics)-1].ConstraintCount - metrics[0].ConstraintCount
		totalDepthIncrease := metrics[len(metrics)-1].MerkleDepth - metrics[0].MerkleDepth
		avgConstraintsPerLevel := float64(totalConstraintIncrease) / float64(totalDepthIncrease)

		fmt.Printf("\n   📊 平均每层 Merkle 验证约束: ≈%.0f 个\n", avgConstraintsPerLevel)
		fmt.Printf("   💡 这是正常的！单个 HTLC 验证只需要沿 Merkle 路径验证\n")
	}

	fmt.Println("\n📋 生成的验证器合约:")
	for _, m := range metrics {
		fmt.Printf("   - BatchUnlockVerifier%d.sol (Merkle 深度: %d)\n", m.BatchSize, m.MerkleDepth)
	}

	fmt.Println("\n⚠️  后续步骤:")
	fmt.Println("   1. Setup 密钥已保存到 zkp/ 目录")
	fmt.Println("   2. 验证器合约已生成")
	fmt.Println("   3. 编译合约: cd blockchain-contracts && npx hardhat compile")
	fmt.Println("   4. 部署合约: npx hardhat run scripts/deploy_multi_verifier.js --network localhost")
}

func calculateDepth(leafCount int) int {
	depth := 0
	n := leafCount
	for n > 1 {
		n = (n + 1) / 2
		depth++
	}
	return depth
}

// calculateTrimmedMean 计算去除最大最小值后的平均值
func calculateTrimmedMean(times []time.Duration) time.Duration {
	if len(times) < 3 {
		// 数据太少，直接返回中位数
		return calculateMedian(times)
	}

	// 复制并排序
	sorted := make([]time.Duration, len(times))
	copy(sorted, times)
	sortDurations(sorted)

	// 去掉最大和最小值
	trimmed := sorted[1 : len(sorted)-1]

	// 计算平均值
	var sum time.Duration
	for _, t := range trimmed {
		sum += t
	}
	return sum / time.Duration(len(trimmed))
}

// calculateMedian 计算中位数
func calculateMedian(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(times))
	copy(sorted, times)
	sortDurations(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

// calculateStdDev 计算标准差
func calculateStdDev(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}

	// 计算平均值
	var sum time.Duration
	for _, t := range times {
		sum += t
	}
	mean := sum / time.Duration(len(times))

	// 计算方差
	var variance float64
	for _, t := range times {
		diff := float64(t - mean)
		variance += diff * diff
	}
	variance /= float64(len(times))

	// 返回标准差
	return time.Duration(math.Sqrt(variance))
}

// getMinMax 获取最小值和最大值
func getMinMax(times []time.Duration) (time.Duration, time.Duration) {
	if len(times) == 0 {
		return 0, 0
	}

	min := times[0]
	max := times[0]
	for _, t := range times {
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}
	return min, max
}

// sortDurations 对时间切片进行排序（冒泡排序）
func sortDurations(times []time.Duration) {
	n := len(times)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if times[i] > times[j] {
				times[i], times[j] = times[j], times[i]
			}
		}
	}
}
