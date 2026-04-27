package actors

import (
	"fmt"
	"math/big"
	"zk-htlc/audit"
	"zk-htlc/circuit"
	"zk-htlc/data"
	"zk-htlc/merkle"
	"zk-htlc/zkp"

	"github.com/consensys/gnark/frontend"
)

// DSP 结构体模拟一个数据服务提供商节点
type DSP struct {
	Name       string
	Pre        *big.Int
	Sn         *big.Int
	H          *big.Int
	zkpHandler *zkp.ZKPHandler

	// 数据审计相关
	StoredData          map[string][]byte             // dataID -> 原始文件数据
	MerkleTrees         map[string]*merkle.MerkleTree // dataID -> Merkle 树
	AuditHandler        *audit.AuditHandler           // 审计证明处理器
	AuditUnlockHandler  *zkp.AuditUnlockHandler       // 审计-解锁联合处理器（新增）
}

// NewDSP 是 DSP 节点的构造函数
func NewDSP(name string, pre, sn, h *big.Int, handler *zkp.ZKPHandler) *DSP {
	fmt.Printf("[ACTOR] 创建 DSP 节点: %s\n", name)
	return &DSP{
		Name:       name,
		Pre:        pre,
		Sn:         sn,
		H:          h,
		zkpHandler: handler,
	}
}

// GenerateUnlockProof 是 DSP 最核心的行为
func (d *DSP) GenerateUnlockProof() (frontend.Circuit, error) {
	fmt.Printf("[%s] 准备为哈希锁 H=%s 生成证明...\n", d.Name, d.H.String()[:10]+"...")

	// 1. 构造 witness 赋值
	assignment := &circuit.UnlockCircuit{
		Pre: d.Pre,
		Sn:  d.Sn,
		H:   d.H,
	}

	return assignment, nil
}

// ReceiveDataPackage DSP 接收数据包并构建 Merkle 树
func (d *DSP) ReceiveDataPackage(pkg *data.DataPackage) (cidf []byte, err error) {
	fmt.Printf("\n[%s] 📥 接收数据包\n", d.Name)
	fmt.Printf("       - DataID: %s\n", pkg.DataID)
	fmt.Printf("       - 文件大小: %d bytes\n", len(pkg.FileData))
	fmt.Printf("       - 分块大小: %d bytes\n", pkg.ChunkSize)

	// 1. 存储原始数据
	if d.StoredData == nil {
		d.StoredData = make(map[string][]byte)
	}
	d.StoredData[pkg.DataID] = pkg.FileData

	// 2. 分块
	chunks := merkle.ChunkData(pkg.FileData, pkg.ChunkSize)
	fmt.Printf("       - 数据块数量: %d\n", len(chunks))

	// 3. 构建 Merkle 树
	mt := merkle.NewMerkleTree(chunks)
	if d.MerkleTrees == nil {
		d.MerkleTrees = make(map[string]*merkle.MerkleTree)
	}
	d.MerkleTrees[pkg.DataID] = mt

	// 4. 获取 CIDF (Merkle 根)
	cidf = mt.GetRoot()
	fmt.Printf("       - Merkle 树已构建\n")
	fmt.Printf("       - CIDF: 0x%x...\n", cidf[:8])

	return cidf, nil
}

// RespondToAudit DSP 响应审计挑战
func (d *DSP) RespondToAudit(challenge *data.AuditChallenge) (*data.AuditResponse, []byte, error) {
	fmt.Printf("\n[%s] 🔍 收到审计挑战\n", d.Name)
	fmt.Printf("       - DataID: %s\n", challenge.DataID)
	fmt.Printf("       - 挑战索引: %d\n", challenge.ChunkIndex)

	// 1. 获取 Merkle 树
	mt, exists := d.MerkleTrees[challenge.DataID]
	if !exists {
		return nil, nil, fmt.Errorf("数据 ID %s 不存在", challenge.DataID)
	}

	// 2. 检查索引有效性
	if challenge.ChunkIndex >= len(mt.Leaves) {
		return nil, nil, fmt.Errorf("索引 %d 超出范围 (总块数: %d)", challenge.ChunkIndex, len(mt.Leaves))
	}

	// 3. 获取数据块和 Merkle 证明
	proofPath, _ := mt.GetProof(challenge.ChunkIndex)
	chunkData := mt.Leaves[challenge.ChunkIndex]
	merkleRoot := mt.GetRoot()

	// 4. 计算数据块哈希（用于公开输入）
	chunkHash := circuit.HashChunk(chunkData)

	fmt.Printf("       - 数据块大小: %d bytes\n", len(chunkData))
	fmt.Printf("       - Merkle 证明路径长度: %d\n", len(proofPath))
	fmt.Printf("       - 块哈希: 0x%x...\n", chunkHash[:8])

	// 5. 构造响应
	response := &data.AuditResponse{
		ChunkIndex:  challenge.ChunkIndex,
		ChunkData:   chunkData,
		MerkleProof: proofPath,
		MerkleRoot:  merkleRoot,
	}

	fmt.Printf("       ✅ 审计响应已构造\n")
	return response, chunkHash, nil
}

