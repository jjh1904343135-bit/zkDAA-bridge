package main

import (
	"fmt"
	"log"
	"math/big"
	"zk-htlc/actors"
	"zk-htlc/circuit"
	"zk-htlc/zkp"

	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
)

//func main() {
	fmt.Println("🔍 证明调试工具\n")

	batchSize := 16

	// Step 1: 生成测试数据
	fmt.Println("Step 1: 生成测试数据...")
	operator := actors.NewOperator()
	if err := operator.GenerateMockBatch(batchSize); err != nil {
		log.Fatal(err)
	}

	tx := operator.TxLocks[0]
	fmt.Printf("✓ Preimage: %x\n", tx.Preimage.Bytes()[:8])
	fmt.Printf("✓ SerialNumber: %x\n", tx.SerialNumber.Bytes()[:8])
	fmt.Printf("✓ 叶子哈希 (H): %x\n", tx.H.Bytes()[:8])
	fmt.Printf("✓ Merkle Root: %x\n", operator.MerkleRoot.Bytes()[:8])

	// Step 2: 手动验证叶子哈希
	fmt.Println("\nStep 2: 验证叶子哈希计算...")
	manualHash := computeLeafHash(tx.Preimage, tx.SerialNumber)
	fmt.Printf("✓ 手动计算: %x\n", manualHash.Bytes()[:8])
	fmt.Printf("✓ TxLock.H: %x\n", tx.H.Bytes()[:8])
	if manualHash.Cmp(tx.H) == 0 {
		fmt.Println("✅ 叶子哈希匹配")
	} else {
		fmt.Println("❌ 叶子哈希不匹配！")
		return
	}

	// Step 3: 验证 Merkle 证明
	fmt.Println("\nStep 3: 验证链下 Merkle 证明...")
	proofElements, err := operator.GetProofPath(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ 证明长度: %d\n", len(proofElements))

	isValid := operator.MerkleTree.VerifyProof(tx.H, 0, proofElements)
	if isValid {
		fmt.Println("✅ 链下 Merkle 验证通过")
	} else {
		fmt.Println("❌ 链下 Merkle 验证失败！")
		return
	}

	// Step 4: 打印证明路径详情
	fmt.Println("\nStep 4: 证明路径详情...")
	for i, elem := range proofElements {
		fmt.Printf("[%d] Hash: %x, IsLeft: %v\n",
			i, elem.Hash.Bytes()[:8], elem.IsLeft)
	}

	// Step 5: 手动计算 Merkle 根
	fmt.Println("\nStep 5: 手动沿路径计算根...")
	computedHash := tx.H
	for i, elem := range proofElements {
		var left, right *big.Int
		if elem.IsLeft {
			left = elem.Hash
			right = computedHash
		} else {
			left = computedHash
			right = elem.Hash
		}

		hFunc := hash.MIMC_BN254.New()
		hFunc.Write(padTo32Bytes(left.Bytes()))
		hFunc.Write(padTo32Bytes(right.Bytes()))
		computedHash = new(big.Int).SetBytes(hFunc.Sum(nil))

		fmt.Printf("[%d] 计算哈希: %x\n", i, computedHash.Bytes()[:8])
	}

	fmt.Printf("✓ 计算的根: %x\n", computedHash.Bytes()[:8])
	fmt.Printf("✓ 实际的根: %x\n", operator.MerkleRoot.Bytes()[:8])
	if computedHash.Cmp(operator.MerkleRoot) == 0 {
		fmt.Println("✅ 手动验证通过")
	} else {
		fmt.Println("❌ 手动验证失败！")
		return
	}

	// Step 6: 测试电路
	fmt.Println("\nStep 6: 测试电路证明...")
	handler, err := zkp.NewBatchZKPHandler(batchSize)
	if err != nil {
		log.Fatal(err)
	}

	depth := handler.GetMerkleDepth()
	proofPath := make([]frontend.Variable, depth)
	helpers := make([]frontend.Variable, depth)
	leafCounts := make([]frontend.Variable, depth)
	leafNumBytes := make([]frontend.Variable, depth)

	for i := 0; i < len(proofElements) && i < depth; i++ {
		proofPath[i] = proofElements[i].Hash
		if proofElements[i].IsLeft {
			helpers[i] = 0
		} else {
			helpers[i] = 1
		}
		leafCounts[i] = proofElements[i].LeafCount
		leafNumBytes[i] = proofElements[i].LeafNumBytes
	}

	// 填充剩余
	for i := len(proofElements); i < depth; i++ {
		proofPath[i] = big.NewInt(0)
		helpers[i] = big.NewInt(0)
		leafCounts[i] = big.NewInt(1)
		leafNumBytes[i] = big.NewInt(32)
	}

	assignment := &circuit.BatchUnlockCircuit{
		Preimage:           tx.Preimage,
		SerialNumber:       tx.SerialNumber,
		TxIndex:            big.NewInt(int64(tx.Index)),
		ProofPath:          proofPath,
		Helpers:            helpers,
		LeafCounts:         leafCounts,
		LeafNumBytes:       leafNumBytes,
		MerkleRoot:         operator.MerkleRoot,
		SerialNumberPublic: tx.SerialNumber,
	}

	proof, err := handler.Prove(assignment)
	if err != nil {
		fmt.Printf("❌ Prove 失败: %v\n", err)
		return
	}
	fmt.Println("✅ Prove 成功")

	err = handler.Verify(proof, []*big.Int{operator.MerkleRoot, tx.SerialNumber})
	if err != nil {
		fmt.Printf("❌ Verify 失败: %v\n", err)
		return
	}
	fmt.Println("✅ Verify 成功")

	fmt.Println("\n🎉 所有链下测试通过！")
	fmt.Println("\n💡 如果链上验证失败，问题可能是：")
	fmt.Println("   1. 验证器合约与当前电路不匹配")
	fmt.Println("   2. 公开输入的顺序在合约中不对")
	fmt.Println("   3. 证明格式化函数有问题")
}

func computeLeafHash(preimage, serialNumber *big.Int) *big.Int {
	hFunc := hash.MIMC_BN254.New()
	hFunc.Write(padTo32Bytes(preimage.Bytes()))
	hFunc.Write(padTo32Bytes(serialNumber.Bytes()))
	return new(big.Int).SetBytes(hFunc.Sum(nil))
}

func padTo32Bytes(data []byte) []byte {
	if len(data) >= 32 {
		return data[len(data)-32:]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(data):], data)
	return padded
}
