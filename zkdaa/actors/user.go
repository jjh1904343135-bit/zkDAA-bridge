package actors

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"zk-htlc/data"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

// 错误定义
var (
	ErrInvalidDataID       = errors.New("invalid data ID")
	ErrDataPackageNotFound = errors.New("data package not found")
	ErrInvalidChunkIndex   = errors.New("invalid chunk index")
	ErrCIDFMismatch        = errors.New("CIDF verification failed")
)

// DSPAInfoPackage 用户发送给 DSPA 的信息包
type DSPAInfoPackage struct {
	Pre *big.Int `json:"pre"`
	Sn  *big.Int `json:"sn"`
	H   *big.Int `json:"h"`
}

// DSPBInfoPackage 用户发送给 DSPB 的信息包
type DSPBInfoPackage struct {
	Pre *big.Int `json:"pre"`
	Sn  *big.Int `json:"sn"`
	H   *big.Int `json:"h"`
}

// User 用户结构体，存储密码学参数和数据包
type User struct {
	// 密码学参数
	Pre_I  *big.Int
	Sn_I   *big.Int
	Pre_II *big.Int
	Sn_II  *big.Int
	Z_256  *big.Int
	H1     *big.Int
	H2     *big.Int

	// 数据审计相关
	mu           sync.RWMutex                 // 保护并发访问
	DataPackages map[string]*data.DataPackage // dataID -> 数据包
}

// BN254 曲线的标量域阶数
var fieldModulus = new(big.Int)

func init() {
	// BN254 曲线的 r (标量域阶数)
	fieldModulus.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
}

// NewUser 创建新用户并生成所有密码学参数
func NewUser() (*User, error) {
	printHeader("用户生成密码学参数 (User Setup)")

	// 1. 生成确定性秘密
	fmt.Println("\n[USER] 📝 生成确定性秘密 (Preimages)...")
	pre_I := generatePreimageFromSecret("user_master_secret_for_DSPA")
	pre_II := generatePreimageFromSecret("user_master_secret_for_DSPB")

	printTruncated("pre_I", pre_I)
	printTruncated("pre_II", pre_II)

	// 2. 生成随机数
	fmt.Println("\n[USER] 🎲 生成随机数...")
	sn_II, err := newRandomFieldElement()
	if err != nil {
		return nil, fmt.Errorf("failed to generate sn_II: %w", err)
	}

	z_256, err := newRandomFieldElement()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Z_256: %w", err)
	}

	// XOR操作后需要确保结果在有限域内
	sn_I_raw := new(big.Int).Xor(sn_II, z_256)
	sn_I := toFieldElement(sn_I_raw)

	printTruncated("sn_II", sn_II)
	printTruncated("Z_256", z_256)
	printTruncated("sn_I (sn_II ⊕ Z_256) mod r", sn_I)

	// 3. 验证 XOR 关系
	if err := verifyXORRelation(sn_I, sn_II, z_256); err != nil {
		return nil, fmt.Errorf("XOR verification failed: %w", err)
	}

	// 4. 计算哈希锁
	fmt.Println("\n[USER] 🔐 基于秘密计算哈希锁 (MiMC Hash)...")
	h1, err := calculateMiMCHash(pre_I, sn_I)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate h1: %w", err)
	}

	h2, err := calculateMiMCHash(pre_II, sn_II)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate h2: %w", err)
	}

	printHashLock("h1 (for DSPA)", h1)
	printHashLock("h2 (for DSPB)", h2)

	fmt.Println("\n✅ 用户参数生成完成!")

	return &User{
		Pre_I:        pre_I,
		Sn_I:         sn_I,
		Pre_II:       pre_II,
		Sn_II:        sn_II,
		Z_256:        z_256,
		H1:           h1,
		H2:           h2,
		DataPackages: make(map[string]*data.DataPackage),
	}, nil
}

