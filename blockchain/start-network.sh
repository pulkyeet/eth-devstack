#!/bin/bash

echo "========================================="
echo "ETHEREUM TESTNET STARTUP"
echo "========================================="

cd ~/eth-devstack/blockchain

# 1. Check current state
echo "1. Checking current state..."
RUNNING=$(docker ps --filter "name=eth-signer" -q | wc -l)
echo "   Running signers: $RUNNING/7"

if [ "$RUNNING" -eq 0 ]; then
  echo "2. Starting Docker containers..."
  docker-compose up -d
  echo "   Waiting 15s for initialization..."
  sleep 15
elif [ "$RUNNING" -lt 7 ]; then
  echo "2. Some containers missing, restarting all..."
  docker-compose restart
  sleep 10
else
  echo "2. All containers running, skipping startup"
fi

# 3. Verify running
echo "3. Container status:"
docker ps --filter "name=eth" --format "table {{.Names}}\t{{.Status}}"

# 4. Check block sync
echo "4. Current block:"
BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc 2>/dev/null || echo "ERROR")
echo "   Block: $BLOCK"

if [ "$BLOCK" == "ERROR" ]; then
  echo "   ⚠️  Geth not responding, waiting 10s..."
  sleep 10
  BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc)
  echo "   Block: $BLOCK"
fi

# 5. Check peer connectivity
echo "5. Peer connectivity:"
for i in {1..7}; do
  PEERS=$(docker exec eth-signer$i geth attach --exec 'net.peerCount' /data/geth.ipc 2>/dev/null || echo "0")
  echo "   Signer$i: $PEERS peers"
done

# 6. If peers low, run mesh script
LOW_PEERS=$(docker exec eth-signer1 geth attach --exec 'net.peerCount' /data/geth.ipc 2>/dev/null)
if [ "$LOW_PEERS" -lt 6 ]; then
  echo "6. Low peer count, rebuilding mesh..."
  ./mesh-peers.sh
else
  echo "6. Peer mesh healthy, skipping"
fi

# 7. Check tx-generator
echo "7. Transaction generator:"
TX_RUNNING=$(docker ps --filter "name=eth-tx-generator" -q | wc -l)
if [ "$TX_RUNNING" -eq 0 ]; then
  echo "   ⚠️  TX generator not running, starting..."
  docker-compose up -d eth-tx-generator
  sleep 5
fi
docker logs --tail 3 eth-tx-generator

# 8. Watch blocks
echo "8. Block production (5 blocks):"
for i in {1..5}; do
  BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc 2>/dev/null)
  TXS=$(docker exec eth-signer1 geth attach --exec 'eth.getBlock("latest").transactions.length' /data/geth.ipc 2>/dev/null)
  echo "   Block $BLOCK: $TXS txs"
  sleep 2
done

echo ""
echo "========================================="
echo "STATUS CHECK COMPLETE ✓"
echo "========================================="