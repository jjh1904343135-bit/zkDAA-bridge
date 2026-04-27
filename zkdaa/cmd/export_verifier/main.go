package main

import (
	"fmt"
	"log"
	"os"
	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	fmt.Println("🔧 生成 BatchUnlockVerifier 合约...")

	// 使用固定的批量大小
	batchSize := 16
	depth := calculateDepth(batchSize)

	fmt.Printf("批量大小: %d\n", batchSize)
	fmt.Printf("Merkle 深度: %d\n", depth)

	// 初始化电路
	dummyCircuit := circuit.BatchUnlockCircuit{
		ProofPath:    make([]frontend.Variable, depth),
		Helpers:      make([]frontend.Variable, depth),
		LeafCounts:   make([]frontend.Variable, depth),
		LeafNumBytes: make([]frontend.Variable, depth),
	}

	// 编译电路
	fmt.Println("📝 编译电路...")
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &dummyCircuit)
	if err != nil {
		log.Fatalf("❌ 编译失败: %v", err)
	}
	fmt.Printf("约束数量: %d\n", ccs.GetNbConstraints())

	// Setup
	fmt.Println("🔑 执行 Setup...")
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatalf("❌ Setup 失败: %v", err)
	}

	// 🔧 新增：保存 Setup 结果到文件
	fmt.Println("💾 保存 Setup 密钥...")

	// 保存 Proving Key
	pkFile, err := os.Create("zkp/batch_pk.bin")
	if err != nil {
		log.Fatalf("❌ 创建 PK 文件失败: %v", err)
	}
	defer pkFile.Close()

	_, err = pk.WriteTo(pkFile)
	if err != nil {
		log.Fatalf("❌ 写入 PK 失败: %v", err)
	}
	pkInfo, _ := os.Stat("zkp/batch_pk.bin")
	fmt.Printf("   ✅ Proving Key 已保存: zkp/batch_pk.bin (%d bytes)\n", pkInfo.Size())

	// 保存 Verifying Key
	vkFile, err := os.Create("zkp/batch_vk.bin")
	if err != nil {
		log.Fatalf("❌ 创建 VK 文件失败: %v", err)
	}
	defer vkFile.Close()

	_, err = vk.WriteTo(vkFile)
	if err != nil {
		log.Fatalf("❌ 写入 VK 失败: %v", err)
	}
	vkInfo, _ := os.Stat("zkp/batch_vk.bin")
	fmt.Printf("   ✅ Verifying Key 已保存: zkp/batch_vk.bin (%d bytes)\n", vkInfo.Size())

	// 导出验证器合约
	outputPath := "blockchain-contracts/contracts/BatchUnlockVerifier.sol"
	fmt.Printf("📤 导出验证器到: %s\n", outputPath)

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("❌ 创建文件失败: %v", err)
	}
	defer f.Close()

	err = vk.ExportSolidity(f)
	if err != nil {
		log.Fatalf("❌ 导出失败: %v", err)
	}

	fmt.Println("✅ 验证器合约生成成功！")
	fmt.Printf("📍 文件位置: %s\n", outputPath)

	// 检查文件大小
	fileInfo, _ := os.Stat(outputPath)
	fmt.Printf("📦 文件大小: %d bytes\n", fileInfo.Size())

	fmt.Println("\n⚠️  重要提示:")
	fmt.Println("   Setup 密钥已保存，请确保:")
	fmt.Println("   1. 重新部署智能合约")
	fmt.Println("   2. 运行测试程序时会自动加载这些密钥")
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
