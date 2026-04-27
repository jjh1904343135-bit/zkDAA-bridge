package zkp

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"time"
	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// BatchZKPHandler 批量解锁的 ZKP 处理器
type BatchZKPHandler struct {
	pk          groth16.ProvingKey
	vk          groth16.VerifyingKey
	ccs         constraint.ConstraintSystem
	batchSize   int
	merkleDepth int

	// 性能指标
	compileTimeMs  float64
	loadKeysTimeMs float64
}

// HandlerMetrics 详细的性能指标
type HandlerMetrics struct {
	CompileTimeMs  float64
	LoadKeysTimeMs float64
}

// BatchMetrics 批量测试性能指标
type BatchMetrics struct {
	BatchSize       int           // 批量大小（16/64/128/256）
	MerkleDepth     int           // Merkle 树深度
	SetupTime       time.Duration // Setup 时间
	ProveTime       time.Duration // Prove 时间
	VerifyTime      time.Duration // Verify 时间
	ConstraintCount int           // 约束数量
	ProofSize       int           // 证明大小（bytes）
}

// NewBatchZKPHandler 根据批量大小初始化
func NewBatchZKPHandler(batchSize int) (*BatchZKPHandler, error) {
	// 计算 Merkle 树深度
	depth := calculateDepth(batchSize)

	// 初始化电路
	dummyCircuit := circuit.BatchUnlockCircuit{
		ProofPath:    make([]frontend.Variable, depth),
		Helpers:      make([]frontend.Variable, depth),
		LeafCounts:   make([]frontend.Variable, depth),
		LeafNumBytes: make([]frontend.Variable, depth),
	}

	// 编译电路
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &dummyCircuit)
	if err != nil {
		return nil, err
	}

	// 🔧 关键修复：根据批量大小加载对应的 Setup
	var pk groth16.ProvingKey
	var vk groth16.VerifyingKey

	// 使用批量大小作为文件名后缀
	pkPath := fmt.Sprintf("zkp/batch_pk_%d.bin", batchSize)
	vkPath := fmt.Sprintf("zkp/batch_vk_%d.bin", batchSize)

	// 检查文件是否存在
	if _, err := os.Stat(pkPath); err == nil {
		// 文件存在，加载
		fmt.Printf("📂 加载批量 %d 的 Setup 密钥...\n", batchSize)

		// 加载 Proving Key
		pkFile, err := os.Open(pkPath)
		if err != nil {
			return nil, fmt.Errorf("打开 PK 文件失败: %w", err)
		}
		pk = groth16.NewProvingKey(ecc.BN254)
		_, err = pk.ReadFrom(pkFile)
		pkFile.Close()
		if err != nil {
			return nil, fmt.Errorf("读取 PK 失败: %w", err)
		}
		fmt.Println("   ✅ Proving Key 加载成功")

		// 加载 Verifying Key
		vkFile, err := os.Open(vkPath)
		if err != nil {
			return nil, fmt.Errorf("打开 VK 文件失败: %w", err)
		}
		vk = groth16.NewVerifyingKey(ecc.BN254)
		_, err = vk.ReadFrom(vkFile)
		vkFile.Close()
		if err != nil {
			return nil, fmt.Errorf("读取 VK 失败: %w", err)
		}
		fmt.Println("   ✅ Verifying Key 加载成功")

	} else {
		// 文件不存在，执行 Setup
		fmt.Printf("⚠️  未找到批量 %d 的 Setup 密钥，执行新的 Setup...\n", batchSize)
		pk, vk, err = groth16.Setup(ccs)
		if err != nil {
			return nil, err
		}

		// 保存新生成的密钥
		fmt.Println("   💾 保存 Setup 密钥...")

		pkFile, err := os.Create(pkPath)
		if err == nil {
			pk.WriteTo(pkFile)
			pkFile.Close()
			fmt.Printf("   ✅ 已保存到: %s\n", pkPath)
		}

		vkFile, err := os.Create(vkPath)
		if err == nil {
			vk.WriteTo(vkFile)
			vkFile.Close()
			fmt.Printf("   ✅ 已保存到: %s\n", vkPath)
		}

		fmt.Println("   ⚠️  警告: 使用新生成的密钥，与链上验证器可能不匹配！")
		fmt.Printf("   对于批量 %d，请重新生成验证器并部署\n", batchSize)
	}

	return &BatchZKPHandler{
		pk:          pk,
		vk:          vk,
		ccs:         ccs,
		batchSize:   batchSize,
		merkleDepth: depth,
	}, nil
}

// Prove 生成证明
func (h *BatchZKPHandler) Prove(assignment *circuit.BatchUnlockCircuit) (groth16.Proof, error) {
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	return groth16.Prove(h.ccs, h.pk, witness)
}

// Verify 验证证明
func (h *BatchZKPHandler) Verify(proof groth16.Proof, publicInputs []*big.Int) error {
	publicWitness, err := frontend.NewWitness(&circuit.BatchUnlockCircuit{
		MerkleRoot:         frontend.Variable(publicInputs[0]),
		SerialNumberPublic: frontend.Variable(publicInputs[1]),
	}, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return err
	}

	return groth16.Verify(proof, h.vk, publicWitness)
}

// GetConstraintCount 返回约束数量
func (h *BatchZKPHandler) GetConstraintCount() int {
	return h.ccs.GetNbConstraints()
}

// GetMerkleDepth 返回 Merkle 树深度
func (h *BatchZKPHandler) GetMerkleDepth() int {
	return h.merkleDepth
}

// GetVerifyingKey 返回验证密钥（用于导出 Solidity 验证器）
func (h *BatchZKPHandler) GetVerifyingKey() groth16.VerifyingKey {
	return h.vk
}

// ExportSolidityVerifier 导出 Solidity 验证器合约
func (h *BatchZKPHandler) ExportSolidityVerifier(w io.Writer) error {
	return h.vk.ExportSolidity(w)
}

// calculateDepth 计算 Merkle 树深度
func calculateDepth(leafCount int) int {
	depth := 0
	n := leafCount
	for n > 1 {
		n = (n + 1) / 2
		depth++
	}
	return depth
}

// GetMetrics 获取详细的性能指标
func (h *BatchZKPHandler) GetMetrics() *HandlerMetrics {
	return &HandlerMetrics{
		CompileTimeMs:  h.compileTimeMs,
		LoadKeysTimeMs: h.loadKeysTimeMs,
	}
}
