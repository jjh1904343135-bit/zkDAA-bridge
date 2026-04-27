package zkp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// MerkleOperationHandler 处理 Merkle 操作的 ZKP
type MerkleOperationHandler struct {
	TreeDepth int

	// 插入操作
	InsertPK groth16.ProvingKey
	InsertVK groth16.VerifyingKey

	// 删除操作
	DeletePK groth16.ProvingKey
	DeleteVK groth16.VerifyingKey

	// 更新操作
	UpdatePK groth16.ProvingKey
	UpdateVK groth16.VerifyingKey
}

// NewMerkleOperationHandler 初始化处理器
func NewMerkleOperationHandler(treeDepth int) (*MerkleOperationHandler, error) {
	handler := &MerkleOperationHandler{TreeDepth: treeDepth}

	fmt.Printf("📦 Initializing Merkle Operation Handler (depth=%d)...\n", treeDepth)

	// 初始化插入电路
	if err := handler.setupInsert(treeDepth); err != nil {
		return nil, fmt.Errorf("setup insert circuit: %w", err)
	}

	// 初始化删除电路
	if err := handler.setupDelete(treeDepth); err != nil {
		return nil, fmt.Errorf("setup delete circuit: %w", err)
	}

	// 初始化更新电路
	if err := handler.setupUpdate(treeDepth); err != nil {
		return nil, fmt.Errorf("setup update circuit: %w", err)
	}

	fmt.Println("✅ Merkle Operation Handler initialized")
	return handler, nil
}

// ========== 插入操作 ==========
func (h *MerkleOperationHandler) setupInsert(depth int) error {
	circuit := &circuit.MerkleInsertCircuit{
		Circuit_merkle: circuit.MerkleProofCircuit{
			Path:        make([]frontend.Variable, depth),
			LeafNum:     make([]frontend.Variable, depth),
			Helper:      make([]frontend.Variable, depth),
			LeafNumByte: make([]frontend.Variable, depth),
		},
		NewNum_byte: make([]frontend.Variable, depth),
	}

	// 尝试加载已有的密钥
	pkPath := fmt.Sprintf("zkp/insert_pk_v3_%d.bin", depth)
	vkPath := fmt.Sprintf("zkp/insert_vk_v3_%d.bin", depth)

	if _, err := os.Stat(pkPath); err == nil {
		// 加载现有密钥
		fmt.Printf("  📁 Loading insert keys (depth=%d)...\n", depth)
		pkFile, _ := os.Open(pkPath)
		h.InsertPK = groth16.NewProvingKey(ecc.BN254)
		_, _ = h.InsertPK.ReadFrom(pkFile)
		pkFile.Close()

		vkFile, _ := os.Open(vkPath)
		h.InsertVK = groth16.NewVerifyingKey(ecc.BN254)
		_, _ = h.InsertVK.ReadFrom(vkFile)
		vkFile.Close()

		fmt.Println("  ✅ Insert keys loaded")
		return nil
	}

	// 生成新密钥
	fmt.Printf("  🔧 Compiling insert circuit (depth=%d)...\n", depth)
	start := time.Now()
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return err
	}
	fmt.Printf("  ⏱️  Compiled in %v\n", time.Since(start))

	fmt.Println("  🔑 Generating insert proving/verifying keys...")
	start = time.Now()
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return err
	}
	h.InsertPK = pk
	h.InsertVK = vk
	fmt.Printf("  ⏱️  Setup completed in %v\n", time.Since(start))

	// 保存密钥
	os.MkdirAll("zkp", 0755)
	pkFile, _ := os.Create(pkPath)
	_, _ = pk.WriteTo(pkFile)
	pkFile.Close()

	vkFile, _ := os.Create(vkPath)
	_, _ = vk.WriteTo(vkFile)
	vkFile.Close()

	fmt.Println("  ✅ Insert circuit setup complete")
	return nil
}

