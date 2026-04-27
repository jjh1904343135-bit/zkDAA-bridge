package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"sort"
	"time"

	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

const chunkSize2 = 1500

var testCases2 = []struct {
	SizeMB int
	Label  string
}{
	{8, "8 MB"},
	{16, "16 MB"},
	{32, "32 MB"},
	{64, "64 MB"},
	{128, "128 MB"},
}

func sizeToDepth2(sizeMB int) int {
	nLeaves := sizeMB * 1024 * 1024 / chunkSize2
	return int(math.Ceil(math.Log2(float64(nLeaves))))
}

// 工具函数

func mimcChain2(inputs ...*big.Int) *big.Int {
	h := hash.MIMC_BN254.New()
	for _, v := range inputs {
		h.Write(to322(v))
	}
	return new(big.Int).SetBytes(h.Sum(nil))
}

func to322(n *big.Int) []byte {
	b := n.Bytes()
	if len(b) >= 32 {
		return b
	}
	p := make([]byte, 32)
	copy(p[32-len(b):], b)
	return p
}

func randBytes2(size int) []byte {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		log.Fatalf("randBytes: %v", err)
	}
	return data
}

func hashBytes2(b []byte) *big.Int {
	h := hash.MIMC_BN254.New()
	h.Write(b)
	return new(big.Int).SetBytes(h.Sum(nil))
}

// Merkle 树（支持任意叶节点的证明路径）

type MerkleTree2 struct {
	Depth  int
	nodes  [][]*big.Int
	leaves []*big.Int
}

func buildTree2(fileData []byte, depth int) *MerkleTree2 {
	numLeaves := 1 << depth
	t := &MerkleTree2{
		Depth:  depth,
		leaves: make([]*big.Int, numLeaves),
	}
	t.nodes = make([][]*big.Int, depth+1)
	t.nodes[0] = make([]*big.Int, numLeaves)

	for i := 0; i < numLeaves; i++ {
		start := i * chunkSize2
		chunk := make([]byte, chunkSize2)
		if start < len(fileData) {
			end := start + chunkSize2
			if end > len(fileData) {
				end = len(fileData)
			}
			copy(chunk, fileData[start:end])
		}
		t.nodes[0][i] = hashBytes2(chunk)
		t.leaves[i] = t.nodes[0][i]
	}
	for level := 1; level <= depth; level++ {
		prev := t.nodes[level-1]
		n := len(prev) / 2
		t.nodes[level] = make([]*big.Int, n)
		nb := big.NewInt(int64(1 << level))
		for i := 0; i < n; i++ {
			t.nodes[level][i] = mimcChain2(nb, prev[2*i], prev[2*i+1])
		}
	}
	return t
}

func (t *MerkleTree2) root() *big.Int { return t.nodes[t.Depth][0] }

// getProofFor 返回任意叶节点的证明路径
func (t *MerkleTree2) getProofFor(leafIdx int) (
	path, helper, leafNum, leafNumByte []frontend.Variable,
) {
	d := t.Depth
	path = make([]frontend.Variable, d)
	helper = make([]frontend.Variable, d)
	leafNum = make([]frontend.Variable, d)
	leafNumByte = make([]frontend.Variable, d)
	idx := leafIdx
	for i := 0; i < d; i++ {
		isRight := idx % 2
		siblingIdx := idx ^ 1
		path[i] = t.nodes[i][siblingIdx]
		helper[i] = big.NewInt(int64(isRight))
		leafNum[i] = big.NewInt(int64(1 << i))
		leafNumByte[i] = big.NewInt(int64(1 << (i + 1)))
		idx >>= 1
	}
	return
}

func lnbSlice2(depth int) []frontend.Variable {
	s := make([]frontend.Variable, depth)
	for i := range s {
		s[i] = big.NewInt(int64(1 << (i + 1)))
	}
	return s
}

// 新根计算（Go 侧，用于填充 witness 的 NewRootHash）

func calcInsertRoot(t *MerkleTree2, newLeaf *big.Int, leafIdx int) *big.Int {
	computed := mimcChain2(newLeaf, t.leaves[leafIdx])
	idx := leafIdx
	for i := 0; i < t.Depth; i++ {
		nb := big.NewInt(int64((1 << (i + 1)) + 1))
		sib := t.nodes[i][idx^1]
		if idx%2 == 1 {
			computed = mimcChain2(nb, sib, computed)
		} else {
			computed = mimcChain2(nb, computed, sib)
		}
		idx >>= 1
	}
	return computed
}

