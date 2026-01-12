# Ethereum Private Testnet - 7-Validator Clique PoA

**Status:** Production-ready, 41 TPS sustained, 2.0s block time

---

## ARCHITECTURE

### Network Topology
```
┌─────────────────────────────────────────────────────────────────┐
│                     TRANSACTION GENERATOR                        │
│         Round-robin distribution to all 7 signers               │
│              40-120 txs/batch, 2s intervals                     │
└────────┬────────────────────────────────────────────────────────┘
         │
    ┌────┴────┬────┬────┬────┬────┬────┐
    │         │    │    │    │    │    │
┌───▼───┐ ┌──▼──┐ ... (7 signers total, full mesh)
│Signer1│ │Signer2│
│:8545  │ │:8546  │  Each signer:
│:30303 │ │:30303 │  - Seals 1/7 of blocks (round-robin)
└───────┘ └───────┘  - 4GB RAM limit, 30M gas/block
                     - HTTP RPC enabled (eth,net,web3,admin,txpool)
                     - Full peer mesh (6 peers each)
```

### Key Specifications

- **Consensus:** Clique Proof of Authority
- **Block time:** 2 seconds (fixed)
- **Gas limit:** 30M per block
- **Chain ID:** 1337
- **Validators:** 7 signers (round-robin sealing)
- **Throughput:** 40-120 txs/block (30-60 TPS sustained)
- **Network:** Private, Docker bridge network

---

## FILE STRUCTURE
```
~/eth-devstack/blockchain/
├── docker-compose.yml          # All services orchestration
├── genesis.json                # Network genesis configuration
├── mesh-peers.sh              # Peer mesh connection script
│
├── setup/
│   ├── accounts.json          # 100 funded accounts (10k ETH each)
│   └── init-signers.sh        # Initialize signer keystores (if needed)
│
├── tx-generator/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go                # Transaction generator
│
└── data/                      # Persistent volumes (created by Docker)
    ├── signer1/
    ├── signer2/
    ├── ...
    └── signer7/
```

---

## PREREQUISITES

- **OS:** Ubuntu 22.04+ or WSL2
- **Docker:** 24.0+
- **Docker Compose:** 2.20+
- **Resources:**
  - 8GB RAM minimum (7 signers × ~200MB + tx-generator)
  - 4 CPU cores minimum
  - 10GB disk space

**Install Docker:**
```bash
# If not already installed:
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker

# Verify:
docker --version
docker-compose --version
```

---

## GENESIS CONFIGURATION

**File:** `~/eth-devstack/blockchain/genesis.json`

**Critical fields:**
```json
{
  "config": {
    "chainId": 1337,
    "clique": {
      "period": 2,      // Block time in seconds
      "epoch": 30000    // Signer voting period
    },
    "londonBlock": 0    // Enable EIP-1559 from genesis
  },
  "difficulty": "1",
  "gasLimit": "30000000",
  "extradata": "0x0000...0000[SIGNER_ADDRESSES]0000...0000",
  "alloc": {
    // 100 pre-funded accounts with 10k ETH each
  }
}
```

**ExtraData Format (Clique-specific):**
- 32 bytes vanity prefix (64 hex chars): `0000...0000`
- Concatenated signer addresses (280 hex chars): `[addr1][addr2]...[addr7]` (no 0x prefix)
- 65 bytes seal suffix (130 hex chars): `0000...0000`
- **Total:** 474 hex characters

**Current signer addresses (in genesis extraData order):**
```
1. 0xc6bcfea104dbc7cb26cc990e33e58febd62fa533
2. 0x529166fe40a2912d369881d8d5c4e63f75e9e1f8
3. 0x52688b264bd5ffa39a3df49d4bcbe91cc561d9fb
4. 0x5f67d2acab6115fd6e337a852ffa7f04e9adb7ed
5. 0x44ec5abdb6e43c42da49d960c31af5d550f413c9
6. 0x0b18dc66d9837e7eed25b6a8cffbe16a82ea1599
7. 0xeef5100728b68504f635150224766861a2888b84
```

