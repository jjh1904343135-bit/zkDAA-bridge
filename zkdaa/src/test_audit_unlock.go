package main

// import (
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"math/big"
// 	"zk-htlc/actors"
// 	_ "zk-htlc/circuit"
// 	"zk-htlc/merkle"
// 	"zk-htlc/zkp"

// 	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
// )

// // 测试新的审计-解锁联合流程
// func main() {
// 	fmt.Println("==============================================")
// 	fmt.Println("   ZK-HTLC 审计-解锁联合电路测试")
// 	fmt.Println("   核心创新: H2 = MiMC(CIDF, Sn_II)")
// 	fmt.Println("==============================================\n")

// 	// ========== 第一阶段：用户准备 ==========
// 	fmt.Println("【阶段 1】用户准备密码学参数...")

// 	// 1. 创建用户
// 	user, err := actors.NewUser()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 2. 用户分发信息给两个 DSP
// 	infoDSPA, infoDSPB := user.DistributeInfo()

// 	fmt.Printf("✅ H1 (for SCA) = %s...\n", infoDSPA.H.String()[:20])
// 	fmt.Printf("✅ Sn_I = %s...\n", infoDSPA.Sn.String()[:20])
// 	fmt.Printf("✅ Sn_II = %s...\n", infoDSPB.Sn.String()[:20])

// 	// ========== 第二阶段：准备数据和 CIDF ==========
// 	fmt.Println("\n【阶段 2】数据迁移准备...")

// 	// 3. 创建测试数据
// 	testDataSize := 8 * 1024 // 8KB
// 	chunkSize := 1024        // 1KB
// 	testData := make([]byte, testDataSize)
// 	rand.Read(testData)

// 	dataID := hex.EncodeToString(testData[:8])
// 	fmt.Printf("✅ 生成测试数据: %d bytes (DataID: %s)\n", testDataSize, dataID)

// 	// 4. 模拟数据在 DSPA，然后迁移到 DSPB
// 	// DSPB 接收数据并构建 Merkle 树
// 	chunks := merkle.ChunkData(testData, chunkSize)
// 	mt := merkle.NewMerkleTree(chunks)
// 	CIDF := mt.GetRoot()

// 	fmt.Printf("✅ DSPB 构建 Merkle 树\n")
// 	fmt.Printf("   - 数据块数量: %d\n", len(chunks))
// 	fmt.Printf("   - CIDF: 0x%x...\n", CIDF[:8])

// 	// ========== 第三阶段：计算 H2 ==========
// 	fmt.Println("\n【阶段 3】计算 H2 = MiMC(CIDF, Sn_II)...")

// 	// 5. 关键创新：用 CIDF 计算 H2
// 	hFunc := mimc.NewMiMC()
// 	hFunc.Write(CIDF)
// 	snIIBytes := infoDSPB.Sn.Bytes()
// 	hFunc.Write(snIIBytes)
// 	H2_bytes := hFunc.Sum(nil)
// 	H2 := new(big.Int).SetBytes(H2_bytes)

// 	fmt.Printf("✅ H2 = MiMC(CIDF, Sn_II)\n")
// 	fmt.Printf("   - H2: %s...\n", H2.String()[:20])

// 	// ========== 第四阶段：初始化 Setup ==========
// 	fmt.Println("\n【阶段 4】初始化审计-解锁联合电路 Setup...")

// 	// 6. 计算 Merkle 树深度
// 	proofDepth := calculateMerkleDepth(len(chunks))
// 	fmt.Printf("   - Merkle 树深度: %d\n", proofDepth)

// 	// 7. 执行 Setup
// 	auditUnlockHandler, err := zkp.NewAuditUnlockHandler(proofDepth)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("✅ Setup 完成，约束数量: %d\n", auditUnlockHandler.GetConstraintCount())

// 	// ========== 第五阶段：创建 DSPB ==========
// 	fmt.Println("\n【阶段 5】创建 DSPB...")

// 	// 8. 创建 DSPB（使用 H2）
// 	dspB := actors.NewDSP("DSPB", nil, infoDSPB.Sn, H2, nil)
// 	dspB.SetAuditUnlockHandler(auditUnlockHandler)

// 	// 9. DSPB 存储数据和 Merkle 树
// 	dspB.StoredData = make(map[string][]byte)
// 	dspB.MerkleTrees = make(map[string]*merkle.MerkleTree)
// 	dspB.StoredData[dataID] = testData
// 	dspB.MerkleTrees[dataID] = mt

