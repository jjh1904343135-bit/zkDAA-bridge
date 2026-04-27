package circuit

import (
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// AuditUnlockCircuit 审计-解锁联合电路
// 用于 SCB：先验证 Merkle 证明得到 CIDF，再验证 H == MiMC(CIDF, Sn)
// 关键创新：用 CIDF 替换原本的 Pre_II，将审计和解锁连接起来
type AuditUnlockCircuit struct {
	// ========== 公开输入 ==========
	ChunkIndex frontend.Variable `gnark:",public"` // 公开：被挑战的索引 i
	ChunkHash  frontend.Variable `gnark:",public"` // 公开：数据块的哈希（预先计算）
	H          frontend.Variable `gnark:",public"` // 公开：哈希锁 H2

	// ========== 审计部分：私有输入 ==========
	ChunkData    []byte              `gnark:"-"`       // 私有：原始数据块 Fi（链下使用）
	ProofPath    []frontend.Variable `gnark:",secret"` // 私有：Merkle 证明路径
	Helpers      []frontend.Variable `gnark:",secret"` // 私有：路径方向（0=左，1=右）
	LeafCounts   []frontend.Variable `gnark:",secret"` // 私有：每层可达叶子数
	LeafNumBytes []frontend.Variable `gnark:",secret"` // 私有：叶子数的字节表示

	// ========== 解锁部分：私有输入 ==========
	Sn frontend.Variable `gnark:",secret"` // 私有：Sn_II（序列号）

}

// Define 定义电路约束
// 验证流程：
// 1. 通过 Merkle 证明计算得到 CIDF
// 2. 验证 H == MiMC(CIDF, Sn)
// 3. 验证 Merkle 路径的索引正确性
func (circuit *AuditUnlockCircuit) Define(api frontend.API) error {
	// 1. 初始化 MiMC 哈希
	mimcHash, _ := mimc.NewMiMC(api)

	// ========== 第一部分：审计验证（计算 CIDF）==========
	// 2. 使用预先计算的数据块哈希作为起点
	computedHash := circuit.ChunkHash

	// 3. 沿 Merkle 路径向上计算，得到根哈希（CIDF）
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

	// 4. 得到 CIDF（Merkle 根）
	CIDF := computedHash

	// ========== 第二部分：解锁验证 ==========
	// 5. 计算 hash = MiMC(CIDF, Sn)
	// 关键：用 CIDF 替换原本的 Pre_II
	mimcHash.Write(CIDF)
	mimcHash.Write(circuit.Sn)
	computedH := mimcHash.Sum()
	mimcHash.Reset()

	// 6. 约束：计算出的哈希必须等于公开的 H
	api.AssertIsEqual(computedH, circuit.H)

	// ========== 第三部分：索引验证 ==========
	// 7. 约束：计算出的索引必须等于挑战的索引 i
	api.AssertIsEqual(calculatedIndex, circuit.ChunkIndex)

	return nil
}

// HashChunkForAuditUnlock 哈希数据块（链下辅助函数）
// 与 AuditCircuit 保持一致
func HashChunkForAuditUnlock(data []byte) []byte {
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(data)
	return hFunc.Sum(nil)
}