func calcModifyRoot(t *MerkleTree2, newLeaf *big.Int, leafIdx int) *big.Int {
	computed := newLeaf
	idx := leafIdx
	for i := 0; i < t.Depth; i++ {
		nb := big.NewInt(int64(1 << (i + 1)))
		sib := t.nodes[i][idx^1]
		if idx%2 == 1 {
			computed = mimcChain2(nb, sib, computed)
		} else {
			computed = mimcChain2(nb, computed, sib)
		}
		idx >>= 1
	}
	return computed
}

func calcDeleteRoot(t *MerkleTree2, leafIdx int) *big.Int {
	computed := t.nodes[0][leafIdx^1]
	idx := leafIdx >> 1
	for i := 1; i < t.Depth; i++ {
		nb := big.NewInt(int64((1 << (i + 1)) - 1))
		sib := t.nodes[i][idx^1]
		if idx%2 == 1 {
			computed = mimcChain2(nb, sib, computed)
		} else {
			computed = mimcChain2(nb, computed, sib)
		}
		idx >>= 1
	}
	return computed
}

// Witness 构造（每个 leafIdx 独立）

func makeInsertW2(t *MerkleTree2, newLeaf *big.Int, leafIdx int) *circuit.MerkleInsertCircuit {
	depth := t.Depth
	path, helper, leafNum, _ := t.getProofFor(leafIdx)
	newRoot := calcInsertRoot(t, newLeaf, leafIdx)
	nnb := make([]frontend.Variable, depth)
	for i := range nnb {
		nnb[i] = big.NewInt(int64((1 << (i + 1)) + 1))
	}
	return &circuit.MerkleInsertCircuit{
		LeafHash: newLeaf, NewRootHash: newRoot,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf: t.leaves[leafIdx], RootHash: t.root(),
			LeafIndex: big.NewInt(int64(leafIdx)),
			Path:      path, LeafNum: leafNum, Helper: helper,
			LeafNumByte: lnbSlice2(depth),
		},
		NewNum_byte: nnb,
	}
}

func makeModifyW2(t *MerkleTree2, newLeaf *big.Int, leafIdx int) *circuit.MerkleUpdateCircuit {
	path, helper, leafNum, leafNumByte := t.getProofFor(leafIdx)
	newRoot := calcModifyRoot(t, newLeaf, leafIdx)
	return &circuit.MerkleUpdateCircuit{
		LeafHash: newLeaf, NewRootHash: newRoot,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf: t.leaves[leafIdx], RootHash: t.root(),
			LeafIndex: big.NewInt(int64(leafIdx)),
			Path:      path, LeafNum: leafNum, Helper: helper, LeafNumByte: leafNumByte,
		},
	}
}

func makeDeleteW2(t *MerkleTree2, leafIdx int) *circuit.MerkleDeleteCircuit {
	depth := t.Depth
	path, helper, leafNum, _ := t.getProofFor(leafIdx)
	newRoot := calcDeleteRoot(t, leafIdx)
	npb := make([]frontend.Variable, depth)
	npb[0] = big.NewInt(0)
	for i := 1; i < depth; i++ {
		npb[i] = big.NewInt(int64((1 << (i + 1)) - 1))
	}
	return &circuit.MerkleDeleteCircuit{
		NewRootHash: newRoot,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf: t.leaves[leafIdx], RootHash: t.root(),
			LeafIndex: big.NewInt(int64(leafIdx)),
			Path:      path, LeafNum: leafNum, Helper: helper,
			LeafNumByte: lnbSlice2(depth),
		},
		NewPath_byte: npb,
	}
}

// 电路编译 & Setup

func makeTpls2(depth int) (
	ins *circuit.MerkleInsertCircuit,
	upd *circuit.MerkleUpdateCircuit,
	del *circuit.MerkleDeleteCircuit,
) {
	mk := func() circuit.MerkleProofCircuit {
		return circuit.MerkleProofCircuit{
			Path: make([]frontend.Variable, depth), LeafNum: make([]frontend.Variable, depth),
			Helper: make([]frontend.Variable, depth), LeafNumByte: make([]frontend.Variable, depth),
		}
	}
	ins = &circuit.MerkleInsertCircuit{Circuit_merkle: mk(), NewNum_byte: make([]frontend.Variable, depth)}
	upd = &circuit.MerkleUpdateCircuit{Circuit_merkle: mk()}
	del = &circuit.MerkleDeleteCircuit{Circuit_merkle: mk(), NewPath_byte: make([]frontend.Variable, depth)}
	return
}

