package merkle

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/hash"
)

// BatchMerkleTree 批量交易的 Merkle 树(基于 big.Int)
type BatchMerkleTree struct {
	Leaves []*big.Int   // 交易哈希(叶子节点)
	Root   *big.Int     // Merkle 根
	Layers [][]*big.Int // 树的所有层
}

// ProofElement Merkle 证明元素
type ProofElement struct {
	Hash         *big.Int // 兄弟节点哈希
	IsLeft       bool     // 兄弟节点是否在左侧(false=在右侧)
	LeafCount    *big.Int // 当前层的叶子数
	LeafNumBytes *big.Int // 叶子数的字节表示
}

// NewBatchMerkleTree 从交易哈希创建批量 Merkle 树
func NewBatchMerkleTree(txHashes []*big.Int) *BatchMerkleTree {
	mt := &BatchMerkleTree{
		Leaves: txHashes,
		Layers: make([][]*big.Int, 0),
	}
	mt.Build()
	return mt
}

// Build 构建 Merkle 树
func (mt *BatchMerkleTree) Build() {
	if len(mt.Leaves) == 0 {
		mt.Root = big.NewInt(0)
		return
	}

	// 第一层: 叶子节点
	currentLevel := make([]*big.Int, len(mt.Leaves))
	copy(currentLevel, mt.Leaves)
	mt.Layers = append(mt.Layers, currentLevel)

	// 逐层向上构建
	for len(currentLevel) > 1 {
		nextLevel := make([]*big.Int, 0)
		leafCount := len(currentLevel)

		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			var right *big.Int

			// 处理奇数节点: 复制左节点
			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			} else {
				right = left
			}

			// 🔧 修复: 使用统一的哈希函数
			parent := hashPairBigIntWithLeafCount(left, right, leafCount)
			nextLevel = append(nextLevel, parent)
		}

		mt.Layers = append(mt.Layers, nextLevel)
		currentLevel = nextLevel
	}

	// 设置根节点
	mt.Root = mt.Layers[len(mt.Layers)-1][0]
}

// GetProof 获取指定索引的 Merkle 证明路径
func (mt *BatchMerkleTree) GetProof(index int) ([]ProofElement, error) {
	if index < 0 || index >= len(mt.Leaves) {
		return nil, fmt.Errorf("索引超出范围: %d (总叶子数: %d)", index, len(mt.Leaves))
	}

	proof := make([]ProofElement, 0)
	currentIndex := index

	// 从叶子层开始,逐层向上
	for level := 0; level < len(mt.Layers)-1; level++ {
		currentLayer := mt.Layers[level]
		leafCount := len(currentLayer)

		// 计算兄弟节点索引
		isRightNode := (currentIndex%2 == 1)
		var siblingIndex int
		var siblingIsLeft bool

		if isRightNode {
			// 当前节点在右侧,兄弟在左侧
			siblingIndex = currentIndex - 1
			siblingIsLeft = true
		} else {
			// 当前节点在左侧,兄弟在右侧
			siblingIndex = currentIndex + 1
			siblingIsLeft = false
		}

		leafCountBig := big.NewInt(int64(leafCount))
		leafNumBytes := leafCountBig

		// 添加兄弟节点到证明路径
		if siblingIndex < len(currentLayer) {
			proof = append(proof, ProofElement{
				Hash:         currentLayer[siblingIndex],
				IsLeft:       siblingIsLeft,
				LeafCount:    leafCountBig,
				LeafNumBytes: leafNumBytes,
			})
		} else {
			// 奇数节点情况: 兄弟节点是自己
			proof = append(proof, ProofElement{
				Hash:         currentLayer[currentIndex],
				IsLeft:       siblingIsLeft,
				LeafCount:    leafCountBig,
				LeafNumBytes: leafNumBytes,
			})
		}

		// 移动到上一层
		currentIndex /= 2
	}

	return proof, nil
}

// VerifyProof 验证 Merkle 证明(链下辅助函数)
func (mt *BatchMerkleTree) VerifyProof(leafHash *big.Int, index int, proof []ProofElement) bool {
	computedHash := leafHash
	currentIndex := index

	for _, elem := range proof {
		var left, right *big.Int

		if elem.IsLeft {
			// 兄弟节点在左侧
			left = elem.Hash
			right = computedHash
		} else {
			// 兄弟节点在右侧
			left = computedHash
			right = elem.Hash
		}

		// 使用 leafCount 计算哈希
		leafCount := int(elem.LeafCount.Int64())
		computedHash = hashPairBigIntWithLeafCount(left, right, leafCount)
		currentIndex /= 2
	}

	return computedHash.Cmp(mt.Root) == 0
}

// 🔧 关键修复: 使用与电路一致的哈希方式
func hashPairBigIntWithLeafCount(left, right *big.Int, leafCount int) *big.Int {
	hFunc := hash.MIMC_BN254.New()

	// 将 big.Int 转换为有限域元素,然后再转换为字节
	// 这样与电路中的 frontend.Variable 行为一致

	leafCountBig := big.NewInt(int64(leafCount))
	leafCountFr := new(fr.Element).SetBigInt(leafCountBig)
	leafCountBytes := leafCountFr.Bytes()
	hFunc.Write(leafCountBytes[:])

	// 左节点
	leftFr := new(fr.Element).SetBigInt(left)
	leftBytes := leftFr.Bytes()
	hFunc.Write(leftBytes[:])

	// 右节点
	rightFr := new(fr.Element).SetBigInt(right)
	rightBytes := rightFr.Bytes()
	hFunc.Write(rightBytes[:])

	// 返回哈希结果
	hashBytes := hFunc.Sum(nil)
	result := new(big.Int).SetBytes(hashBytes)

	// 确保结果在有限域内
	resultFr := new(fr.Element).SetBigInt(result)
	return resultFr.BigInt(new(big.Int))
}

// GetDepth 获取树的深度
func (mt *BatchMerkleTree) GetDepth() int {
	return len(mt.Layers) - 1
}

// GetLeafCount 获取叶子数量
func (mt *BatchMerkleTree) GetLeafCount() int {
	return len(mt.Leaves)
}
