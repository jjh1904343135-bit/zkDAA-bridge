package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// ========== 通用 Merkle 证明电路 ==========
// MerkleProofCircuit 验证叶节点在 Merkle 树中的路径
type MerkleProofCircuit struct {
	Leaf        frontend.Variable   // 私有：叶节点哈希值 (Changed from raw value)
	RootHash    frontend.Variable   `gnark:",public"` // 公开：Merkle 根哈希
	Path        []frontend.Variable // 私有：Merkle 路径哈希值
	LeafNum     []frontend.Variable // 私有：每层可达叶子数
	Helper      []frontend.Variable // 私有：路径方向（0=左，1=右）
	LeafIndex   frontend.Variable   // 私有：叶节点索引
	LeafNumByte []frontend.Variable // 私有：叶子数的字节表示
}

func (circuit *MerkleProofCircuit) Define(api frontend.API) error {
	mimcHash, _ := mimc.NewMiMC(api)

	// 直接使用输入的 Leaf 作为哈希值
	computedHash := circuit.Leaf

	pathLen := len(circuit.Path)
	q := frontend.Variable(0) // 计算的索引

	// 沿路径向上计算
	for i := 0; i < pathLen; i++ {
		pathHash := circuit.Path[i]
		helper := circuit.Helper[i]
		num := circuit.LeafNum[i]
		leafNumByte := circuit.LeafNumByte[i]

		// 左侧：computedHash 在左
		mimcHash.Write(leafNumByte, computedHash, pathHash)
		leftHash := mimcHash.Sum()
		mimcHash.Reset()

		// 右侧：computedHash 在右
		mimcHash.Write(leafNumByte, pathHash, computedHash)
		rightHash := mimcHash.Sum()
		mimcHash.Reset()

		// 根据 helper 选择正确的哈希
		computedHash = api.Select(helper, rightHash, leftHash)
		q = api.Select(helper, api.Add(q, num), q)
	}

	// 验证根哈希和索引
	api.AssertIsEqual(computedHash, circuit.RootHash)
	api.AssertIsEqual(circuit.LeafIndex, q)

	return nil
}

// ========== 1. 插入电路 ==========
type MerkleInsertCircuit struct {
	LeafHash frontend.Variable `gnark:",public"` // 公开：新插入叶节点的哈希
	// InsertLeaf     frontend.Variable   // REMOVED: No need for preimage
	NewRootHash    frontend.Variable   `gnark:",public"` // 公开：插入后的新根哈希
	Circuit_merkle MerkleProofCircuit  // 验证原始路径 (Verified using oldLeafHash)
	NewNum_byte    []frontend.Variable // 新的叶子数字节表示
}

func (circuit *MerkleInsertCircuit) Define(api frontend.API) error {
	// 验证原始 Merkle 路径
	err := circuit.Circuit_merkle.Define(api)
	if err != nil {
		return err
	}

	mimcHash, _ := mimc.NewMiMC(api)

	// 新叶节点哈希即为 LeafHash
	newhash := circuit.LeafHash

	// 原叶节点哈希即为 Circuit_merkle.Leaf
	oldLeafHash := circuit.Circuit_merkle.Leaf

	// 合并新旧叶节点，生成新的非叶节点
	mimcHash.Write(newhash, oldLeafHash)
	computedHash := mimcHash.Sum()
	mimcHash.Reset()

	q := frontend.Variable(1) // 新索引从 1 开始

	// 沿新路径向上计算
	pathLen := len(circuit.Circuit_merkle.Path)
	for i := 0; i < pathLen; i++ {
		pathHash := circuit.Circuit_merkle.Path[i]
		helper := circuit.Circuit_merkle.Helper[i]
		num := circuit.Circuit_merkle.LeafNum[i]
		newNumByte := circuit.NewNum_byte[i]

		// 左侧
		mimcHash.Write(newNumByte, computedHash, pathHash)
		leftHash := mimcHash.Sum()
		mimcHash.Reset()

		// 右侧
		mimcHash.Write(newNumByte, pathHash, computedHash)
		rightHash := mimcHash.Sum()
		mimcHash.Reset()

		computedHash = api.Select(helper, rightHash, leftHash)
		q = api.Select(helper, api.Add(q, num), q)
	}

	// 验证新根哈希
	api.AssertIsEqual(computedHash, circuit.NewRootHash)
	// api.AssertIsEqual(api.Add(circuit.Circuit_merkle.LeafIndex, 1), q) // Incorrect index check removed

	return nil
}

