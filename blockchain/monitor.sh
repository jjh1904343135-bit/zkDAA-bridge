#!/bin/bash
# monitor.sh

while true; do
  clear
  echo "==================== 网络状态 ===================="
  PENDING=$(geth attach node1/geth.ipc --exec "eth.pendingTransactions.length" 2>/dev/null)
  BLOCK=$(geth attach node1/geth.ipc --exec "eth.blockNumber" 2>/dev/null)
  PEERS=$(geth attach node1/geth.ipc --exec "admin.peers.length" 2>/dev/null)
  
  echo "Pending 交易: $PENDING"
  echo "当前区块: $BLOCK"
  echo "连接节点: $PEERS"
  echo "时间: $(date +%H:%M:%S)"
  
  sleep 2
done