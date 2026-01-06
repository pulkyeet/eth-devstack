# Private Ethereum Testnet

Clique PoA network with 2 signer nodes.

## Quick Start
```bash
# Start network
docker compose up -d

# View logs
docker compose logs -f

# Stop network
docker compose down
```

## Network Details

- **Chain ID:** 1337
- **Block Time:** 15 seconds
- **Consensus:** Clique PoA
- **RPC:** http://localhost:8545
- **WebSocket:** ws://localhost:8546

## Pre-funded Accounts

See `accounts.txt` for full list.

Account 1:
- Address: `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
- Private Key: `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
- Balance: 100,000 ETH

Account 2:
- Address: `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
- Private Key: `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
- Balance: 100,000 ETH

## Testing
```bash
# Check current block
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Get balance
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"eth_getBalance",
    "params":["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "latest"],
    "id":1
  }'
```
