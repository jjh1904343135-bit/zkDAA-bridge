package zkp

import (
	"fmt"
	"os"
	"time"

	"zk-htlc/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// ZKPHandler 封装了 ZKP 的所有核心组件
type ZKPHandler struct {
	R1CS         constraint.ConstraintSystem
	ProvingKey   groth16.ProvingKey
	VerifyingKey groth16.VerifyingKey
}

// NewZKPHandler 执行一次性的可信设置，生成证明和验证密钥
// 如果密钥文件已存在，则直接加载；否则重新生成
func NewZKPHandler() (*ZKPHandler, error) {
	fmt.Println("\n[ZKP] 正在执行一次性可信设置 (Setup)...")
	startTime := time.Now()

	// 🔑 尝试从文件加载已保存的密钥
	if keysExist() {
		fmt.Println("[ZKP] 🔍 检测到已保存的密钥文件，正在加载...")
		return loadKeysFromFile()
	}

	// 如果密钥不存在，执行完整的 Setup
	fmt.Println("[ZKP] ⚙️  未找到密钥文件，执行完整 Setup...")

	// 1. 实例化我们的电路
	var unlockCircuit circuit.UnlockCircuit

	// 2. 编译电路 -> R1CS
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &unlockCircuit)
	if err != nil {
		return nil, fmt.Errorf("电路编译失败: %w", err)
	}
	fmt.Printf("[ZKP] 电路编译完成。约束数量: %d\n", ccs.GetNbConstraints())

	// 3. 运行 Groth16 Setup -> Proving Key & Verifying Key
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, fmt.Errorf("Groth16 Setup 失败: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("[ZKP] ✅ 可信设置完成, 耗时: %v\n", duration)

	fmt.Println("[ZKP] ⚠️  警告: 密钥未保存到文件！")
	fmt.Println("[ZKP]     请先运行 'go run tools/gen_verifier.go' 生成并保存密钥")

	return &ZKPHandler{
		R1CS:         ccs,
		ProvingKey:   pk,
		VerifyingKey: vk,
	}, nil
}

// keysExist 检查密钥文件是否存在
func keysExist() bool {
	_, err1 := os.Stat("keys/proving.key")
	_, err2 := os.Stat("keys/verifying.key")
	_, err3 := os.Stat("keys/circuit.r1cs")
	return err1 == nil && err2 == nil && err3 == nil
}

// loadKeysFromFile 从文件加载密钥
func loadKeysFromFile() (*ZKPHandler, error) {
	startTime := time.Now()

	// 加载 R1CS
	r1csFile, err := os.Open("keys/circuit.r1cs")
	if err != nil {
		return nil, fmt.Errorf("打开 circuit.r1cs 失败: %w", err)
	}
	defer r1csFile.Close()

	ccs := groth16.NewCS(ecc.BN254)
	_, err = ccs.ReadFrom(r1csFile)
	if err != nil {
		return nil, fmt.Errorf("读取 circuit.r1cs 失败: %w", err)
	}
	fmt.Printf("[ZKP] ✅ R1CS 加载成功，约束数量: %d\n", ccs.GetNbConstraints())

	// 加载 Proving Key
	pkFile, err := os.Open("keys/proving.key")
	if err != nil {
		return nil, fmt.Errorf("打开 proving.key 失败: %w", err)
	}
	defer pkFile.Close()

	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(pkFile)
	if err != nil {
		return nil, fmt.Errorf("读取 proving.key 失败: %w", err)
	}
	fmt.Println("[ZKP] ✅ Proving Key 加载成功")

	// 加载 Verifying Key
	vkFile, err := os.Open("keys/verifying.key")
	if err != nil {
		return nil, fmt.Errorf("打开 verifying.key 失败: %w", err)
	}
	defer vkFile.Close()

	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vk.ReadFrom(vkFile)
	if err != nil {
		return nil, fmt.Errorf("读取 verifying.key 失败: %w", err)
	}
	fmt.Println("[ZKP] ✅ Verifying Key 加载成功")

	duration := time.Since(startTime)
	fmt.Printf("[ZKP] ✅ 密钥加载完成, 耗时: %v\n", duration)

	return &ZKPHandler{
		R1CS:         ccs,
		ProvingKey:   pk,
		VerifyingKey: vk,
	}, nil
}

// Prove 生成一个零知识证明
func (h *ZKPHandler) Prove(assignment frontend.Circuit) (groth16.Proof, error) {
	fmt.Println("[ZKP] 正在为解锁交易生成证明...")
	startTime := time.Now()

	// 创建 witness
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("创建 witness 失败: %w", err)
	}

	// 调用 Groth16 的 Prove 函数
	proof, err := groth16.Prove(h.R1CS, h.ProvingKey, witness)
	if err != nil {
		return nil, fmt.Errorf("生成证明失败: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("[ZKP] ✅ 证明生成成功, 耗时: %v\n", duration)
	return proof, nil
}

// ExportVerifier 导出 Verifier.sol 合约（公开方法，供 tools 使用）
func (h *ZKPHandler) ExportVerifier(outputPath string) error {
	fmt.Printf("[ZKP] 正在导出 Verifier.sol 到 %s...\n", outputPath)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	err = h.VerifyingKey.ExportSolidity(file)
	if err != nil {
		return fmt.Errorf("导出 Solidity 失败: %w", err)
	}

	fmt.Println("[ZKP] ✅ Verifier.sol 导出成功!")
	return nil
}

// GetConstraintCount 返回电路的约束数量
func (h *ZKPHandler) GetConstraintCount() int {
	return h.R1CS.GetNbConstraints()
}

// exportVerifier 是一个内部辅助函数（已弃用，使用 ExportVerifier 代替）
func exportVerifier(vk groth16.VerifyingKey) {
	fmt.Println("[ZKP] 正在导出 Verifier.sol 合约...")
	file, err := os.Create("Verifier.sol")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = vk.ExportSolidity(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("[ZKP] ✅ Verifier.sol 导出成功!")
}
