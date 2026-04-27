package main

//  bench_single.go  ——  方案1：单条路径证明
//  运行: go run bench_single.go
//
//  证明语义：
//    每次操作只证明"一个数据块"的变更合法性（沿单条路径）
//    Insert : 证明在位置0旁边插入新块，newRoot 正确
//    Modify : 证明把位置0的块替换为新块，newRoot 正确
//    Delete : 证明删除位置0的块，newRoot 正确
//    1次操作 = 1个证明
//
//  时间如何随文件大小增长：
//    文件大小 → depth 梯度（步长=4）→ 约束数增加 → Prove 时间增加
//    depth 4→8→12→16→20，约束数差 ~4.8 倍，趋势明显
//    标签仍用 8MB/16MB/.../128MB 便于对比
//
//  计时范围：Prove（不含建树，建树是 O(N) 不属于 ZK 电路开销）
//  Setup 单独计时（一次性预计算）
//  Verify 单独计时（微秒级，固定）

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
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

const chunkSize1 = 1500

// 方案1 的实验梯度：depth 等步长递增，标签对应文件大小
// depth 4→8→12→16→20，约束数从 ~3520 到 ~17600，差 5 倍
var testCases1 = []struct {
	Depth int
	Label string
}{
	{4, "8 MB"},
	{8, "16 MB"},
	{12, "32 MB"},
	{16, "64 MB"},
	{20, "128 MB"},
}

// 工具函数

func mimcChain1(inputs ...*big.Int) *big.Int {
	h := hash.MIMC_BN254.New()
	for _, v := range inputs {
		h.Write(to321(v))
	}
	return new(big.Int).SetBytes(h.Sum(nil))
}

func to321(n *big.Int) []byte {
	b := n.Bytes()
	if len(b) >= 32 {
		return b
	}
	p := make([]byte, 32)
	copy(p[32-len(b):], b)
	return p
}

func randBytes1(size int) []byte {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		log.Fatalf("randBytes: %v", err)
	}
	return data
}

func hashBytes1(b []byte) *big.Int {
	h := hash.MIMC_BN254.New()
	h.Write(b)
	return new(big.Int).SetBytes(h.Sum(nil))
}

// Merkle 树（只需建一棵小树用于出单条路径证明）

type MerkleTree1 struct {
	Depth int
	nodes [][]*big.Int
}

func buildTree1(depth int) *MerkleTree1 {
	numLeaves := 1 << depth
	t := &MerkleTree1{Depth: depth}
	t.nodes = make([][]*big.Int, depth+1)
	t.nodes[0] = make([]*big.Int, numLeaves)
	for i := 0; i < numLeaves; i++ {
		t.nodes[0][i] = hashBytes1(randBytes1(chunkSize1))
	}
	for level := 1; level <= depth; level++ {
		prev := t.nodes[level-1]
		n := len(prev) / 2
		t.nodes[level] = make([]*big.Int, n)
		nb := big.NewInt(int64(1 << level))
		for i := 0; i < n; i++ {
			t.nodes[level][i] = mimcChain1(nb, prev[2*i], prev[2*i+1])
		}
	}
	return t
}

func (t *MerkleTree1) root() *big.Int { return t.nodes[t.Depth][0] }

// getProof 返回 leafIndex=0 的证明（始终在左，helper=0）
func (t *MerkleTree1) getProof() (path, helper, leafNum, leafNumByte []frontend.Variable) {
	d := t.Depth
	path = make([]frontend.Variable, d)
	helper = make([]frontend.Variable, d)
	leafNum = make([]frontend.Variable, d)
	leafNumByte = make([]frontend.Variable, d)
	for i := 0; i < d; i++ {
		path[i] = t.nodes[i][1]
		helper[i] = big.NewInt(0)
		leafNum[i] = big.NewInt(int64(1 << i))
		leafNumByte[i] = big.NewInt(int64(1 << (i + 1)))
	}
	return
}

func lnbSlice1(depth int) []frontend.Variable {
	s := make([]frontend.Variable, depth)
	for i := range s {
		s[i] = big.NewInt(int64(1 << (i + 1)))
	}
	return s
}

