const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("\n╔═══════════════════════════════════════════════════╗");
  console.log("║       批量数据迁移合约部署                        ║");
  console.log("╚═══════════════════════════════════════════════════╝\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 部署账户:", deployer.address);
  console.log("💰 账户余额:", hre.ethers.formatEther(await hre.ethers.provider.getBalance(deployer.address)), "ETH\n");

  // 检查验证器合约是否存在
  const verifierPath = path.join(__dirname, "../contracts/BatchUnlockVerifier.sol");
  if (!fs.existsSync(verifierPath)) {
    console.error("❌ 错误: BatchUnlockVerifier.sol 不存在！");
    console.error("请先运行: go run gen_verifier2.go");
    process.exit(1);
  }
  console.log("✅ 找到验证器合约:", verifierPath);

  // 1. 部署验证器合约
  console.log("\n[1/2] 部署 BatchUnlockVerifier...");
  const Verifier = await hre.ethers.getContractFactory("BatchUnlockVerifier");
  const verifier = await Verifier.deploy();
  await verifier.waitForDeployment();
  const verifierAddress = await verifier.getAddress();
  console.log("✅ BatchUnlockVerifier 部署到:", verifierAddress);

  // 2. 部署主合约
  console.log("\n[2/2] 部署 BatchDataMigration...");
  const BatchDataMigration = await hre.ethers.getContractFactory("BatchDataMigration");
  const batchContract = await BatchDataMigration.deploy(verifierAddress);
  await batchContract.waitForDeployment();
  const batchAddress = await batchContract.getAddress();
  console.log("✅ BatchDataMigration 部署到:", batchAddress);

  // 3. 保存部署信息
  const deployInfo = {
    network: hre.network.name,
    chainId: (await hre.ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    contracts: {
      BatchUnlockVerifier: verifierAddress,
      BatchDataMigration: batchAddress
    },
    timestamp: new Date().toISOString()
  };

  const deployPath = path.join(__dirname, "../deployments.json");
  fs.writeFileSync(deployPath, JSON.stringify(deployInfo, null, 2));
  console.log("\n📁 部署信息已保存到:", deployPath);

  // 4. 验证合约
  console.log("\n🔍 验证合约功能...");
  try {
    // 测试提交根
    const testRoot = "0x" + "1".repeat(64);
    const tx = await batchContract.submitBatchRoot(testRoot);
    await tx.wait();
    
    const storedRoot = await batchContract.currentBatchRoot();
    console.log("✅ 合约功能正常 - 测试根:", storedRoot);
  } catch (error) {
    console.log("⚠️  合约验证警告:", error.message);
  }

  console.log("\n╔═══════════════════════════════════════════════════╗");
  console.log("║              部署成功！                           ║");
  console.log("╚═══════════════════════════════════════════════════╝");
  console.log("\n📍 合约地址:");
  console.log("   BatchUnlockVerifier:", verifierAddress);
  console.log("   BatchDataMigration: ", batchAddress);
  console.log("\n💡 在 Go 代码中使用此地址:");
  console.log(`   contractAddr := common.HexToAddress("${batchAddress}")`);
  console.log("");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ 部署失败:", error);
    process.exit(1);
  });
