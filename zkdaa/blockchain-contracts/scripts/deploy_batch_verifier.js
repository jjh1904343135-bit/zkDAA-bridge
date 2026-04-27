const hre = require("hardhat");

async function main() {
  console.log("\n🚀 部署批量解锁验证器和批量迁移合约...\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 部署账户:", deployer.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 账户余额:", hre.ethers.formatEther(balance), "ETH\n");

  // Step 1: 部署 BatchUnlockVerifier
  console.log("⏳ [1/2] 部署 BatchUnlockVerifier...");
  const Verifier = await hre.ethers.getContractFactory("BatchUnlockVerifier");
  const verifier = await Verifier.deploy();
  await verifier.waitForDeployment();
  const verifierAddress = await verifier.getAddress();
  console.log("✅ BatchUnlockVerifier 部署成功:", verifierAddress);

  // Step 2: 部署 BatchDataMigration（使用新的验证器）
  console.log("\n⏳ [2/2] 部署 BatchDataMigration...");
  const BatchDataMigration = await hre.ethers.getContractFactory("BatchDataMigration");
  const migration = await BatchDataMigration.deploy(verifierAddress);
  await migration.waitForDeployment();
  const migrationAddress = await migration.getAddress();
  console.log("✅ BatchDataMigration 部署成功:", migrationAddress);

  // 保存地址
  const fs = require('fs');
  const addresses = {
    network: hre.network.name,
    verifier: verifierAddress,
    batchDataMigration: migrationAddress,
    deployer: deployer.address,
    timestamp: new Date().toISOString()
  };
  
  fs.writeFileSync(
    'deployed_batch_addresses.json',
    JSON.stringify(addresses, null, 2)
  );
  
  console.log("\n💾 地址已保存到 deployed_batch_addresses.json");

  // 打印使用说明
  console.log("\n" + "=".repeat(60));
  console.log("📋 部署完成！下一步操作:");
  console.log("=".repeat(60));
  console.log("\n1. 更新 Go 代码中的合约地址:");
  console.log("   编辑 cmd/onchain_test/main.go 第 69 行:");
  console.log(`   contractAddr := common.HexToAddress("${migrationAddress}")`);
  
  console.log("\n2. 运行链上测试:");
  console.log("   go run cmd/onchain_test/main.go");
  
  console.log("\n" + "=".repeat(60) + "\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ 部署失败:", error);
    process.exit(1);
  });