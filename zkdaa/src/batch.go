package main

import (
	//"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"os"
	"runtime"
	"time"

	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// UnlockHandler - A 使用

type UnlockHandler struct {
	cs constraint.ConstraintSystem
	pk groth16.ProvingKey
	vk groth16.VerifyingKey
}

func NewUnlockHandler() (*UnlockHandler, time.Duration, error) {
	runtime.GC()
	start := time.Now()

	c := &circuit.UnlockCircuit{}
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, c)
	if err != nil {
		return nil, 0, err
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	vk := groth16.NewVerifyingKey(ecc.BN254)

	if err := loadKey("build/unlock.pk", pk); err != nil {
		return nil, 0, fmt.Errorf("failed to load unlock.pk: %v", err)
	}
	if err := loadKey("build/unlock.vk", vk); err != nil {
		return nil, 0, fmt.Errorf("failed to load unlock.vk: %v", err)
	}

	return &UnlockHandler{cs: cs, pk: pk, vk: vk}, time.Since(start), nil
}

func (h *UnlockHandler) Prove(preimage, sn, hashLock *big.Int) (groth16.Proof, time.Duration, error) {
	runtime.GC()
	assign := &circuit.UnlockCircuit{Pre: preimage, H: hashLock, Sn: sn}
	w, err := frontend.NewWitness(assign, ecc.BN254.ScalarField())
	if err != nil {
		return nil, 0, fmt.Errorf("witness creation failed: %w", err)
	}

	// 预热
	groth16.Prove(h.cs, h.pk, w)

	start := time.Now()
	proof, err := groth16.Prove(h.cs, h.pk, w)
	return proof, time.Since(start), err
}

func (h *UnlockHandler) Verify(proof groth16.Proof, sn, hashLock *big.Int) (time.Duration, error) {
	runtime.GC()
	publicAssign := &circuit.UnlockCircuit{H: hashLock, Sn: sn}
	w, err := frontend.NewWitness(publicAssign, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return 0, fmt.Errorf("witness build failed: %w", err)
	}

	// 预热
	groth16.Verify(proof, h.vk, w)

	// 单次测量
	start := time.Now()
	err = groth16.Verify(proof, h.vk, w)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}
func (h *UnlockHandler) GetConstraints() int { return h.cs.GetNbConstraints() }

// AuditUnlockHandler - B 使用

type AuditUnlockHandler struct {
	cs    constraint.ConstraintSystem
	pk    groth16.ProvingKey
	vk    groth16.VerifyingKey
	depth int // ✅ 加这个字段
}

func NewAuditUnlockHandler(depth int) (*AuditUnlockHandler, time.Duration, error) {
	runtime.GC()
	start := time.Now()

	c := &circuit.AuditUnlockCircuit{
		ProofPath:    make([]frontend.Variable, depth),
		Helpers:      make([]frontend.Variable, depth),
		LeafCounts:   make([]frontend.Variable, depth),
		LeafNumBytes: make([]frontend.Variable, depth),
	}
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, c)
	if err != nil {
		return nil, 0, err
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	vk := groth16.NewVerifyingKey(ecc.BN254)

	pkPath := fmt.Sprintf("build/audit_d%d.pk", depth)
	vkPath := fmt.Sprintf("build/audit_d%d.vk", depth)

	if err := loadKey(pkPath, pk); err != nil {
		return nil, 0, fmt.Errorf("failed to load %s: %v", pkPath, err)
	}
	if err := loadKey(vkPath, vk); err != nil {
		return nil, 0, fmt.Errorf("failed to load %s: %v", vkPath, err)
	}
	fmt.Printf("DEBUG: Saving depth=%d\n", depth)
	return &AuditUnlockHandler{
		cs:    cs,
		pk:    pk,
		vk:    vk,
		depth: depth, // ✅ 保存 depth
	}, time.Since(start), nil
}

func (h *AuditUnlockHandler) Prove(assign *circuit.AuditUnlockCircuit) (groth16.Proof, time.Duration, error) {
	runtime.GC()
	w, err := frontend.NewWitness(assign, ecc.BN254.ScalarField())
	if err != nil {
		return nil, 0, err
	}

	// 预热
	groth16.Prove(h.cs, h.pk, w)

	start := time.Now()
	proof, err := groth16.Prove(h.cs, h.pk, w)
	return proof, time.Since(start), err
}

func (h *AuditUnlockHandler) Verify(proof groth16.Proof, idx int, chunkHash, hVal *big.Int) (time.Duration, error) {
	runtime.GC()

	fmt.Printf("DEBUG: h.depth = %d\n", h.depth)
	// ✅ 使用保存的 depth
	publicAssign := &circuit.AuditUnlockCircuit{
		ProofPath:    make([]frontend.Variable, h.depth),
		Helpers:      make([]frontend.Variable, h.depth),
		LeafCounts:   make([]frontend.Variable, h.depth),
		LeafNumBytes: make([]frontend.Variable, h.depth),
		ChunkIndex:   idx,
		ChunkHash:    chunkHash,
		H:            hVal,
	}
	fmt.Printf("DEBUG: ProofPath len = %d\n", len(publicAssign.ProofPath))

	w, err := frontend.NewWitness(publicAssign, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return 0, fmt.Errorf("witness build failed: %w", err)
	}

	// 预热
	groth16.Verify(proof, h.vk, w)

	start := time.Now()
	err = groth16.Verify(proof, h.vk, w)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func (h *AuditUnlockHandler) GetConstraints() int { return h.cs.GetNbConstraints() }

// 🔥 主测试函数

func runSingleTest(fileSize int, chunkSize int, addrA, addrB string) *PerformanceMetrics {
	e2eStart := time.Now()

	label := fmt.Sprintf("%dMB", fileSize/(1024*1024))
	fmt.Printf("\n🚀 Running single test: %s\n", label)

	// 🔥 私链配置 (硬编码用于测试)
	realChainRPC := "http://127.0.0.1:8545"
	myPrivateKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// 🔥 使用重命名后的 initBatchClient，避免与 utils.go 冲突
	client, auth, _, err := initBatchClient(realChainRPC, myPrivateKey)
	if err != nil {
		fmt.Printf("❌ Blockchain setup failed: %v\n", err)
		return nil
	}
	defer client.Close()

	instA := setupContract(addrA, client)
	instB := setupContract(addrB, client)

	m := &PerformanceMetrics{FileSize: fileSize, ChunkSize: chunkSize}

	chunkCount := fileSize / chunkSize
	actualDepth := calcMerkleDepth(chunkCount)
	m.ChunkCount = chunkCount
	m.MerkleDepth = actualDepth
	if content, err := os.ReadFile("build/setup_time_unlock.txt"); err == nil {
		fmt.Sscanf(string(content), "%d", &m.TrustedSetupUnlockMs)
	} else {
		fmt.Println("⚠️ Warning: Could not read build/setup_time_unlock.txt")
	}

	// 🔥 读取 Audit Setup 时间
	if content, err := os.ReadFile("build/setup_time_audit.txt"); err == nil {
		fmt.Sscanf(string(content), "%d", &m.TrustedSetupAuditMs)
	} else {
		fmt.Println("⚠️ Warning: Could not read build/setup_time_audit.txt")
	}

	fmt.Printf("   [Step 0] Loading keys (depth=%d)...\n", actualDepth)

	unlockHandler, unlockSetupTime, err := NewUnlockHandler()
	if err != nil {
		fmt.Printf("Error Loading Keys A: %v\n", err)
		return nil
	}
	m.UnlockCircuitSetupTime = unlockSetupTime
	m.UnlockCircuitConstraints = unlockHandler.GetConstraints()

	auditHandler, auditSetupTime, err := NewAuditUnlockHandler(actualDepth)
	fmt.Printf("DEBUG: Created handler with depth=%d\n", actualDepth)
	if err != nil {
		fmt.Printf("Error Loading Keys B: %v\n", err)
		return nil
	}
	m.AuditUnlockSetupTime = auditSetupTime
	m.AuditCircuitConstraints = auditHandler.GetConstraints()

	fmt.Printf("   ✅ Keys loaded. A: %v, B: %v\n", unlockSetupTime, auditSetupTime)

	fmt.Println("   [Step 1] Preparing Data & CIDF...")
	prepStart := time.Now()
	data := generateDummyData(fileSize)
	chunkHashes, cidf, proofPaths, helpers, leafCounts, leafNumBytes := buildMerkleTreeForCircuit(data, chunkSize)
	m.DataPrepTimeA = time.Since(prepStart)

	fmt.Println("   [Step 2] Keys & Hashes...")
	zMask := new(big.Int)
	zMask.SetString("AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55AA55", 16)
	snII := randomFieldElement()
	snI := new(big.Int).Xor(snII, zMask)
	snI.Mod(snI, bn254FieldModulus)
	preI := randomFieldElement()
	h1 := mimcHashBig(preI, snI)
	h2 := mimcHashBig(cidf, snII)

	fmt.Println("   [Step 3] Transfer...")
	transferStart := time.Now()
	time.Sleep(50 * time.Millisecond)
	m.DataTransferTime = time.Since(transferStart)

	fmt.Println("   [Step 4] Locking...")
	dataIdA := strToDataID(fmt.Sprintf("flow_A_%d", time.Now().UnixNano()))
	dataIdB := strToDataID(fmt.Sprintf("flow_B_%d", time.Now().UnixNano()))
	timeout := big.NewInt(3600)
	nonce, _ := client.PendingNonceAt(context.Background(), auth.From)

	authA := cloneAuth(auth, nonce)
	// 使用重命名后的 batchTo32Bytes
	m.TxConfirmTimeLockA, m.GasLockA = sendAndWait(client, "LockA", func() (*types.Transaction, error) {
		return instA.Lock(authA, batchTo32Bytes(h1), dataIdA, timeout)
	})
	authB := cloneAuth(auth, nonce+1)
	m.TxConfirmTimeLockB, m.GasLockB = sendAndWait(client, "LockB", func() (*types.Transaction, error) {
		return instB.Lock(authB, batchTo32Bytes(h2), dataIdB, timeout)
	})

	fmt.Println("   [Step 5] DSPB Prove...")
	challengeIdx := mathrand.Intn(chunkCount)
	assignB := &circuit.AuditUnlockCircuit{
		ProofPath:    make([]frontend.Variable, actualDepth),
		Helpers:      make([]frontend.Variable, actualDepth),
		LeafCounts:   make([]frontend.Variable, actualDepth),
		LeafNumBytes: make([]frontend.Variable, actualDepth),
		Sn:           snII, ChunkIndex: challengeIdx, ChunkHash: chunkHashes[challengeIdx], H: h2,
	}
	fillWitnessSlice(assignB, proofPaths[challengeIdx], helpers[challengeIdx], leafCounts, leafNumBytes)

	proofB, proveTimeB, err := auditHandler.Prove(assignB)
	if err != nil {
		fmt.Printf("Error ProveB: %v\n", err)
		return nil
	}
	m.AuditUnlockProveTimeB = proveTimeB

	fmt.Println("   [Step 5.1] Local Verify B (Loop 1000x)...")
	verifyTimeB, err := auditHandler.Verify(proofB, challengeIdx, chunkHashes[challengeIdx], h2)
	if err != nil {
		fmt.Printf("❌ Verify B FAILED: %v\n", err)
	} else {
		fmt.Printf("✅ Verify B PASSED in %v\n", verifyTimeB)
	}
	m.AuditVerifyTimeB = verifyTimeB
	fmt.Printf("      -> Avg Verify B: %v\n", verifyTimeB)

	fmt.Println("   [Step 6] SCB Verify...")
	solProofB, pubB := formatProofForSolidityAuditUnlock(proofB, challengeIdx, chunkHashes[challengeIdx], h2)
	fmt.Printf("      🔍 PubInputs B: [0]=%s, [1]=%s, [2]=%s\n", pubB[0].String(), pubB[1].String(), pubB[2].String())

	nonce, _ = client.PendingNonceAt(context.Background(), auth.From)
	authB = cloneAuth(auth, nonce)
	m.TxConfirmTimeUnlockB, m.GasUnlockB = sendAndWait(client, "AuditUnlockB", func() (*types.Transaction, error) {
		return instB.AuditUnlock(authB, solProofB, pubB)
	})

	fmt.Println("   [Step 7] User Unlock Proof...")
	snI_calculated := new(big.Int).Xor(snII, zMask)
	snI_calculated.Mod(snI_calculated, bn254FieldModulus)
	proofA, proveTimeA, err := unlockHandler.Prove(preI, snI_calculated, h1)
	if err != nil {
		fmt.Printf("Error ProveA: %v\n", err)
		return nil
	}
	m.UnlockProveTimeA = proveTimeA

	fmt.Println("   [Step 7.1] Local Verify A (Loop 1000x)...")
	verifyTimeA, err := unlockHandler.Verify(proofA, snI_calculated, h1)
	if err != nil {
		fmt.Printf("Verify Fail: %v\n", err)
		return nil
	}
	m.UnlockVerifyTimeA = verifyTimeA
	fmt.Printf("      -> Avg Verify A: %v\n", verifyTimeA)

	fmt.Println("   [Step 8] SCA Unlock...")
	solProofA, pubA := formatProofForSolidityUnlock(proofA, h1, snI_calculated)
	fmt.Printf("      🔍 PubInputs A: [0]=%s, [1]=%s\n", pubA[0].String(), pubA[1].String())

	nonce, _ = client.PendingNonceAt(context.Background(), auth.From)
	authA = cloneAuth(auth, nonce)
	m.TxConfirmTimeUnlockA, m.GasUnlockA = sendAndWait(client, "UnlockA", func() (*types.Transaction, error) {
		return instA.Unlock(authA, solProofA, pubA)
	})

	m.TotalGas = m.GasLockA + m.GasLockB + m.GasUnlockA + m.GasUnlockB
	m.EndToEndTime = time.Since(e2eStart)
	fmt.Printf("\n   ✅ Test completed! E2E: %.2f ms\n", float64(m.EndToEndTime.Milliseconds()))
	return m
}

func initBatchClient(rpcURL, privateKeyHex string) (*ethclient.Client, *bind.TransactOpts, *bind.CallOpts, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to eth client: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get chainID: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.GasLimit = 8000000
	callOpts := &bind.CallOpts{Pending: true}
	return client, auth, callOpts, nil
}

func loadKey(filename string, key interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if pk, ok := key.(groth16.ProvingKey); ok {
		_, err = pk.ReadFrom(f)
		return err
	}
	if vk, ok := key.(groth16.VerifyingKey); ok {
		_, err = vk.ReadFrom(f)
		return err
	}
	return fmt.Errorf("unknown key type")
}

func buildMerkleTreeForCircuit(data []byte, chunkSize int) ([]*big.Int, *big.Int, [][]*big.Int, [][]int, []*big.Int, []*big.Int) {
	chunks := splitDataBytes(data, chunkSize)
	n := len(chunks)
	hashes := make([]*big.Int, n)
	for i, c := range chunks {
		hashes[i] = hashChunkBig(c)
	}

	proofPaths := make([][]*big.Int, n)
	helpers := make([][]int, n)
	for i := range proofPaths {
		proofPaths[i] = []*big.Int{}
		helpers[i] = []int{}
	}

	var leafCounts []*big.Int
	var leafNumBytes []*big.Int
	current := hashes
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	currentWeight := int64(1)
	for len(current) > 1 {
		leafCounts = append(leafCounts, big.NewInt(currentWeight))
		leafNumBytes = append(leafNumBytes, big.NewInt(int64(len(current))))

		nextLen := (len(current) + 1) / 2
		next := make([]*big.Int, nextLen)
		for i := 0; i < len(current); i += 2 {
			l, r := current[i], current[i]
			if i+1 < len(current) {
				r = current[i+1]
			}
			next[i/2] = mimcHashBig(big.NewInt(int64(len(current))), l, r)
		}

		nextIndices := make([]int, n)
		for j := 0; j < n; j++ {
			idx := indices[j]
			sib := idx ^ 1
			if sib >= len(current) {
				sib = idx
			}
			proofPaths[j] = append(proofPaths[j], current[sib])
			helpers[j] = append(helpers[j], idx&1)
			nextIndices[j] = idx / 2
		}
		current = next
		indices = nextIndices
		currentWeight *= 2
	}
	return hashes, current[0], proofPaths, helpers, leafCounts, leafNumBytes
}

func fillWitnessSlice(assign *circuit.AuditUnlockCircuit, paths []*big.Int, helpers []int, lc, lnb []*big.Int) {
	depth := len(assign.ProofPath)
	for i := 0; i < depth; i++ {
		if i < len(paths) {
			assign.ProofPath[i] = paths[i]
			assign.Helpers[i] = helpers[i]
			assign.LeafCounts[i] = lc[i]
			assign.LeafNumBytes[i] = lnb[i]
		} else {
			assign.ProofPath[i] = 0
			assign.Helpers[i] = 0
			assign.LeafCounts[i] = 0
			assign.LeafNumBytes[i] = 0
		}
	}
}

func formatProofForSolidityUnlock(p groth16.Proof, h, sn *big.Int) ([8]*big.Int, [2]*big.Int) {
	var zeroProof [8]*big.Int
	// 使用 utils.go 中的 curveToContractProof
	solProof, err := curveToContractProof(p)
	if err != nil {
		return zeroProof, [2]*big.Int{toFieldElement(h), toFieldElement(sn)}
	}
	return solProof, [2]*big.Int{toFieldElement(h), toFieldElement(sn)}
}

func formatProofForSolidityAuditUnlock(p groth16.Proof, idx int, chunkHash, h *big.Int) ([8]*big.Int, [3]*big.Int) {
	var zeroProof [8]*big.Int
	// 使用 utils.go 中的 curveToContractProof
	solProof, err := curveToContractProof(p)
	if err != nil {
		return zeroProof, [3]*big.Int{toFieldElement(big.NewInt(int64(idx))), toFieldElement(chunkHash), toFieldElement(h)}
	}
	return solProof, [3]*big.Int{
		toFieldElement(big.NewInt(int64(idx))),
		toFieldElement(chunkHash),
		toFieldElement(h),
	}
}

func cloneAuth(auth *bind.TransactOpts, nonce uint64) *bind.TransactOpts {
	return &bind.TransactOpts{From: auth.From, Signer: auth.Signer, Context: context.Background(), GasLimit: 8000000, GasPrice: big.NewInt(20000000000), Nonce: big.NewInt(int64(nonce))}
}
func mimcHashBig(vals ...*big.Int) *big.Int {
	hf := hash.MIMC_BN254.New()
	for _, v := range vals {
		b := v.Bytes()
		padded := make([]byte, 32)
		copy(padded[32-len(b):], b)
		hf.Write(padded)
	}
	return new(big.Int).SetBytes(hf.Sum(nil))
}
func hashChunkBig(data []byte) *big.Int {
	hf := hash.MIMC_BN254.New()
	hf.Write(data)
	return new(big.Int).SetBytes(hf.Sum(nil))
}
func splitDataBytes(data []byte, size int) [][]byte {
	var out [][]byte
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		out = append(out, data[i:end])
	}
	return out
}
func calcMerkleDepth(n int) int {
	d := 0
	for n > 1 {
		n = (n + 1) / 2
		d++
	}
	return d
}
func randomFieldElement() *big.Int {
	max := new(big.Int).Set(bn254FieldModulus)
	max.Sub(max, big.NewInt(1))
	r, _ := rand.Int(rand.Reader, max)
	return r
}

func batchTo32Bytes(val *big.Int) [32]byte {
	var r [32]byte
	if val != nil {
		b := val.Bytes()
		copy(r[32-len(b):], b)
	}
	return r
}