// SetAuditHandler 设置审计处理器
func (d *DSP) SetAuditHandler(handler *audit.AuditHandler) {
	d.AuditHandler = handler
	fmt.Printf("[%s] 审计处理器已设置\n", d.Name)
}

// SetAuditUnlockHandler 设置审计-解锁联合处理器
func (d *DSP) SetAuditUnlockHandler(handler *zkp.AuditUnlockHandler) {
	d.AuditUnlockHandler = handler
	fmt.Printf("[%s] 审计-解锁联合处理器已设置\n", d.Name)
}

// GenerateAuditUnlockProof DSPB 生成审计-解锁联合证明
// 关键创新：用 CIDF 替换 Pre_II，实现 H2 = MiMC(CIDF, Sn_II)
func (d *DSP) GenerateAuditUnlockProof(
	dataID string,
	chunkIndex int,
	sn *big.Int, // Sn_II
) (*circuit.AuditUnlockCircuit, []byte, []byte, error) {
	fmt.Printf("\n[%s] 🔐 生成审计-解锁联合证明\n", d.Name)
	fmt.Printf("       - DataID: %s\n", dataID)
	fmt.Printf("       - 挑战索引: %d\n", chunkIndex)

	// 1. 获取 Merkle 树
	mt, exists := d.MerkleTrees[dataID]
	if !exists {
		return nil, nil, nil, fmt.Errorf("数据 ID %s 不存在", dataID)
	}

	// 2. 检查索引有效性
	if chunkIndex >= len(mt.Leaves) {
		return nil, nil, nil, fmt.Errorf("索引 %d 超出范围 (总块数: %d)", chunkIndex, len(mt.Leaves))
	}

	// 3. 获取 Merkle 证明和数据
	proofPath, helpers := mt.GetProof(chunkIndex)
	chunkData := mt.Leaves[chunkIndex]
	merkleRoot := mt.GetRoot() // 这就是 CIDF

	// 4. 计算数据块哈希
	chunkHash := circuit.HashChunkForAuditUnlock(chunkData)

	fmt.Printf("       - CIDF: 0x%x...\n", merkleRoot[:8])
	fmt.Printf("       - 数据块大小: %d bytes\n", len(chunkData))
	fmt.Printf("       - Merkle 路径长度: %d\n", len(proofPath))
	fmt.Printf("       - 块哈希: 0x%x...\n", chunkHash[:8])

	// 5. 构造 witness 赋值
	assignment := &circuit.AuditUnlockCircuit{
		ChunkData:    chunkData,
		ProofPath:    make([]frontend.Variable, len(proofPath)),
		Helpers:      make([]frontend.Variable, len(helpers)),
		LeafCounts:   make([]frontend.Variable, len(mt.PathCounts)),
		LeafNumBytes: make([]frontend.Variable, len(mt.PathBytes)),
		Sn:           sn, // Sn_II
		ChunkIndex:   chunkIndex,
		ChunkHash:    chunkHash,
		H:            d.H, // H2 = MiMC(CIDF, Sn_II)
	}

	// 6. 填充 Merkle 路径相关数据
	for i := 0; i < len(proofPath); i++ {
		assignment.ProofPath[i] = proofPath[i]
		assignment.Helpers[i] = helpers[i]
		assignment.LeafCounts[i] = mt.PathCounts[i]
		assignment.LeafNumBytes[i] = mt.PathBytes[i]
	}

	fmt.Printf("       ✅ 联合证明 witness 已构造\n")
	return assignment, chunkHash, merkleRoot, nil
}
