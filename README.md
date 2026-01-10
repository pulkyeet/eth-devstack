# Ethereum DevStack

A complete, vertically-integrated Ethereum development environment featuring a private testnet, block explorer backend with real-time indexing, and multi-chain support.

## 🎯 Project Overview

This is a production-architecture blockchain infrastructure built from scratch to demonstrate:
- **Blockchain Engineering**: Private Ethereum testnet, multi-chain indexing
- **Backend Systems**: Real-time data synchronization, REST APIs, SSE streaming
- **Database Design**: Complex PostgreSQL schema with proper indexing
- **System Architecture**: Microservices, event processing, multi-chain abstraction

**Tech Stack:** Go, PostgreSQL, Geth (go-ethereum), Docker, Fiber

---

## 🏗️ Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                     REST API (Fiber)                        │
│  /blocks  /transactions  /addresses  /stats  /stream/blocks │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                   Indexer Service (Go)                      │
│  • Multi-chain sync    • Log parsing (ERC20)                │
│  • Reorg handling      • Address tracking                   │
└───────────────────────────┬─────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
   ┌────┴─────┐      ┌──────┴──────┐     ┌─────┴──────┐
   │PostgreSQL│      │ Local Geth  │     │  Sepolia   │
   │  (Index) │      │  (Testnet)  │     │ (Optional) │
   └──────────┘      └─────────────┘     └────────────┘
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 16
- Docker & Docker Compose
- `golang-migrate` CLI

### 1. Clone & Setup
```bash
git clone https://github.com/pulkyeet/eth-devstack.git
cd eth-devstack
```

### 2. Start PostgreSQL
```bash
sudo systemctl start postgresql

# Create database and user
sudo -u postgres psql << EOF
CREATE USER eth_user WITH PASSWORD 'eth_pass_dev_only';
CREATE DATABASE ethereum_explorer OWNER eth_user;
GRANT ALL PRIVILEGES ON DATABASE ethereum_explorer TO eth_user;
EOF
```

### 3. Run Migrations
```bash
cd backend

# Install migrate if needed
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate

# Run migrations
make migrate-up
```

### 4. Start Local Testnet
```bash
cd blockchain
docker-compose up -d

# Verify blocks are being produced
docker logs blockchain-signer1-1 --tail 20
```

### 5. Start Backend Services
```bash
cd backend

# Terminal 1: Start Indexer
go run cmd/indexer/main.go

# Terminal 2: Start API
go run cmd/api/main.go
```

### 6. Test It
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get latest blocks
curl http://localhost:8080/api/v1/blocks?chain_id=1337&limit=5

# Network stats
curl http://localhost:8080/api/v1/stats?chain_id=1337

# Real-time block stream
curl -N http://localhost:8080/api/v1/stream/blocks?chain_id=1337
```

---

## 📊 Features

### ✅ Phase 1: Private Testnet
- 2-node Clique PoA network
- 15-second block time
- Pre-funded test accounts
- Persistent data storage

### ✅ Phase 2: Block Explorer Backend
- **Real-time indexing** of blocks, transactions, logs
- **Multi-chain support** (local, Sepolia, mainnet ready)
- **ERC20 token detection** and transfer tracking
- **Address activity tracking** with balance updates
- **RESTful API** with pagination and search
- **Server-Sent Events (SSE)** for real-time block streaming
- **Comprehensive stats** endpoint

### 🚧 Coming Soon (Phase 3-6)
- Next.js frontend (block explorer UI)
- HD wallet with transaction signing
- Smart contracts (ERC20, Staking, DEX)
- dApp frontends

---

## 📁 Repository Structure
```
eth-devstack/
├── blockchain/              # Local Geth testnet
│   ├── genesis.json
│   ├── docker-compose.yml
│   └── accounts.txt
├── backend/                 # Go services
│   ├── cmd/
│   │   ├── api/            # REST API server
│   │   └── indexer/        # Blockchain indexer
│   ├── internal/
│   │   ├── blockchain/     # Chain abstraction layer
│   │   ├── database/       # Data access layer
│   │   ├── indexer/        # Sync logic
│   │   ├── api/            # HTTP handlers
│   │   └── models/         # Data models
│   ├── Makefile
│   └── README.md
└── docs/
    └── MASTER_BLUEPRINT.md # Complete project spec
