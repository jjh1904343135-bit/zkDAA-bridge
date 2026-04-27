// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IVerifier {
    /// @notice 验证 Groth16 证明
    /// @param proof 包含 8 个元素的证明数组 [Ax, Ay, Bx0, Bx1, By0, By1, Cx, Cy]
    /// @param input 包含 2 个公开输入的数组 [h, sn]
    function verifyProof(
        uint256[8] calldata proof,
        uint256[2] calldata input
    ) external view;
}