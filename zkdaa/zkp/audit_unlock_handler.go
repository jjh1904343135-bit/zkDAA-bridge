package zkp

import (
	"fmt"
	"os"
	"time"

	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// AuditUnlockHandler 审计-解锁联合处理器
type AuditUnlockHandler struct {
	R1CS         constraint.ConstraintSystem
	ProvingKey   groth16.ProvingKey
	VerifyingKey groth16.VerifyingKey
}

// NewAuditUnlockHandler 初始化审计-解锁处理器（执行 Setup）
func NewAuditUnlockHandler(proofDepth int) (*AuditUnlockHandler, error) {
	fmt.Println("\n[AUDIT-UNLOCK] 执行审计-解锁联合电路 Setup...")
	startTime := time.Now()

	// 检查是否存在已保存的密钥
	pkPath := "zkp/audit_unlock_pk.bin"
	vkPath := "zkp/audit_unlock_vk.bin"

	var pk groth16.ProvingKey
	var vk groth16.VerifyingKey
	var ccs constraint.ConstraintSystem

	// 1. 创建电路模板
	auditUnlockCircuit := &circuit.AuditUnlockCircuit{
		ProofPath:    make([]frontend.Variable, proofDepth),
		Helpers:      make([]frontend.Variable, proofDepth),
		LeafCounts:   make([]frontend.Variable, proofDepth),
		LeafNumBytes: make([]frontend.Variable, proofDepth),
	}

	// 2. 编译电路
	var err error
	ccs, err = frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, auditUnlockCircuit)
	if err != nil {
		return nil, fmt.Errorf("审计-解锁电路编译失败: %w", err)
	}
	fmt.Printf("[AUDIT-UNLOCK] 电路编译完成。约束数量: %d\n", ccs.GetNbConstraints())

	// 3. 尝试加载已保存的密钥
	if _, err := os.Stat(pkPath); err == nil {
		fmt.Println("[AUDIT-UNLOCK] 📂 加载已保存的密钥...")

		// 加载 Proving Key
		pkFile, err := os.Open(pkPath)
		if err == nil {
			pk = groth16.NewProvingKey(ecc.BN254)
			_, err = pk.ReadFrom(pkFile)
			pkFile.Close()
			if err == nil {
				fmt.Println("[AUDIT-UNLOCK]    ✅ Proving Key 加载成功")
			}
		}

		// 加载 Verifying Key
		vkFile, err := os.Open(vkPath)
		if err == nil {
			vk = groth16.NewVerifyingKey(ecc.BN254)
			_, err = vk.ReadFrom(vkFile)
			vkFile.Close()
			if err == nil {
				fmt.Println("[AUDIT-UNLOCK]    ✅ Verifying Key 加载成功")
			}
		}
	}

	// 4. 如果加载失败，执行新的 Setup
	if pk == nil || vk == nil {
		fmt.Println("[AUDIT-UNLOCK] ⚙️  执行新的 Setup...")
		pk, vk, err = groth16.Setup(ccs)
		if err != nil {
			return nil, fmt.Errorf("Setup 失败: %w", err)
		}

		// 保存密钥
		fmt.Println("[AUDIT-UNLOCK] 💾 保存 Setup 密钥...")
		pkFile, err := os.Create(pkPath)
		if err == nil {
			pk.WriteTo(pkFile)
			pkFile.Close()
			fmt.Printf("[AUDIT-UNLOCK]    ✅ 已保存到: %s\n", pkPath)
		}

		vkFile, err := os.Create(vkPath)
		if err == nil {
			vk.WriteTo(vkFile)
			vkFile.Close()
			fmt.Printf("[AUDIT-UNLOCK]    ✅ 已保存到: %s\n", vkPath)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("[AUDIT-UNLOCK] ✅ Setup 完成, 耗时: %v\n", duration)

	return &AuditUnlockHandler{
		R1CS:         ccs,
		ProvingKey:   pk,
		VerifyingKey: vk,
	}, nil
}

// GenerateProof 生成审计-解锁联合证明
func (h *AuditUnlockHandler) GenerateProof(
	assignment *circuit.AuditUnlockCircuit,
) (groth16.Proof, error) {
	fmt.Printf("[AUDIT-UNLOCK] 生成联合证明...\n")
	startTime := time.Now()

	// 创建 witness
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("创建 witness 失败: %w", err)
	}

	// 生成证明
	proof, err := groth16.Prove(h.R1CS, h.ProvingKey, witness)
	if err != nil {
		return nil, fmt.Errorf("生成证明失败: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("[AUDIT-UNLOCK] ✅ 联合证明生成成功, 耗时: %v\n", duration)
	return proof, nil
}

// Verify 验证审计-解锁证明
func (h *AuditUnlockHandler) Verify(
	proof groth16.Proof,
	chunkIndex int,
	chunkHash []byte,
	H []byte,
) error {
	fmt.Println("[AUDIT-UNLOCK] 验证联合证明...")
	startTime := time.Now()

	// 构造公开输入
	publicAssignment := &circuit.AuditUnlockCircuit{
		ChunkIndex: chunkIndex,
		ChunkHash:  chunkHash,
		H:          H,
	}

	publicWitness, err := frontend.NewWitness(publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return fmt.Errorf("创建公开 witness 失败: %w", err)
	}

	// 验证
	err = groth16.Verify(proof, h.VerifyingKey, publicWitness)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("[AUDIT-UNLOCK] ❌ 证明验证失败: %v\n", err)
		return err
	}

	fmt.Printf("[AUDIT-UNLOCK] ✅ 证明验证成功, 耗时: %v\n", duration)
	return nil
}

// GetConstraintCount 返回约束数量
func (h *AuditUnlockHandler) GetConstraintCount() int {
	return h.R1CS.GetNbConstraints()
}

// ExportVerifier 导出 Solidity 验证器
func (h *AuditUnlockHandler) ExportVerifier(outputPath string) error {
	fmt.Printf("[AUDIT-UNLOCK] 正在导出 Verifier.sol 到 %s...\n", outputPath)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	err = h.VerifyingKey.ExportSolidity(file)
	if err != nil {
		return fmt.Errorf("导出 Solidity 失败: %w", err)
	}

	fmt.Println("[AUDIT-UNLOCK] ✅ Verifier.sol 导出成功!")
	return nil
}
