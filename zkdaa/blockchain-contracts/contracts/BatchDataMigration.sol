// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./IVerifier.sol";

contract BatchDataMigration {
    // 支持多个批量大小的验证器
    mapping(uint256 => IVerifier) public verifiers;
    
    // 当前批量的 Merkle 根（按批量大小分别存储）
    mapping(uint256 => bytes32) public batchRoots;
    mapping(uint256 => uint256) public batchTimestamps;
    
    // 防双花（全局，跨所有批量大小）
    mapping(uint256 => bool) public usedSerialNumbers;
    
    address public owner;
    
    event BatchRootSubmitted(uint256 indexed batchSize, bytes32 indexed root, uint256 timestamp);
    event Unlocked(uint256 indexed batchSize, uint256 indexed serialNumber, address indexed unlocker);
    event VerifierSet(uint256 indexed batchSize, address verifier);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    /**
     * @dev 设置指定批量大小的验证器
     * @param batchSize 批量大小 (16, 64, 128, 256)
     * @param verifierAddress 验证器合约地址
     */
    function setVerifier(uint256 batchSize, address verifierAddress) external onlyOwner {
        require(verifierAddress != address(0), "Verifier address cannot be zero");
        require(
            batchSize == 16 || batchSize == 64 || batchSize == 128 || batchSize == 256,
            "Invalid batch size"
        );
        
        verifiers[batchSize] = IVerifier(verifierAddress);
        emit VerifierSet(batchSize, verifierAddress);
    }
    
    /**
     * @dev 批量设置所有验证器
     */
    function setAllVerifiers(
        address verifier16,
        address verifier64,
        address verifier128,
        address verifier256
    ) external onlyOwner {
        require(verifier16 != address(0), "Verifier16 cannot be zero");
        require(verifier64 != address(0), "Verifier64 cannot be zero");
        require(verifier128 != address(0), "Verifier128 cannot be zero");
        require(verifier256 != address(0), "Verifier256 cannot be zero");
        
        verifiers[16] = IVerifier(verifier16);
        verifiers[64] = IVerifier(verifier64);
        verifiers[128] = IVerifier(verifier128);
        verifiers[256] = IVerifier(verifier256);
        
        emit VerifierSet(16, verifier16);
        emit VerifierSet(64, verifier64);
        emit VerifierSet(128, verifier128);
        emit VerifierSet(256, verifier256);
    }
    
    /**
     * @dev Operator 提交批量 Merkle 根
     * @param batchSize 批量大小
     * @param root 批量交易的 Merkle 根
     */
    function submitBatchRoot(uint256 batchSize, bytes32 root) external {
        require(root != bytes32(0), "Root cannot be zero");
        require(address(verifiers[batchSize]) != address(0), "Verifier not set for this batch size");
        
        batchRoots[batchSize] = root;
        batchTimestamps[batchSize] = block.timestamp;
        
        emit BatchRootSubmitted(batchSize, root, block.timestamp);
    }
    
    /**
     * @dev DSP 解锁（根据批量大小选择验证器）
     * @param batchSize 批量大小
     * @param proof Groth16 证明
     * @param publicInputs [merkleRoot, serialNumber]
     */
    function unlock(
        uint256 batchSize,
        uint256[8] calldata proof,
        uint256[2] calldata publicInputs
    ) external {
        // 检查验证器是否已设置
        IVerifier verifier = verifiers[batchSize];
        require(address(verifier) != address(0), "Verifier not set for this batch size");
        
        // 提取公开输入
        bytes32 claimedRoot = bytes32(publicInputs[0]);
        uint256 serialNumber = publicInputs[1];
        
        // 1. 验证提交的 root 是否正确
        require(claimedRoot == batchRoots[batchSize], "Invalid Merkle root");
        
        // 2. 防止双花
        require(!usedSerialNumbers[serialNumber], "Serial number already used");
        
        // 3. 使用对应批量大小的验证器进行链上验证
        verifier.verifyProof(proof, publicInputs);
        
        // 4. 标记已使用
        usedSerialNumbers[serialNumber] = true;
        
        emit Unlocked(batchSize, serialNumber, msg.sender);
    }
    
    // 查询接口
    function getBatchRoot(uint256 batchSize) external view returns (bytes32) {
        return batchRoots[batchSize];
    }
    
    function getVerifier(uint256 batchSize) external view returns (address) {
        return address(verifiers[batchSize]);
    }
}
