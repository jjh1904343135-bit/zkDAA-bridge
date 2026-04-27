#!/usr/bin/env bash
set -e

# === 参数检查 ===
if [ -z "$1" ]; then
  echo "用法: $0 <节点数量N，例如: 10 或 20>"
  exit 1
fi

NODES="$1"

# === 一些路径（根据你现在的情况设定，如果不一样自己改一下） ===
PRIVCHAIN_DIR="/root/privchainA"
ZKHTLC_DIR="/root/zk-htlc"
CONTRACTS_DIR="$ZKHTLC_DIR/blockchain-contracts"

# Hardhat / Go 共用的私钥对应地址（你现在用的那个）
HARDHAT_ADDR="0x279D2b1aB48736BDA48ECFd9e90BaBbCb9bCAb7E"

# TODO: 把这里换成你 node1 账号的密码（setup_clique_network.sh 里创建账号时用的那个）
NODE1_PASS="passwd1"

# 1. 杀掉旧 geth
echo "👉 停掉旧的 geth 进程..."
pkill -f geth 2>/dev/null || true

# 2. 重建 N 节点 Clique 网络
cd "$PRIVCHAIN_DIR"

echo "👉 重新初始化 $NODES 个签名节点..."
./setup_clique_network.sh "$NODES"

echo "👉 启动 $NODES 个签名节点..."
./start_clique_network.sh "$NODES"

echo "⏳ 等待节点启动和互相连接 (10s)..."
sleep 10

echo "👉 查看当前 peer 数:"
geth attach --exec 'net.peerCount' node1/geth.ipc || true

# 3. 给 Hardhat / Go 账号转钱
echo "  给部署+实验账号 $HARDHAT_ADDR 转 100 ETH"


geth attach node1/geth.ipc <<EOF
console.log("coinbase =", eth.coinbase);
console.log("coinbase balance(before) =", eth.getBalance(eth.coinbase));
console.log("target balance(before)   =", eth.getBalance("$HARDHAT_ADDR"));
personal.unlockAccount(eth.coinbase, "$NODE1_PASS", 0);
var txHash = eth.sendTransaction({
  from: eth.coinbase,
  to: "$HARDHAT_ADDR",
  value: web3.toWei(100, "ether")
});
console.log("sendTransaction txHash =", txHash);
miner.start();
admin.sleepBlocks(1);
miner.stop();
console.log("target balance(after)    =", eth.getBalance("$HARDHAT_ADDR"));
exit;
EOF

# 4. 在当前这条链上用 Hardhat 部署合约

echo "  用 Hardhat 在 privchain 上重新部署合约       "


cd "$CONTRACTS_DIR"

# 这里会把输出同时写到 /tmp/deploy_output.txt，方便后面 grep 地址
npx hardhat run --network privchain scripts/deploy_privchain.js | tee /tmp/deploy_output.txt

# 5. 从输出里提取 Verifier / SC-A / SC-B 地址
# 根据你脚本里的中文日志做的简单提取，如果你日志有微调，下面的 grep 关键字也要跟着改一下
VERIFIER_ADDR=$(grep -i "Verifier 合约成功部署到地址" /tmp/deploy_output.txt | awk '{print $NF}')
CONTRACT_A_ADDR=$(grep -i "DataMigration 合约 A" /tmp/deploy_output.txt | awk '{print $NF}')
CONTRACT_B_ADDR=$(grep -i "DataMigration 合约 B" /tmp/deploy_output.txt | awk '{print $NF}')

echo
echo "==============================================="
echo "  当前 N = $NODES 链上的合约地址如下：        "
echo "==============================================="
echo "Verifier     : $VERIFIER_ADDR"
echo "Contract A   : $CONTRACT_A_ADDR"
echo "Contract B   : $CONTRACT_B_ADDR"
echo "==============================================="
echo
echo "👉 请把上面的 Contract A / Contract B 地址，"
echo "   手动写入 Go 的 main.go 中："
echo
echo "   contractA_Addr := \"$CONTRACT_A_ADDR\""
echo "   contractB_Addr := \"$CONTRACT_B_ADDR\""
echo
echo "在 $ZKHTLC_DIR 目录执行类似命令："
echo
echo "   cd $ZKHTLC_DIR"
echo "   go run main.go -size 8388608 -chunk 1024 -output metrics_${NODES}nodes.json"
echo

