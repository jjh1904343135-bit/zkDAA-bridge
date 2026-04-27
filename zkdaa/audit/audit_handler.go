package audit

import (
	"fmt"
	"time"
	"zk-htlc/circuit"
	"zk-htlc/merkle"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// AuditHandler 处理数据审计的 ZKP
type AuditHandler struct {
	R1CS         constraint.ConstraintSystem
	ProvingKey   groth16.ProvingKey
	VerifyingKey groth16.VerifyingKey
}

// NewAuditHandler 初始化审计处理器（执行 Setup）
func NewAuditHandler(proofDepth int) (*AuditHandler, error) {
	fmt.Println("\n[AUDIT] 执行审计电路 Setup...")
	startTime := time.Now()

	// 1. 创建电路模板
	auditCircuit := &circuit.AuditCircuit{
		ProofPath:    make([]frontend.Variable, proofDepth),
		Helpers:      make([]frontend.Variable, proofDepth),
		LeafCounts:   make([]frontend.Variable, proofDepth),
		LeafNumBytes: make([]frontend.Variable, proofDepth),
	}

	// 2. 编译电路
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, auditCircuit)
	if err != nil {
		return nil, fmt.Errorf("审计电路编译失败: %w", err)
	}
	fmt.Printf("[AUDIT] 电路编译完成。约束数量: %d\n", ccs.GetNbConstraints())

	// 3. Setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, fmt.Errorf("Setup 失败: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("[AUDIT] ✅ Setup 完成, 耗时: %v\n", duration)

	return &AuditHandler{
		R1CS:         ccs,
		ProvingKey:   pk,
		VerifyingKey: vk,
	}, nil
}

// GenerateAuditProof DSP 生成审计证明
func (h *AuditHandler) GenerateAuditProof(
	mt *merkle.MerkleTree,
	chunkIndex int,
) (groth16.Proof, error) {
	fmt.Printf("[AUDIT] 生成索引 %d 的审计证明...\n", chunkIndex)
	startTime := time.Now()

	// 1. 获取 Merkle 证明
	proofPath, helpers := mt.GetProof(chunkIndex)
	merkleRoot := mt.GetRoot()
	chunkData := mt.Leaves[chunkIndex]

	// 2. 计算数据块哈希
	chunkHash := circuit.HashChunk(chunkData)

	// 3. 构造 witness
	assignment := &circuit.AuditCircuit{
		ChunkData:    chunkData,
		ProofPath:    make([]frontend.Variable, len(proofPath)),
		Helpers:      make([]frontend.Variable, len(helpers)),
		LeafCounts:   make([]frontend.Variable, len(mt.PathCounts)),
		LeafNumBytes: make([]frontend.Variable, len(mt.PathBytes)),
		MerkleRoot:   merkleRoot,
		ChunkIndex:   chunkIndex,
		ChunkHash:    chunkHash,
	}

	for i := 0; i < len(proofPath); i++ {
		assignment.ProofPath[i] = proofPath[i]
		assignment.Helpers[i] = helpers[i]
		assignment.LeafCounts[i] = mt.PathCounts[i]
		assignment.LeafNumBytes[i] = mt.PathBytes[i]
	}

	// 4. 创建 witness
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("创建 witness 失败: %w", err)
	}

	// 5. 生成证明
	proof, err := groth16.Prove(h.R1CS, h.ProvingKey, witness)
	if err != nil {
		return nil, fmt.Errorf("生成证明失败: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("[AUDIT] ✅ 审计证明生成成功, 耗时: %v\n", duration)
	return proof, nil
}

// VerifyAuditProof 用户验证审计证明
func (h *AuditHandler) VerifyAuditProof(
	proof groth16.Proof,
	merkleRoot []byte,
	chunkIndex int,
	chunkHash []byte,
) error {
	fmt.Println("[AUDIT] 验证审计证明...")
	startTime := time.Now()

	// 构造公开输入
	publicAssignment := &circuit.AuditCircuit{
		MerkleRoot: merkleRoot,
		ChunkIndex: chunkIndex,
		ChunkHash:  chunkHash,
	}

	publicWitness, err := frontend.NewWitness(publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return fmt.Errorf("创建公开 witness 失败: %w", err)
	}

	// 验证
	err = groth16.Verify(proof, h.VerifyingKey, publicWitness)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("[AUDIT] ❌ 审计证明验证失败: %v\n", err)
		return err
	}

	fmt.Printf("[AUDIT] ✅ 审计证明验证成功, 耗时: %v\n", duration)
	return nil
}

// GetConstraintCount 返回审计电路的约束数量
func (h *AuditHandler) GetConstraintCount() int {
	return h.R1CS.GetNbConstraints()
}
