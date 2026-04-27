require("@nomicfoundation/hardhat-toolbox");

module.exports = {
  solidity: "0.8.28",
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      chainId: 31337,
    },
    // 你的 geth 私链
    privchain: {
      url: "http://127.0.0.1:8545",
      chainId: 12345, // 对应 setup.sh 里的 chainId
      // 🔥 [修改] 这里换成上帝私钥 (对应 setup.sh 里给钱的那个)
      accounts: [
        "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
      ],
    },
  },
};