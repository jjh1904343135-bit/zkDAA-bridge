const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("\n╔═══════════════════════════════════════════════════════════╗");
  console.log("║     部署支持多验证器的批量数据迁移系统                   ║");
  console.log("╚═══════════════════════════════════════════════════════════╝\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 部署账户:", deployer.address);
  console.log("💰 账户余额:", hre.ethers.formatEther(await hre.ethers.provider.getBalance(deployer.address)), "ETH\n");

  const batchSizes = [16, 64, 128, 256];
  const verifierAddrs = {};

  // 步骤 1: 部署所有验证器
  console.log("=" .repeat(60));
  console.log("步骤 1: 部署所有验证器");
  console.log("=".repeat(60) + "\n");
  
  for (const size of batchSizes) {
    console.log(`[${size}] 部署 BatchUnlockVerifier${size}...`);
    
    const VerifierFactory = await hre.ethers.getContractFactory(`BatchUnlockVerifier${size}`);
    const verifier = await VerifierFactory.deploy();
    await verifier.waitForDeployment();
    const verifierAddr = await verifier.getAddress();
    
    verifierAddrs[size] = verifierAddr;
    console.log(`      ✅ 地址: ${verifierAddr}\n`);
  }

  // 步骤 2: 部署主合约
  console.log("=".repeat(60));
  console.log("步骤 2: 部署 BatchDataMigration 主合约");
  console.log("=".repeat(60) + "\n");
  
  const BatchDataMigration = await hre.ethers.getContractFactory("BatchDataMigration");
  const mainContract = await BatchDataMigration.deploy();
  await mainContract.waitForDeployment();
  const mainAddr = await mainContract.getAddress();
  
  console.log("✅ BatchDataMigration 部署到:", mainAddr, "\n");

  // 步骤 3: 配置验证器
  console.log("=".repeat(60));
  console.log("步骤 3: 配置所有验证器");
  console.log("=".repeat(60) + "\n");
  
  console.log("调用 setAllVerifiers()...");
  const tx = await mainContract.setAllVerifiers(
    verifierAddrs[16],
    verifierAddrs[64],
    verifierAddrs[128],
    verifierAddrs[256]
  );
  await tx.wait();
  console.log("✅ 所有验证器配置完成\n");

  // 验证配置
  console.log("验证配置:");
  for (const size of batchSizes) {
    const configuredAddr = await mainContract.getVerifier(size);
    const match = configuredAddr.toLowerCase() === verifierAddrs[size].toLowerCase();
    console.log(`  批量 ${size}: ${match ? "✅" : "❌"} ${configuredAddr}`);
  }

  // 保存部署信息
  const deployInfo = {
    network: hre.network.name,
    chainId: (await hre.ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    mainContract: mainAddr,
    verifiers: verifierAddrs,
    timestamp: new Date().toISOString()
  };

  const deployPath = path.join(__dirname, "../multi_verifier_deployment.json");
  fs.writeFileSync(deployPath, JSON.stringify(deployInfo, null, 2));
  
  console.log("\n📁 部署信息已保存到:", deployPath);

  console.log("\n╔═══════════════════════════════════════════════════════════╗");
  console.log("║                   部署完成！                              ║");
  console.log("╚═══════════════════════════════════════════════════════════╝");
  
  console.log("\n📍 合约地址:");
  console.log(`   主合约: ${mainAddr}`);
  console.log("\n   验证器:");
  for (const [size, addr] of Object.entries(verifierAddrs)) {
    console.log(`     批量 ${size}: ${addr}`);
  }
  
  console.log("\n💡 下一步:");
  console.log(`   1. 更新 Go 代码中的主合约地址: ${mainAddr}`);
  console.log("   2. 更新 SubmitBatchRoot 和 Unlock 调用，传入 batchSize 参数");
  console.log("   3. 运行测试: go run cmd/batch_onchain/main.go -batch");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ 部署失败:", error);
    process.exit(1);
  });