---

## DOCKER COMPOSE CONFIGURATION

**File:** `~/eth-devstack/blockchain/docker-compose.yml`

**Key service configuration (Signer1 example):**
```yaml
services:
  signer1:
    image: ethereum/client-go:v1.13.5
    container_name: eth-signer1
    command:
      - --networkid=1337
      - --datadir=/data
      - --syncmode=full
      - --gcmode=archive
      - --http
      - --http.addr=0.0.0.0
      - --http.port=8545
      - --http.api=eth,net,web3,personal,admin,txpool,debug
      - --http.corsdomain=*
      - --ws
      - --ws.addr=0.0.0.0
      - --ws.port=8546
      - --ws.api=eth,net,web3
      - --ws.origins=*
      - --allow-insecure-unlock
      - --mine
      - --miner.etherbase=0xc6bcfea104dbc7cb26cc990e33e58febd62fa533
      - --unlock=0xc6bcfea104dbc7cb26cc990e33e58febd62fa533
      - --password=/data/password.txt
      - --nodiscover
      - --maxpeers=25
      - --netrestrict=172.25.0.0/16
      - --miner.recommit=500ms
    ports:
      - "8545:8545"  # HTTP RPC
      - "30303:30303"  # P2P (not exposed externally)
    volumes:
      - ./data/signer1:/data
      - ./genesis.json:/genesis.json
    networks:
      eth-network:
        ipv4_address: 172.25.0.11
    mem_limit: 4g
    restart: unless-stopped
```

**Critical flags explained:**
- `--mine --miner.etherbase=<addr> --unlock=<addr>`: Enable mining as this signer
- `--http.api=...admin,txpool...`: Enable admin APIs for peer management
- `--nodiscover --netrestrict`: Private network (no public discovery)
- `--maxpeers=25`: Allow sufficient peer connections
- `--miner.recommit=500ms`: Fast block attempts (Geth enforces 1s minimum)
- `--allow-insecure-unlock`: Required for automated unlocking (dev only)

**Port mapping:**
```
Signer1: 8545 → 8545 (HTTP RPC)
Signer2: 8546 → 8545
Signer3: 8547 → 8545
...
Signer7: 8551 → 8545
```

**Network:**
```yaml
networks:
  eth-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.25.0.0/16
```

---

## TRANSACTION GENERATOR

**File:** `~/eth-devstack/blockchain/tx-generator/main.go`

**Architecture:**
- Connects to all 7 signers via HTTP RPC
- Round-robin distribution (tx N goes to signer N%7)
- Fresh nonce fetch per tx (no local caching)
- Random amounts: 0.001-1000 ETH
- 100 funded accounts (loaded from accounts.json)
- 2s batch intervals

**Key implementation details:**
```go
// Connect to all 7 signers
clients := make([]*ethclient.Client, 7)
for i := 0; i < 7; i++ {
    url := fmt.Sprintf("http://signer%d:8545", i+1)
    clients[i], _ = ethclient.Dial(url)
}

// Round-robin distribution
signerIndex := txNum % 7
client := clients[signerIndex]

// CRITICAL: Fetch fresh nonce every time
nonce, _ := client.PendingNonceAt(ctx, fromAddr)

// Build & sign tx
tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, nil)
signedTx, _ := types.SignTx(tx, signer, privateKey)

// Send to specific signer
client.SendTransaction(ctx, signedTx)
```

**Why no goroutines?**
- Original implementation: 40-120 goroutines/2s with 20-semaphore bottleneck
- Caused memory spikes, OOM kills, inconsistent throughput
- Simple loop: More predictable, easier to debug, sufficient throughput

---

## SETUP FROM SCRATCH

### Step 1: Project Structure
```bash
mkdir -p ~/eth-devstack/blockchain/{setup,tx-generator}
cd ~/eth-devstack/blockchain
```

