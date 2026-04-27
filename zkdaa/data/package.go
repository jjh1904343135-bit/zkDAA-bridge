package data

import "time"

// DataPackage 表示用户发送给 DSP 的数据包
type DataPackage struct {
	DataID    string // 数据唯一标识
	TargetIP  string // 目标 DSP 的 IP 地址 (对于 DSPA，是 DSPB 的 IP)
	FileData  []byte // 原始文件数据
	ChunkSize int    // 分块大小（例如 1KB）
}

// TransferInfo 记录数据传输信息
type TransferInfo struct {
	From        string // 发送方 (User/DSPA)
	To          string // 接收方 (DSPA/DSPB)
	DataID      string
	TotalChunks int  // 总块数
	Transferred bool // 是否已传输完成
}

// AuditChallenge 用户发起的审计挑战
type AuditChallenge struct {
	DataID     string
	ChunkIndex int   // 被挑战的数据块索引
	Timestamp  int64 // 挑战时间戳
}

// AuditResponse DSP 的审计响应
type AuditResponse struct {
	ChunkIndex  int      // 被挑战块的索引
	ChunkData   []byte   // 原始数据块
	MerkleProof [][]byte // Merkle 证明路径
	MerkleRoot  []byte   // 文件的 Merkle 根 (CIDF)
}

// NewDataPackage 创建新的数据包
func NewDataPackage(dataID string, fileData []byte, targetIP string, chunkSize int) *DataPackage {
	return &DataPackage{
		DataID:    dataID,
		TargetIP:  targetIP,
		FileData:  fileData,
		ChunkSize: chunkSize,
	}
}

// NewAuditChallenge 创建新的审计挑战
func NewAuditChallenge(dataID string, chunkIndex int) *AuditChallenge {
	return &AuditChallenge{
		DataID:     dataID,
		ChunkIndex: chunkIndex,
		Timestamp:  time.Now().Unix(),
	}
}
