const { ethers } = require("hardhat");

async function main() {
  // 🔥 1. 从环境变量获取深度 (Shell 脚本传进来的)
  const depth = process.env.DEPTH;
  if (!depth) {
    throw new Error("❌ 环境变量 DEPTH 未设置！请通过 DEPTH=13 npx hardhat run ... 运行");
  }

  const verifierName = `AuditVerifier_d${depth}`;
  console.log(`\n⚙️  正在部署针对 Depth=${depth} 的合约...`);
  console.log(`   目标合约名: ${verifierName}`);

  const [deployer] = await ethers.getSigners();

  // 2. 部署 UnlockVerifier (固定)
  const UnlockVerifier = await ethers.getContractFactory("UnlockVerifier");
  const unlockVerifier = await UnlockVerifier.deploy();
  await unlockVerifier.waitForDeployment();
  const unlockVerifierAddress = unlockVerifier.target;

  // 3. 部署 AuditVerifier (动态名字!)
  const AuditVerifier = await ethers.getContractFactory(verifierName);
  const auditVerifier = await AuditVerifier.deploy();
  await auditVerifier.waitForDeployment();
  const auditVerifierAddress = auditVerifier.target;

  // 4. 部署 DataMigration (需要传入 Verifiers)
  const DataMigration = await ethers.getContractFactory("DataMigration");
  
  // 部署 Contract A
  const contractA = await DataMigration.deploy(unlockVerifierAddress, auditVerifierAddress);
  await contractA.waitForDeployment();
  const addrA = contractA.target;

  // 部署 Contract B
  const contractB = await DataMigration.deploy(unlockVerifierAddress, auditVerifierAddress);
  await contractB.waitForDeployment();
  const addrB = contractB.target;

  // 🔥 5. 打印特定格式的日志，供 Shell 脚本 grep 抓取
  // Shell 脚本里写的是: grep "DataMigration deployed at:"
  // 所以这里必须输出得一模一样，或者用特定分隔符
  console.log("========================================");
  console.log(`DataMigration deployed at: ${addrA} ${addrB}`); 
  console.log("========================================");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});