### Step 2: Create Genesis File

Create `genesis.json` with the content provided above (or generate new one with your signer addresses).

**To generate new genesis with new signers:**
```bash
# Generate 7 new accounts:
for i in {1..7}; do
  geth account new --datadir ./data/signer$i --password <(echo "password")
done

# Extract addresses and build extraData:
# (Manual process - concat addresses without 0x prefix)
```

### Step 3: Create accounts.json

Create `setup/accounts.json` with 100 funded accounts (or copy existing).

### Step 4: Create docker-compose.yml

Copy the full docker-compose.yml (shown above). Ensure:
- All 7 signers have unique `miner.etherbase` addresses
- All ports mapped correctly
- Volume paths correct
- Network IPs sequential

### Step 5: Create Transaction Generator

**`tx-generator/Dockerfile`:**
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o tx-generator main.go
CMD ["./tx-generator"]
```

**`tx-generator/go.mod`:**
```go
module tx-generator

go 1.21

require github.com/ethereum/go-ethereum v1.13.5
```

**`tx-generator/main.go`:**
(Copy your current working implementation)

### Step 6: Initialize Signers
```bash
# Start only signers (not tx-generator yet):
docker-compose up -d signer1 signer2 signer3 signer4 signer5 signer6 signer7

# Wait for initialization (10 seconds):
sleep 10

# Verify all started:
docker ps --filter "name=eth-signer"
```

### Step 7: Create Peer Mesh

**`mesh-peers.sh`:**
```bash
#!/bin/bash
set -e

echo "=== Creating peer mesh for 7 signers ==="

# Get enodes from all signers
ENODES=()
for i in {1..7}; do
  ENODE=$(docker exec eth-signer$i geth attach --exec 'admin.nodeInfo.enode' /data/geth.ipc | tr -d '"')
  # Replace 127.0.0.1 with container name
  ENODE=$(echo $ENODE | sed "s/127.0.0.1/signer$i/g")
  ENODES[$i]=$ENODE
  echo "Signer$i: ${ENODES[$i]}"
done

echo ""
echo "=== Connecting peers ==="

# Connect each signer to all others (full mesh)
for i in {1..7}; do
  for j in {1..7}; do
    if [ $i -ne $j ]; then
      docker exec eth-signer$i geth attach --exec "admin.addPeer('${ENODES[$j]}')" /data/geth.ipc > /dev/null
    fi
  done
  echo "✓ Signer$i connected to 6 peers"
done

echo ""
echo "=== Verifying mesh ==="

for i in {1..7}; do
  PEER_COUNT=$(docker exec eth-signer$i geth attach --exec 'net.peerCount' /data/geth.ipc)
  echo "Signer$i: $PEER_COUNT peers"
done

echo ""
echo "=== Mesh complete ==="
```

Make executable and run:
```bash
chmod +x mesh-peers.sh
./mesh-peers.sh
```

**Expected output:**
```
✓ Signer1 connected to 6 peers
✓ Signer2 connected to 6 peers
...
Signer1: 6 peers
Signer2: 6 peers
...
```

### Step 8: Start Transaction Generator
```bash
docker-compose up -d eth-tx-generator

# Watch logs:
docker logs -f eth-tx-generator

# Expected output:
# 📦 Batch: 87 transactions
#    ✅ Sent 87/87
# 📦 Batch: 92 transactions
#    ✅ Sent 92/92
```

### Step 9: Verify Network Health

Run the health check script (see Testing section below).

---

## ESSENTIAL COMMANDS

### Start/Stop Network
```bash
# Start all services:
cd ~/eth-devstack/blockchain
docker-compose up -d

# After start, ALWAYS run:
./mesh-peers.sh

# Stop all:
docker-compose down

# Full clean restart (wipes data):
docker-compose down -v
docker-compose up -d
./mesh-peers.sh
```

### Monitor Services
```bash
# Check all containers:
docker ps --filter "name=eth"