// DistributeInfo 将信息分发给两个 DSP
func (u *User) DistributeInfo() (*DSPAInfoPackage, *DSPBInfoPackage) {
	printHeader("用户分发信息包 (Info Distribution)")

	dspaPackage := &DSPAInfoPackage{
		Pre: new(big.Int).Set(u.Pre_I), // 使用副本防止外部修改
		Sn:  new(big.Int).Set(u.Sn_I),
		H:   new(big.Int).Set(u.H1),
	}

	printPackageInfo("DSPA", dspaPackage.Pre, dspaPackage.Sn, dspaPackage.H)

	dspbPackage := &DSPBInfoPackage{
		Pre: new(big.Int).Set(u.Pre_II),
		Sn:  new(big.Int).Set(u.Sn_II),
		H:   new(big.Int).Set(u.H2),
	}

	printPackageInfo("DSPB", dspbPackage.Pre, dspbPackage.Sn, dspbPackage.H)

	return dspaPackage, dspbPackage
}

// SendDataPackage 向 DSP 发送数据包
func (u *User) SendDataPackage(dataID string, fileData []byte, targetIP string, chunkSize int) (*data.DataPackage, error) {
	if dataID == "" {
		return nil, ErrInvalidDataID
	}

	if len(fileData) == 0 {
		return nil, errors.New("file data is empty")
	}

	if chunkSize <= 0 {
		return nil, errors.New("chunk size must be positive")
	}

	pkg := data.NewDataPackage(dataID, fileData, targetIP, chunkSize)

	u.mu.Lock()
	u.DataPackages[dataID] = pkg
	u.mu.Unlock()

	fmt.Printf("\n[USER] 📦 向 DSP 发送数据包\n")
	fmt.Printf("       - DataID: %s\n", dataID)
	fmt.Printf("       - 文件大小: %d bytes\n", len(fileData))
	fmt.Printf("       - 目标IP: %s\n", targetIP)
	fmt.Printf("       - 分块大小: %d bytes\n", chunkSize)
	fmt.Printf("       - 分块数量: %d\n", (len(fileData)+chunkSize-1)/chunkSize)

	return pkg, nil
}

// InitiateAudit 发起审计挑战
func (u *User) InitiateAudit(dataID string, chunkIndex int) (*data.AuditChallenge, error) {
	if dataID == "" {
		return nil, ErrInvalidDataID
	}

	u.mu.RLock()
	pkg, exists := u.DataPackages[dataID]
	u.mu.RUnlock()

	if !exists {
		return nil, ErrDataPackageNotFound
	}

	// 验证分块索引的有效性
	maxChunks := (len(pkg.FileData) + pkg.ChunkSize - 1) / pkg.ChunkSize
	if chunkIndex < 0 || chunkIndex >= maxChunks {
		return nil, fmt.Errorf("%w: index %d, max %d", ErrInvalidChunkIndex, chunkIndex, maxChunks-1)
	}

	challenge := data.NewAuditChallenge(dataID, chunkIndex)

	fmt.Printf("\n[USER] 🔍 发起审计挑战\n")
	fmt.Printf("       - DataID: %s\n", dataID)
	fmt.Printf("       - 挑战索引: %d/%d\n", chunkIndex, maxChunks-1)
	fmt.Printf("       - 时间戳: %d\n", challenge.Timestamp)

	return challenge, nil
}

// VerifyAuditResponse 验证审计响应
func (u *User) VerifyAuditResponse(response *data.AuditResponse, expectedCIDF []byte) error {
	if response == nil {
		return errors.New("response is nil")
	}

	fmt.Printf("\n[USER] 🔐 验证审计响应\n")
	fmt.Printf("       - 块索引: %d\n", response.ChunkIndex)
	fmt.Printf("       - CIDF: 0x%s...\n", hex.EncodeToString(response.MerkleRoot[:min(8, len(response.MerkleRoot))]))

	// 验证 CIDF
	if len(response.MerkleRoot) != len(expectedCIDF) {
		return fmt.Errorf("%w: length mismatch (got %d, expected %d)",
			ErrCIDFMismatch, len(response.MerkleRoot), len(expectedCIDF))
	}

	// 使用常量时间比较防止时序攻击
	if !bytesEqual(response.MerkleRoot, expectedCIDF) {
		return ErrCIDFMismatch
	}

	fmt.Println("       ✅ 审计响应有效")
	return nil
}

