const { ethers } = require("ethers");

// 常见的验证器错误
const errors = [
    "ProofInvalid()",
    "InvalidProof()",
    "VerificationFailed()",
];

errors.forEach(err => {
    const selector = ethers.id(err).slice(0, 10);
    console.log(`${err}: ${selector}`);
});

// 你遇到的错误
console.log("\n你的错误: 0x7fcdd1f4");