# Check specific signer logs:
docker logs -f eth-signer1

# Check tx-generator logs:
docker logs -f eth-tx-generator

# Live block production:
watch -n 2 'docker exec eth-signer1 geth attach --exec "eth.blockNumber" /data/geth.ipc'
```

### Query Blockchain
```bash
# Current block number:
docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc

# Get latest block:
docker exec eth-signer1 geth attach --exec 'eth.getBlock("latest")' /data/geth.ipc

# Check account balance:
docker exec eth-signer1 geth attach --exec 'eth.getBalance("0xYOUR_ADDRESS")' /data/geth.ipc

# Txpool status:
docker exec eth-signer1 geth attach --exec 'txpool.status' /data/geth.ipc

# Peer count:
docker exec eth-signer1 geth attach --exec 'net.peerCount' /data/geth.ipc

# Authorized signers:
docker exec eth-signer1 geth attach --exec 'clique.getSigners()' /data/geth.ipc
```

### Send Manual Transaction
```bash
# Unlock account first:
docker exec eth-signer1 geth attach --exec 'personal.unlockAccount(eth.accounts[0], "password", 300)' /data/geth.ipc

# Send transaction:
docker exec eth-signer1 geth attach --exec 'eth.sendTransaction({from: eth.accounts[0], to: "0xRECIPIENT", value: web3.toWei(1, "ether")})' /data/geth.ipc
```

---

## TESTING & VERIFICATION

### Complete Health Check Script

Create `health-check.sh`:
```bash
#!/bin/bash
set -e

echo "=== ETHEREUM TESTNET HEALTH CHECK ==="
echo ""

# 1. Container Status
echo "=== CONTAINER STATUS ==="
docker ps --filter "name=eth-signer" --format "table {{.Names}}\t{{.Status}}"
echo ""

# 2. Block Sync
echo "=== BLOCK SYNC ==="
for i in {1..7}; do
  BLOCK=$(docker exec eth-signer$i geth attach --exec 'eth.blockNumber' /data/geth.ipc)
  echo "Signer$i: Block $BLOCK"
done
echo ""

# 3. Peer Connectivity
echo "=== PEER CONNECTIVITY ==="
for i in {1..7}; do
  PEERS=$(docker exec eth-signer$i geth attach --exec 'net.peerCount' /data/geth.ipc)
  echo "Signer$i: $PEERS peers (expected: 6)"
done
echo ""

# 4. Mining Status
echo "=== MINING STATUS ==="
for i in {1..7}; do
  MINING=$(docker exec eth-signer$i geth attach --exec 'eth.mining' /data/geth.ipc)
  echo "Signer$i: mining=$MINING"
done
echo ""

# 5. Txpool Status
echo "=== TXPOOL STATUS ==="
for i in {1..7}; do
  PENDING=$(docker exec eth-signer$i geth attach --exec 'txpool.status.pending' /data/geth.ipc)
  echo "Signer$i: $PENDING pending txs"
done
echo ""

# 6. Throughput Test (60 seconds)
echo "=== THROUGHPUT TEST (60 seconds) ==="
START_BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc)
echo "Start block: $START_BLOCK"
sleep 60
END_BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc)
echo "End block: $END_BLOCK"

BLOCKS=$((END_BLOCK - START_BLOCK))
echo "Blocks produced: $BLOCKS (expected: ~30)"

# Count transactions
TOTAL_TXS=0
for ((i=START_BLOCK+1; i<=END_BLOCK; i++)); do
  TX_COUNT=$(docker exec eth-signer1 geth attach --exec "eth.getBlock($i).transactions.length" /data/geth.ipc)
  TOTAL_TXS=$((TOTAL_TXS + TX_COUNT))
done

echo "Total transactions: $TOTAL_TXS"
echo "Average TPS: $((TOTAL_TXS / 60))"
echo "Average txs/block: $((TOTAL_TXS / BLOCKS))"
echo ""