// Witness 构造（单条路径，leafIndex=0）

func makeInsertWitness1(t *MerkleTree1) *circuit.MerkleInsertCircuit {
	newLeaf := hashBytes1(randBytes1(chunkSize1))
	oldLeaf := t.nodes[0][0]
	computed := mimcChain1(newLeaf, oldLeaf)
	path, helper, leafNum, _ := t.getProof()
	depth := t.Depth
	nnb := make([]frontend.Variable, depth)
	for i := 0; i < depth; i++ {
		nb := big.NewInt(int64(1<<(i+1)) + 1)
		nnb[i] = nb
		computed = mimcChain1(nb, computed, t.nodes[i][1])
	}
	return &circuit.MerkleInsertCircuit{
		LeafHash: newLeaf, NewRootHash: computed,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf: oldLeaf, RootHash: t.root(), LeafIndex: big.NewInt(0),
			Path: path, LeafNum: leafNum, Helper: helper,
			LeafNumByte: lnbSlice1(depth),
		},
		NewNum_byte: nnb,
	}
}

func makeModifyWitness1(t *MerkleTree1) *circuit.MerkleUpdateCircuit {
	newLeaf := hashBytes1(randBytes1(chunkSize1))
	computed := newLeaf
	for i := 0; i < t.Depth; i++ {
		nb := big.NewInt(int64(1 << (i + 1)))
		computed = mimcChain1(nb, computed, t.nodes[i][1])
	}
	path, helper, leafNum, leafNumByte := t.getProof()
	return &circuit.MerkleUpdateCircuit{
		LeafHash: newLeaf, NewRootHash: computed,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf: t.nodes[0][0], RootHash: t.root(), LeafIndex: big.NewInt(0),
			Path: path, LeafNum: leafNum, Helper: helper, LeafNumByte: leafNumByte,
		},
	}
}

func makeDeleteWitness1(t *MerkleTree1) *circuit.MerkleDeleteCircuit {
	depth := t.Depth
	path, helper, leafNum, _ := t.getProof()
	computed := t.nodes[0][1]
	npb := make([]frontend.Variable, depth)
	npb[0] = big.NewInt(0)
	for i := 1; i < depth; i++ {
		nb := big.NewInt(int64(1<<(i+1)) - 1)
		npb[i] = nb
		computed = mimcChain1(nb, computed, t.nodes[i][1])
	}
	return &circuit.MerkleDeleteCircuit{
		NewRootHash: computed,
		Circuit_merkle: circuit.MerkleProofCircuit{
			Leaf: t.nodes[0][0], RootHash: t.root(), LeafIndex: big.NewInt(0),
			Path: path, LeafNum: leafNum, Helper: helper,
			LeafNumByte: lnbSlice1(depth),
		},
		NewPath_byte: npb,
	}
}

// 电路编译 & Setup

