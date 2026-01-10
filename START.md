# Complete Startup Guide - Every Time You Boot

Save this as `~/eth-devstack/START.md`

---

## 1. Start PostgreSQL

```bash
# Check status
sudo systemctl status postgresql

# If not running, start it
sudo systemctl start postgresql

# Verify it's up
psql -U eth_user -d ethereum_explorer -h localhost -c "SELECT COUNT(*) FROM blocks WHERE chain_id=1337;"
# Password: eth_pass_dev_only
```

**Expected:** Should show block count (e.g., 932+)

---

## 2. Start Geth Testnet (Blockchain)

```bash
cd ~/eth-devstack/blockchain

# Start both signer nodes
docker-compose up -d

# Verify they're running
docker ps | grep signer

# Check logs - should see "Successfully sealed new block"
docker logs eth-signer1 --tail 20

# Test RPC is accessible
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

**Expected:** 
- 2 containers running (signer1, signer2)
- Logs show blocks being sealed every ~15 seconds
- RPC returns block number in hex (e.g., `{"jsonrpc":"2.0","id":1,"result":"0x3a5"}`)

---

## 3. Start Backend Indexer

```bash
# Open new terminal (Terminal 1)
cd ~/eth-devstack/backend

# Start indexer
go run cmd/indexer/main.go
```

**Expected output:**
```
INFO Starting indexer service
INFO Starting chain indexer chain_id=1337
INFO Syncing blocks chain_id=1337 from=932 to=1032
```

**Leave this terminal open.** Indexer will continuously sync new blocks.

**To verify it's working:**
```bash
# In another terminal, watch blocks increase
watch -n 5 'psql -U eth_user -d ethereum_explorer -h localhost -t -c "SELECT MAX(block_number) FROM blocks WHERE chain_id=1337;"'
```

---

## 4. Start Backend API

```bash
# Open new terminal (Terminal 2)
cd ~/eth-devstack/backend

# Start API server
go run cmd/api/main.go
```

**Expected output:**
```
INFO Starting Ethereum Explorer API port=8080
INFO Database connected
INFO API Server started successfully
INFO Starting API server port=8080
```

**Leave this terminal open.** API will serve requests on port 8080.

**To verify it's working:**
```bash
# In another terminal
curl http://localhost:8080/api/v1/health
```

Should return: `{"success":true,"data":{"status":"healthy","database":"connected"},...}`

---

## 5. Quick Health Check (All Services)

```bash
# Check all services with one script
cat > ~/eth-devstack/check-services.sh << 'EOF'
#!/bin/bash

echo "=== Service Health Check ==="
echo ""

# PostgreSQL
echo -n "PostgreSQL: "
sudo systemctl is-active postgresql

# Geth
echo -n "Geth Signers: "
docker ps --filter "name=signer" --format "{{.Names}}" | wc -l | awk '{print $1 " containers"}'

# Indexer
echo -n "Indexer: "
pgrep -f "cmd/indexer/main.go" > /dev/null && echo "Running" || echo "NOT RUNNING"

# API
echo -n "API: "
pgrep -f "cmd/api/main.go" > /dev/null && echo "Running" || echo "NOT RUNNING"

# Test API endpoint
echo -n "API Health: "
curl -s http://localhost:8080/api/v1/health | grep -q "healthy" && echo "OK" || echo "FAILED"

# Latest block
echo -n "Latest Block in DB: "
psql -U eth_user -d ethereum_explorer -h localhost -t -c "SELECT MAX(block_number) FROM blocks WHERE chain_id=1337;" 2>/dev/null | xargs

echo ""
echo "=== All checks complete ==="
EOF

chmod +x ~/eth-devstack/check-services.sh

# Run it
~/eth-devstack/check-services.sh
```

**Expected output:**
```
=== Service Health Check ===

PostgreSQL: active
Geth Signers: 2 containers
Indexer: Running
API: Running
API Health: OK
Latest Block in DB: 1045

=== All checks complete ===
```

---

## 6. Test All Endpoints

```bash
# Health
curl http://localhost:8080/api/v1/health