type compiled2 struct {
	Depth    int
	Label    string
	InsertCS constraint.ConstraintSystem
	InsertPK groth16.ProvingKey
	InsertVK groth16.VerifyingKey
	ModifyCS constraint.ConstraintSystem
	ModifyPK groth16.ProvingKey
	ModifyVK groth16.VerifyingKey
	DeleteCS constraint.ConstraintSystem
	DeletePK groth16.ProvingKey
	DeleteVK groth16.VerifyingKey
	SetupMs  int64
}

func setup2(depth int, label string) *compiled2 {
	fmt.Printf("  ⚙️  %s (depth=%d) Setup...\n", label, depth)
	ins, upd, del := makeTpls2(depth)
	cs := &compiled2{Depth: depth, Label: label}
	var err error
	t0 := time.Now()
	cs.InsertCS, err = frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, ins)
	if err != nil {
		log.Fatalf("compile Insert: %v", err)
	}
	cs.InsertPK, cs.InsertVK, err = groth16.Setup(cs.InsertCS)
	if err != nil {
		log.Fatalf("setup Insert: %v", err)
	}
	cs.ModifyCS, err = frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, upd)
	if err != nil {
		log.Fatalf("compile Modify: %v", err)
	}
	cs.ModifyPK, cs.ModifyVK, err = groth16.Setup(cs.ModifyCS)
	if err != nil {
		log.Fatalf("setup Modify: %v", err)
	}
	cs.DeleteCS, err = frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, del)
	if err != nil {
		log.Fatalf("compile Delete: %v", err)
	}
	cs.DeletePK, cs.DeleteVK, err = groth16.Setup(cs.DeleteCS)
	if err != nil {
		log.Fatalf("setup Delete: %v", err)
	}
	cs.SetupMs = time.Since(t0).Milliseconds()
	fmt.Printf("     ✅ %.3f s\n", float64(cs.SetupMs)/1000.0)
	return cs
}

// 核心：为每个 chunk 循环出证明

type batchResult2 struct {
	NumChunks    int
	BuildMs      int64   // 建树时间
	TotalProveMs int64   // Σ Prove（N 次）
	TotalTimeMs  int64   // BuildMs + TotalProveMs
	AvgProveMs   float64 // 平均每个 chunk 的 Prove
	AvgVerifyMs  float64
}

func proveAllChunks2(
	fileData []byte,
	depth int,
	op string,
	ccs constraint.ConstraintSystem,
	pk groth16.ProvingKey,
	vk groth16.VerifyingKey,
) (batchResult2, error) {

	numLeaves := 1 << depth
	res := batchResult2{NumChunks: numLeaves}

	// ① 建树（计入总时间）
	t0 := time.Now()
	tree := buildTree2(fileData, depth)
	res.BuildMs = time.Since(t0).Milliseconds()

	// ② 准备每个 chunk 的新叶（Modify/Insert 用）
	var newLeaves []*big.Int
	var sharedNewLeaf *big.Int
	switch op {
	case "modify":
		newLeaves = make([]*big.Int, numLeaves)
		for i := range newLeaves {
			newLeaves[i] = hashBytes2(randBytes2(chunkSize2))
		}
	case "insert":
		// 所有位置插入同一个新 chunk（测性能用）
		sharedNewLeaf = hashBytes2(randBytes2(chunkSize2))
	}

	// ③ 循环出证明
	fmt.Printf("    为 %d 个 chunk 生成证明...\n", numLeaves)
	var totalVerifyUs int64

	for i := 0; i < numLeaves; i++ {
		var assignment frontend.Circuit
		switch op {
		case "insert":
			assignment = makeInsertW2(tree, sharedNewLeaf, i)
		case "modify":
			assignment = makeModifyW2(tree, newLeaves[i], i)
		case "delete":
			assignment = makeDeleteW2(tree, i)
		}

		fullW, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
		if err != nil {
			return res, fmt.Errorf("chunk %d NewWitness: %w", i, err)
		}
		pubW, err := fullW.Public()
		if err != nil {
			return res, fmt.Errorf("chunk %d PublicWitness: %w", i, err)
		}

		tp := time.Now()
		proof, err := groth16.Prove(ccs, pk, fullW)
		proveMs := time.Since(tp).Milliseconds()
		if err != nil {
			return res, fmt.Errorf("chunk %d Prove: %w", i, err)
		}

		tv := time.Now()
		if err = groth16.Verify(proof, vk, pubW); err != nil {
			return res, fmt.Errorf("chunk %d Verify: %w", i, err)
		}
		totalVerifyUs += time.Since(tv).Microseconds()
		res.TotalProveMs += proveMs

		if (i+1)%500 == 0 || i == numLeaves-1 {
			fmt.Printf("      进度 %d/%d (%.0f%%)\n",
				i+1, numLeaves, float64(i+1)*100/float64(numLeaves))
		}
	}

	res.TotalTimeMs = res.BuildMs + res.TotalProveMs
	res.AvgProveMs = float64(res.TotalProveMs) / float64(numLeaves)
	res.AvgVerifyMs = float64(totalVerifyUs) / float64(numLeaves) / 1000.0
	return res, nil
}