// ProveInsert 生成插入操作的 ZK 证明
func (h *MerkleOperationHandler) ProveInsert(
	oldLeaf []byte,
	newLeaf []byte,
	oldRoot []byte,
	proofPath [][]byte,
	helpers []int,
	leafCounts []int,
	leafNumBytes [][]byte,
	leafIndex int,
	newRoot []byte,
) (groth16.Proof, error) {

	fmt.Println("  🔐 Generating insert proof...")
	start := time.Now()

	// 构造 witness
	assignment := &circuit.MerkleInsertCircuit{
		LeafHash:    hashChunk(newLeaf), // Use Hash
		NewRootHash: newRoot,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf:        hashChunk(oldLeaf), // Use Hash
			RootHash:    oldRoot,
			LeafIndex:   leafIndex,
			Path:        make([]frontend.Variable, len(proofPath)),
			LeafNum:     make([]frontend.Variable, len(leafCounts)),
			Helper:      make([]frontend.Variable, len(helpers)),
			LeafNumByte: make([]frontend.Variable, len(leafNumBytes)),
		},
		NewNum_byte: make([]frontend.Variable, len(leafNumBytes)),
	}

	// 填充路径数据
	for i := 0; i < len(proofPath); i++ {
		assignment.Circuit_merkle.Path[i] = proofPath[i]
		assignment.Circuit_merkle.LeafNum[i] = leafCounts[i]
		assignment.Circuit_merkle.Helper[i] = helpers[i]
		assignment.Circuit_merkle.LeafNumByte[i] = leafNumBytes[i]
		assignment.NewNum_byte[i] = calculateNewLeafNumByte(leafCounts[i], 1)
	}

	// 生成 witness
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("create witness: %w", err)
	}

	// 编译电路（临时）
	circuit := &circuit.MerkleInsertCircuit{
		Circuit_merkle: circuit.MerkleProofCircuit{
			Path:        make([]frontend.Variable, len(proofPath)),
			LeafNum:     make([]frontend.Variable, len(leafCounts)),
			Helper:      make([]frontend.Variable, len(helpers)),
			LeafNumByte: make([]frontend.Variable, len(leafNumBytes)),
		},
		NewNum_byte: make([]frontend.Variable, len(leafNumBytes)),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, fmt.Errorf("compile circuit: %w", err)
	}

	// 生成证明
	proof, err := groth16.Prove(r1cs, h.InsertPK, witness)
	if err != nil {
		return nil, fmt.Errorf("generate proof: %w", err)
	}

	fmt.Printf("  ⏱️  Insert proof generated in %v\n", time.Since(start))
	return proof, nil
}

// VerifyInsert 验证插入证明
func (h *MerkleOperationHandler) VerifyInsert(proof groth16.Proof, leafHash, newRoot []byte) error {
	publicWitness := &circuit.MerkleInsertCircuit{
		LeafHash:    leafHash,
		NewRootHash: newRoot,
	}

	witness, err := frontend.NewWitness(publicWitness, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return fmt.Errorf("create public witness: %w", err)
	}

	err = groth16.Verify(proof, h.InsertVK, witness)
	if err != nil {
		return fmt.Errorf("verify failed: %w", err)
	}

	return nil
}

