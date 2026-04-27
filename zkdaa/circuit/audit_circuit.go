package circuit

import (
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// AuditCircuit 验证 Merkle 证明的电路
// 约束：给定数据块 Fi 和证明路径 Ωi，验证其根哈希等于 CIDF
type AuditCircuit struct {
	// 私有输入
	ChunkData    []byte              `gnark:"-"`       // 私有：原始数据块 Fi（链下使用）
	ProofPath    []frontend.Variable `gnark:",secret"` // 私有：Merkle 证明路径
	Helpers      []frontend.Variable `gnark:",secret"` // 私有：路径方向（0=左，1=右）
	LeafCounts   []frontend.Variable `gnark:",secret"` // 私有：每层可达叶子数
	LeafNumBytes []frontend.Variable `gnark:",secret"` // 私有：叶子数的字节表示

	// 公开输入
	MerkleRoot frontend.Variable `gnark:",public"` // 公开：CIDF
	ChunkIndex frontend.Variable `gnark:",public"` // 公开：被挑战的索引 i
	ChunkHash  frontend.Variable `gnark:",public"` // 公开：数据块的哈希（预先计算）
}

// Define 定义电路约束
func (circuit *AuditCircuit) Define(api frontend.API) error {
	// 1. 初始化 MiMC 哈希
	mimcHash, _ := mimc.NewMiMC(api)

	// 2. 使用预先计算的数据块哈希
	computedHash := circuit.ChunkHash

	// 3. 沿 Merkle 路径向上计算
	pathLen := len(circuit.ProofPath)
	calculatedIndex := frontend.Variable(0) // 用于验证索引

	for i := 0; i < pathLen; i++ {
		pathHash := circuit.ProofPath[i]
		helper := circuit.Helpers[i]
		leafCount := circuit.LeafCounts[i]
		leafNumByte := circuit.LeafNumBytes[i]

		// 计算两种可能的哈希（左右顺序）
		// 当 helper 为 0 时：computedHash 在左，pathHash 在右
		mimcHash.Write(leafNumByte, computedHash, pathHash)
		leftHash := mimcHash.Sum()
		mimcHash.Reset()

		// 当 helper 为 1 时：pathHash 在左，computedHash 在右
		mimcHash.Write(leafNumByte, pathHash, computedHash)
		rightHash := mimcHash.Sum()
		mimcHash.Reset()

		// 根据 helper 选择正确的哈希
		computedHash = api.Select(helper, rightHash, leftHash)

		// 更新计算的索引
		calculatedIndex = api.Select(helper,
			api.Add(calculatedIndex, leafCount),
			calculatedIndex)
	}

	// 4. 约束：计算出的哈希必须等于 CIDF
	api.AssertIsEqual(computedHash, circuit.MerkleRoot)

	// 5. 约束：计算出的索引必须等于挑战的索引 i
	api.AssertIsEqual(calculatedIndex, circuit.ChunkIndex)

	return nil
}

// HashChunk 哈希数据块（链下辅助函数）
func HashChunk(data []byte) []byte {
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(data)
	return hFunc.Sum(nil)
}