// 	fmt.Printf("✅ DSPB 已准备完毕\n")

// 	// ========== 第六阶段：SCB 发起审计挑战 ==========
// 	fmt.Println("\n【阶段 6】SCB 发起审计挑战...")

// 	// 10. 随机选择一个数据块进行审计
// 	challengeIndex := 3
// 	fmt.Printf("   - 挑战索引: %d\n", challengeIndex)

// 	// ========== 第七阶段：DSPB 生成联合证明 ==========
// 	fmt.Println("\n【阶段 7】DSPB 生成审计-解锁联合证明...")

// 	// 11. DSPB 生成联合证明 witness
// 	assignment, chunkHash, cidfResult, err := dspB.GenerateAuditUnlockProof(
// 		dataID,
// 		challengeIndex,
// 		infoDSPB.Sn, // Sn_II
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 验证 CIDF 一致性
// 	if hex.EncodeToString(CIDF) != hex.EncodeToString(cidfResult) {
// 		panic("CIDF 不一致！")
// 	}

// 	fmt.Printf("✅ 联合证明 witness 构造完成\n")

// 	// 12. 生成证明
// 	fmt.Println("\n【阶段 8】生成零知识证明...")
// 	proof, err := auditUnlockHandler.GenerateProof(assignment)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// ========== 第八阶段：SCB 验证证明 ==========
// 	fmt.Println("\n【阶段 9】SCB 验证证明...")

// 	// 13. SCB 验证证明
// 	err = auditUnlockHandler.Verify(proof, challengeIndex, chunkHash, H2.Bytes())
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("\n✅ 验证通过！SCB 解锁成功！")
// 	fmt.Printf("   - CIDF 已验证: 0x%x...\n", CIDF[:8])
// 	fmt.Printf("   - SCB 公开 Sn_II: %s...\n", infoDSPB.Sn.String()[:20])

// 	// ========== 第九阶段：SCA 解锁 ==========
// 	fmt.Println("\n【阶段 10】SCA 监听到 Sn_II，准备解锁...")

// 	// 14. 验证 XOR 关系
// 	snI := infoDSPA.Sn
// 	snII := infoDSPB.Sn
// 	z256 := user.Z_256

// 	xorResult := new(big.Int).Xor(snI, snII)
// 	if xorResult.Cmp(z256) == 0 {
// 		fmt.Println("✅ XOR 验证通过: Sn_I ⊕ Sn_II = Z_256")
// 	} else {
// 		fmt.Println("❌ XOR 验证失败")
// 		return
// 	}

// 	// 15. SCA 可以使用原有的 UnlockCircuit 解锁
// 	fmt.Println("✅ SCA 可以进行解锁（使用原有 UnlockCircuit）")
// 	fmt.Printf("   - H1 = MiMC(Pre_I, Sn_I)\n")

// 	// ========== 总结 ==========
// 	fmt.Println("\n==============================================")
// 	fmt.Println("          🎉 测试成功！")
// 	fmt.Println("==============================================")
// 	fmt.Println("\n流程总结:")
// 	fmt.Println("1. ✅ 用户生成 H1 = MiMC(Pre_I, Sn_I) 发给 SCA")
// 	fmt.Println("2. ✅ 用户生成 H2 = MiMC(CIDF, Sn_II) 发给 SCB")
// 	fmt.Println("3. ✅ 数据迁移到 DSPB，构建 Merkle 树得到 CIDF")
// 	fmt.Println("4. ✅ SCB 发起审计挑战")
// 	fmt.Println("5. ✅ DSPB 生成审计-解锁联合证明")
// 	fmt.Println("6. ✅ SCB 验证通过，公开 Sn_II")
// 	fmt.Println("7. ✅ SCA 验证 XOR 关系，准备解锁")
// 	fmt.Println("\n核心创新:")
// 	fmt.Println("  H2 = MiMC(CIDF, Sn_II) ← 用 CIDF 替换 Pre_II")
// 	fmt.Println("  将审计和解锁完美连接！")
// 	fmt.Println("==============================================\n")
// }

// // calculateMerkleDepth 计算 Merkle 树深度
// func calculateMerkleDepth(leafCount int) int {
// 	depth := 0
// 	n := leafCount
// 	for n > 1 {
// 		n = (n + 1) / 2
// 		depth++
// 	}
// 	return depth
// }
