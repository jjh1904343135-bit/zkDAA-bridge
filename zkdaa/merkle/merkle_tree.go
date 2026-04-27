package merkle

import (
	"bytes"
	"encoding/binary"

	"github.com/consensys/gnark-crypto/hash"
)

// MerkleTree 表示一棵 Merkle 树
type MerkleTree struct {
	Leaves     [][]byte   // 原始数据块（叶子节点）
	TreeLayers [][][]byte // 树的所有层（从叶子到根）
	PathCounts []int      // 每层节点可达叶子数
	PathBytes  [][]byte   // 每层路径的字节表示
}

// NewMerkleTree 从数据块创建 Merkle 树
func NewMerkleTree(chunks [][]byte) *MerkleTree {
	mt := &MerkleTree{Leaves: chunks}
	mt.Build()
	return mt
}

// Build 构建 Merkle 树
func (mt *MerkleTree) Build() {
	// 1. 确保叶子数为偶数
	if len(mt.Leaves)%2 != 0 {
		mt.Leaves = append(mt.Leaves, mt.Leaves[len(mt.Leaves)-1])
	}

	// 2. 计算路径参数
	depth := calculateDepth(len(mt.Leaves))
	mt.PathCounts, mt.PathBytes = calculateLeafNodesInPath(depth)

	// 3. 哈希叶子节点层
	hashedLeaves := make([][]byte, len(mt.Leaves))
	for i, leaf := range mt.Leaves {
		hashedLeaves[i] = hashFunction(leaf)
	}
	mt.TreeLayers = append(mt.TreeLayers, hashedLeaves)

	// 4. 逐层构建树
	currentLevel := hashedLeaves
	levelIndex := 0
	for len(currentLevel) > 1 {
		var newLevel [][]byte
		for i := 0; i < len(currentLevel); i += 2 {
			// 合并节点：pathByte + left + right
			combined := append(append(mt.PathBytes[levelIndex],
				currentLevel[i]...),
				currentLevel[i+1]...)
			combinedHash := hashFunction(combined)
			newLevel = append(newLevel, combinedHash)
		}

		// 处理奇数节点
		if len(newLevel)%2 != 0 && len(newLevel) != 1 {
			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}

		mt.TreeLayers = append(mt.TreeLayers, newLevel)
		currentLevel = newLevel
		levelIndex++
	}
}

// GetRoot 获取 Merkle 根哈希 (CIDF)
func (mt *MerkleTree) GetRoot() []byte {
	if len(mt.TreeLayers) == 0 {
		return nil
	}
	return mt.TreeLayers[len(mt.TreeLayers)-1][0]
}

// GetProof 获取指定索引的 Merkle 证明路径
func (mt *MerkleTree) GetProof(chunkIndex int) (proof [][]byte, helpers []int) {
	proof = [][]byte{}
	helpers = []int{}
	layerSize := len(mt.Leaves)

	for level := 0; level < len(mt.TreeLayers)-1; level++ {
		isRightNode := (chunkIndex%2 == 1)
		var siblingIndex int

		if isRightNode {
			siblingIndex = chunkIndex - 1
			helpers = append(helpers, 1) // 当前节点在右侧
		} else {
			siblingIndex = chunkIndex + 1
			helpers = append(helpers, 0) // 当前节点在左侧
		}

		if siblingIndex < layerSize {
			proof = append(proof, mt.TreeLayers[level][siblingIndex])
		}

		chunkIndex /= 2
		layerSize = len(mt.TreeLayers[level+1])
	}

	return proof, helpers
}

// hashFunction 使用 MiMC 哈希
func hashFunction(data []byte) []byte {
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(data)
	return hFunc.Sum(nil)
}

// calculateDepth 计算树的深度
func calculateDepth(leafCount int) int {
	depth := 0
	n := leafCount
	for n > 1 {
		n = (n + 1) / 2
		depth++
	}
	return depth
}

// calculateLeafNodesInPath 计算路径中每个节点可到达的叶子数
func calculateLeafNodesInPath(depth int) ([]int, [][]byte) {
	leafCounts := make([]int, depth)
	leafCountBytes := make([][]byte, depth)
	currentLeaves := 1

	for i := 0; i < depth; i++ {
		leafCounts[i] = currentLeaves

		// 将叶子数转换为字节（用于哈希）
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, int32(currentLeaves*2))
		leafCountBytes[i] = hashFunction(buf.Bytes())

		currentLeaves *= 2
	}

	return leafCounts, leafCountBytes
}

// ChunkData 将数据分块
func ChunkData(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := make([]byte, chunkSize)
		copy(chunk, data[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}
