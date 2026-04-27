#!/usr/bin/env bash
set -e

if [ -z "$1" ]; then
  echo "用法: $0 <节点个数>"
  exit 1
fi

N=$1
echo "👉 [Setup] 正在重置环境 (节点数: $N)..."

# 1. 清理
pkill -f "geth" 2>/dev/null || true
pkill -f "bootnode" 2>/dev/null || true
sleep 2
rm -rf node* genesis.json

addresses=()

# 2. 创建节点账户
for i in $(seq 1 "$N"); do
  mkdir -p "node$i"
  echo "passwd$i" > "node$i/password.txt"
  addr=$(geth --datadir "node$i" account new --password "node$i/password.txt" 2>&1 | grep -o '0x[0-9a-fA-F]\+' | head -n1)
  addresses+=("$addr")
done

# 3. 生成 Genesis
export ADDRS="${addresses[*]}"
export NODE_COUNT="$N"

python3 - <<'PYEOF' > genesis.json
import os, json, sys

addrs = os.environ["ADDRS"].split()
node_count = int(os.environ["NODE_COUNT"])
clean = [a.lower().replace("0x","") for a in addrs]
extra = "0x" + "00"*32 + "".join(clean) + "00"*65

if node_count <= 10:
    period = 2
    gas_limit = 50_000_000
elif node_count <= 50:
    period = 3
    gas_limit = 60_000_000
else:
    period = 5
    gas_limit = 80_000_000

genesis = {
  "config": {
    "chainId": 12345,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "clique": {
      "period": period,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": hex(gas_limit),
  "extraData": extra,
  "alloc": {}
}

# 1. 节点账户
for addr in addrs:
  genesis["alloc"][addr] = {"balance": "0x3635C9ADC5DEA00000"}

balance_wei = 1_000_000 * 10**18 

# 2. 🔥 [账户 A] 上帝账户 (用于部署合约 & Go程序)
# 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
genesis["alloc"]["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"] = {"balance": hex(balance_wei)}

# 3. 🔥 [账户 B] 背景压测账户 (用于 send 程序)
# 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
genesis["alloc"]["0x70997970C51812dc3A010C7d01b50e0d17dc79C8"] = {"balance": hex(balance_wei)}

print(json.dumps(genesis, indent=2))
PYEOF

# 4. 初始化
for i in $(seq 1 "$N"); do
  geth --datadir "node$i" init genesis.json > /dev/null 2>&1
done

echo "✅ Setup 完成 (双账户已充值)"