// 统计（7次取中位数）

func median7b(results []batchResult2) batchResult2 {
	s := make([]batchResult2, len(results))
	copy(s, results)
	sort.Slice(s, func(i, j int) bool { return s[i].TotalTimeMs < s[j].TotalTimeMs })
	return s[len(s)/2]
}

// 结果结构

type runRec2 struct {
	BuildMs      int64   `json:"build_ms"`
	TotalProveMs int64   `json:"total_prove_ms"`
	TotalTimeMs  int64   `json:"total_time_ms"`
	AvgProveMs   float64 `json:"avg_prove_ms"`
}

type opResult2 struct {
	Operation     string    `json:"operation"`
	SizeMB        int       `json:"size_mb"`
	Label         string    `json:"label"`
	Depth         int       `json:"depth"`
	NumChunks     int       `json:"num_chunks"`
	SetupMs       int64     `json:"setup_ms"`
	Runs          []runRec2 `json:"runs"`
	TotalMedianMs float64   `json:"total_median_ms"`
	AvgProveMs    float64   `json:"avg_prove_ms"`
	AvgVerifyMs   float64   `json:"avg_verify_ms"`
}

// 主函数

func main() {
	const N = 7

	fmt.Printf("%-12s %-8s %-10s\n", "文件大小", "depth", "chunk 数")
	for _, tc := range testCases2 {
		d := sizeToDepth2(tc.SizeMB)
		fmt.Printf("%-12s %-8d %d\n", tc.Label, d, 1<<d)
	}
	fmt.Println()

	var allResults []opResult2

	for _, tc := range testCases2 {
		depth := sizeToDepth2(tc.SizeMB)
		numChunks := 1 << depth

		fmt.Printf("📁 %s (depth=%d, chunks=%d)\n", tc.Label, depth, numChunks)

		cs := setup2(depth, tc.Label)

		fmt.Printf("  📄 生成 %s 随机文件...\n", tc.Label)
		fileData := randBytes2(tc.SizeMB * 1024 * 1024)
		fmt.Printf("  ✅ 就绪\n\n")

		type opDef struct {
			name string
			op   string
			ccs  constraint.ConstraintSystem
			pk   groth16.ProvingKey
			vk   groth16.VerifyingKey
		}
		ops := []opDef{
			{"Insert（插入）", "insert", cs.InsertCS, cs.InsertPK, cs.InsertVK},
			{"Modify（修改）", "modify", cs.ModifyCS, cs.ModifyPK, cs.ModifyVK},
			{"Delete（删除）", "delete", cs.DeleteCS, cs.DeletePK, cs.DeleteVK},
		}

		for _, op := range ops {
			fmt.Printf("  ▶ %s\n", op.name)
			batchRuns := make([]batchResult2, N)
			ok := true

			for i := 0; i < N; i++ {
				fmt.Printf("    第 %d/%d 次\n", i+1, N)
				br, err := proveAllChunks2(fileData, depth, op.op, op.ccs, op.pk, op.vk)
				if err != nil {
					fmt.Printf("    ❌ %v\n", err)
					ok = false
					break
				}
				batchRuns[i] = br
				fmt.Printf("    ✅ 总时间=%.2f s  平均/chunk=%.3f ms\n",
					float64(br.TotalTimeMs)/1000.0, br.AvgProveMs)
			}
			if !ok {
				fmt.Printf("  ⚠️  跳过\n\n")
				continue
			}

			med := median7b(batchRuns)
			recs := make([]runRec2, N)
			for i, br := range batchRuns {
				recs[i] = runRec2{
					BuildMs: br.BuildMs, TotalProveMs: br.TotalProveMs,
					TotalTimeMs: br.TotalTimeMs, AvgProveMs: br.AvgProveMs,
				}
			}
			r := opResult2{
				Operation: op.name, SizeMB: tc.SizeMB, Label: tc.Label,
				Depth: depth, NumChunks: numChunks, SetupMs: cs.SetupMs,
				Runs: recs, TotalMedianMs: float64(med.TotalTimeMs),
				AvgProveMs: med.AvgProveMs, AvgVerifyMs: med.AvgVerifyMs,
			}
			allResults = append(allResults, r)
			fmt.Printf("  ✅ 总时间中位数=%.2f s  平均/chunk=%.3f ms\n\n",
				r.TotalMedianMs/1000.0, r.AvgProveMs)
		}
	}

	if b, err := json.MarshalIndent(allResults, "", "  "); err == nil {
		_ = os.WriteFile("result_full.json", b, 0644)
	}
	genMD2(allResults)
	fmt.Println("💾 result_full.json")
	fmt.Println("📝 report_full.md")
	fmt.Println("✅ 完成！")
}