// GetDataPackage 获取数据包（线程安全）
func (u *User) GetDataPackage(dataID string) (*data.DataPackage, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	pkg, exists := u.DataPackages[dataID]
	if !exists {
		return nil, ErrDataPackageNotFound
	}

	return pkg, nil
}

// --- 密码学辅助函数 ---

// calculateMiMCHash 使用 MiMC 哈希函数计算哈希值
func calculateMiMCHash(pre, sn *big.Int) (*big.Int, error) {
	if pre == nil || sn == nil {
		return nil, errors.New("nil input to hash function")
	}

	// 1. 转换为有限域元素
	preField := toFieldElement(pre)
	snField := toFieldElement(sn)

	// 2. 转换为 fr.Element
	var preFr, snFr fr.Element
	preFr.SetBigInt(preField)
	snFr.SetBigInt(snField)

	// 3. 初始化 MiMC 哈希
	mimcHasher := mimc.NewMiMC()

	// 4. 写入数据
	preBytes := preFr.Bytes()
	if _, err := mimcHasher.Write(preBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to write pre to hasher: %w", err)
	}

	snBytes := snFr.Bytes()
	if _, err := mimcHasher.Write(snBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to write sn to hasher: %w", err)
	}

	// 5. 获取哈希结果
	hashBytes := mimcHasher.Sum(nil)

	// 6. 转换回 big.Int
	var hashFr fr.Element
	hashFr.SetBytes(hashBytes)

	return hashFr.BigInt(new(big.Int)), nil
}

// toFieldElement 确保值在 BN254 标量域内
func toFieldElement(val *big.Int) *big.Int {
	if val == nil {
		return big.NewInt(0)
	}

	result := new(big.Int)
	result.Mod(val, fieldModulus)
	return result
}

// generatePreimageFromSecret 从字符串种子生成确定性前像
func generatePreimageFromSecret(secret string) *big.Int {
	hash := sha256.Sum256([]byte(secret))
	val := new(big.Int).SetBytes(hash[:])
	return toFieldElement(val)
}

// newRandomFieldElement 生成随机有限域元素
func newRandomFieldElement() (*big.Int, error) {
	// 生成 256 位随机数
	max := new(big.Int).Lsh(big.NewInt(1), 256)
	val, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}

	return toFieldElement(val), nil
}

// verifyXORRelation 验证 XOR 关系（考虑有限域模运算）
func verifyXORRelation(sn_I, sn_II, z_256 *big.Int) error {
	computed_raw := new(big.Int).Xor(sn_II, z_256)
	computed := toFieldElement(computed_raw)

	if computed.Cmp(sn_I) != 0 {
		return errors.New("XOR relation verification failed")
	}

	fmt.Println("       ✅ XOR 关系验证通过: sn_I = (sn_II ⊕ Z_256) mod r")
	return nil
}

// --- 工具函数 ---

// bytesEqual 常量时间字节数组比较
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// min 返回两个整数的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- 打印辅助函数 ---

// printHeader 打印格式化的标题
func printHeader(title string) {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Printf("║  %-54s  ║\n", title)
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
}

// printTruncated 打印截断的大整数
func printTruncated(name string, val *big.Int) {
	s := val.String()
	if len(s) > 20 {
		s = s[:20] + "..."
	}
	fmt.Printf("       ✓ %-8s %s\n", name+":", s)
}

// printHashLock 打印哈希锁
func printHashLock(name string, hash *big.Int) {
	bytes := hash.Bytes()
	truncated := bytes
	if len(bytes) > 8 {
		truncated = bytes[:8]
	}
	fmt.Printf("       ✓ %s: 0x%s...\n", name, hex.EncodeToString(truncated))
}

// printPackageInfo 打印信息包详情
func printPackageInfo(target string, pre, sn, h *big.Int) {
	fmt.Printf("\n[USER] 📦 -> [%s] 发送信息包\n", target)
	printTruncated("Pre", pre)
	printTruncated("Sn", sn)
	printHashLock("H", h)
}
