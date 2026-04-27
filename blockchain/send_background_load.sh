#!/usr/bin/env bash


# 1. 获取节点数量 (默认为 20，测试时可传入 40, 60, 80, 100)
NODE_COUNT=${1:-20}

# 2.20节点 -> 200 pending; 100节点 -> 1000 pending
TARGET_PENDING=$((NODE_COUNT * 10)) 

# RPC配置
RPC_URL="http://127.0.0.1:8545"
PRIV_KEY="c0d166aa26dbbed533c2561f2d4f8d28ec4565c0d284ea5f36c15f465b529063"

REFILL_TPS=100


echo "[INIT] 正在编译发送器..."
if [ -f "go.mod" ]; then go mod tidy; fi
go build -o send_load send_load.go
if [ $? -ne 0 ]; then echo "❌ 编译失败"; exit 1; fi

echo "[$(date +%H:%M:%S)] 🧪 论文实验环境启动 (节点数: $NODE_COUNT)"
echo "   📋 实验约束: 10 x NodeCount = $TARGET_PENDING 笔 Pending 交易"
echo "   📡 RPC 接口: $RPC_URL"

while true; do
  # 1. 获取当前 pending 数量
  HEX_COUNT=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByNumber","params":["pending"],"id":1}' \
    "$RPC_URL" | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

  if [ -z "$HEX_COUNT" ] || [ "$HEX_COUNT" == "null" ]; then
      CURRENT=0
  else
      CURRENT=$(printf "%d" "$HEX_COUNT")
  fi

  # 2. 状态检查与补充
  THRESHOLD=$((TARGET_PENDING * 90 / 100))
  
  if [ "$CURRENT" -lt "$THRESHOLD" ]; then
    # 计算缺口
    TO_SEND=$((TARGET_PENDING - CURRENT))
    
    # 稍微多补一点(例如+20笔)以抵消脚本执行期间被打包的交易，维持水位
    TO_SEND=$((TO_SEND + 20))
    
    echo "[$(date +%H:%M:%S)] 📉 水位过低 ($CURRENT < $TARGET_PENDING) -> 补充 $TO_SEND 笔 (速率: $REFILL_TPS TPS)"
    
    # 调用 Go 程序平滑补充
    ./send_load \
      -rpc "$RPC_URL" \
      -key "$PRIV_KEY" \
      -count "$TO_SEND" \
      -tps "$REFILL_TPS" &
      
    # 等待补充完成的大致时间，避免重复检测
    SLEEP_TIME=$(echo "$TO_SEND / $REFILL_TPS" | bc)
    if [ "$SLEEP_TIME" -lt 2 ]; then SLEEP_TIME=2; fi
    sleep $SLEEP_TIME

  else
    echo -ne "[$(date +%H:%M:%S)] ✅ 环境达标: $CURRENT / $TARGET_PENDING Pending (符合 1:10 比例) \r"
    sleep 1
  fi
done