# 7. Resource Usage
echo "=== RESOURCE USAGE ==="
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | grep signer
echo ""

echo "=== HEALTH CHECK COMPLETE ==="
echo ""
echo "Expected healthy values:"
echo "  - All signers on same block number"
echo "  - All signers have 6 peers"
echo "  - All signers mining=true"
echo "  - ~30 blocks in 60s (2s block time)"
echo "  - 30-60 TPS sustained"
echo "  - 60-90 txs/block average"
echo "  - <200MB RAM per signer"
```

Make executable:
```bash
chmod +x health-check.sh
./health-check.sh
```

### Stress Test (5 minutes)
```bash
START_BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc)
START_TIME=$(date +%s)

echo "Starting 5-minute stress test..."
sleep 300

END_BLOCK=$(docker exec eth-signer1 geth attach --exec 'eth.blockNumber' /data/geth.ipc)
END_TIME=$(date +%s)

BLOCKS=$((END_BLOCK - START_BLOCK))
DURATION=$((END_TIME - START_TIME))

TOTAL_TXS=0
for ((i=START_BLOCK+1; i<=END_BLOCK; i++)); do
  TX_COUNT=$(docker exec eth-signer1 geth attach --exec "eth.getBlock($i).transactions.length" /data/geth.ipc)
  TOTAL_TXS=$((TOTAL_TXS + TX_COUNT))
done

echo "Duration: ${DURATION}s"
echo "Blocks: $BLOCKS (expected: ~150)"
echo "Total txs: $TOTAL_TXS"
echo "Average TPS: $((TOTAL_TXS / DURATION))"
echo "Block time: $((DURATION / BLOCKS))s (expected: 2s)"
```

**Expected results:**
- ~150 blocks in 301s
- ~10,000-15,000 total transactions
- 35-50 TPS sustained
- 2.0s block time

---

## TROUBLESHOOTING

### Problem: Signers not producing blocks

**Symptoms:**
- Block number not increasing
- `eth.mining = false`

**Fix:**
```bash
# Check if signers are in genesis:
docker exec eth-signer1 geth attach --exec 'clique.getSigners()' /data/geth.ipc

# Check if signer addresses match genesis:
for i in {1..7}; do
  docker exec eth-signer$i geth attach --exec 'eth.coinbase' /data/geth.ipc
done

# If mismatch: regenerate genesis or fix signer keystores
```

### Problem: Low peer count

**Symptoms:**
- `net.peerCount < 6`

**Fix:**
```bash
# Re-run peer mesh script:
./mesh-peers.sh

# Verify connectivity:
for i in {1..7}; do
  docker exec eth-signer$i geth attach --exec 'admin.peers.length' /data/geth.ipc
done
```

### Problem: No transactions in blocks

**Symptoms:**
- Blocks empty or very few transactions
- Txpool empty on all signers

**Fix:**
```bash
# Check tx-generator status:
docker logs eth-tx-generator

# Restart tx-generator:
docker-compose restart eth-tx-generator

# Check if accounts have funds:
docker exec eth-signer1 geth attach --exec 'eth.getBalance("0xFIRST_ACCOUNT_FROM_JSON")' /data/geth.ipc
```

### Problem: Inconsistent block times

**Symptoms:**
- Some blocks take >3s
- Blocks come in bursts

**Fix:**
```bash
# Check for network forks (all should have same block hash):
for i in {1..7}; do
  LATEST=$(docker exec eth-signer$i geth attach --exec 'eth.blockNumber' /data/geth.ipc)
  HASH=$(docker exec eth-signer$i geth attach --exec "eth.getBlock($LATEST).hash" /data/geth.ipc)
  echo "Signer$i block $LATEST: $HASH"
done

# If different hashes: network forked, restart required
```

### Problem: High memory usage

**Symptoms:**
- Signers using >1GB RAM
- OOM kills

**Fix:**
```bash
# Check current usage:
docker stats --no-stream | grep signer

# Reduce txpool size in docker-compose.yml:
# Add flags: --txpool.globalslots=2048 --txpool.globalqueue=512

# Restart with new limits:
docker-compose down
docker-compose up -d
```

### Problem: Mesh peers script fails

**Symptoms:**
- Error: "enode not found"
- Connections fail

**Fix:**
```bash
# Ensure all signers are fully started (wait 30s after docker-compose up)
sleep 30

# Check if geth IPC is accessible:
docker exec eth-signer1 geth attach --exec 'admin.nodeInfo' /data/geth.ipc

# Manually add peers one by one:
ENODE=$(docker exec eth-signer2 geth attach --exec 'admin.nodeInfo.enode' /data/geth.ipc)
docker exec eth-signer1 geth attach --exec "admin.addPeer('$ENODE')" /data/geth.ipc
```

---

## PERFORMANCE METRICS

### Observed Performance (Production)

**Hardware:**
- CPU: Intel i5-12600K (16 threads)
- RAM: 7.6GB available
- Storage: SSD

**Network Performance:**
```
Block time:        2.0s ±0.1s (consistent)
Throughput:        40-120 txs/block
TPS (sustained):   35-60 TPS
Gas per block:     1.5-2.5M / 30M limit (5-8% utilization)
Resource/signer:   170-240MB RAM, 10-20% CPU
Blocks/5min:       150 blocks (exactly on target)
Latency:           <50ms RPC response time
```

### Scaling Characteristics

**Current bottleneck:** Transaction generator (single-threaded loop)

**Theoretical maximum:**
- Gas limit: 30M
- Gas per tx: 21,000 (simple transfer)
- Max txs/block: ~1,428
- Max TPS: ~714 (at 2s blocks)

**Practical limits:**
- Network propagation: ~500ms between peers
- RPC processing: ~50ms per tx submission
- Realistic max: ~200-300 TPS before propagation delays

**To increase throughput:**
1. Parallelize tx-generator (goroutines with proper rate limiting)
2. Increase gas limit (requires genesis change + full restart)
3. Add more signers (diminishing returns after 10-15)

---

## CLIQUE-SPECIFIC NOTES

### Why miner = 0x00 in Blocks?

**Clique doesn't use the `miner` field.** The actual signer is recovered from:
- Block's `extraData` field contains ECDSA signature (last 65 bytes)
- Signature is verified against chain's authorized signers
- Signer address is recovered via ECDSA from signature

**This is normal and expected.** Block explorers like Etherscan also show miner=0x00 for Clique chains.

### Block Distribution

With 7 signers, each seals ~14.3% of blocks (1/7). Clique enforces round-robin:
- Signer can't seal consecutive blocks
- Must wait for (N/2 + 1) other signers to seal before next turn
- Prevents single signer takeover

**To verify distribution:** Count unique addresses in block extraData signatures (requires ECDSA recovery, not trivial via RPC).

### Txpool Imbalance

Round-robin tx distribution causes temporary txpool imbalances:
1. Tx-generator sends tx to signer3
2. Signer3 gets tx immediately (0ms)
3. Other signers get tx via gossip (100-500ms)
4. If signer3 seals next block (within 2s), it includes the tx before others see it

**This is normal.** Txs don't get stuck; they're just mined by the recipient signer faster.

---

## INTEGRATION WITH INDEXER/EXPLORER

Your backend indexer can connect to any signer's RPC endpoint:
```
http://localhost:8545  (Signer1)
http://localhost:8546  (Signer2)
...
http://localhost:8551  (Signer7)
```

**All expose standard Ethereum JSON-RPC:**
- `eth_blockNumber` - Get current block
- `eth_getBlockByNumber` - Get block details
- `eth_getTransactionReceipt` - Get tx receipt
- `eth_getLogs` - Get event logs
- `eth_getBalance` - Get address balance
- etc.

**Example indexer connection (Go):**
```go
import "github.com/ethereum/go-ethereum/ethclient"