// Markdown 报告

func genMD2(results []opResult2) {
	f, _ := os.Create("report_full.md")
	defer f.Close()
	ops := []string{"Insert（插入）", "Modify（修改）", "Delete（删除）"}
	w := func(s string, a ...interface{}) { fmt.Fprintf(f, s, a...) }

	w("## 文件大小与 Chunk 数\n\n| 文件大小 | depth | chunk 数 |\n|:---:|:---:|:---:|\n")
	for _, tc := range testCases2 {
		d := sizeToDepth2(tc.SizeMB)
		w("| %s | %d | %d |\n", tc.Label, d, 1<<d)
	}
	w("\n")

	hdr, sep := "| 操作 |", "|:---:|"
	for _, tc := range testCases2 {
		hdr += " " + tc.Label + " |"
		sep += ":---:|"
	}

	// 总时间
	w("## 总时间（BuildTree + N×Prove，中位数，秒）\n\n%s\n%s\n", hdr, sep)
	for _, op := range ops {
		row := "| **" + op + "** |"
		for _, tc := range testCases2 {
			if v := find2(results, op, tc.SizeMB); v != nil {
				row += fmt.Sprintf(" %.2f |", v.TotalMedianMs/1000.0)
			} else {
				row += " N/A |"
			}
		}
		w("%s\n", row)
	}
	w("\n")

	// 平均每 chunk
	w("## 平均每 Chunk 的 Prove 时间（中位数，ms）\n\n%s\n%s\n", hdr, sep)
	for _, op := range ops {
		row := "| **" + op + "** |"
		for _, tc := range testCases2 {
			if v := find2(results, op, tc.SizeMB); v != nil {
				row += fmt.Sprintf(" %.3f ms |", v.AvgProveMs)
			} else {
				row += " N/A |"
			}
		}
		w("%s\n", row)
	}
	w("\n")

	// 详细
	w("## 详细数据\n\n")
	for _, tc := range testCases2 {
		d := sizeToDepth2(tc.SizeMB)
		w("### %s（depth=%d，%d chunks）\n\n", tc.Label, d, 1<<d)
		w("| 操作 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | **中位数(s)** | **avg/chunk(ms)** |\n")
		w("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|\n")
		for _, op := range ops {
			v := find2(results, op, tc.SizeMB)
			if v == nil {
				continue
			}
			w("| %s |", op)
			for _, r := range v.Runs {
				w(" %.2f |", float64(r.TotalTimeMs)/1000.0)
			}
			w(" **%.2f** | **%.3f** |\n", v.TotalMedianMs/1000.0, v.AvgProveMs)
		}
		w("\n")
	}

	// CSV
	w("## CSV\n\n```csv\noperation,label,depth,num_chunks,total_median_s,avg_prove_ms,avg_verify_ms\n")
	for _, r := range results {
		w("%s,%s,%d,%d,%.4f,%.4f,%.4f\n",
			r.Operation, r.Label, r.Depth, r.NumChunks,
			r.TotalMedianMs/1000.0, r.AvgProveMs, r.AvgVerifyMs)
	}
	w("```\n")
}

func find2(results []opResult2, op string, sizeMB int) *opResult2 {
	for i := range results {
		if results[i].Operation == op && results[i].SizeMB == sizeMB {
			return &results[i]
		}
	}
	return nil
}
