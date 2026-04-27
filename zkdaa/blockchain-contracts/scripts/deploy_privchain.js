const { ethers } = require("hardhat");

async function main() {
  // 1. 获取环境变量中的 DEPTH
  const depth = process.env.DEPTH;
  if (!depth) {
    throw new Error("❌ 环境变量 DEPTH 未设置! 请通过 DEPTH=13 npx hardhat run ... 运行");
  }

  const [deployer] = await ethers.getSigners();
  console.log("正在使用账户:", deployer.address, "进行部署...");

  // ====================================================
  // 2. 部署 UnlockVerifier (固定名称)
  // ====================================================
  console.log("准备部署 UnlockVerifier ...");
  const UnlockVerifier = await ethers.getContractFactory("UnlockVerifier");
  const unlockVerifier = await UnlockVerifier.deploy();
  await unlockVerifier.waitForDeployment();
  const unlockAddress = await unlockVerifier.getAddress();
  console.log("✅ UnlockVerifier 部署地址:", unlockAddress);

  // ====================================================
  // 3. 部署 AuditVerifier (动态名称: AuditVerifier_d{depth})
  // ====================================================
  const auditVerifierName = `AuditVerifier_d${depth}`;
  console.log(`准备部署 ${auditVerifierName} ...`);
  
  // 检查 Artifact 是否存在，防止报错
  try {
    await ethers.getContractFactory(auditVerifierName);
  } catch (e) {
    throw new Error(`❌ 找不到合约 ${auditVerifierName}。请检查 artifacts 是否已编译，或 depth 是否正确。`);
  }

  const AuditVerifier = await ethers.getContractFactory(auditVerifierName);
  const auditVerifier = await AuditVerifier.deploy();
  await auditVerifier.waitForDeployment();
  const auditAddress = await auditVerifier.getAddress();
  console.log(`✅ ${auditVerifierName} 部署地址:`, auditAddress);

  // ====================================================
  // 4. 部署 DataMigration 合约 (传入两个验证器地址)
  // ====================================================
  const DataMigration = await ethers.getContractFactory("DataMigration");

  // 部署 SC-A
  console.log("\n准备部署 DataMigration A (SC-A) ...");
  // 注意：这里传入了两个参数 (UnlockVerifier, AuditVerifier)
  const contractA = await DataMigration.deploy(unlockAddress, auditAddress);
  await contractA.waitForDeployment();
  const addrA = await contractA.getAddress();
  console.log("✅ SC-A 部署成功:", addrA);

  // 部署 SC-B
  console.log("准备部署 DataMigration B (SC-B) ...");
  const contractB = await DataMigration.deploy(unlockAddress, auditAddress);
  await contractB.waitForDeployment();
  const addrB = await contractB.getAddress();
  console.log("✅ SC-B 部署成功:", addrB);

  // ====================================================
  // 5. 输出关键日志 (供 Bash 脚本 grep 使用)
  // ====================================================
  console.log("==================================================");
  // 🔥这一行非常重要，run_batch_tests.sh 依赖这句话提取地址🔥
  console.log(`DataMigration deployed at: ${addrA} ${addrB}`);
  console.log("==================================================");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });