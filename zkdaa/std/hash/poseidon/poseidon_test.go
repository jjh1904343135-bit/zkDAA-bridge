package poseidon

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

// TestPoseidon tests the Gnark Poseidon implementation against Iden3's Go implementation on all the test
// vectors outlined in the original paper's reference repository, which can be found here: https://extgit.iaik.tugraz.at/krypto/hadeshash/-/tree/master/code
//
// The actual test vectors are outlined here: https://extgit.iaik.tugraz.at/krypto/hadeshash/-/blob/master/code/test_vectors.txt
// We have included more for the sake of robustness.
// Note that our implementation is focused on the 3-input variant with an x^5 S-box, so not all the test vectors apply.
func TestPoseidon(t *testing.T) {
	var datasize = 63488000
	var dag_size = 256000
	var leafnum = datasize / dag_size
	var block_size = 64
	//var chunksize = dag_size / block_size
	mod := ecc.BN254.ScalarField()
	var leaves = make([][]byte, leafnum)
	for i := 0; i < len(leaves); i++ {
		leaf, _ := rand.Int(rand.Reader, mod)
		//assert.NoError(err)
		b := leaf.Bytes()
		if len(b) < dag_size {
			// 创建一个32字节的数组，并在前面填充零
			padded := make([]byte, dag_size)
			copy(padded[dag_size-len(b):], b)
			b = padded
		}
		leaves[i] = b
	}

	circuit_input := make([]frontend.Variable, dag_size/block_size)
	var result = make([][]byte, dag_size/block_size)
	var result_int = make([]*big.Int, dag_size/block_size)
	var index = 0
	for i := 0; i < len(leaves[0]); i += block_size {
		end := i + block_size
		if end > len(leaves[0]) {
			end = len(leaves[0])
		}
		result[index] = leaves[0][i:end]
		//bigInt := new(big.Int)
		//bytes := bigInt.SetBytes(result[index])
		bytes := int64(binary.BigEndian.Uint64(result[index]))

		circuit_input[index] = bytes
		result_int[index] = big.NewInt(bytes)

		index++
	}

	tests := map[string]struct {
		gnarkPoseidonInput     [4000]frontend.Variable
		referencePoseidonInput []*big.Int
	}{
		"happy path: basic input": {
			gnarkPoseidonInput:     [4000]frontend.Variable(circuit_input),
			referencePoseidonInput: result_int,
		},
		//"official test vector: poseidonperm_x5_254_3": {
		//	gnarkPoseidonInput: [17]frontend.Variable{0, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
		//	referencePoseidonInput: []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(2),
		//		big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(2)},
		//},
		//"zero vector": {
		//	gnarkPoseidonInput:     [3]frontend.Variable{0, 0, 0},
		//	referencePoseidonInput: []*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
		//},
		//"larger inputs": {
		//	gnarkPoseidonInput:     [3]frontend.Variable{129048, 990217, 2234383333},
		//	referencePoseidonInput: []*big.Int{big.NewInt(129048), big.NewInt(990217), big.NewInt(2234383333)},
		//},
		//"decreasing vector inputs": {
		//	gnarkPoseidonInput:     [3]frontend.Variable{10000000, 10000, 100},
		//	referencePoseidonInput: []*big.Int{big.NewInt(10000000), big.NewInt(10000), big.NewInt(100)},
		//},
	}

	for name, testCase := range tests {

		assert := test.NewAssert(t)
		var circuit circuitPoseidon

		// Compute reference hash to test against
		//分批计算
		var referenceHash *big.Int
		var input = make([]*big.Int, 16)
		var dirty bool
		//var index_circuit = 0
		for j := 0; j < 16; j++ {
			input[j] = new(big.Int)

		}
		k := 0
		for i := 0; i < len(testCase.referencePoseidonInput); i++ {
			dirty = true
			input[k] = testCase.referencePoseidonInput[i]
			if k == 15 {
				referenceHash, _ = poseidon.Hash(input)
				dirty = false
				//if err != nil {
				//	return nil, err
				//}
				input = make([]*big.Int, 16)
				input[0] = referenceHash
				for j := 1; j < 16; j++ {
					input[j] = new(big.Int)
				}
				k = 1
			} else {
				k++
			}
		}
		if dirty {
			var final_input = make([]*big.Int, k)
			for i := 0; i < k; i++ {
				final_input[i] = input[i]
			}
			// we haven't hashed something in the main sponge loop and need to do hash here
			referenceHash, _ = poseidon.Hash(final_input)
		}

		//if err != nil {
		//	t.Fatal(err, "Failed to compute reference poseidon hash for test case: ", name)
		//}
		t.Logf("Reference hash: %s", referenceHash.String())

		// Generate poseidon hash using gnark implementation
		assert.ProverSucceeded(&circuit, &circuitPoseidon{
			A:    testCase.gnarkPoseidonInput,
			Hash: referenceHash,
		}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

		// Ensure output correctly compiles
		_r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
		if err != nil {
			t.Fatal(err, "Failed to compile computed poseidon hash for test case: ", name)
		}

		// Sanity check and debugging support for internal variables
		internal, secret, public := _r1cs.GetNbVariables()
		t.Logf("Public, secret, internal %v, %v, %v\n", public, secret, internal)
	}
}

// --- Test Helpers ---

type circuitPoseidon struct {
	A    [4000]frontend.Variable `gnark:",public"`
	Hash frontend.Variable       `gnark:",public"`
}

func (t *circuitPoseidon) Define(api frontend.API) error {

	var out frontend.Variable
	var input = make([]frontend.Variable, 16)
	var dirty bool
	for j := 0; j < 16; j++ {
		input[j] = new(frontend.Variable)
	}
	k := 0
	for i := 0; i < len(t.A); i++ {
		dirty = true
		input[k] = t.A[i]
		if k == 15 {
			out = Poseidon(api, input)
			//out = PoseidonEx(api, input, 0, 1)
			dirty = false
			//if err != nil {
			//	return nil, err
			//}
			input = make([]frontend.Variable, 16)
			input[0] = out
			for j := 1; j < 16; j++ {
				input[j] = new(frontend.Variable)
			}
			k = 1
		} else {
			k++
		}
	}
	if dirty {
		var final_input = make([]frontend.Variable, k)
		for i := 0; i < k; i++ {
			final_input[i] = input[i]
		}
		// we haven't hashed something in the main sponge loop and need to do hash here
		out = Poseidon(api, final_input)
	}

	//hash := Poseidon(api, t.A[:])
	api.Println(t.Hash)
	api.Println(out)
	api.AssertIsEqual(out, t.Hash)
	return nil
}
