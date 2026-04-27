#!/usr/bin/env bash

# 自动探测节点数
N=$(find . -maxdepth 1 -type d -name "node*" | wc -l)
echo "🚀 [Start] 正在启动 $N 个节点..."

# 1. 启动 Node1 (作为 Bootnode + RPC)
# 注意：开启了 --http.api personal,admin 以便后续控制
echo "   - Starting Node 1 (RPC Port 8545)..."
NODE1_ADDR=$(cat node1/keystore/* | grep -o '"address":"[^"]*"' | cut -d'"' -f4)

nohup geth \
  --datadir node1 \
  --networkid 12345 \
  --port 30301 \
  --nat "extip:127.0.0.1" \
  --authrpc.port 8551 \
  --http --http.addr "0.0.0.0" --http.port 8545 --http.corsdomain "*" --http.api "eth,net,web3,personal,admin,miner,txpool" \
  --unlock "0x$NODE1_ADDR" --password node1/password.txt --mine --miner.etherbase "0x$NODE1_ADDR" \
  --allow-insecure-unlock \
  --maxpeers 200 \
  --txpool.globalslots 10000 --txpool.globalqueue 5000 \
  --syncmode "full" --gcmode archive \
  > node1.log 2>&1 &

# 等待 Node1 启动并生成 IPC
sleep 3
if [ ! -S "node1/geth.ipc" ]; then
    echo "❌ Node 1 启动失败 (未找到 IPC)"
    exit 1
fi

# 获取 Node1 的 Enode URL
ENODE=$(geth attach --exec "admin.nodeInfo.enode" node1/geth.ipc | tr -d '"')
# 将 [::] 替换为 127.0.0.1 确保连接顺畅
ENODE=${ENODE/\[\:\:\]/127.0.0.1}
echo "🔗 Bootnode Enode: $ENODE"

# 2. 启动其他节点 (Node 2 ~ N)
for i in $(seq 2 "$N"); do
  NODE_ADDR=$(cat node$i/keystore/* | grep -o '"address":"[^"]*"' | cut -d'"' -f4)
  PORT=$((30300 + i))
  AUTH_PORT=$((8550 + i))
  
  nohup geth \
    --datadir "node$i" \
    --networkid 12345 \
    --port "$PORT" \
    --nat "extip:127.0.0.1" \
    --authrpc.port "$AUTH_PORT" \
    --unlock "0x$NODE_ADDR" --password "node$i/password.txt" --mine --miner.etherbase "0x$NODE_ADDR" \
    --allow-insecure-unlock \
    --bootnodes "$ENODE" \
    --maxpeers 200 \
    --txpool.globalslots 10000 --txpool.globalqueue 5000 \
    --syncmode "full" \
    > "node$i.log" 2>&1 &
    
    # 稍微错开启动时间，避免 CPU 瞬时爆炸
    sleep 0.1
done

echo "✅ 所有 $N 个节点已后台启动！"
echo "⏳ 等待 10秒让节点互联..."
sleep 10

# 3. 检查连接数
PEER_COUNT=$(geth attach --exec "admin.peers.length" node1/geth.ipc | tr -d '\n')
echo "📊 当前 Node1 连接数: $PEER_COUNT / $((N-1))"

# 如果连接数太少，尝试手动互联 (双保险)
if [ "$PEER_COUNT" -eq 0 ] && [ "$N" -gt 1 ]; then
    echo "⚠️ 连接数异常，尝试强制互联..."
    for i in $(seq 2 "$N"); do
        geth attach --exec "admin.addPeer('$ENODE')" node$i/geth.ipc >/dev/null 2>&1
    done
fi