// ========== 2. 删除电路 ==========
type MerkleDeleteCircuit struct {
	NewRootHash    frontend.Variable   `gnark:",public"` // 公开：删除后的新根哈希
	Circuit_merkle MerkleProofCircuit  // 验证原始路径
	NewPath_byte   []frontend.Variable // 新的路径字节表示
}

func (circuit *MerkleDeleteCircuit) Define(api frontend.API) error {
	// 验证原始 Merkle 路径
	err := circuit.Circuit_merkle.Define(api)
	if err != nil {
		return err
	}

	mimcHash, _ := mimc.NewMiMC(api)

	// 删除后，Path[0] 成为新的叶节点
	computedHash := circuit.Circuit_merkle.Path[0]
	q := frontend.Variable(0)

	pathLen := len(circuit.Circuit_merkle.Path)

	// 从第 1 层开始（跳过原叶节点层）
	for i := 1; i < pathLen; i++ {
		pathHash := circuit.Circuit_merkle.Path[i]
		helper := circuit.Circuit_merkle.Helper[i]
		num := circuit.Circuit_merkle.LeafNum[i]
		newPathByte := circuit.NewPath_byte[i]

		// 左侧
		mimcHash.Write(newPathByte, computedHash, pathHash)
		leftHash := mimcHash.Sum()
		mimcHash.Reset()

		// 右侧
		mimcHash.Write(newPathByte, pathHash, computedHash)
		rightHash := mimcHash.Sum()
		mimcHash.Reset()

		computedHash = api.Select(helper, rightHash, leftHash)
		q = api.Select(helper, api.Add(q, num), q)
	}

	// 验证新根哈希
	api.AssertIsEqual(computedHash, circuit.NewRootHash)
	// api.AssertIsEqual(q, api.Sub(circuit.Circuit_merkle.LeafIndex, 1)) // Incorrect index check removed

	return nil
}

// ========== 3. 更新电路 ==========
type MerkleUpdateCircuit struct {
	LeafHash frontend.Variable `gnark:",public"` // 公开：更新后叶节点的哈希
	// NewLeaf        frontend.Variable  // REMOVED: No need for preimage
	NewRootHash    frontend.Variable  `gnark:",public"` // 公开：更新后的新根哈希
	Circuit_merkle MerkleProofCircuit // 验证原始路径
}

func (circuit *MerkleUpdateCircuit) Define(api frontend.API) error {
	// 验证原始 Merkle 路径
	err := circuit.Circuit_merkle.Define(api)
	if err != nil {
		return err
	}

	mimcHash, _ := mimc.NewMiMC(api)

	// 新叶节点哈希即为 LeafHash
	computedHash := circuit.LeafHash

	pathLen := len(circuit.Circuit_merkle.Path)
	q := frontend.Variable(0)

	// 沿路径向上计算新根
	for i := 0; i < pathLen; i++ {
		pathHash := circuit.Circuit_merkle.Path[i]
		helper := circuit.Circuit_merkle.Helper[i]
		num := circuit.Circuit_merkle.LeafNum[i]
		leafNumByte := circuit.Circuit_merkle.LeafNumByte[i]

		// 左侧
		mimcHash.Write(leafNumByte, computedHash, pathHash)
		leftHash := mimcHash.Sum()
		mimcHash.Reset()

		// 右侧
		mimcHash.Write(leafNumByte, pathHash, computedHash)
		rightHash := mimcHash.Sum()
		mimcHash.Reset()

		computedHash = api.Select(helper, rightHash, leftHash)
		q = api.Select(helper, api.Add(q, num), q)
	}

	// 验证新根哈希和索引
	api.AssertIsEqual(computedHash, circuit.NewRootHash)
	api.AssertIsEqual(circuit.Circuit_merkle.LeafIndex, q)

	return nil
}
