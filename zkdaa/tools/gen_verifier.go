package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	depthPtr := flag.Int("depth", 10, "merkle tree depth")
	flag.Parse()
	depth := *depthPtr

	// 确保 build 目录存在
	if _, err := os.Stat("build"); os.IsNotExist(err) {
		os.Mkdir("build", 0755)
	}

	fmt.Printf("[1/2] Generating UnlockVerifier (A)...\n")
	generateUnlockVerifier()

	fmt.Printf("\n[2/2] Generating AuditVerifier (B) with depth=%d...\n", depth)
	generateAuditVerifier(depth)
	
	fmt.Println("\n✅ All Setup Completed!")
}

func generateUnlockVerifier() {
	var c circuit.UnlockCircuit
	
	// 1. 编译
	startCompile := time.Now()
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &c)
	if err != nil { panic(err) }
	compileTime := time.Since(startCompile)
	fmt.Printf("   Compiling circuit... Done! (%v)\n", compileTime)
	fmt.Printf("   Constraints: %d\n", ccs.GetNbConstraints())

	// 2. Setup
	startSetup := time.Now()
	pk, vk, err := groth16.Setup(ccs)
	if err != nil { panic(err) }
	setupTime := time.Since(startSetup)
	fmt.Printf("   Running Groth16 Setup... Done! (%v)\n", setupTime)

	// 🔥 写入 Unlock Setup 时间
	writeTimeToFile("build/setup_time_unlock.txt", setupTime)

	saveKeys("unlock", pk, vk)
	saveVerifier("UnlockVerifier", vk)
}

func generateAuditVerifier(depth int) {
	c := circuit.AuditUnlockCircuit{
		ProofPath:    make([]frontend.Variable, depth),
		Helpers:      make([]frontend.Variable, depth),
		LeafCounts:   make([]frontend.Variable, depth),
		LeafNumBytes: make([]frontend.Variable, depth),
	}

	// 1. 编译
	startCompile := time.Now()
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &c)
	if err != nil { panic(err) }
	compileTime := time.Since(startCompile)
	fmt.Printf("   Compiling circuit... Done! (%v)\n", compileTime)
	fmt.Printf("   Constraints: %d\n", ccs.GetNbConstraints())

	// 2. Setup
	startSetup := time.Now()
	pk, vk, err := groth16.Setup(ccs)
	if err != nil { panic(err) }
	setupTime := time.Since(startSetup)
	fmt.Printf("   Running Groth16 Setup... Done! (%v)\n", setupTime)

	// 🔥 写入 Audit Setup 时间
	writeTimeToFile("build/setup_time_audit.txt", setupTime)

	saveKeys(fmt.Sprintf("audit_d%d", depth), pk, vk)
	saveVerifier(fmt.Sprintf("AuditVerifier_d%d", depth), vk)
}

// 通用写入函数
func writeTimeToFile(filename string, d time.Duration) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Warning: failed to save setup time: %v\n", err)
		return
	}
	defer f.Close()
	// 写入毫秒数
	f.WriteString(fmt.Sprintf("%d", d.Milliseconds()))
}

func saveKeys(name string, pk groth16.ProvingKey, vk groth16.VerifyingKey) {
	pkFile, _ := os.Create("build/" + name + ".pk")
	defer pkFile.Close()
	pk.WriteTo(pkFile)

	vkFile, _ := os.Create("build/" + name + ".vk")
	defer vkFile.Close()
	vk.WriteTo(vkFile)
}

func saveVerifier(contractName string, vk groth16.VerifyingKey) {
	var buf bytes.Buffer
	if err := vk.ExportSolidity(&buf); err != nil {
		panic(err)
	}
	content := buf.String()
	newContent := strings.Replace(content, "contract Verifier", "contract "+contractName, 1)

	filePath := "blockchain-contracts/contracts/" + contractName + ".sol"
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	f.WriteString(newContent)
	fmt.Printf("   Saved to %s (renamed to contract %s)\n", filePath, contractName)
}