client, err := ethclient.Dial("http://localhost:8545")
if err != nil {
    log.Fatal(err)
}

// Start syncing from genesis:
for blockNum := 0; ; blockNum++ {
    block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
    if err != nil {
        // Caught up to chain head, wait for next block
        time.Sleep(2 * time.Second)
        continue
    }
    
    // Index block into database
    indexBlock(block)
}
```

**No changes needed to your existing indexer if it:**
- Supports standard Ethereum RPC
- Handles Clique blocks (miner=0x00 is fine)
- Works with ChainID 1337

---

## SECURITY NOTES (DEV ENVIRONMENT ONLY)

**This setup is FOR DEVELOPMENT ONLY and has intentional security weaknesses:**

1. **`--allow-insecure-unlock`** - Allows HTTP unlock (never use in production)
2. **`--http.corsdomain=*`** - Allows any origin (should restrict in production)
3. **Plain text passwords** - Keystores unlocked with `password.txt`
4. **No TLS** - All RPC over plain HTTP
5. **No authentication** - RPC endpoints open to anyone on host
6. **Private keys in volumes** - Accessible to anyone with Docker access

**For production:**
- Use hardware wallets or remote signers
- Enable TLS for RPC
- Implement authentication (JWT, API keys)
- Use firewall rules to restrict access
- Store keystores encrypted with HSM

---

## MAINTENANCE

### Backup Blockchain Data
```bash
# Stop network:
docker-compose down

# Backup all signer data:
tar -czf blockchain-backup-$(date +%Y%m%d).tar.gz ./data/

# Restart:
docker-compose up -d
./mesh-peers.sh
```

### Restore from Backup
```bash
docker-compose down -v  # Wipe existing data
tar -xzf blockchain-backup-YYYYMMDD.tar.gz
docker-compose up -d
./mesh-peers.sh
```

### Clean Restart (Full Wipe)
```bash
docker-compose down -v
rm -rf ./data/
docker-compose up -d
./mesh-peers.sh
```

**Note:** This resets to genesis. All transactions lost.

### Log Rotation
```bash
# Signers logs can grow large. Rotate with:
docker-compose logs --tail=1000 signer1 > signer1.log
docker-compose restart signer1

# Or configure Docker daemon log rotation in /etc/docker/daemon.json:
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

---

## NEXT STEPS

1. **Document current state** ✓ (this file)
2. **Phase 3:** Backend - PostgreSQL + Go indexer
3. **Phase 4:** REST API for block explorer
4. **Phase 5:** Next.js frontend

**Your indexer should work immediately with:**
```bash
RPC_URL=http://localhost:8545
CHAIN_ID=1337
```

---

## TROUBLESHOOTING CHECKLIST

Before asking for help, verify:

- [ ] All 7 signers running (`docker ps`)
- [ ] All signers on same block number
- [ ] All signers have 6 peers (`net.peerCount`)
- [ ] All signers mining (`eth.mining = true`)
- [ ] Tx-generator running and sending txs
- [ ] Blocks producing every 2 seconds
- [ ] Txpool has pending transactions
- [ ] No errors in logs (`docker logs eth-signer1`)

---

## REFERENCE

**Official Geth Documentation:**
- Clique: https://geth.ethereum.org/docs/fundamentals/private-network
- RPC API: https://geth.ethereum.org/docs/interacting-with-geth/rpc

**Docker Compose:**
- https://docs.docker.com/compose/

**go-ethereum:**
- https://github.com/ethereum/go-ethereum
- https://pkg.go.dev/github.com/ethereum/go-ethereum

---

## CHANGELOG

**2025-01-11:** Initial documentation - 7-validator Clique PoA network operational, 41 TPS sustained.

---

**Network Status:** ✅ PRODUCTION READY
**Maintainer:** [Your Name]
**Last Updated:** 2025-01-11