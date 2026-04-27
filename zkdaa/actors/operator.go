package actors

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"zk-htlc/merkle"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/hash"
)

// Operator 批量锁定的运营商
type Operator struct {
	TxLocks    []TxLock                // 收集到的锁定交易
	MerkleTree *merkle.BatchMerkleTree // 构建的 Merkle 树
	MerkleRoot *big.Int                // 提交到合约的根
}

// TxLock 单个锁定交易
type TxLock struct {
	Preimage     *big.Int // preimage
	SerialNumber *big.Int // serial number
	H            *big.Int // hash(preimage, sn)
	Index        int      // 在批量中的索引
}

// NewOperator 创建运营商
func NewOperator() *Operator {
	return &Operator{
		TxLocks: make([]TxLock, 0),
	}
}

// CollectLock 收集一个锁定交易
func (op *Operator) CollectLock(preimage, serialNumber, h *big.Int) {
	op.TxLocks = append(op.TxLocks, TxLock{
		Preimage:     preimage,
		SerialNumber: serialNumber,
		H:            h,
		Index:        len(op.TxLocks),
	})
}

// BuildMerkleTree 构建批量 Merkle 树
func (op *Operator) BuildMerkleTree() error {
	if len(op.TxLocks) == 0 {
		return fmt.Errorf("没有交易可构建 Merkle 树")
	}

	// 使用每个 TxLock 的 H 作为叶子
	leaves := make([]*big.Int, len(op.TxLocks))
	for i, tx := range op.TxLocks {
		leaves[i] = tx.H
	}

	// 构建批量 Merkle 树
	op.MerkleTree = merkle.NewBatchMerkleTree(leaves)
	op.MerkleRoot = op.MerkleTree.Root
	return nil
}

// GetProofPath 获取某个索引的 Merkle 路径
func (op *Operator) GetProofPath(index int) ([]merkle.ProofElement, error) {
	if op.MerkleTree == nil {
		return nil, fmt.Errorf("Merkle 树未构建")
	}
	return op.MerkleTree.GetProof(index)
}

// GenerateMockBatch 生成模拟批量交易(用于测试)
func (op *Operator) GenerateMockBatch(batchSize int) error {
	for i := 0; i < batchSize; i++ {
		// 生成随机 preimage 和 serial number
		preimage, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
		serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

		// 🔧 修复: 使用与电路一致的 MiMC 哈希
		h := hashMiMC(preimage, serialNumber)

		op.CollectLock(preimage, serialNumber, h)
	}

	return op.BuildMerkleTree()
}

// 需要保证链下计算与电路中的行为完全一致
func hashMiMC(preimage, serialNumber *big.Int) *big.Int {
	hFunc := hash.MIMC_BN254.New()

	// 将 big.Int 转换为有限域元素,然后再转换为字节
	// 这样与电路中的 frontend.Variable 行为一致
	preimageFr := new(fr.Element).SetBigInt(preimage)
	preimageBytes := preimageFr.Bytes()
	hFunc.Write(preimageBytes[:])

	serialNumberFr := new(fr.Element).SetBigInt(serialNumber)
	serialNumberBytes := serialNumberFr.Bytes()
	hFunc.Write(serialNumberBytes[:])

	result := hFunc.Sum(nil)
	hash := new(big.Int).SetBytes(result)

	// 确保结果在有限域内
	hashFr := new(fr.Element).SetBigInt(hash)
	return hashFr.BigInt(new(big.Int))
}
