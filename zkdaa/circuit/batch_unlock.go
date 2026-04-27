package circuit

import (
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// BatchUnlockCircuit 批量解锁电路
// 证明：我知道 preimage 和 serial_number，且我的 TxLock 在批量 Merkle 树中
type BatchUnlockCircuit struct {
	// 私有输入
	Preimage     frontend.Variable   `gnark:",secret"` // preimage
	SerialNumber frontend.Variable   `gnark:",secret"` // serial number
	TxIndex      frontend.Variable   `gnark:",secret"` // 在批量中的索引
	ProofPath    []frontend.Variable `gnark:",secret"` // Merkle 路径
	Helpers      []frontend.Variable `gnark:",secret"` // 路径方向（0=左，1=右）
	LeafCounts   []frontend.Variable `gnark:",secret"` // 每层可达叶子数
	LeafNumBytes []frontend.Variable `gnark:",secret"` // 叶子数本身（用于哈希）

	// 公开输入
	MerkleRoot         frontend.Variable `gnark:",public"` // 合约存储的批量根
	SerialNumberPublic frontend.Variable `gnark:",public"` // 公开的 serial number
}

// Define 定义电路约束
func (circuit *BatchUnlockCircuit) Define(api frontend.API) error {
	mimcHash, _ := mimc.NewMiMC(api)

	// 1. 约束：计算 h = hash(preimage, serial_number)
	mimcHash.Write(circuit.Preimage, circuit.SerialNumber)
	h := mimcHash.Sum()
	mimcHash.Reset()

	// 2. 计算 TxLock 哈希（作为 Merkle 叶子）
	// 简化版：直接用 h 作为叶子，实际可添加更多字段
	txHash := h

	// 3. 沿 Merkle 路径向上计算，验证 txHash 在树中
	computedHash := txHash
	calculatedIndex := frontend.Variable(0)

	pathLen := len(circuit.ProofPath)
	for i := 0; i < pathLen; i++ {
		pathHash := circuit.ProofPath[i]
		helper := circuit.Helpers[i]
		leafCount := circuit.LeafCounts[i]
		leafNumByte := circuit.LeafNumBytes[i]

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

		// 🔧 修复：更新计算的索引
		// 如果 helper=1（当前在右侧），索引 += leafCount/2（左子树节点数）
		// 如果 helper=0（当前在左侧），索引不变
		leftSubtreeSize := api.Div(leafCount, 2)
		calculatedIndex = api.Select(helper,
			api.Add(calculatedIndex, leftSubtreeSize),
			calculatedIndex)
	}

	// 4. 约束：计算出的根 == 合约存储的批量根
	api.AssertIsEqual(computedHash, circuit.MerkleRoot)

	// 5. 约束：计算出的索引 == 私有输入的索引
	api.AssertIsEqual(calculatedIndex, circuit.TxIndex)

	// 6. 约束：私有的 serial_number == 公开的值
	api.AssertIsEqual(circuit.SerialNumber, circuit.SerialNumberPublic)

	return nil
}

// HashTxLock 计算 TxLock 交易哈希 - 链下辅助函数
func HashTxLock(preimage, serialNumber []byte) []byte {
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(preimage)
	hFunc.Write(serialNumber)
	return hFunc.Sum(nil)
}