```

---

## 🗄️ Database Schema

**Core tables:**
- `chains` - Multi-chain configuration
- `blocks` - Indexed blockchain blocks
- `transactions` - Transaction history with receipts
- `transaction_logs` - Event logs (ERC20 transfers, etc.)
- `addresses` - Address metadata and activity
- `tokens` - ERC20/721/1155 token registry
- `token_transfers` - Token transfer events
- `token_balances` - Current token holdings

**Optimizations:**
- Composite indexes on (chain_id, block_number)
- Partial indexes for active chains
- Materialized views for stats (future)

---

## 🔌 API Endpoints

### Health & Chains
- `GET /api/v1/health` - Service health check
- `GET /api/v1/chains` - List supported chains

### Blocks
- `GET /api/v1/blocks` - Paginated block list
- `GET /api/v1/blocks/:id` - Block by number or hash

### Transactions
- `GET /api/v1/transactions` - Transaction list
- `GET /api/v1/transactions/:hash` - Transaction details

### Addresses
- `GET /api/v1/addresses/:address` - Address info
- `GET /api/v1/addresses/:address/transactions` - Address history
- `GET /api/v1/addresses/:address/tokens` - Token balances

### Stats & Search
- `GET /api/v1/stats?chain_id=1337` - Network statistics
- `GET /api/v1/search?q=<query>` - Universal search

### Real-time
- `GET /api/v1/stream/blocks` - SSE block stream

---

## 🧪 Testing
```bash
cd backend

# Unit tests
make test

# Integration tests (API must be running)
go test -v ./test/...

# Coverage report
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

---

## 🛠️ Development

### Makefile Commands
```bash
make help              # Show all commands
make build             # Build binaries
make run               # Start API server
make run-indexer       # Start indexer
make test              # Run unit tests
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make clean             # Clean build artifacts
```

### Configuration

Edit `backend/.env`:
```bash
API_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ethereum_explorer
DB_USER=eth_user
DB_PASSWORD=eth_pass_dev_only
```

### Adding New Chains

Edit `backend/internal/config/chains.json`:
```json
{
  "chain_id": 11155111,
  "name": "Ethereum Sepolia",
  "rpc_endpoint": "https://rpc.sepolia.org",
  "block_time_seconds": 12,
  "is_active": true
}
```

Restart indexer to begin syncing.

---

## 🎓 Technical Highlights

### Multi-Chain Architecture
- Chain abstraction layer with RPC failover
- Parallel sync workers (configurable)
- Per-chain sync status tracking

### Real-Time Processing
- SSE streaming for live block updates
- Event log parsing (Transfer, Approval, etc.)
- Address balance tracking

### Production Patterns
- Structured logging (Zap)
- Database connection pooling
- Graceful shutdown
- CORS support
- Error handling middleware

### Optimizations
- Batch block processing (100 blocks/batch)
- Composite database indexes
- Minimal RPC calls (receipt caching)

---

## 📈 Performance

**Current metrics** (local testnet):
- **Indexer sync speed:** ~50 blocks/second
- **API response time:** <100ms for simple queries
- **Database size:** ~50MB per 10,000 blocks
- **Memory usage:** ~100MB (indexer), ~50MB (API)

---

## 🔒 Security Notes

- This is a **development environment**
- Never use in production without security hardening
- Private keys in Phase 4 wallet are **educational only**
- Always use hardware wallets for real funds

---

## 📚 Documentation

- [Master Blueprint](docs/MASTER_BLUEPRINT.md) - Complete project specification
- [Startup Guide](START.md) - Daily workflow

---

## 🤝 Contributing

This is a portfolio project, but suggestions welcome via issues!

---

## 📝 License

MIT

---

## 🎯 Roadmap

- [x] Phase 1: Private Testnet
- [x] Phase 2: Block Explorer Backend
- [ ] Phase 3: Block Explorer Frontend
- [ ] Phase 4: Wallet Backend (Educational)
- [ ] Phase 5: Wallet Frontend
- [ ] Phase 6: Smart Contracts (ERC20, Staking, DEX)
- [ ] Phase 7: dApp Frontends
- [ ] Phase 8: Testing, CI/CD, Documentation

---

## 📧 Contact

**Pulkit** - [GitHub](https://github.com/pulkyeet)

**Project:** [github.com/pulkyeet/eth-devstack](https://github.com/pulkyeet/eth-devstack)