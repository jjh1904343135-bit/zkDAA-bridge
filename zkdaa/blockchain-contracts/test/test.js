// test/DataMigration.test.js
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("DataMigration 合约测试", function () {
  let verifier, dataMigration;
  let owner, dspA, dspB;
  let dataId, hashLock, timeout;

  beforeEach(async function () {
    // 获取测试账户
    [owner, dspA, dspB] = await ethers.getSigners();

    // 部署 Verifier 合约（模拟版本，总是返回 true）
    const Verifier = await ethers.getContractFactory("Verifier");
    verifier = await Verifier.deploy();
    await verifier.deployed();

    // 部署 DataMigration 合约
    const DataMigration = await ethers.getContractFactory("DataMigration");
    dataMigration = await DataMigration.deploy(verifier.address);
    await dataMigration.deployed();

    // 准备测试数据
    dataId = ethers.utils.formatBytes32String("test_data_123");
    hashLock = ethers.utils.formatBytes32String("test_hash_lock");
    timeout = 3600; // 1 小时
  });

  describe("部署验证", function () {
    it("应该正确设置 Verifier 地址", async function () {
      expect(await dataMigration.verifier()).to.equal(verifier.address);
    });
  });

  describe("Lock 功能", function () {
    it("应该成功锁定哈希", async function () {
      await expect(
        dataMigration.connect(dspA).lock(hashLock, dataId, timeout)
      )
        .to.emit(dataMigration, "Locked")
        .withArgs(hashLock, dspA.address, dataId, await getTimeout(timeout));
    });

    it("不应该允许重复锁定相同的哈希", async function () {
      await dataMigration.connect(dspA).lock(hashLock, dataId, timeout);
      
      await expect(
        dataMigration.connect(dspB).lock(hashLock, dataId, timeout)
      ).to.be.revertedWith("Lock already exists");
    });

    it("不应该允许锁定零哈希", async function () {
      const zeroHash = ethers.constants.HashZero;
      
      await expect(
        dataMigration.connect(dspA).lock(zeroHash, dataId, timeout)
      ).to.be.revertedWith("Hash cannot be zero");
    });

    it("应该正确存储锁定信息", async function () {
      await dataMigration.connect(dspA).lock(hashLock, dataId, timeout);
      
      const lockDetails = await dataMigration.activeLocks(hashLock);
      expect(lockDetails.locker).to.equal(dspA.address);
      expect(lockDetails.dataId).to.equal(dataId);
    });
  });

  describe("Unlock 功能", function () {
    beforeEach(async function () {
      // 先锁定
      await dataMigration.connect(dspA).lock(hashLock, dataId, timeout);
    });

    it("应该成功解锁（使用模拟 Verifier）", async function () {
      // 构造模拟的证明参数
      const proof_a = [1, 2];
      const proof_b = [[3, 4], [5, 6]];
      const proof_c = [7, 8];
      const publicInputs = [hashLock, 12345]; // [h, sn]

      await expect(
        dataMigration.connect(dspB).unlock(proof_a, proof_b, proof_c, publicInputs)
      )
        .to.emit(dataMigration, "Unlocked")
        .withArgs(hashLock);
    });

    it("解锁后应该删除锁定记录", async function () {
      const proof_a = [1, 2];
      const proof_b = [[3, 4], [5, 6]];
      const proof_c = [7, 8];
      const publicInputs = [hashLock, 12345];

      await dataMigration.connect(dspB).unlock(proof_a, proof_b, proof_c, publicInputs);
      
      const lockDetails = await dataMigration.activeLocks(hashLock);
      expect(lockDetails.locker).to.equal(ethers.constants.AddressZero);
    });

    it("不应该允许解锁不存在的哈希", async function () {
      const nonExistentHash = ethers.utils.formatBytes32String("non_existent");
      const proof_a = [1, 2];
      const proof_b = [[3, 4], [5, 6]];
      const proof_c = [7, 8];
      const publicInputs = [nonExistentHash, 12345];

      await expect(
        dataMigration.connect(dspB).unlock(proof_a, proof_b, proof_c, publicInputs)
      ).to.be.revertedWith("Lock does not exist");
    });
  });

  describe("Reclaim 功能", function () {
    beforeEach(async function () {
      await dataMigration.connect(dspA).lock(hashLock, dataId, timeout);
    });

    it("应该在超时后允许原锁定者回收", async function () {
      // 增加区块链时间
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");

      await expect(
        dataMigration.connect(dspA).reclaim(hashLock)
      )
        .to.emit(dataMigration, "Reclaimed")
        .withArgs(hashLock);
    });

    it("不应该在超时前允许回收", async function () {
      await expect(
        dataMigration.connect(dspA).reclaim(hashLock)
      ).to.be.revertedWith("Timeout not reached");
    });

    it("不应该允许非锁定者回收", async function () {
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");

      await expect(
        dataMigration.connect(dspB).reclaim(hashLock)
      ).to.be.revertedWith("Only locker can reclaim");
    });

    it("回收后应该删除锁定记录", async function () {
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");

      await dataMigration.connect(dspA).reclaim(hashLock);
      
      const lockDetails = await dataMigration.activeLocks(hashLock);
      expect(lockDetails.locker).to.equal(ethers.constants.AddressZero);
    });
  });

  describe("Gas 消耗测试", function () {
    it("应该记录 lock 操作的 Gas 消耗", async function () {
      const tx = await dataMigration.connect(dspA).lock(hashLock, dataId, timeout);
      const receipt = await tx.wait();
      
      console.log("      Lock Gas 消耗:", receipt.gasUsed.toString());
      expect(receipt.gasUsed).to.be.lt(100000); // 应该小于 100k gas
    });

    it("应该记录 unlock 操作的 Gas 消耗", async function () {
      await dataMigration.connect(dspA).lock(hashLock, dataId, timeout);
      
      const proof_a = [1, 2];
      const proof_b = [[3, 4], [5, 6]];
      const proof_c = [7, 8];
      const publicInputs = [hashLock, 12345];

      const tx = await dataMigration.connect(dspB).unlock(proof_a, proof_b, proof_c, publicInputs);
      const receipt = await tx.wait();
      
      console.log("      Unlock Gas 消耗:", receipt.gasUsed.toString());
      // 真实的 Groth16 验证约需 260k gas，这里是模拟版本
    });
  });

  // 辅助函数
  async function getTimeout(duration) {
    const block = await ethers.provider.getBlock("latest");
    return block.timestamp + duration;
  }
});
