const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("\n╔═══════════════════════════════════════════════════════════╗");
  console.log("║       部署所有批量大小的验证器和主合约                   ║");
  console.log("╚═══════════════════════════════════════════════════════════╝\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 部署账户:", deployer.address);
  console.log("💰 账户余额:", hre.ethers.formatEther(await hre.ethers.provider.getBalance(deployer.address)), "ETH\n");

  const batchSizes = [16, 64, 128, 256];
  const deployments = {};

  // 部署所有验证器
  for (const size of batchSizes) {
    console.log(`\n[${size}] 部署 BatchUnlockVerifier${size}...`);
    
    const VerifierFactory = await hre.ethers.getContractFactory(`BatchUnlockVerifier${size}`);
    const verifier = await VerifierFactory.deploy();
    await verifier.waitForDeployment();
    const verifierAddr = await verifier.getAddress();
    
    console.log(`✅ BatchUnlockVerifier${size} 部署到:`, verifierAddr);
    
    deployments[`verifier_${size}`] = verifierAddr;
  }

  // 部署主合约（使用批量 16 的验证器作为默认）
  console.log("\n[主合约] 部署 BatchDataMigration...");
  const BatchDataMigration = await hre.ethers.getContractFactory("BatchDataMigration");
  const mainContract = await BatchDataMigration.deploy(deployments.verifier_16);
  await mainContract.waitForDeployment();
  const mainAddr = await mainContract.getAddress();
  
  console.log("✅ BatchDataMigration 部署到:", mainAddr);
  deployments.main_contract = mainAddr;

  // 保存部署信息
  const deployInfo = {
    network: hre.network.name,
    chainId: (await hre.ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    contracts: deployments,
    timestamp: new Date().toISOString()
  };

  const deployPath = path.join(__dirname, "../all_deployments.json");
  fs.writeFileSync(deployPath, JSON.stringify(deployInfo, null, 2));
  
  console.log("\n📁 部署信息已保存到:", deployPath);

  console.log("\n╔═══════════════════════════════════════════════════════════╗");
  console.log("║                   部署完成！                              ║");
  console.log("╚═══════════════════════════════════════════════════════════╝");
  
  console.log("\n📍 合约地址汇总:");
  for (const [key, addr] of Object.entries(deployments)) {
    console.log(`   ${key}: ${addr}`);
  }
  
  console.log("\n💡 下一步:");
  console.log("   1. 更新 Go 代码中的合约地址");
  console.log(`   2. 主合约地址: ${mainAddr}`);
  console.log("   3. 运行测试: go run cmd/batch_onchain/main.go -batch");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ 部署失败:", error);
    process.exit(1);
  });