// ProveDelete 生成删除操作的 ZK 证明
func (h *MerkleOperationHandler) ProveDelete(
	oldLeaf []byte,
	oldRoot []byte,
	proofPath [][]byte,
	helpers []int,
	leafCounts []int,
	leafNumBytes [][]byte,
	leafIndex int,
	newRoot []byte,
) (groth16.Proof, error) {

	fmt.Println("  🔐 Generating delete proof...")
	start := time.Now()

	assignment := &circuit.MerkleDeleteCircuit{
		NewRootHash: newRoot,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf:        hashChunk(oldLeaf), // Use Hash
			RootHash:    oldRoot,
			LeafIndex:   leafIndex,
			Path:        make([]frontend.Variable, len(proofPath)),
			LeafNum:     make([]frontend.Variable, len(leafCounts)),
			Helper:      make([]frontend.Variable, len(helpers)),
			LeafNumByte: make([]frontend.Variable, len(leafNumBytes)),
		},
		NewPath_byte: make([]frontend.Variable, len(proofPath)),
	}

	for i := 0; i < len(proofPath); i++ {
		assignment.Circuit_merkle.Path[i] = proofPath[i]
		assignment.Circuit_merkle.LeafNum[i] = leafCounts[i]
		assignment.Circuit_merkle.Helper[i] = helpers[i]
		assignment.Circuit_merkle.LeafNumByte[i] = leafNumBytes[i]

		// For deletion, NewPath_byte[i] logic:
		// Based on calculateNewRootAfterDelete, we assumed NewPathByte[i] is H(newCount).
		// Let's verify standard binary MerkleDeleteCircuit expectations.
		// If implementation mirrors Insert, then NewPath_byte[i] is indeed the hash of the new count.
		newCount := leafCounts[i] - 1
		assignment.NewPath_byte[i] = calculateNewLeafNumByte(newCount, 0) // Delta 0 effectively, just calculating hash of newCount
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("create witness: %w", err)
	}

	circuit := &circuit.MerkleDeleteCircuit{
		Circuit_merkle: circuit.MerkleProofCircuit{
			Path:        make([]frontend.Variable, len(proofPath)),
			LeafNum:     make([]frontend.Variable, len(leafCounts)),
			Helper:      make([]frontend.Variable, len(helpers)),
			LeafNumByte: make([]frontend.Variable, len(leafNumBytes)),
		},
		NewPath_byte: make([]frontend.Variable, len(leafNumBytes)),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, fmt.Errorf("compile circuit: %w", err)
	}

	proof, err := groth16.Prove(r1cs, h.DeletePK, witness)
	if err != nil {
		return nil, fmt.Errorf("generate proof: %w", err)
	}

	fmt.Printf("  ⏱️  Delete proof generated in %v\n", time.Since(start))
	return proof, nil
}

// ========== 删除操作 ==========
func (h *MerkleOperationHandler) setupDelete(depth int) error {
	circuit := &circuit.MerkleDeleteCircuit{
		Circuit_merkle: circuit.MerkleProofCircuit{
			Path:        make([]frontend.Variable, depth),
			LeafNum:     make([]frontend.Variable, depth),
			Helper:      make([]frontend.Variable, depth),
			LeafNumByte: make([]frontend.Variable, depth),
		},
		NewPath_byte: make([]frontend.Variable, depth),
	}

	pkPath := fmt.Sprintf("zkp/delete_pk_v3_%d.bin", depth)
	vkPath := fmt.Sprintf("zkp/delete_vk_v3_%d.bin", depth)

	if _, err := os.Stat(pkPath); err == nil {
		fmt.Printf("  📁 Loading delete keys (depth=%d)...\n", depth)
		pkFile, _ := os.Open(pkPath)
		h.DeletePK = groth16.NewProvingKey(ecc.BN254)
		_, _ = h.DeletePK.ReadFrom(pkFile)
		pkFile.Close()

		vkFile, _ := os.Open(vkPath)
		h.DeleteVK = groth16.NewVerifyingKey(ecc.BN254)
		_, _ = h.DeleteVK.ReadFrom(vkFile)
		vkFile.Close()

		fmt.Println("  ✅ Delete keys loaded")
		return nil
	}

	fmt.Printf("  🔧 Compiling delete circuit (depth=%d)...\n", depth)
	start := time.Now()
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return err
	}
	fmt.Printf("  ⏱️  Compiled in %v\n", time.Since(start))

	fmt.Println("  🔑 Generating delete proving/verifying keys...")
	start = time.Now()
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return err
	}
	h.DeletePK = pk
	h.DeleteVK = vk
	fmt.Printf("  ⏱️  Setup completed in %v\n", time.Since(start))

	os.MkdirAll("zkp", 0755)
	pkFile, _ := os.Create(pkPath)
	_, _ = pk.WriteTo(pkFile)
	pkFile.Close()

	vkFile, _ := os.Create(vkPath)
	_, _ = vk.WriteTo(vkFile)
	vkFile.Close()

	fmt.Println("  ✅ Delete circuit setup complete")
	return nil
}

