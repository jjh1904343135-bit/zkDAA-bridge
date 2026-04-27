#!/usr/bin/env bash
# 诊断节点连接问题

echo "时间: $(date)"
echo

# 1. 检查进程
echo "📊 运行的 geth 进程:"
GETH_COUNT=$(ps aux | grep "geth --datadir" | grep -v grep | wc -l)
echo "   总数: $GETH_COUNT"
if [ "$GETH_COUNT" -eq 0 ]; then
  echo "   ❌ 没有 geth 进程在运行！"
  exit 1
fi
echo

# 2. 检查 IPC 文件
echo "🔌 IPC 文件状态:"
for i in {1..5}; do
  if [ -S "node$i/geth.ipc" ]; then
    echo "   ✅ node$i/geth.ipc 存在"
    ls -lh "node$i/geth.ipc" | awk '{print "      权限: " $1 " 所有者: " $3}'
  else
    echo "   ❌ node$i/geth.ipc 不存在"
  fi
done
echo

# 3. 检查 node1 的 bootnode 信息
echo "🔗 node1 enode:"
if [ -S "node1/geth.ipc" ]; then
  ENODE=$(geth attach --exec "admin.nodeInfo.enode" node1/geth.ipc 2>/dev/null | tr -d '"')
  echo "   $ENODE"
  
  # 检查 IP 地址
  if [[ "$ENODE" == *"[::]"* ]] || [[ "$ENODE" == *"0.0.0.0"* ]]; then
  fi
else
  echo "   ❌ 无法连接 node1"
fi
echo

# 4. 检查其他节点的 enode（前3个）
echo "📡 其他节点的 enode:"
for i in {2..4}; do
  if [ -S "node$i/geth.ipc" ]; then
    NODE_ENODE=$(geth attach --exec "admin.nodeInfo.enode" "node$i/geth.ipc" 2>/dev/null | tr -d '"')
    if [ -n "$NODE_ENODE" ]; then
      echo "   node$i: ${NODE_ENODE:0:50}..."
    else
      echo "   node$i: 无法获取 enode"
    fi
  else
    echo "   node$i: IPC 不存在"
  fi
done
echo

# 5. 检查 node1 的 peers
echo "👥 node1 连接的 peers:"
if [ -S "node1/geth.ipc" ]; then
  PEERS=$(geth attach --exec "admin.peers.length" node1/geth.ipc 2>/dev/null)
  echo "   连接数: $PEERS"
  
  if [ "$PEERS" = "0" ]; then
    echo "   ⚠️  没有连接的节点"
    echo
    echo "   📋 检查 node1 的监听地址:"
    geth attach --exec "admin.nodeInfo.listenAddr" node1/geth.ipc 2>/dev/null
  fi
fi
echo

# 6. 检查端口
echo "🌐 端口监听状态:"
for port in 30301 30302 30303; do
  if netstat -tuln 2>/dev/null | grep -q ":$port "; then
    echo "   ✅ 端口 $port 正在监听"
  elif ss -tuln 2>/dev/null | grep -q ":$port "; then
    echo "   ✅ 端口 $port 正在监听 (ss)"
  else
    echo "   ❌ 端口 $port 未监听"
  fi
done
echo

# 7. 检查最近的日志错误
echo "📝 node2 最近的日志 (最后10行):"
tail -10 node2.log | grep -i "error\|warn\|fail" || echo "   无明显错误"
echo

echo "📝 node3 最近的日志 (最后10行):"
tail -10 node3.log | grep -i "error\|warn\|fail" || echo "   无明显错误"
echo

# 8. 尝试手动连接测试
echo "🔧 尝试手动添加 peer (node2 -> node1)..."
if [ -S "node2/geth.ipc" ] && [ -n "$ENODE" ]; then
  FIXED_ENODE=$(echo "$ENODE" | sed 's/@\[::\]/@127.0.0.1/g' | sed 's/@0\.0\.0\.0/@127.0.0.1/g')
  RESULT=$(geth attach --exec "admin.addPeer('$FIXED_ENODE')" node2/geth.ipc 2>/dev/null)
  echo "   结果: $RESULT"
  
  sleep 2
  PEERS=$(geth attach --exec "admin.peers.length" node1/geth.ipc 2>/dev/null)
  echo "   node1 连接数: $PEERS"
fi
echo
