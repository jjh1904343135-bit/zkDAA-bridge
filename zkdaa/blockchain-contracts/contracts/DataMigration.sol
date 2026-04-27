// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// ==========================================
// 1. 定义两个不同的验证器接口
// ==========================================

// 对应 DSPA 的 UnlockCircuit (2 个公共输入: H, Sn)
interface IUnlockVerifier {
    // 🔥 修正：移除 returns (bool)，Gnark Verifier 验证失败会直接 revert
    function verifyProof(
        uint256[8] calldata proof,
        uint256[2] calldata pubInputs
    ) external view;
}

// 对应 DSPB 的 AuditUnlockCircuit (3 个公共输入: ChunkIndex, ChunkHash, H)
interface IAuditVerifier {
    // 🔥 修正：移除 returns (bool)
    function verifyProof(
        uint256[8] calldata proof,
        uint256[3] calldata pubInputs
    ) external view;
}

contract DataMigration {
    // --- 状态变量 ---
    
    // 分别存储两个验证器实例
    IUnlockVerifier public immutable unlockVerifier;
    IAuditVerifier public immutable auditVerifier;

    struct LockDetails {
        address locker;
        bytes32 dataId;
        uint256 timeout;
    }

    mapping(bytes32 => LockDetails) public activeLocks;

    // --- 事件 ---
    event Locked(bytes32 indexed h, address indexed locker, bytes32 indexed dataId, uint256 timeout);
    event Unlocked(bytes32 indexed h, string mechanic); 
    event Reclaimed(bytes32 indexed h);

    // --- 构造函数 ---
    constructor(address _unlockVerifier, address _auditVerifier) {
        require(_unlockVerifier != address(0), "Unlock Verifier cannot be zero");
        require(_auditVerifier != address(0), "Audit Verifier cannot be zero");
        
        unlockVerifier = IUnlockVerifier(_unlockVerifier);
        auditVerifier = IAuditVerifier(_auditVerifier);
    }

    // --- 通用锁定函数 ---
    function lock(bytes32 _h, bytes32 _dataId, uint256 _timeoutDuration) external {
        require(activeLocks[_h].locker == address(0), "Lock already exists");
        require(_h != bytes32(0), "Hash cannot be zero");

        uint256 timeoutTimestamp = block.timestamp + _timeoutDuration;

        activeLocks[_h] = LockDetails({
            locker: msg.sender,
            dataId: _dataId,
            timeout: timeoutTimestamp
        });

        emit Locked(_h, msg.sender, _dataId, timeoutTimestamp);
    }

    // --- DSPA 专用解锁函数 (2 Inputs: H, Sn) ---
    function unlock(
        uint256[8] calldata proof,
        uint256[2] calldata publicInputs
    ) external {
        // 在 UnlockCircuit 中，H 通常是第 0 个输入
        bytes32 h = bytes32(publicInputs[0]);

        // 1. 检查锁是否存在
        require(activeLocks[h].locker != address(0), "Lock does not exist");
        
        // 2. 调用 UnlockVerifier 验证
        // 🔥 修正：直接调用，如果失败 Verifier 会 revert "ProofInvalid"
        // 这里的 verifyProof 没有返回值，不能放在 require 里面
        unlockVerifier.verifyProof(proof, publicInputs);

        // 3. 删除锁
        delete activeLocks[h];

        // 4. 发出事件
        emit Unlocked(h, "DSPA_Unlock");
    }

    // --- DSPB 专用审计解锁函数 (3 Inputs: Index, Hash, H) ---
    function auditUnlock(
        uint256[8] calldata proof,
        uint256[3] calldata publicInputs
    ) external {
        // 🔥 关键确认：根据 Go 代码 batch.go 发送的顺序 [Index, ChunkHash, H]
        // H 位于数组的最后一个位置 (索引 2)
        bytes32 h = bytes32(publicInputs[2]);

        // 1. 检查锁是否存在
        require(activeLocks[h].locker != address(0), "Lock does not exist");

        // 2. 调用 AuditVerifier 验证
        // 🔥 修正：直接调用，验证失败会自动 revert
        auditVerifier.verifyProof(proof, publicInputs);

        // 3. 删除锁
        delete activeLocks[h];

        // 4. 发出事件
        emit Unlocked(h, "DSPB_AuditUnlock");
    }

    // --- 取回函数 ---
    function reclaim(bytes32 _h) external {
        LockDetails memory currentLock = activeLocks[_h];

        require(currentLock.locker != address(0), "Lock does not exist");
        require(currentLock.locker == msg.sender, "Only locker can reclaim");
        require(block.timestamp >= currentLock.timeout, "Timeout not reached");

        delete activeLocks[_h];

        emit Reclaimed(_h);
    }
}