// ========== 更新操作 ==========
func (h *MerkleOperationHandler) setupUpdate(depth int) error {
	circuit := &circuit.MerkleUpdateCircuit{
		Circuit_merkle: circuit.MerkleProofCircuit{
			Path:        make([]frontend.Variable, depth),
			LeafNum:     make([]frontend.Variable, depth),
			Helper:      make([]frontend.Variable, depth),
			LeafNumByte: make([]frontend.Variable, depth),
		},
	}

	pkPath := fmt.Sprintf("zkp/update_pk_v3_%d.bin", depth)
	vkPath := fmt.Sprintf("zkp/update_vk_v3_%d.bin", depth)

	if _, err := os.Stat(pkPath); err == nil {
		fmt.Printf("  📁 Loading update keys (depth=%d)...\n", depth)
		pkFile, _ := os.Open(pkPath)
		h.UpdatePK = groth16.NewProvingKey(ecc.BN254)
		_, _ = h.UpdatePK.ReadFrom(pkFile)
		pkFile.Close()

		vkFile, _ := os.Open(vkPath)
		h.UpdateVK = groth16.NewVerifyingKey(ecc.BN254)
		_, _ = h.UpdateVK.ReadFrom(vkFile)
		vkFile.Close()

		fmt.Println("  ✅ Update keys loaded")
		return nil
	}

	fmt.Printf("  🔧 Compiling update circuit (depth=%d)...\n", depth)
	start := time.Now()
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return err
	}
	fmt.Printf("  ⏱️  Compiled in %v\n", time.Since(start))

	fmt.Println("  🔑 Generating update proving/verifying keys...")
	start = time.Now()
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return err
	}
	h.UpdatePK = pk
	h.UpdateVK = vk
	fmt.Printf("  ⏱️  Setup completed in %v\n", time.Since(start))

	os.MkdirAll("zkp", 0755)
	pkFile, _ := os.Create(pkPath)
	_, _ = pk.WriteTo(pkFile)
	pkFile.Close()

	vkFile, _ := os.Create(vkPath)
	_, _ = vk.WriteTo(vkFile)
	vkFile.Close()

	fmt.Println("  ✅ Update circuit setup complete")
	return nil
}

// ProveUpdate 生成更新操作的 ZK 证明
func (h *MerkleOperationHandler) ProveUpdate(
	oldLeaf []byte,
	newLeaf []byte,
	oldRoot []byte,
	proofPath [][]byte,
	helpers []int,
	leafCounts []int,
	leafNumBytes [][]byte,
	leafIndex int,
	newRoot []byte,
) (groth16.Proof, error) {

	fmt.Println("  🔐 Generating update proof...")
	start := time.Now()

	assignment := &circuit.MerkleUpdateCircuit{
		LeafHash:    hashChunk(newLeaf), // Use Hash
		NewRootHash: newRoot,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf:        hashChunk(oldLeaf), // Use Hash
			RootHash:    oldRoot,
			LeafIndex:   leafIndex,
			Path:        make([]frontend.Variable, len(proofPath)),
			LeafNum:     make([]frontend.Variable, len(leafCounts)),
			Helper:      make([]frontend.Variable, len(helpers)),
			LeafNumByte: make([]frontend.Variable, len(leafNumBytes)),
		},
	}

	for i := 0; i < len(proofPath); i++ {
		assignment.Circuit_merkle.Path[i] = proofPath[i]
		assignment.Circuit_merkle.LeafNum[i] = leafCounts[i]
		assignment.Circuit_merkle.Helper[i] = helpers[i]
		assignment.Circuit_merkle.LeafNumByte[i] = leafNumBytes[i]
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("create witness: %w", err)
	}

	// 编译电路（临时）
	circuit := &circuit.MerkleUpdateCircuit{
		Circuit_merkle: circuit.MerkleProofCircuit{
			Path:        make([]frontend.Variable, len(proofPath)),
			LeafNum:     make([]frontend.Variable, len(leafCounts)),
			Helper:      make([]frontend.Variable, len(helpers)),
			LeafNumByte: make([]frontend.Variable, len(leafNumBytes)),
		},
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, fmt.Errorf("compile circuit: %w", err)
	}

	proof, err := groth16.Prove(r1cs, h.UpdatePK, witness)
	if err != nil {
		return nil, fmt.Errorf("generate proof: %w", err)
	}

	fmt.Printf("  ⏱️  Update proof generated in %v\n", time.Since(start))
	return proof, nil
}

// ========== 辅助函数 ==========
func hashChunk(data []byte) []byte {
	// 使用 MiMC 哈希
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(data)
	return hFunc.Sum(nil)
}

func calculateNewLeafNumByte(oldNum int, delta int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int32((oldNum+delta)*2))
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(buf.Bytes())
	return hFunc.Sum(nil)
}