# Chains
curl http://localhost:8080/api/v1/chains

# Latest blocks
curl http://localhost:8080/api/v1/blocks?chain_id=1337&limit=3

# Single block
curl http://localhost:8080/api/v1/blocks/1?chain_id=1337

# Transactions
curl http://localhost:8080/api/v1/transactions?chain_id=1337&limit=3

# Stats
curl http://localhost:8080/api/v1/stats?chain_id=1337

# Search
curl "http://localhost:8080/api/v1/search?q=1&chain_id=1337"

# SSE Stream (leave running for 30 seconds to see blocks)
curl -N http://localhost:8080/api/v1/stream/blocks?chain_id=1337
```

---

## 7. Common Issues & Fixes

### Issue: "Connection refused" on PostgreSQL

```bash
# Start PostgreSQL
sudo systemctl start postgresql

# Check status
sudo systemctl status postgresql
```

### Issue: "Connection refused" on API (port 8080)

```bash
# Check if API is running
pgrep -f "cmd/api/main.go"

# If not, start it
cd ~/eth-devstack/backend
go run cmd/api/main.go
```

### Issue: Indexer not syncing new blocks

```bash
# Check if Geth is producing blocks
docker logs blockchain-signer1-1 --tail 20

# Restart indexer
pkill -f "cmd/indexer/main.go"
cd ~/eth-devstack/backend
go run cmd/indexer/main.go
```

### Issue: Docker containers not running

```bash
cd ~/eth-devstack/blockchain

# Stop and restart
docker-compose down
docker-compose up -d

# Check logs
docker logs blockchain-signer1-1 --tail 50
```

### Issue: Too many Go processes running

```bash
# Kill all
pkill -f "cmd/api/main.go"
pkill -f "cmd/indexer/main.go"

# Restart clean
cd ~/eth-devstack/backend
go run cmd/indexer/main.go &
go run cmd/api/main.go &
```

---

## 8. Stop Everything (Shutdown)

```bash
# Stop Go processes (Ctrl+C in their terminals, or:)
pkill -f "cmd/api/main.go"
pkill -f "cmd/indexer/main.go"

# Stop Docker containers
cd ~/eth-devstack/blockchain
docker-compose down

# Optionally stop PostgreSQL (not recommended, usually leave running)
# sudo systemctl stop postgresql
```

---

## 9. Daily Workflow Summary

**Morning startup (5 commands):**

```bash
# 1. Start PostgreSQL (if not running)
sudo systemctl start postgresql

# 2. Start Geth
cd ~/eth-devstack/blockchain && docker-compose up -d

# 3. Start Indexer (Terminal 1)
cd ~/eth-devstack/backend && go run cmd/indexer/main.go

# 4. Start API (Terminal 2)
cd ~/eth-devstack/backend && go run cmd/api/main.go

# 5. Check all services
~/eth-devstack/check-services.sh
```

**That's it. You're ready to work.**

---

## 10. Port Reference

| Service | Port | URL |
|---------|------|-----|
| PostgreSQL | 5432 | localhost:5432 |
| Geth RPC | 8545 | http://localhost:8545 |
| Geth WS | 8546 | ws://localhost:8546 |
| API | 8080 | http://localhost:8080/api/v1 |

---

## 11. Useful Aliases (Optional)

Add to `~/.bashrc`:

```bash
# Eth DevStack aliases
alias eth-start-db='sudo systemctl start postgresql'
alias eth-start-geth='cd ~/eth-devstack/blockchain && docker-compose up -d'
alias eth-start-indexer='cd ~/eth-devstack/backend && go run cmd/indexer/main.go'
alias eth-start-api='cd ~/eth-devstack/backend && go run cmd/api/main.go'
alias eth-check='~/eth-devstack/check-services.sh'
alias eth-stop='pkill -f "cmd/api/main.go"; pkill -f "cmd/indexer/main.go"; cd ~/eth-devstack/blockchain && docker-compose down'
```

Then: `source ~/.bashrc`

**Now you can just type:** `eth-check` to verify everything!

---

Save this file and you'll never need to ask again. 🚀