const hre = require("hardhat");

async function main() {
    // 检查两个系统
    const accounts = await hre.ethers.getSigners();
    console.log("检查部署的合约...\n");
    
    // 获取所有部署的合约地址
    const provider = hre.ethers.provider;
    const latestBlock = await provider.getBlockNumber();
    
    console.log("当前区块高度:", latestBlock);
    
    // 列出最近的合约部署
    for (let i = Math.max(0, latestBlock - 20); i <= latestBlock; i++) {
        const block = await provider.getBlock(i);
        if (block && block.transactions) {
            for (const txHash of block.transactions) {
                const receipt = await provider.getTransactionReceipt(txHash);
                if (receipt && receipt.contractAddress) {
                    const code = await provider.getCode(receipt.contractAddress);
                    console.log(`\n区块 ${i}: 合约部署到 ${receipt.contractAddress}`);
                    console.log(`  代码大小: ${code.length} bytes`);
                    
                    // 尝试识别合约类型
                    if (code.length > 50000) {
                        console.log("  可能是: 验证器合约（代码较大）");
                    }
                }
            }
        }
    }
}

main().catch(console.error);