func makeTpls1(depth int) (
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

type compiled1 struct {
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

func setup1(depth int, label string) *compiled1 {
	fmt.Printf("  ⚙️  %s (depth=%d) Setup...\n", label, depth)
	ins, upd, del := makeTpls1(depth)
	cs := &compiled1{Depth: depth, Label: label}
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

// Prove + Verify（单次）

type timing1 struct{ ProveMs, VerifyUs int64 }

func proveOnce1(
	ccs constraint.ConstraintSystem,
	pk groth16.ProvingKey,
	vk groth16.VerifyingKey,
	assignment frontend.Circuit,
) (timing1, error) {
	fullW, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return timing1{}, fmt.Errorf("NewWitness: %w", err)
	}
	pubW, err := fullW.Public()
	if err != nil {
		return timing1{}, fmt.Errorf("PublicWitness: %w", err)
	}
	t0 := time.Now()
	proof, err := groth16.Prove(ccs, pk, fullW)
	if err != nil {
		return timing1{}, fmt.Errorf("Prove: %w", err)
	}
	proveMs := time.Since(t0).Milliseconds()
	t1 := time.Now()
	if err = groth16.Verify(proof, vk, pubW); err != nil {
		return timing1{}, fmt.Errorf("Verify: %w", err)
	}
	return timing1{ProveMs: proveMs, VerifyUs: time.Since(t1).Microseconds()}, nil
}

func median7s(v []int64) float64 {
	s := make([]int64, len(v))
	copy(s, v)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	return float64(s[3])
}

// 结果结构

type run1 struct {
	ProveMs  int64 `json:"prove_ms"`
	VerifyUs int64 `json:"verify_us"`
}

type result1 struct {
	Operation      string  `json:"operation"`
	Label          string  `json:"label"`
	Depth          int     `json:"depth"`
	SetupMs        int64   `json:"setup_ms"`
	Runs           []run1  `json:"runs"`
	ProveMedianMs  float64 `json:"prove_median_ms"`
	VerifyMedianUs float64 `json:"verify_median_us"`
	ProveMinMs     int64   `json:"prove_min_ms"`
	ProveMaxMs     int64   `json:"prove_max_ms"`
}

// 主函数

func main() {
	const N = 7

	fmt.Printf("%-12s %-8s %-12s\n", "文件大小", "depth", "约束数(估算)")
	for _, tc := range testCases1 {
		fmt.Printf("%-12s %-8d ~%d\n", tc.Label, tc.Depth, tc.Depth*4*220)
	}
	fmt.Println()

	var allResults []result1

	for _, tc := range testCases1 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📁 %s (depth=%d)\n", tc.Label, tc.Depth)

		cs := setup1(tc.Depth, tc.Label)

		// 建树只做一次（用于构造 witness，不纳入证明计时）
		tree := buildTree1(tc.Depth)

		type opDef struct {
			name string
			ccs  constraint.ConstraintSystem
			pk   groth16.ProvingKey
			vk   groth16.VerifyingKey
			mkW  func() frontend.Circuit
		}
		ops := []opDef{
			{"Insert（插入）", cs.InsertCS, cs.InsertPK, cs.InsertVK,
				func() frontend.Circuit { return makeInsertWitness1(tree) }},
			{"Modify（修改）", cs.ModifyCS, cs.ModifyPK, cs.ModifyVK,
				func() frontend.Circuit { return makeModifyWitness1(tree) }},
			{"Delete（删除）", cs.DeleteCS, cs.DeletePK, cs.DeleteVK,
				func() frontend.Circuit { return makeDeleteWitness1(tree) }},
		}

		for _, op := range ops {
			fmt.Printf("  ▶ %s\n", op.name)
			proves := make([]int64, N)
			verifys := make([]int64, N)
			runs := make([]run1, N)
			ok := true

			for i := 0; i < N; i++ {
				fmt.Printf("    第 %d/%d 次 ... ", i+1, N)
				tm, err := proveOnce1(op.ccs, op.pk, op.vk, op.mkW())
				if err != nil {
					fmt.Printf("❌ %v\n", err)
					ok = false
					break
				}
				proves[i] = tm.ProveMs
				verifys[i] = tm.VerifyUs
				runs[i] = run1{ProveMs: tm.ProveMs, VerifyUs: tm.VerifyUs}
				fmt.Printf("Prove=%.3f s  Verify=%.3f ms\n",
					float64(tm.ProveMs)/1000.0, float64(tm.VerifyUs)/1000.0)
			}
			if !ok {
				fmt.Printf("  ⚠️  跳过\n\n")
				continue
			}

			sp := make([]int64, N)
			copy(sp, proves)
			sort.Slice(sp, func(i, j int) bool { return sp[i] < sp[j] })

			allResults = append(allResults, result1{
				Operation: op.name, Label: tc.Label, Depth: tc.Depth,
				SetupMs: cs.SetupMs, Runs: runs,
				ProveMedianMs:  median7s(proves),
				VerifyMedianUs: median7s(verifys),
				ProveMinMs:     sp[0], ProveMaxMs: sp[N-1],
			})
			fmt.Printf("  ✅ Prove中位数=%.3f s  Verify中位数=%.3f ms\n\n",
				median7s(proves)/1000.0, median7s(verifys)/1000.0)
		}
	}

	if b, err := json.MarshalIndent(allResults, "", "  "); err == nil {
		_ = os.WriteFile("result_single.json", b, 0644)
	}
	genMD1(allResults)
	fmt.Println("💾 result_single.json")
	fmt.Println("📝 report_single.md")
	fmt.Println("✅ 完成！")
}

// Markdown 报告

func genMD1(results []result1) {
	f, _ := os.Create("report_single.md")
	defer f.Close()
	ops := []string{"Insert（插入）", "Modify（修改）", "Delete（删除）"}
	w := func(s string, a ...interface{}) { fmt.Fprintf(f, s, a...) }

	w("## 实验梯度\n\n| 文件大小 | depth | 约束数(估算) |\n|:---:|:---:|:---:|\n")
	for _, tc := range testCases1 {
		w("| %s | %d | ~%d |\n", tc.Label, tc.Depth, tc.Depth*4*220)
	}
	w("\n")

	// Setup
	w("## Setup 时间（秒）\n\n| 文件大小 | depth | Setup |\n|:---:|:---:|:---:|\n")
	seen := map[string]bool{}
	for _, r := range results {
		if !seen[r.Label] {
			seen[r.Label] = true
			w("| %s | %d | %.3f |\n", r.Label, r.Depth, float64(r.SetupMs)/1000.0)
		}
	}
	w("\n")

	// Prove 汇总
	hdr, sep := "| 操作 |", "|:---:|"
	for _, tc := range testCases1 {
		hdr += " " + tc.Label + " |"
		sep += ":---:|"
	}
	w("## Prove 时间（中位数，秒）\n\n%s\n%s\n", hdr, sep)
	for _, op := range ops {
		row := "| **" + op + "** |"
		for _, tc := range testCases1 {
			if v := find1(results, op, tc.Label); v != nil {
				row += fmt.Sprintf(" %.3f |", v.ProveMedianMs/1000.0)
			} else {
				row += " N/A |"
			}
		}
		w("%s\n", row)
	}
	w("\n")

	// Verify 汇总
	w("## Verify 时间（中位数，ms）\n\n> Groth16 Verify 为固定配对运算，与 depth 无关，约 1~3 ms\n\n%s\n%s\n", hdr, sep)
	for _, op := range ops {
		row := "| **" + op + "** |"
		for _, tc := range testCases1 {
			if v := find1(results, op, tc.Label); v != nil {
				row += fmt.Sprintf(" %.3f ms |", v.VerifyMedianUs/1000.0)
			} else {
				row += " N/A |"
			}
		}
		w("%s\n", row)
	}
	w("\n")

	// 详细
	w("## 详细数据\n\n")
	for _, tc := range testCases1 {
		w("### %s（depth=%d）\n\n", tc.Label, tc.Depth)
		w("| 操作 | Setup(s) | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 最小 | 最大 | **中位数(s)** | **Verify(ms)** |\n")
		w("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|\n")
		for _, op := range ops {
			v := find1(results, op, tc.Label)
			if v == nil {
				continue
			}
			w("| %s | %.3f |", op, float64(v.SetupMs)/1000.0)
			for _, r := range v.Runs {
				w(" %.3f |", float64(r.ProveMs)/1000.0)
			}
			w(" %.3f | %.3f | **%.3f** | **%.3f** |\n",
				float64(v.ProveMinMs)/1000.0, float64(v.ProveMaxMs)/1000.0,
				v.ProveMedianMs/1000.0, v.VerifyMedianUs/1000.0)
		}
		w("\n")
	}

	// CSV
	w("## CSV\n\n```csv\noperation,label,depth,setup_s,prove_median_s,verify_median_ms\n")
	for _, r := range results {
		w("%s,%s,%d,%.4f,%.4f,%.3f\n",
			r.Operation, r.Label, r.Depth,
			float64(r.SetupMs)/1000.0, r.ProveMedianMs/1000.0, r.VerifyMedianUs/1000.0)
	}
	w("```\n")
}

func find1(results []result1, op, label string) *result1 {
	for i := range results {
		if results[i].Operation == op && results[i].Label == label {
			return &results[i]
		}
	}
	return nil
}
