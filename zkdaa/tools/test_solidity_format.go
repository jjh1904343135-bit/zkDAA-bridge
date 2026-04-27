package main

import (
	"fmt"
	"math/big"
	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
)

//func main() {
	// 使用简单的测试值
	testMerkleRoot := big.NewInt(12345678)
	testSerialNumber := big.NewInt(999)

	assignment := &circuit.BatchUnlockCircuit{
		Preimage:           big.NewInt(111),
		SerialNumber:       testSerialNumber,
		TxIndex:            big.NewInt(0),
		ProofPath:          []frontend.Variable{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4)},
		Helpers:            []frontend.Variable{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		LeafCounts:         []frontend.Variable{big.NewInt(16), big.NewInt(8), big.NewInt(4), big.NewInt(2)},
		LeafNumBytes:       []frontend.Variable{big.NewInt(16), big.NewInt(8), big.NewInt(4), big.NewInt(2)},
		MerkleRoot:         testMerkleRoot,
		SerialNumberPublic: testSerialNumber,
	}

	fmt.Println("赋值的公开字段:")
	fmt.Printf("  MerkleRoot: %s\n", testMerkleRoot.String())
	fmt.Printf("  SerialNumberPublic: %s\n", testSerialNumber.String())

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(fmt.Sprintf("创建 witness 失败: %v", err))
	}

	publicWitness, err := witness.Public()
	if err != nil {
		panic(fmt.Sprintf("提取公开 witness 失败: %v", err))
	}

	publicVector, ok := publicWitness.Vector().(fr.Vector)
	if !ok {
		panic("类型断言失败")
	}

	fmt.Printf("\n公开输入向量长度: %d\n", len(publicVector))

	for i := 0; i < len(publicVector); i++ {
		val := new(big.Int)
		publicVector[i].BigInt(val)
		fmt.Printf("  publicVector[%d]: %s\n", i, val.String())
	}

	// 检查是否匹配
	fmt.Println("\n预期:")
	fmt.Println("  publicVector[0] = 1 (固定值)")
	fmt.Printf("  publicVector[1] = %s (MerkleRoot)\n", testMerkleRoot.String())
	fmt.Printf("  publicVector[2] = %s (SerialNumberPublic)\n", testSerialNumber.String())
}
