# MASTER BLUEPRINT: Vertically Integrated Ethereum Stack (REVISED)

**Version 2.0 - Critical Security & Quality Focus**

---

## OVERVIEW & PHILOSOPHY

### Project Goal
Build a production-architecture Ethereum development environment demonstrating:
- **Smart contract expertise** (50% of effort)
- **Full-stack blockchain engineering** (30%)
- **System design & architecture** (20%)

### Core Principles
1. **Security First** - Document security considerations, even for dev features
2. **Test Everything** - >80% coverage on backend, >95% on contracts
3. **Quality Over Speed** - No timeline pressure, build it right
4. **Production Patterns** - Use real-world patterns, document trade-offs
5. **Multi-Chain from Day 1** - Config-driven, works with any EVM chain

### What This Proves to Employers
- Deep Solidity knowledge (gas optimization, security, testing)
- Backend systems (indexing, APIs, databases)
- Frontend integration (Web3, signing, state management)
- DevOps basics (Docker, CI/CD)
- System design thinking (documented decisions)

---

## I. FINAL ARCHITECTURE

```
┌──────────────────────────────────────────────────────────────┐
│                    FRONTEND (Single Next.js App)              │
│                                                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  /explorer  │  │   /wallet   │  │   /dapps    │          │
│  │  (Route)    │  │   (Route)   │  │   (Route)   │          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
│         └─────────────────┴─────────────────┘                │
│                           │                                   │
│                  Shared Components                            │
│         (NetworkSwitcher, ConnectButton, etc.)                │
└───────────────────────────┬──────────────────────────────────┘
                            │
┌───────────────────────────┴──────────────────────────────────┐
│                    APPLICATION LAYER                          │
│                                                                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │            API Gateway (Go/Fiber)                     │   │
│  │  - CORS  - Logging  - Error Handling  - Health       │   │
│  └────────┬──────────────────────────────┬──────────────┘   │
│           │                               │                   │
│  ┌────────┴──────────┐         ┌────────┴──────────┐        │
│  │  Explorer Service │         │  Wallet Service    │        │
│  │  - Block queries  │         │  - Tx signing      │        │
│  │  - Tx queries     │         │  - Balance track   │        │
│  │  - Search         │         │  - Multi-chain     │        │
│  │  - Real-time SSE  │         │  DEV ONLY ⚠️       │        │
│  └────────┬──────────┘         └────────┬───────────┘        │
└───────────┼─────────────────────────────┼────────────────────┘
            │                              │
┌───────────┴──────────────────────────────┴────────────────────┐
│                        DATA LAYER                              │
│                                                                 │
│  ┌──────────────────┐        ┌────────────────────────────┐   │
│  │   PostgreSQL     │        │   Redis (Later)            │   │
│  │   - Blocks       │        │   - Cache                  │   │
│  │   - Transactions │        │   - Rate limits            │   │
│  │   - Addresses    │        └────────────────────────────┘   │
│  │   - Tokens       │                                          │
│  │   - Wallet keys  │        ⚠️ SECURITY WARNING:              │
│  │     (encrypted)  │        Storing encrypted keys           │
│  └──────────────────┘        server-side is DEV ONLY          │
│                              Production: client-side signing   │
└────────────────────────┬───────────────────────────────────────┘
                         │
┌────────────────────────┴───────────────────────────────────────┐
│                    BLOCKCHAIN LAYER                             │
│                                                                 │
│  ┌───────────────────────────────────────────────────────┐    │
│  │              Indexer Service (Go)                      │    │
│  │  - Multi-chain sync (parallel workers)                │    │
│  │  - Block/tx/log parsing                               │    │
│  │  - Reorg handling                                     │    │
│  │  - Tested (unit + integration)                        │    │
│  └────────────┬──────────────────────────────────────────┘    │
│               │                                                │
│  ┌────────────┴──────────────────────────────────────────┐    │
│  │         Chain Abstraction Layer                        │    │
│  │  - ChainConfig (from JSON)                            │    │
│  │  - RPC client pool with failover                      │    │
│  │  - Gas price oracles (EIP-1559 / Legacy / Fixed)     │    │
│  └────────────┬──────────────────────────────────────────┘    │
│               │                                                │
│  ┌────────────┴────────────┬──────────────┬─────────────┐     │
│  │ Local Testnet           │ Ethereum     │ Polygon     │     │
│  │ (Geth - Clique PoA)     │ Sepolia      │ Mumbai      │     │
│  │ - 2 signer nodes        │ (Testnet)    │ (Testnet)   │     │
│  └─────────────────────────┴──────────────┴─────────────┘     │
└─────────────────────────────────────────────────────────────────┘

                    ┌──────────────────────┐
                    │  Smart Contracts     │
                    │  (Solidity/Foundry)  │
                    │  - Advanced ERC20    │
                    │  - Staking (complex) │
                    │  - DEX (AMM)         │
                    │                      │
                    │  >95% test coverage  │
                    │  Gas optimized       │
                    │  Natspec documented  │
                    └──────────────────────┘
```

---

## II. TECH STACK (FINALIZED)

### Backend
```
Language:        Go 1.21+
Framework:       Fiber v2 (fast, Express-like)
Ethereum:        go-ethereum v1.13+
Database:        PostgreSQL 16
Cache:           Redis 7 (Phase 8+)
Testing:         testify, httptest, table-driven tests
Logging:         zap (structured, leveled)
Config:          viper (env + file)
Migrations:      golang-migrate
```

### Frontend (Single App)
```
Framework:       Next.js 14 (App Router)
Language:        TypeScript (strict mode)
Web3:            viem + wagmi (v2)
UI:              TailwindCSS + shadcn/ui
State:           Zustand (lightweight)
Data Fetching:   TanStack Query (React Query v5)
Charts:          recharts
Forms:           react-hook-form + zod
Testing:         Vitest (optional, document if skipped)
```

### Smart Contracts
```
Language:        Solidity ^0.8.20
Framework:       Foundry (forge, cast, anvil)
Testing:         Forge (Solidity tests) + Fuzzing
Linting:         solhint
Formatting:      forge fmt
Security:        slither, aderyn
Gas Reports:     forge test --gas-report
```

### Infrastructure
```
Containers:      Docker + Docker Compose
Nodes:           Geth official images
Reverse Proxy:   Nginx (Phase 8+)
Monitoring:      Prometheus + Grafana (Phase 8+)
CI/CD:           GitHub Actions
Deployment:      Vercel (frontend), VPS/Railway (backend - optional)
```

---

## III. REPOSITORY STRUCTURE (MONOREPO)

```
ethereum-project/                    # Root monorepo
├── README.md                        # Project overview
├── docker-compose.yml               # All services
├── .gitignore
├── .github/
│   └── workflows/
│       ├── backend-tests.yml        # Run Go tests on PR
│       ├── contract-tests.yml       # Run Forge tests on PR
│       └── frontend-build.yml       # Build Next.js on PR
│
├── docs/
│   ├── ARCHITECTURE.md              # System design
│   ├── DECISIONS.md                 # Architecture Decision Records
│   ├── API.md                       # API documentation
│   ├── SECURITY.md                  # Security considerations
│   └── DEPLOYMENT.md                # Deployment guide
│
├── blockchain/                      # Local testnet
│   ├── README.md
│   ├── genesis.json
│   ├── docker-compose.yml
│   ├── accounts.txt                 # Pre-funded accounts
│   ├── signer1/                     # Node 1 data
│   ├── signer2/                     # Node 2 data
│   └── scripts/
│       └── faucet.go                # Fund test accounts
│
├── backend/                         # Go services
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile                     # Helper commands
│   ├── .env.example
│   │
│   ├── cmd/                         # Binaries
│   │   ├── api/
│   │   │   └── main.go
│   │   ├── indexer/
│   │   │   └── main.go
│   │   └── migrate/                 # DB migrations CLI
│   │       └── main.go
│   │
│   ├── internal/                    # Private packages
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   ├── config_test.go      # Test everything
│   │   │   ├── chains.go
│   │   │   └── chains.json
│   │   │
│   │   ├── blockchain/              # Chain abstraction
│   │   │   ├── client.go
│   │   │   ├── client_test.go
│   │   │   ├── manager.go
│   │   │   ├── manager_test.go
│   │   │   ├── gas_oracle.go
│   │   │   └── gas_oracle_test.go
│   │   │
│   │   ├── indexer/                 # Blockchain sync
│   │   │   ├── service.go
│   │   │   ├── service_test.go
│   │   │   ├── block_processor.go
│   │   │   ├── block_processor_test.go
│   │   │   ├── tx_processor.go
│   │   │   ├── tx_processor_test.go
│   │   │   ├── reorg_handler.go
│   │   │   └── reorg_handler_test.go
│   │   │
│   │   ├── database/                # Data access
│   │   │   ├── postgres.go
│   │   │   ├── postgres_test.go
│   │   │   ├── migrations/
│   │   │   │   ├── 001_initial.up.sql
│   │   │   │   └── 001_initial.down.sql
│   │   │   └── queries/            # SQL queries
│   │   │       ├── blocks.sql
│   │   │       └── transactions.sql
│   │   │
│   │   ├── models/                  # Data models
│   │   │   ├── block.go
│   │   │   ├── transaction.go
│   │   │   ├── address.go
│   │   │   └── token.go
│   │   │
│   │   ├── api/                     # HTTP layer
│   │   │   ├── server.go
│   │   │   ├── server_test.go      # Integration tests
│   │   │   ├── handlers/
│   │   │   │   ├── blocks.go
│   │   │   │   ├── blocks_test.go
│   │   │   │   ├── transactions.go
│   │   │   │   ├── transactions_test.go
│   │   │   │   ├── addresses.go
│   │   │   │   └── search.go
│   │   │   ├── middleware/
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── recovery.go     # Panic recovery
│   │   │   │   └── errors.go       # Error handling
│   │   │   └── routes.go
│   │   │
│   │   ├── wallet/                  # ⚠️ DEV ONLY
│   │   │   ├── README.md           # Security warnings
│   │   │   ├── service.go
│   │   │   ├── service_test.go
│   │   │   ├── hd_wallet.go
│   │   │   ├── hd_wallet_test.go
│   │   │   ├── keystore.go         # Encryption
│   │   │   ├── keystore_test.go
│   │   │   ├── signer.go
│   │   │   ├── signer_test.go
│   │   │   └── nonce_manager.go
│   │   │
│   │   └── utils/
│   │       ├── logger.go            # Zap setup
│   │       ├── errors.go            # Custom errors
│   │       └── helpers.go
│   │
│   └── test/                        # Integration tests
│       ├── api_test.go              # Full API tests
│       └── testdata/
│
├── frontend/                        # Single Next.js app
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.js
│   ├── tailwind.config.ts
│   ├── .env.local.example
│   │
│   ├── public/
│   │   └── icons/                   # Chain icons
│   │
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx           # Root layout
│   │   │   ├── page.tsx             # Home/landing
│   │   │   │
│   │   │   ├── explorer/            # Route group
│   │   │   │   ├── layout.tsx
│   │   │   │   ├── page.tsx         # Explorer home
│   │   │   │   ├── blocks/
│   │   │   │   │   ├── page.tsx     # Block list
│   │   │   │   │   └── [id]/
│   │   │   │   │       └── page.tsx # Block detail
│   │   │   │   ├── tx/
│   │   │   │   │   └── [hash]/
│   │   │   │   │       └── page.tsx
│   │   │   │   └── address/
│   │   │   │       └── [address]/
│   │   │   │           └── page.tsx
│   │   │   │
│   │   │   ├── wallet/              # Route group
│   │   │   │   ├── layout.tsx
│   │   │   │   ├── page.tsx         # Wallet home
│   │   │   │   ├── create/
│   │   │   │   ├── import/
│   │   │   │   ├── send/
│   │   │   │   └── dashboard/
│   │   │   │
│   │   │   └── dapps/               # Route group
│   │   │       ├── layout.tsx
│   │   │       ├── page.tsx         # dApps home
│   │   │       ├── token/
│   │   │       ├── staking/
│   │   │       └── dex/
│   │   │
│   │   ├── components/              # Shared components
│   │   │   ├── layout/
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   └── Footer.tsx
│   │   │   ├── web3/
│   │   │   │   ├── ConnectButton.tsx
│   │   │   │   ├── NetworkSwitcher.tsx
│   │   │   │   └── AddressDisplay.tsx
│   │   │   ├── ui/                  # shadcn components
│   │   │   └── common/
│   │   │       ├── LoadingSpinner.tsx
│   │   │       └── ErrorBoundary.tsx
│   │   │
│   │   ├── lib/
│   │   │   ├── wagmi.ts             # wagmi config
│   │   │   ├── api.ts               # API client
│   │   │   ├── chains.ts            # Chain configs
│   │   │   ├── contracts.ts         # Contract ABIs
│   │   │   └── utils.ts
│   │   │
│   │   ├── hooks/
│   │   │   ├── useChain.ts
│   │   │   ├── useContracts.ts
│   │   │   └── useBalance.ts
│   │   │
│   │   ├── context/
│   │   │   └── ChainContext.tsx
│   │   │
│   │   └── types/
│   │       ├── api.ts
│   │       └── contracts.ts
│   │
│   └── tests/                       # Optional
│       └── setup.ts
│
├── contracts/                       # Solidity
│   ├── foundry.toml
│   ├── remappings.txt
│   ├── .env.example
│   │
│   ├── src/
│   │   ├── Token.sol                # Advanced ERC20
│   │   ├── Staking.sol              # Complex staking
│   │   └── DEX.sol                  # AMM with LP tokens
│   │
│   ├── test/
│   │   ├── Token.t.sol              # Unit tests
│   │   ├── TokenFuzz.t.sol          # Fuzz tests
│   │   ├── Staking.t.sol
│   │   ├── StakingFuzz.t.sol
│   │   ├── DEX.t.sol
│   │   └── DEXFuzz.t.sol
│   │
│   ├── script/
│   │   ├── Deploy.s.sol             # Deployment
│   │   └── Interact.s.sol           # Interaction scripts
│   │
│   ├── lib/                         # Dependencies
│   │   └── openzeppelin-contracts/
│   │
│   └── gas-reports/                 # Gas optimization docs
│       └── optimizations.md
│
└── monitoring/                      # Phase 8+
    ├── prometheus/
    │   └── prometheus.yml
    └── grafana/
        └── dashboards/
```

---

## IV. DATABASE SCHEMA (COMPLETE & INDEXED)

```sql
-- ============================================================================
-- EXTENSIONS
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- CHAIN METADATA
-- ============================================================================

CREATE TABLE chains (
    chain_id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    short_name VARCHAR(20) NOT NULL UNIQUE,
    native_symbol VARCHAR(10) NOT NULL,
    rpc_endpoint TEXT NOT NULL,
    ws_endpoint TEXT,
    block_time_seconds INT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    is_testnet BOOLEAN DEFAULT false,
    last_indexed_block BIGINT DEFAULT 0,
    icon_url TEXT,
    explorer_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_chains_active ON chains(is_active) WHERE is_active = true;

-- ============================================================================
-- BLOCKCHAIN DATA
-- ============================================================================

CREATE TABLE blocks (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    block_number BIGINT NOT NULL,
    hash VARCHAR(66) NOT NULL,
    parent_hash VARCHAR(66) NOT NULL,
    nonce VARCHAR(18),
    sha3_uncles VARCHAR(66),
    miner VARCHAR(42) NOT NULL,
    state_root VARCHAR(66),
    transactions_root VARCHAR(66),
    receipts_root VARCHAR(66),
    difficulty NUMERIC(78, 0),
    total_difficulty NUMERIC(78, 0),
    size BIGINT,
    gas_limit BIGINT NOT NULL,
    gas_used BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    extra_data TEXT,
    mix_hash VARCHAR(66),
    base_fee_per_gas NUMERIC(78, 0), -- EIP-1559
    tx_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, block_number),
    UNIQUE(chain_id, hash)
);

CREATE INDEX idx_blocks_chain_number ON blocks(chain_id, block_number DESC);
CREATE INDEX idx_blocks_chain_timestamp ON blocks(chain_id, timestamp DESC);
CREATE INDEX idx_blocks_chain_miner ON blocks(chain_id, miner);
CREATE INDEX idx_blocks_hash ON blocks(hash); -- For quick hash lookups

-- ============================================================================

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    hash VARCHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    block_hash VARCHAR(66) NOT NULL,
    transaction_index INT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42), -- NULL for contract creation
    value NUMERIC(78, 0) NOT NULL DEFAULT 0,
    gas BIGINT NOT NULL,
    gas_price NUMERIC(78, 0), -- Legacy
    max_fee_per_gas NUMERIC(78, 0), -- EIP-1559
    max_priority_fee_per_gas NUMERIC(78, 0), -- EIP-1559
    input TEXT,
    nonce BIGINT NOT NULL,
    transaction_type INT DEFAULT 0, -- 0: legacy, 1: EIP-2930, 2: EIP-1559
    
    -- Receipt data
    status INT, -- 1: success, 0: failed, NULL: pending
    gas_used BIGINT,
    cumulative_gas_used BIGINT,
    effective_gas_price NUMERIC(78, 0),
    contract_address VARCHAR(42), -- If contract creation
    logs_bloom TEXT,
    
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, hash)
);

CREATE INDEX idx_tx_chain_block ON transactions(chain_id, block_number DESC);
CREATE INDEX idx_tx_chain_from ON transactions(chain_id, from_address, timestamp DESC);
CREATE INDEX idx_tx_chain_to ON transactions(chain_id, to_address, timestamp DESC) WHERE to_address IS NOT NULL;
CREATE INDEX idx_tx_hash ON transactions(hash);
CREATE INDEX idx_tx_chain_timestamp ON transactions(chain_id, timestamp DESC);
CREATE INDEX idx_tx_status ON transactions(chain_id, status) WHERE status IS NOT NULL;

-- ============================================================================

CREATE TABLE transaction_logs (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    transaction_hash VARCHAR(66) NOT NULL,
    log_index INT NOT NULL,
    address VARCHAR(42) NOT NULL, -- Contract that emitted
    data TEXT,
    topic0 VARCHAR(66), -- Event signature
    topic1 VARCHAR(66),
    topic2 VARCHAR(66),
    topic3 VARCHAR(66),
    block_number BIGINT NOT NULL,
    block_hash VARCHAR(66) NOT NULL,
    transaction_index INT NOT NULL,
    removed BOOLEAN DEFAULT false, -- For reorgs
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, transaction_hash, log_index)
);

CREATE INDEX idx_logs_chain_address ON transaction_logs(chain_id, address);
CREATE INDEX idx_logs_chain_topic0 ON transaction_logs(chain_id, topic0);
CREATE INDEX idx_logs_chain_block ON transaction_logs(chain_id, block_number DESC);
CREATE INDEX idx_logs_topics ON transaction_logs(chain_id, topic0, topic1) WHERE topic1 IS NOT NULL;

-- ============================================================================
-- ADDRESS DATA
-- ============================================================================

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    address VARCHAR(42) NOT NULL,
    balance NUMERIC(78, 0) DEFAULT 0,
    nonce BIGINT DEFAULT 0,
    is_contract BOOLEAN DEFAULT false,
    contract_creator VARCHAR(42),
    creation_tx_hash VARCHAR(66),
    code_hash VARCHAR(66),
    tx_count BIGINT DEFAULT 0,
    first_seen_block BIGINT,
    last_seen_block BIGINT,
    first_seen_at TIMESTAMP,
    last_seen_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, address)
);

CREATE INDEX idx_addresses_chain ON addresses(chain_id, address);
CREATE INDEX idx_addresses_chain_balance ON addresses(chain_id, balance DESC);
CREATE INDEX idx_addresses_contract ON addresses(chain_id, is_contract) WHERE is_contract = true;
CREATE INDEX idx_addresses_updated ON addresses(updated_at DESC);

-- ============================================================================
-- TOKEN DATA (ERC20/ERC721/ERC1155)
-- ============================================================================

CREATE TABLE tokens (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    address VARCHAR(42) NOT NULL,
    type VARCHAR(20) NOT NULL, -- ERC20, ERC721, ERC1155
    name VARCHAR(255),
    symbol VARCHAR(50),
    decimals INT, -- ERC20
    total_supply NUMERIC(78, 0), -- ERC20
    holder_count BIGINT DEFAULT 0,
    transfer_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, address),
    CHECK (type IN ('ERC20', 'ERC721', 'ERC1155'))
);

CREATE INDEX idx_tokens_chain ON tokens(chain_id);
CREATE INDEX idx_tokens_chain_type ON tokens(chain_id, type);
CREATE INDEX idx_tokens_symbol ON tokens(chain_id, symbol);

-- ============================================================================

CREATE TABLE token_transfers (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    transaction_hash VARCHAR(66) NOT NULL,
    log_index INT NOT NULL,
    token_address VARCHAR(42) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    value NUMERIC(78, 0), -- ERC20/ERC1155
    token_id NUMERIC(78, 0), -- ERC721/ERC1155
    block_number BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, transaction_hash, log_index)
);

CREATE INDEX idx_token_transfers_chain_token ON token_transfers(chain_id, token_address, timestamp DESC);
CREATE INDEX idx_token_transfers_chain_from ON token_transfers(chain_id, from_address, timestamp DESC);
CREATE INDEX idx_token_transfers_chain_to ON token_transfers(chain_id, to_address, timestamp DESC);
CREATE INDEX idx_token_transfers_block ON token_transfers(chain_id, block_number DESC);

-- ============================================================================

CREATE TABLE token_balances (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    token_address VARCHAR(42) NOT NULL,
    holder_address VARCHAR(42) NOT NULL,
    balance NUMERIC(78, 0) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, token_address, holder_address)
);

CREATE INDEX idx_token_balances_chain_token ON token_balances(chain_id, token_address, balance DESC);
CREATE INDEX idx_token_balances_chain_holder ON token_balances(chain_id, holder_address);
CREATE INDEX idx_token_balances_updated ON token_balances(updated_at DESC);

-- ============================================================================
-- WALLET DATA (⚠️ DEV ONLY - See SECURITY.md)
-- ============================================================================

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address VARCHAR(42) UNIQUE NOT NULL,
    encrypted_private_key TEXT NOT NULL, -- AES-256-GCM encrypted
    encryption_iv VARCHAR(32) NOT NULL, -- Initialization vector
    encryption_tag VARCHAR(32) NOT NULL, -- Auth tag for GCM
    name VARCHAR(100),
    derivation_path VARCHAR(100), -- m/44'/60'/0'/0/N
    is_imported BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_wallets_address ON wallets(address);
CREATE INDEX idx_wallets_created ON wallets(created_at DESC);

-- ============================================================================

CREATE TABLE wallet_balances (
    id BIGSERIAL PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    balance NUMERIC(78, 0) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(wallet_id, chain_id)
);

CREATE INDEX idx_wallet_balances_wallet ON wallet_balances(wallet_id);
CREATE INDEX idx_wallet_balances_updated ON wallet_balances(last_updated DESC);

-- ============================================================================

CREATE TABLE wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    tx_hash VARCHAR(66) NOT NULL,
    direction VARCHAR(10) NOT NULL, -- 'sent' or 'received'
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    value NUMERIC(78, 0) NOT NULL,
    gas_used BIGINT,
    gas_price NUMERIC(78, 0),
    status INT NOT NULL, -- 1: success, 0: failed, 2: pending
    block_number BIGINT,
    timestamp TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (direction IN ('sent', 'received'))
);

CREATE INDEX idx_wallet_tx_wallet_chain ON wallet_transactions(wallet_id, chain_id, timestamp DESC);
CREATE INDEX idx_wallet_tx_hash ON wallet_transactions(chain_id, tx_hash);
CREATE INDEX idx_wallet_tx_status ON wallet_transactions(wallet_id, status);

-- ============================================================================
-- SYSTEM METADATA
-- ============================================================================

CREATE TABLE sync_status (
    chain_id BIGINT PRIMARY KEY REFERENCES chains(chain_id) ON DELETE CASCADE,
    last_synced_block BIGINT NOT NULL DEFAULT 0,
    latest_block BIGINT NOT NULL DEFAULT 0, -- From chain
    is_syncing BOOLEAN DEFAULT false,
    sync_rate FLOAT, -- blocks/second
    last_sync_time TIMESTAMP,
    error_count INT DEFAULT 0,
    last_error TEXT,
    last_error_time TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sync_status_syncing ON sync_status(is_syncing) WHERE is_syncing = true;

-- ============================================================================

CREATE TABLE reorgs (
    id SERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    old_block_number BIGINT NOT NULL,
    old_block_hash VARCHAR(66) NOT NULL,
    new_block_hash VARCHAR(66) NOT NULL,
    depth INT NOT NULL, -- How many blocks reorged
    detected_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reorgs_chain ON reorgs(chain_id, detected_at DESC);

-- ============================================================================
-- ANALYTICS (Materialized Views - Refresh periodically)
-- ============================================================================

CREATE MATERIALIZED VIEW daily_stats AS
SELECT 
    chain_id,
    DATE(timestamp) as date,
    COUNT(*) as tx_count,
    COUNT(DISTINCT from_address) as unique_senders,
    COUNT(DISTINCT COALESCE(to_address, contract_address)) as unique_receivers,
    SUM(value) as total_value_transferred,
    AVG(gas_price) as avg_gas_price,
    SUM(gas_used) as total_gas_used,
    AVG(gas_used) as avg_gas_used
FROM transactions
WHERE status = 1 -- Only successful transactions
GROUP BY chain_id, DATE(timestamp);

CREATE UNIQUE INDEX idx_daily_stats_chain_date ON daily_stats(chain_id, date DESC);

-- Manual refresh: REFRESH MATERIALIZED VIEW CONCURRENTLY daily_stats;

-- ============================================================================

CREATE MATERIALIZED VIEW hourly_stats AS
SELECT 
    chain_id,
    DATE_TRUNC('hour', timestamp) as hour,
    COUNT(*) as tx_count,
    AVG(gas_price) as avg_gas_price,
    SUM(gas_used) as total_gas_used
FROM transactions
WHERE status = 1
GROUP BY chain_id, DATE_TRUNC('hour', timestamp);

CREATE UNIQUE INDEX idx_hourly_stats_chain_hour ON hourly_stats(chain_id, hour DESC);

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_chains_updated_at BEFORE UPDATE ON chains
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tokens_updated_at BEFORE UPDATE ON tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- INITIAL DATA
-- ============================================================================

INSERT INTO chains (chain_id, name, short_name, native_symbol, rpc_endpoint, block_time_seconds, is_testnet, is_active) VALUES
(1337, 'Local Testnet', 'local', 'ETH', 'http://localhost:8545', 15, true, true),
(11155111, 'Ethereum Sepolia', 'sepolia', 'ETH', 'https://rpc.sepolia.org', 12, true, false),
(1, 'Ethereum Mainnet', 'ethereum', 'ETH', 'https://eth.llamarpc.com', 12, false, false),
(80002, 'Polygon Amoy', 'amoy', 'MATIC', 'https://rpc-amoy.polygon.technology', 2, true, false);

-- ============================================================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================================================

COMMENT ON TABLE chains IS 'Supported EVM chains configuration';
COMMENT ON TABLE blocks IS 'Indexed blockchain blocks';
COMMENT ON TABLE transactions IS 'Indexed blockchain transactions with receipts';
COMMENT ON TABLE transaction_logs IS 'Event logs emitted by transactions';
COMMENT ON TABLE addresses IS 'Address metadata and activity tracking';
COMMENT ON TABLE tokens IS 'ERC20/721/1155 token contracts';
COMMENT ON TABLE token_transfers IS 'Token transfer events';
COMMENT ON TABLE token_balances IS 'Current token balances per holder';
COMMENT ON TABLE wallets IS '⚠️ DEV ONLY: Encrypted wallet keys (NOT for production)';
COMMENT ON TABLE sync_status IS 'Indexer synchronization status per chain';
COMMENT ON TABLE reorgs IS 'Detected blockchain reorganizations';

COMMENT ON COLUMN wallets.encrypted_private_key IS 'AES-256-GCM encrypted private key. NEVER log or expose.';
COMMENT ON COLUMN wallets.encryption_iv IS 'Unique initialization vector for AES-GCM';
COMMENT ON COLUMN wallets.encryption_tag IS 'Authentication tag for AES-GCM';
```

---

## V. API SPECIFICATION (REST + SSE)

### Base Configuration

```
Base URL: http://localhost:8080/api/v1
Content-Type: application/json
```

### Standard Response Format

```json
{
  "success": true,
  "data": { /* ... */ },
  "error": null,
  "meta": {
    "timestamp": "2024-01-07T12:00:00Z",
    "chain_id": 1337,
    "version": "1.0.0"
  }
}
```

### Error Response Format

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "INVALID_ADDRESS",
    "message": "Address format is invalid",
    "details": {
      "address": "0xinvalid",
      "expected_format": "0x[40 hex chars]"
    }
  },
  "meta": {
    "timestamp": "2024-01-07T12:00:00Z",
    "request_id": "uuid-here"
  }
}
```

### Error Codes

```
INVALID_ADDRESS          - Malformed Ethereum address
INVALID_HASH             - Malformed transaction/block hash
CHAIN_NOT_FOUND          - Chain ID not configured
RESOURCE_NOT_FOUND       - Block/tx/address not found
DATABASE_ERROR           - Internal database error
RPC_ERROR                - Blockchain RPC error
RATE_LIMIT_EXCEEDED      - Too many requests
INTERNAL_SERVER_ERROR    - Unexpected error
```

---

### 1. Health & Status

#### `GET /health`
Health check endpoint

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "database": "connected",
    "redis": "connected",
    "uptime_seconds": 3600
  }
}
```

#### `GET /chains`
List all configured chains

**Response:**
```json
{
  "data": {
    "chains": [
      {
        "chain_id": 1337,
        "name": "Local Testnet",
        "short_name": "local",
        "native_symbol": "ETH",
        "block_time_seconds": 15,
        "is_active": true,
        "is_testnet": true,
        "icon_url": "/icons/ethereum.svg"
      }
    ]
  }
}
```

#### `GET /chains/:chainId/status`
Sync status for a chain

**Response:**
```json
{
  "data": {
    "chain_id": 1337,
    "last_synced_block": 12450,
    "latest_block": 12450,
    "is_syncing": false,
    "sync_percentage": 100.0,
    "blocks_behind": 0,
    "sync_rate": 45.5,
    "last_sync_time": "2024-01-07T12:00:00Z"
  }
}
```

---

### 2. Blocks

#### `GET /blocks`
List recent blocks

**Query Parameters:**
```
chain_id   (int)     Chain ID (default: 1337)
page       (int)     Page number (default: 1)
limit      (int)     Items per page (default: 20, max: 100)
sort       (string)  Sort order: 'asc' or 'desc' (default: 'desc')
```

**Response:**
```json
{
  "data": {
    "blocks": [
      {
        "block_number": 12450,
        "hash": "0xabc...",
        "parent_hash": "0xdef...",
        "timestamp": "2024-01-07T12:00:00Z",
        "miner": "0x123...",
        "tx_count": 15,
        "gas_used": 3500000,
        "gas_limit": 8000000,
        "base_fee_per_gas": "15000000000"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 12450,
      "total_pages": 623
    }
  }
}
```

#### `GET /blocks/:blockId`
Get block details (by number or hash)

**Path Parameters:**
```
blockId  (string)  Block number or hash
```

**Query Parameters:**
```
chain_id       (int)   Chain ID
include_txs    (bool)  Include full transactions (default: false)
```

**Response:**
```json
{
  "data": {
    "block_number": 12450,
    "hash": "0xabc...",
    "parent_hash": "0xdef...",
    "timestamp": "2024-01-07T12:00:00Z",
    "miner": "0x123...",
    "difficulty": "1",
    "total_difficulty": "12450",
    "size": 2048,
    "gas_used": 3500000,
    "gas_limit": 8000000,
    "base_fee_per_gas": "15000000000",
    "nonce": "0x0000000000000000",
    "extra_data": "0x...",
    "tx_count": 15,
    "transactions": [
      "0xhash1...",
      "0xhash2..."
    ]
  }
}
```

---

### 3. Transactions

#### `GET /transactions`
List recent transactions

**Query Parameters:**
```
chain_id       (int)     Chain ID
page           (int)     Page number
limit          (int)     Items per page
block_number   (int)     Filter by block
address        (string)  Filter by from/to address
status         (string)  'success' or 'failed'
```

**Response:**
```json
{
  "data": {
    "transactions": [
      {
        "hash": "0x123...",
        "block_number": 12450,
        "timestamp": "2024-01-07T12:00:00Z",
        "from": "0xfrom...",
        "to": "0xto...",
        "value": "1000000000000000000",
        "gas_used": 21000,
        "gas_price": "20000000000",
        "status": 1
      }
    ],
    "pagination": { /* ... */ }
  }
}
```

#### `GET /transactions/:hash`
Get transaction details

**Query Parameters:**
```
chain_id  (int)  Chain ID
```

**Response:**
```json
{
  "data": {
    "hash": "0x123...",
    "block_number": 12450,
    "block_hash": "0xabc...",
    "timestamp": "2024-01-07T12:00:00Z",
    "from": "0xfrom...",
    "to": "0xto...",
    "value": "1000000000000000000",
    "gas": 21000,
    "gas_price": "20000000000",
    "gas_used": 21000,
    "effective_gas_price": "20000000000",
    "nonce": 5,
    "transaction_index": 2,
    "transaction_type": 0,
    "input": "0x",
    "status": 1,
    "contract_address": null,
    "logs": [
      {
        "log_index": 0,
        "address": "0xcontract...",
        "topics": ["0xtopic0...", "0xtopic1..."],
        "data": "0x..."
      }
    ]
  }
}
```

---

### 4. Addresses

#### `GET /addresses/:address`
Get address details

**Query Parameters:**
```
chain_id  (int)  Chain ID
```

**Response:**
```json
{
  "data": {
    "address": "0x123...",
    "balance": "5000000000000000000",
    "nonce": 10,
    "is_contract": false,
    "tx_count": 25,
    "first_seen_block": 100,
    "last_seen_block": 12450,
    "first_seen_at": "2024-01-01T00:00:00Z",
    "last_seen_at": "2024-01-07T12:00:00Z",
    "token_balances": [
      {
        "token_address": "0xtoken...",
        "token_name": "MyToken",
        "token_symbol": "MTK",
        "token_type": "ERC20",
        "balance": "1000000000000000000",
        "decimals": 18
      }
    ]
  }
}
```

#### `GET /addresses/:address/transactions`
Get transactions for an address

**Query Parameters:**
```
chain_id    (int)     Chain ID
page        (int)     Page number
limit       (int)     Items per page
direction   (string)  'sent', 'received', or 'all' (default: 'all')
```

**Response:**
```json
{
  "data": {
    "address": "0x123...",
    "transactions": [ /* same format as /transactions */ ],
    "pagination": { /* ... */ }
  }
}
```

#### `GET /addresses/:address/tokens`
Get token balances

**Query Parameters:**
```
chain_id     (int)     Chain ID
token_type   (string)  Filter: 'ERC20', 'ERC721', 'ERC1155'
```

**Response:**
```json
{
  "data": {
    "tokens": [
      {
        "token_address": "0xtoken...",
        "name": "MyToken",
        "symbol": "MTK",
        "type": "ERC20",
        "balance": "1000000000000000000",
        "decimals": 18
      }
    ]
  }
}
```

---

### 5. Tokens

#### `GET /tokens`
List tokens

**Query Parameters:**
```
chain_id  (int)     Chain ID
type      (string)  'ERC20', 'ERC721', 'ERC1155'
page      (int)     Page number
limit     (int)     Items per page
sort_by   (string)  'holder_count', 'transfer_count', 'created_at'
```

#### `GET /tokens/:address`
Get token details

**Response:**
```json
{
  "data": {
    "address": "0xtoken...",
    "type": "ERC20",
    "name": "MyToken",
    "symbol": "MTK",
    "decimals": 18,
    "total_supply": "1000000000000000000000000",
    "holder_count": 150,
    "transfer_count": 5000,
    "contract_creator": "0xcreator...",
    "creation_tx": "0xtx..."
  }
}
```

#### `GET /tokens/:address/transfers`
Token transfer history

**Query Parameters:**
```
chain_id  (int)     Chain ID
page      (int)     Page number
limit     (int)     Items per page
from      (string)  Filter by sender
to        (string)  Filter by receiver
```

#### `GET /tokens/:address/holders`
Token holders

**Query Parameters:**
```
chain_id      (int)     Chain ID
page          (int)     Page number
limit         (int)     Items per page
min_balance   (string)  Minimum balance filter
```

---

### 6. Search

#### `GET /search`
Universal search

**Query Parameters:**
```
q         (string)  Query (address, tx hash, block number/hash)
chain_id  (int)     Chain ID
```

**Response:**
```json
{
  "data": {
    "type": "address",  // or "transaction", "block"
    "result": { /* entity data */ }
  }
}
```

**Search Logic:**
1. If starts with `0x` and 66 chars → Block hash or Tx hash
2. If starts with `0x` and 42 chars → Address
3. If numeric → Block number
4. Else → Error

---

### 7. Statistics

#### `GET /stats`
Network statistics

**Query Parameters:**
```
chain_id  (int)  Chain ID
```

**Response:**
```json
{
  "data": {
    "chain_id": 1337,
    "latest_block": 12450,
    "avg_block_time": 15.2,
    "total_transactions": 185000,
    "total_addresses": 5000,
    "active_addresses_24h": 250,
    "tps_24h": 12.5,
    "gas_price": {
      "low": "15000000000",
      "medium": "20000000000",
      "high": "25000000000"
    }
  }
}
```

#### `GET /stats/daily`
Daily aggregated stats

**Query Parameters:**
```
chain_id     (int)     Chain ID
from_date    (string)  Start date (YYYY-MM-DD)
to_date      (string)  End date (YYYY-MM-DD)
```

**Response:**
```json
{
  "data": {
    "stats": [
      {
        "date": "2024-01-07",
        "tx_count": 5000,
        "unique_senders": 150,
        "unique_receivers": 200,
        "total_value": "500000000000000000000",
        "avg_gas_price": "20000000000",
        "total_gas_used": 100000000
      }
    ]
  }
}
```

---

### 8. Real-time (Server-Sent Events)

#### `GET /stream/blocks`
Stream new blocks

**Query Parameters:**
```
chain_id  (int)  Chain ID
```

**Headers:**
```
Accept: text/event-stream
```

**Response Stream:**
```
event: block
data: {"block_number": 12451, "hash": "0x...", "timestamp": "...", ...}

event: block
data: {"block_number": 12452, ...}
```

#### `GET /stream/transactions`
Stream new transactions

**Query Parameters:**
```
chain_id  (int)     Chain ID
address   (string)  Filter by address (optional)
```

**Response Stream:**
```
event: transaction
data: {"hash": "0x...", "from": "0x...", "to": "0x...", ...}
```

---

### 9. Wallet API (⚠️ DEV ONLY)

**⚠️ SECURITY WARNING:**
These endpoints store private keys server-side (encrypted). This is for DEVELOPMENT/LEARNING only. Production wallets should NEVER store private keys on servers. Use client-side signing (MetaMask, WalletConnect, etc.).

#### `POST /wallet/create`
Create new HD wallet

**Request:**
```json
{
  "name": "My Wallet",
  "password": "secure_password_min_12_chars"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "address": "0x123...",
    "mnemonic": "word1 word2 ... word12",
    "warning": "⚠️ SAVE YOUR MNEMONIC SECURELY - Cannot be recovered!"
  }
}
```

**Security Notes:**
- Password hashed with Argon2id before use as encryption key
- Private key encrypted with AES-256-GCM
- Mnemonic shown ONCE, never stored
- Session timeout after 30 minutes

#### `POST /wallet/import`
Import wallet from mnemonic or private key

**Request:**
```json
{
  "method": "mnemonic",  // or "private_key"
  "mnemonic": "word1 word2 ...",  // or "private_key": "0x..."
  "password": "secure_password",
  "name": "Imported Wallet"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "address": "0x123...",
    "name": "Imported Wallet"
  }
}
```

#### `GET /wallet/:walletId/balance`
Get balances across all chains

**Headers:**
```
Authorization: Bearer <session_token>
```

**Response:**
```json
{
  "data": {
    "wallet_id": "uuid",
    "address": "0x123...",
    "balances": [
      {
        "chain_id": 1337,
        "chain_name": "Local Testnet",
        "balance": "5000000000000000000",
        "symbol": "ETH"
      }
    ]
  }
}
```

#### `POST /wallet/send`
Send transaction

**Request:**
```json
{
  "wallet_id": "uuid",
  "chain_id": 1337,
  "to": "0xto...",
  "value": "1000000000000000000",
  "password": "secure_password",
  "gas_limit": 21000,       // optional
  "gas_price": "20000000000" // optional
}
```

**Response:**
```json
{
  "data": {
    "tx_hash": "0x123...",
    "status": "pending"
  }
}
```

**Security Flow:**
1. Verify password → decrypt private key (in memory)
2. Sign transaction → send to blockchain
3. Clear private key from memory
4. Never log private key or password

#### `POST /wallet/sign`
Sign arbitrary message (EIP-191)

**Request:**
```json
{
  "wallet_id": "uuid",
  "message": "Sign this message",
  "password": "secure_password"
}
```

**Response:**
```json
{
  "data": {
    "message": "Sign this message",
    "signature": "0xsignature...",
    "message_hash": "0xhash..."
  }
}
```

#### `GET /wallet/:walletId/transactions`
Wallet transaction history

**Query Parameters:**
```
chain_id  (int)  Filter by chain (optional)
page      (int)  Page number
limit     (int)  Items per page
```

---

## VI. CONFIGURATION FILES

### Backend Configuration

**`backend/internal/config/chains.json`:**
```json
{
  "chains": [
    {
      "chain_id": 1337,
      "name": "Local Testnet",
      "short_name": "local",
      "native_symbol": "ETH",
      "rpc_endpoint": "http://localhost:8545",
      "ws_endpoint": "ws://localhost:8546",
      "block_time_seconds": 15,
      "is_testnet": true,
      "is_active": true,
      "explorer_url": "",
      "icon_url": "/icons/ethereum.svg",
      "gas_price_oracle": "legacy",
      "supports_eip1559": false,
      "backup_rpc_endpoints": []
    },
    {
      "chain_id": 11155111,
      "name": "Ethereum Sepolia",
      "short_name": "sepolia",
      "native_symbol": "ETH",
      "rpc_endpoint": "${SEPOLIA_RPC_URL}",
      "ws_endpoint": "${SEPOLIA_WS_URL}",
      "block_time_seconds": 12,
      "is_testnet": true,
      "is_active": false,
      "explorer_url": "https://sepolia.etherscan.io",
      "icon_url": "/icons/ethereum.svg",
      "gas_price_oracle": "eip1559",
      "supports_eip1559": true,
      "backup_rpc_endpoints": [
        "https://rpc.sepolia.org",
        "https://ethereum-sepolia.publicnode.com"
      ]
    },
    {
      "chain_id": 1,
      "name": "Ethereum Mainnet",
      "short_name": "ethereum",
      "native_symbol": "ETH",
      "rpc_endpoint": "${ETHEREUM_RPC_URL}",
      "ws_endpoint": "${ETHEREUM_WS_URL}",
      "block_time_seconds": 12,
      "is_testnet": false,
      "is_active": false,
      "explorer_url": "https://etherscan.io",
      "icon_url": "/icons/ethereum.svg",
      "gas_price_oracle": "eip1559",
      "supports_eip1559": true,
      "backup_rpc_endpoints": [
        "https://eth.llamarpc.com",
        "https://rpc.ankr.com/eth"
      ]
    },
    {
      "chain_id": 80002,
      "name": "Polygon Amoy",
      "short_name": "amoy",
      "native_symbol": "MATIC",
      "rpc_endpoint": "${POLYGON_AMOY_RPC_URL}",
      "ws_endpoint": "${POLYGON_AMOY_WS_URL}",
      "block_time_seconds": 2,
      "is_testnet": true,
      "is_active": false,
      "explorer_url": "https://amoy.polygonscan.com",
      "icon_url": "/icons/polygon.svg",
      "gas_price_oracle": "eip1559",
      "supports_eip1559": true
    }
  ],
  "default_chain_id": 1337
}
```

**`backend/.env.example`:**
```bash
# ============================================================================
# SERVER
# ============================================================================

API_PORT=8080
API_HOST=0.0.0.0
ENVIRONMENT=development  # development, staging, production

# ============================================================================
# DATABASE
# ============================================================================

DB_HOST=localhost
DB_PORT=5432
DB_NAME=ethereum_explorer
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Connection string format (alternative to above)
# DATABASE_URL=postgresql://user:pass@localhost:5432/dbname?sslmode=disable

# ============================================================================
# REDIS (Phase 8+)
# ============================================================================

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# ============================================================================
# CHAIN RPC ENDPOINTS (Public chains - optional)
# ============================================================================

# Ethereum Sepolia
SEPOLIA_RPC_URL=https://rpc.sepolia.org
SEPOLIA_WS_URL=

# Ethereum Mainnet
ETHEREUM_RPC_URL=https://eth.llamarpc.com
ETHEREUM_WS_URL=

# Polygon Amoy
POLYGON_AMOY_RPC_URL=https://rpc-amoy.polygon.technology
POLYGON_AMOY_WS_URL=

# ============================================================================
# WALLET ENCRYPTION (⚠️ DEV ONLY)
# ============================================================================

# Generate with: openssl rand -hex 32
WALLET_ENCRYPTION_KEY=your-32-byte-hex-key-here-change-this-immediately

# Argon2id parameters (for password hashing)
ARGON2_TIME=1
ARGON2_MEMORY=64
ARGON2_THREADS=4
ARGON2_KEY_LENGTH=32

# ============================================================================
# INDEXER
# ============================================================================

INDEXER_BATCH_SIZE=100           # Blocks per batch
INDEXER_WORKERS=3                # Parallel workers per chain
INDEXER_START_BLOCK=0            # Start from genesis
INDEXER_REORG_CHECK_DEPTH=12     # How many blocks back to check for reorgs

# ============================================================================
# API FEATURES
# ============================================================================

# Rate Limiting (Phase 8+)
RATE_LIMIT_ENABLED=false
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_BURST=20

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourfrontend.vercel.app
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# ============================================================================
# LOGGING
# ============================================================================

LOG_LEVEL=debug                  # debug, info, warn, error
LOG_FORMAT=json                  # json, console
LOG_OUTPUT=stdout                # stdout, file
LOG_FILE_PATH=./logs/app.log

# ============================================================================
# MONITORING (Phase 8+)
# ============================================================================

PROMETHEUS_ENABLED=false
PROMETHEUS_PORT=9090
PROMETHEUS_PATH=/metrics

# ============================================================================
# SECURITY
# ============================================================================

# JWT for wallet sessions (⚠️ DEV ONLY)
JWT_SECRET=your-jwt-secret-here-change-this
JWT_EXPIRY=30m

# API Key for admin endpoints (optional)
ADMIN_API_KEY=

# ============================================================================
# MISC
# ============================================================================

# Graceful shutdown timeout
SHUTDOWN_TIMEOUT=30s

# Request timeout
REQUEST_TIMEOUT=30s
```

### Frontend Configuration

**`frontend/.env.local.example`:**
```bash
# API Backend
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Chain Configs (can override chains.ts)
NEXT_PUBLIC_DEFAULT_CHAIN_ID=1337

# WalletConnect (optional)
NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID=

# Analytics (optional)
NEXT_PUBLIC_GA_TRACKING_ID=

# Environment
NEXT_PUBLIC_ENV=development
```

**`frontend/src/lib/chains.ts`:**
```typescript
export interface Chain {
  chainId: number;
  name: string;
  shortName: string;
  nativeSymbol: string;
  rpcUrl: string;
  explorerUrl?: string;
  iconUrl: string;
  isTestnet: boolean;
}

export const chains: Chain[] = [
  {
    chainId: 1337,
    name: 'Local Testnet',
    shortName: 'local',
    nativeSymbol: 'ETH',
    rpcUrl: 'http://localhost:8545',
    iconUrl: '/icons/ethereum.svg',
    isTestnet: true,
  },
  {
    chainId: 11155111,
    name: 'Ethereum Sepolia',
    shortName: 'sepolia',
    nativeSymbol: 'ETH',
    rpcUrl: 'https://rpc.sepolia.org',
    explorerUrl: 'https://sepolia.etherscan.io',
    iconUrl: '/icons/ethereum.svg',
    isTestnet: true,
  },
  {
    chainId: 1,
    name: 'Ethereum',
    shortName: 'ethereum',
    nativeSymbol: 'ETH',
    rpcUrl: 'https://eth.llamarpc.com',
    explorerUrl: 'https://etherscan.io',
    iconUrl: '/icons/ethereum.svg',
    isTestnet: false,
  },
  {
    chainId: 80002,
    name: 'Polygon Amoy',
    shortName: 'amoy',
    nativeSymbol: 'MATIC',
    rpcUrl: 'https://rpc-amoy.polygon.technology',
    explorerUrl: 'https://amoy.polygonscan.com',
    iconUrl: '/icons/polygon.svg',
    isTestnet: true,
  },
];

export const getChainById = (chainId: number): Chain | undefined => {
  return chains.find((c) => c.chainId === chainId);
};

export const defaultChain = chains[0];
```

---

## VII. PHASE-BY-PHASE IMPLEMENTATION

---

## **PHASE 1: PRIVATE TESTNET** ✅ COMPLETE

See previous setup - you've already done this!

**Deliverables:**
- ✅ 2-node Clique PoA network
- ✅ Blocks producing every 15s
- ✅ RPC accessible at localhost:8545
- ✅ Pre-funded test accounts
- ✅ Docker Compose setup
- ✅ Documentation

---

## **PHASE 2: BLOCK EXPLORER BACKEND**

### Goal
PostgreSQL database + Go indexer syncing blockchain data + REST API

### Phase 2A: Database Setup

**Tasks:**

1. **Add PostgreSQL to docker-compose**
2. **Run database migrations**
3. **Verify schema created**

**Testing:**
- Connect to PostgreSQL
- Verify all tables exist
- Check indexes created

### Phase 2B: Go Backend Structure

**Tasks:**

1. **Initialize Go module**
2. **Setup project structure**
3. **Install dependencies**
4. **Create base configuration**
5. **Setup structured logging (zap)**

**Testing:**
- `go build` compiles without errors
- Config loads from env
- Logger outputs JSON format

### Phase 2C: Blockchain Client Layer

**Tasks:**

1. **Implement ChainManager**
   - Load chains from config
   - Create RPC clients
   - Handle failover

2. **Implement ChainClient wrapper**
   - Get block by number/hash
   - Get transaction + receipt
   - Get latest block number
   - Estimate gas

3. **Write unit tests**
   - Mock RPC responses
   - Test error handling
   - Test failover logic

**Testing:**
- Unit tests pass
- Can connect to local geth
- Failover works (stop primary, uses backup)

### Phase 2D: Database Access Layer

**Tasks:**

1. **Implement database.DB wrapper**
   - Connection pooling
   - Health check
   - Transaction helpers

2. **Implement data access methods**
   - InsertBlock
   - InsertTransaction
   - InsertLog
   - GetBlock (by number/hash)
   - GetTransaction
   - GetAddress
   - Update address balances

3. **Write unit tests**
   - Use test database
   - Test CRUD operations
   - Test constraints (unique, foreign keys)

**Testing:**
- All tests pass
- Can insert/query data
- Constraints enforced

### Phase 2E: Indexer Service

**Tasks:**

1. **Implement core indexer**
   - Sync loop (poll for new blocks)
   - Parse blocks → database
   - Parse transactions + receipts
   - Parse logs (event detection)
   - Update address balances

2. **Implement reorg handler**
   - Detect reorgs (block hash mismatch)
   - Delete orphaned data
   - Re-sync canonical blocks

3. **Add monitoring/metrics**
   - Blocks/sec sync rate
   - Current sync status
   - Error tracking

4. **Write tests**
   - Mock blockchain responses
   - Test full sync flow
   - Test reorg handling

**Testing:**
- Start indexer → syncs from genesis
- Stop indexer → resume from last block
- Artificially create reorg → handles correctly
- Performance: >50 blocks/sec on local network

### Phase 2F: REST API

**Tasks:**

1. **Setup Fiber server**
   - Routes
   - Middleware (CORS, logging, recovery)
   - Error handling

2. **Implement all handlers**
   - Health check
   - Chains
   - Blocks (list, get)
   - Transactions (list, get)
   - Addresses (get, txs, tokens)
   - Search
   - Stats

3. **Implement SSE for real-time**
   - Stream new blocks
   - Stream new transactions

4. **Write integration tests**
   - Test all endpoints
   - Test error cases
   - Test pagination
   - Test search logic

**Testing:**
- All endpoints return 200
- Pagination works
- Search finds data
- SSE streams blocks in real-time
- Error responses formatted correctly

### Phase 2G: End-to-End Testing

**Tasks:**

1. **Integration test suite**
   - Start testnet
   - Start indexer
   - Wait for sync
   - Query API
   - Verify data matches chain

2. **Load testing**
   - Send 100 transactions to testnet
   - Verify all indexed correctly
   - Check API performance under load

**Success Criteria:**
- Indexer syncs 1000 blocks in <20 seconds
- API responds <100ms for simple queries
- 100% of transactions indexed correctly
- No memory leaks (run for 1 hour)

**Deliverables:**
- ✅ PostgreSQL with full schema
- ✅ Indexer syncing blockchain → database
- ✅ REST API serving all endpoints
- ✅ Real-time SSE for blocks/txs
- ✅ >80% test coverage on Go code
- ✅ Makefile with helper commands
- ✅ Documentation (API.md)

---

## **PHASE 3: BLOCK EXPLORER FRONTEND**

### Goal
Single Next.js app with `/explorer` route displaying blockchain data

### Phase 3A: Project Setup

**Tasks:**

1. **Initialize Next.js 14 (App Router)**
2. **Install dependencies**
   - viem, wagmi, @tanstack/react-query
   - TailwindCSS, shadcn/ui
   - recharts (for charts)
3. **Setup project structure**
4. **Configure TypeScript (strict mode)**
5. **Setup wagmi config**

**Testing:**
- `npm run dev` starts without errors
- Can navigate to routes
- Hot reload works

### Phase 3B: Shared Components & Context

**Tasks:**

1. **Create ChainContext**
   - Selected chain state
   - Switch chain function

2. **Create shared components**
   - Header (with nav)
   - NetworkSwitcher (dropdown)
   - Footer
   - LoadingSpinner
   - ErrorBoundary

3. **Create API client**
   - Wrapper around fetch
   - Type-safe (TypeScript interfaces)
   - Error handling

**Testing:**
- Can switch networks in UI
- API client fetches data
- Components render correctly

### Phase 3C: Explorer Routes

**Tasks:**

1. **Explorer home (`/explorer`)**
   - Recent blocks list
   - Recent transactions list
   - Network stats cards

2. **Block list (`/explorer/blocks`)**
   - Paginated block list
   - Search bar

3. **Block detail (`/explorer/blocks/[id]`)**
   - Full block info
   - Transaction list
   - Link to parent/child blocks

4. **Transaction detail (`/explorer/tx/[hash]`)**
   - Full transaction info
   - Event logs
   - Link to block and addresses

5. **Address detail (`/explorer/address/[address]`)**
   - Balance
   - Transaction history
   - Token balances

6. **Search (`/explorer/search`)**
   - Universal search
   - Results page

**Testing:**
- All routes load
- Data displays correctly
- Links work
- Pagination works
- Search works

### Phase 3D: Real-time Updates

**Tasks:**

1. **Implement SSE or polling**
   - New blocks appear automatically
   - Block number updates in header

2. **Add notification system**
   - Toast for new blocks (optional)

**Testing:**
- New blocks appear without refresh
- Performance doesn't degrade over time

### Phase 3E: Polish & Responsive Design

**Tasks:**

1. **Make responsive (mobile-friendly)**
2. **Add loading states**
3. **Add error states**
4. **Add empty states**
5. **Improve typography, spacing**
6. **Add tooltips for technical terms**

**Testing:**
- Works on mobile
- Looks professional
- All states handled

### Phase 3F: Deployment

**Tasks:**

1. **Build for production**
2. **Deploy to Vercel**
3. **Configure env variables**
4. **Test deployed version**

**Deliverables:**
- ✅ Functional block explorer UI
- ✅ All pages implemented
- ✅ Real-time updates
- ✅ Responsive design
- ✅ Deployed to Vercel
- ✅ Works with local backend

---

## **PHASE 4: WALLET BACKEND** (⚠️ DEV ONLY)

### Goal
HD wallet service with transaction signing (with security warnings)

### Phase 4A: Security Documentation

**Before writing code, create `backend/internal/wallet/README.md`:**

```markdown
# Wallet Service (⚠️ DEVELOPMENT ONLY)

## Security Warning

**This wallet implementation stores encrypted private keys on the server. This is INSECURE and should NEVER be used in production.**

### Why This is Unsafe

1. **Server Compromise**: If server is hacked, all keys exposed
2. **Database Breach**: Even encrypted, keys are vulnerable
3. **Admin Access**: Server admins can access keys
4. **No Hardware Security**: Software encryption < hardware wallets

### Production Alternatives

- **Client-side signing**: Keys never leave browser (MetaMask, WalletConnect)
- **Hardware wallets**: Ledger, Trezor
- **MPC wallets**: Multi-party computation (no single key)
- **Smart contract wallets**: Social recovery (Argent, Safe)

### Why We Built This

**Learning purposes:**
- Understanding HD wallets (BIP39/BIP44)
- Key derivation
- Transaction signing
- Encryption best practices

**Use cases:**
- Local development
- Testing
- Demonstrating concepts

## Architecture

[Document the rest...]
```

**Create `docs/SECURITY.md`:**

```markdown
# Security Considerations

## Wallet Security (Critical)

### Current Implementation (DEV ONLY)

The wallet service stores encrypted private keys server-side using:
- Argon2id for password hashing
- AES-256-GCM for encryption
- Unique IVs per key
- Authentication tags

**This is still fundamentally insecure** because:
1. Keys exist on server (even if encrypted)
2. Server compromise = all keys compromised
3. Requires trusting server operator

### Production Recommendations

**Never store private keys server-side. Period.**

Recommended approaches:
1. **Client-side wallets** (MetaMask, WalletConnect)
2. **Hardware wallets** (Ledger, Trezor)
3. **Smart contract wallets** (Safe, Argent - social recovery)
4. **MPC wallets** (Fireblocks - distributed key generation)

[Continue with other security topics...]
```

### Phase 4B: Cryptography Setup

**Tasks:**

1. **Implement password hashing (Argon2id)**
   - High time/memory cost
   - Unique salt per user

2. **Implement key encryption (AES-256-GCM)**
   - Unique IV per encryption
   - Authentication tags
   - Constant-time comparison

3. **Write extensive tests**
   - Test encryption/decryption
   - Test wrong password fails
   - Test tampered ciphertext fails
   - Test key derivation

**Testing:**
- Encrypt/decrypt works
- Wrong password fails gracefully
- Timing attacks prevented (constant-time ops)

### Phase 4C: HD Wallet Implementation

**Tasks:**

1. **Implement BIP39 (mnemonic generation)**
   - 12/24 word mnemonic
   - Validate mnemonic

2. **Implement BIP32/BIP44 (key derivation)**
   - Standard Ethereum path: `m/44'/60'/0'/0/N`
   - Derive multiple accounts

3. **Write tests**
   - Generate mnemonic → same as MetaMask
   - Derive accounts → match MetaMask
   - Test derivation path variations

**Testing:**
- Same mnemonic in MetaMask and code → same addresses
- Can derive 100 accounts without issues

### Phase 4D: Transaction Signing

**Tasks:**

1. **Implement transaction builder**
   - Build legacy transactions
   - Build EIP-1559 transactions
   - Proper RLP encoding

2. **Implement signing**
   - Sign with EIP-155 (replay protection)
   - Proper v/r/s values

3. **Implement nonce management**
   - Track pending nonces
   - Handle concurrent requests
   - Retry on nonce errors

4. **Write tests**
   - Build & sign transaction
   - Verify signature valid
   - Test nonce incrementing

**Testing:**
- Signed transaction broadcasts successfully
- Nonce management handles 10 concurrent sends

### Phase 4E: Wallet Service & API

**Tasks:**

1. **Implement WalletService**
   - CreateWallet
   - ImportWallet
   - GetBalance (multi-chain)
   - SendTransaction
   - SignMessage

2. **Implement API handlers**
   - POST /wallet/create
   - POST /wallet/import
   - GET /wallet/:id/balance
   - POST /wallet/send
   - POST /wallet/sign
   - GET /wallet/:id/transactions

3. **Implement session management (JWT)**
   - Generate token on wallet unlock
   - Expire after 30 minutes
   - Refresh token mechanism

4. **Write integration tests**
   - Full flow: create → send → verify

**Testing:**
- Can create wallet via API
- Can send transaction
- Sessions expire correctly
- Wrong password rejected

### Phase 4F: Security Audit & Documentation

**Tasks:**

1. **Self-audit checklist**
   - No private keys in logs ✓
   - No passwords in logs ✓
   - Constant-time comparisons ✓
   - Secure random generation ✓
   - SQL injection prevented ✓
   - HTTPS only (in production) ✓

2. **Document security considerations**
   - What we did right
   - What's still insecure
   - Production alternatives

**Deliverables:**
- ✅ HD wallet backend functional
- ✅ Can create/import wallets
- ✅ Can send transactions
- ✅ Encrypted storage
- ✅ Security warnings everywhere
- ✅ Tests passing (>80% coverage)
- ✅ Documentation complete

---

## **PHASE 5: WALLET FRONTEND**

### Goal
`/wallet` route in Next.js app for wallet management

### Phase 5A: Wallet Creation Flow

**Tasks:**

1. **Create wallet page (`/wallet/create`)**
   - Password input (strength indicator)
   - Confirm password
   - Create button

2. **Mnemonic display page**
   - Show 12 words
   - Copy button
   - Multiple warnings
   - Checkbox: "I have saved my mnemonic"
   - Continue button

**Testing:**
- Can create wallet
- Mnemonic shown once
- Warnings clear

### Phase 5B: Wallet Import Flow

**Tasks:**

1. **Import page (`/wallet/import`)**
   - Mnemonic input (12 word fields)
   - Or: Private key input
   - Password setup
   - Import button

**Testing:**
- Can import wallet
- Same mnemonic → same address

### Phase 5C: Wallet Dashboard

**Tasks:**

1. **Dashboard page (`/wallet/dashboard`)**
   - Balance cards (per chain)
   - Recent transactions
   - Quick actions (send/receive)
   - Network switcher

2. **Account management**
   - List derived accounts
   - Add new account
   - Switch active account

**Testing:**
- Shows correct balances
- Switches chains correctly

### Phase 5D: Send Transaction Flow

**Tasks:**

1. **Send page (`/wallet/send`)**
   - Recipient input (with validation)
   - Amount input
   - Gas settings (simple/advanced)
   - Password confirmation
   - Send button

2. **Transaction confirmation modal**
   - Review details
   - Estimated fee
   - Confirm/Cancel

3. **Transaction status**
   - Pending state
   - Link to explorer
   - Success/failure notification

**Testing:**
- Can send transaction
- Gas estimation works
- Wrong password rejected
- Transaction appears in explorer

### Phase 5E: Receive & History

**Tasks:**

1. **Receive modal**
   - Show address
   - QR code
   - Copy button

2. **Transaction history**
   - List all transactions
   - Filter by chain
   - Link to explorer

**Testing:**
- QR code scannable
- History accurate

### Phase 5F: Security & Polish

**Tasks:**

1. **Session management**
   - Auto-lock after 30 min
   - Lock button
   - Unlock prompt

2. **Security warnings**
   - Banner: "This is a dev wallet"
   - Link to security docs

3. **Polish**
   - Responsive design
   - Loading states
   - Error handling

**Deliverables:**
- ✅ Functional wallet UI
- ✅ Can create/import
- ✅ Can send transactions
- ✅ Multi-chain support
- ✅ Security warnings visible
- ✅ Deployed to Vercel

---

## **PHASE 6: SMART CONTRACTS**

### Goal
3 production-grade contracts with >95% test coverage

### Phase 6A: Foundry Setup

**Tasks:**

1. **Initialize Foundry project**
   - `forge init`
   - Install OpenZeppelin
   - Configure remappings

2. **Setup testing environment**
   - Create test helpers
   - Setup mock contracts

**Testing:**
- Example test passes
- Can import OpenZeppelin

### Phase 6B: Contract 1 - Advanced ERC20

**Requirements:**
- Standard ERC20 (transfer, approve, etc.)
- Mintable (owner only)
- Burnable
- Pausable (owner only)
- Max supply cap
- Token locking (time-based)
- Snapshot (for governance - future)

**Implementation:**

```solidity
// src/Token.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MyToken
 * @notice Advanced ERC20 implementation with minting, burning, pausing, and time-lock features
 * @dev Inherits from OpenZeppelin's battle-tested contracts
 */
contract MyToken is ERC20, ERC20Burnable, ERC20Pausable, Ownable {
    /// @notice Maximum total supply (1 million tokens)
    uint256 public constant MAX_SUPPLY = 1_000_000 * 10**18;
    
    /// @notice Mapping of addresses to their token unlock timestamp
    mapping(address => uint256) public lockedUntil;
    
    /// @notice Emitted when tokens are locked for an address
    event TokensLocked(address indexed account, uint256 until);
    
    /**
     * @notice Contract constructor
     * @dev Mints initial supply to deployer
     */
    constructor() ERC20("MyToken", "MTK") Ownable(msg.sender) {
        _mint(msg.sender, 100_000 * 10**18); // 100k initial supply
    }
    
    /**
     * @notice Mint new tokens (owner only)
     * @param to Recipient address
     * @param amount Amount to mint
     * @dev Cannot exceed MAX_SUPPLY
     */
    function mint(address to, uint256 amount) public onlyOwner {
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        _mint(to, amount);
    }
    
    /**
     * @notice Pause all token transfers (owner only)
     * @dev Useful for emergency situations
     */
    function pause() public onlyOwner {
        _pause();
    }
    
    /**
     * @notice Unpause token transfers (owner only)
     */
    function unpause() public onlyOwner {
        _unpause();
    }
    
    /**
     * @notice Lock tokens for an address until a specific timestamp
     * @param account Address to lock tokens for
     * @param until Unix timestamp when tokens unlock
     * @dev Owner only. Useful for vesting, team tokens, etc.
     */
    function lockTokens(address account, uint256 until) public onlyOwner {
        require(until > block.timestamp, "Must be future timestamp");
        lockedUntil[account] = until;
        emit TokensLocked(account, until);
    }
    
    /**
     * @notice Internal update hook that enforces pausing and locking
     * @dev Overrides both ERC20 and ERC20Pausable
     */
    function _update(address from, address to, uint256 value)
        internal
        override(ERC20, ERC20Pausable)
    {
        // Check if sender's tokens are locked
        if (from != address(0)) {
            require(block.timestamp >= lockedUntil[from], "Tokens are locked");
        }
        
        super._update(from, to, value);
    }
}
```

**Testing:**

```solidity
// test/Token.t.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Token.sol";

contract TokenTest is Test {
    MyToken token;
    address owner = address(1);
    address user1 = address(2);
    address user2 = address(3);
    
    function setUp() public {
        vm.prank(owner);
        token = new MyToken();
    }
    
    function testInitialSupply() public {
        assertEq(token.totalSupply(), 100_000 * 10**18);
        assertEq(token.balanceOf(owner), 100_000 * 10**18);
    }
    
    function testMint() public {
        vm.prank(owner);
        token.mint(user1, 1000 * 10**18);
        assertEq(token.balanceOf(user1), 1000 * 10**18);
    }
    
    function testMintExceedsMaxSupply() public {
        vm.prank(owner);
        vm.expectRevert("Exceeds max supply");
        token.mint(user1, 1_000_000 * 10**18); // Would exceed cap
    }
    
    function testNonOwnerCannotMint() public {
        vm.prank(user1);
        vm.expectRevert();
        token.mint(user1, 1000);
    }
    
    function testPause() public {
        vm.prank(owner);
        token.transfer(user1, 1000);
        
        vm.prank(owner);
        token.pause();
        
        vm.prank(owner);
        vm.expectRevert();
        token.transfer(user2, 100);
    }
    
    function testUnpause() public {
        vm.prank(owner);
        token.pause();
        
        vm.prank(owner);
        token.unpause();
        
        vm.prank(owner);
        token.transfer(user1, 1000);
        assertEq(token.balanceOf(user1), 1000);
    }
    
    function testLockTokens() public {
        // Give tokens to user1
        vm.prank(owner);
        token.transfer(user1, 1000 * 10**18);
        
        // Lock for 1 day
        uint256 unlockTime = block.timestamp + 1 days;
        vm.prank(owner);
        token.lockTokens(user1, unlockTime);
        
        // Try to transfer - should fail
        vm.prank(user1);
        vm.expectRevert("Tokens are locked");
        token.transfer(user2, 100);
        
        // Fast forward past lock time
        vm.warp(unlockTime + 1);
        
        // Now should work
        vm.prank(user1);
        token.transfer(user2, 100);
        assertEq(token.balanceOf(user2), 100);
    }
    
    function testBurn() public {
        vm.prank(owner);
        token.burn(1000 * 10**18);
        assertEq(token.totalSupply(), 99_000 * 10**18);
    }
    
    function testBurnFrom() public {
        // Owner approves user1 to burn
        vm.prank(owner);
        token.approve(user1, 1000 * 10**18);
        
        // User1 burns owner's tokens
        vm.prank(user1);
        token.burnFrom(owner, 500 * 10**18);
        
        assertEq(token.totalSupply(), 99_500 * 10**18);
        assertEq(token.allowance(owner, user1), 500 * 10**18);
    }
    
    function testTransfer() public {
        vm.prank(owner);
        token.transfer(user1, 1000);
        
        assertEq(token.balanceOf(owner), 100_000 * 10**18 - 1000);
        assertEq(token.balanceOf(user1), 1000);
    }
    
    function testApproveAndTransferFrom() public {
        vm.prank(owner);
        token.approve(user1, 1000);
        
        vm.prank(user1);
        token.transferFrom(owner, user2, 500);
        
        assertEq(token.balanceOf(user2), 500);
        assertEq(token.allowance(owner, user1), 500);
    }
}

// Fuzz tests
contract TokenFuzzTest is Test {
    MyToken token;
    address owner = address(1);
    
    function setUp() public {
        vm.prank(owner);
        token = new MyToken();
    }
    
    function testFuzzMint(address to, uint256 amount) public {
        vm.assume(to != address(0));
        vm.assume(amount <= token.MAX_SUPPLY() - token.totalSupply());
        
        vm.prank(owner);
        token.mint(to, amount);
        
        assertEq(token.balanceOf(to), amount);
    }
    
    function testFuzzTransfer(address to, uint256 amount) public {
        vm.assume(to != address(0) && to != owner);
        vm.assume(amount <= token.balanceOf(owner));
        
        vm.prank(owner);
        token.transfer(to, amount);
        
        assertEq(token.balanceOf(to), amount);
    }
}
```

**Tasks:**

1. Implement token contract
2. Write comprehensive tests (aim for 100%)
3. Run gas report
4. Document gas optimizations
5. Run security tools (slither)

**Gas Optimization Notes:**

```markdown
# Gas Optimizations - MyToken

## Storage Packing
- Used uint256 for all values (EVM word size)
- Mappings don't benefit from packing

## Optimizations Applied
1. **Constant MAX_SUPPLY**: Saves gas vs storage variable
2. **Inherited implementations**: Use OpenZeppelin's optimized code
3. **Single _update override**: Reduces bytecode size

## Gas Report
```
| Function        | Gas (avg) | Gas (optimized) | Savings |
|-----------------|-----------|-----------------|---------|
| transfer        | 51,000    | 51,000          | -       |
| mint            | 48,000    | 48,000          | -       |
| burn            | 28,000    | 28,000          | -       |
| lockTokens      | 45,000    | 45,000          | -       |
```

## Trade-offs
- Chose clarity over micro-optimizations
- OpenZeppelin contracts prioritize security
- Extra features (locking) add gas cost but provide value
```

**Testing:**
- All tests pass
- Fuzz tests pass (100+ runs)
- Coverage >95%
- Gas usage documented

### Phase 6C: Contract 2 - Staking

**Requirements:**
- Stake ERC20 tokens
- Earn rewards over time
- Multiple reward tokens (optional)
- Slashing for early withdrawal (penalty)
- Emergency withdrawal (forfeit rewards)
- Time-weighted rewards
- Minimum stake duration

**Implementation:**

```solidity
// src/Staking.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title Staking
 * @notice Stake tokens to earn rewards over time
 * @dev Uses reward-per-token calculation for fair distribution
 */
contract Staking is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;
    
    /// @notice Token to stake
    IERC20 public immutable stakingToken;
    
    /// @notice Token to reward (can be same as staking token)
    IERC20 public immutable rewardToken;
    
    /// @notice Reward rate (tokens per second across all stakers)
    uint256 public rewardRate = 100 * 10**18; // 100 tokens per second
    
    /// @notice Last time rewards were calculated
    uint256 public lastUpdateTime;
    
    /// @notice Accumulated reward per token staked
    uint256 public rewardPerTokenStored;
    
    /// @notice User's snapshot of rewardPerToken when they last claimed
    mapping(address => uint256) public userRewardPerTokenPaid;
    
    /// @notice Unclaimed rewards for each user
    mapping(address => uint256) public rewards;
    
    /// @notice Amount staked by each user
    mapping(address => uint256) public stakedBalance;
    
    /// @notice Timestamp when user staked
    mapping(address => uint256) public stakeTimestamp;
    
    /// @notice Total amount staked across all users
    uint256 public totalStaked;
    
    /// @notice Minimum stake duration (to prevent flash loan attacks)
    uint256 public constant MIN_STAKE_DURATION = 1 days;
    
    /// @notice Penalty for early withdrawal (10% = 1000)
    uint256 public earlyWithdrawalPenalty = 1000; // 10%
    uint256 public constant PENALTY_DENOMINATOR = 10000;
    
    event Staked(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount, uint256 penalty);
    event RewardsClaimed(address indexed user, uint256 amount);
    event RewardRateUpdated(uint256 newRate);
    event EmergencyWithdraw(address indexed user, uint256 amount);
    
    /**
     * @notice Constructor
     * @param _stakingToken Token to stake
     * @param _rewardToken Token to reward
     */
    constructor(address _stakingToken, address _rewardToken) Ownable(msg.sender) {
        stakingToken = IERC20(_stakingToken);
        rewardToken = IERC20(_rewardToken);
        lastUpdateTime = block.timestamp;
    }
    
    /**
     * @notice Calculate current reward per token
     * @return Current accumulated reward per token
     */
    function rewardPerToken() public view returns (uint256) {
        if (totalStaked == 0) {
            return rewardPerTokenStored;
        }
        
        return rewardPerTokenStored + 
            (((block.timestamp - lastUpdateTime) * rewardRate * 1e18) / totalStaked);
    }
    
    /**
     * @notice Calculate earned rewards for an account
     * @param account Address to check
     * @return Amount of rewards earned
     */
    function earned(address account) public view returns (uint256) {
        return ((stakedBalance[account] * 
            (rewardPerToken() - userRewardPerTokenPaid[account])) / 1e18) + 
            rewards[account];
    }
    
    /**
     * @notice Modifier to update rewards before state changes
     * @param account Address to update rewards for
     */
    modifier updateReward(address account) {
        rewardPerTokenStored = rewardPerToken();
        lastUpdateTime = block.timestamp;
        
        if (account != address(0)) {
            rewards[account] = earned(account);
            userRewardPerTokenPaid[account] = rewardPerTokenStored;
        }
        _;
    }
    
    /**
     * @notice Stake tokens
     * @param amount Amount to stake
     */
    function stake(uint256 amount) external nonReentrant updateReward(msg.sender) {
        require(amount > 0, "Cannot stake 0");
        
        totalStaked += amount;
        stakedBalance[msg.sender] += amount;
        stakeTimestamp[msg.sender] = block.timestamp;
        
        stakingToken.safeTransferFrom(msg.sender, address(this), amount);
        emit Staked(msg.sender, amount);
    }
    
    /**
     * @notice Withdraw staked tokens
     * @param amount Amount to withdraw
     * @dev Penalty applied if withdrawn before MIN_STAKE_DURATION
     */
    function withdraw(uint256 amount) external nonReentrant updateReward(msg.sender) {
        require(amount > 0, "Cannot withdraw 0");
        require(stakedBalance[msg.sender] >= amount, "Insufficient balance");
        
        uint256 penalty = 0;
        
        // Apply penalty if withdrawing before minimum duration
        if (block.timestamp < stakeTimestamp[msg.sender] + MIN_STAKE_DURATION) {
            penalty = (amount * earlyWithdrawalPenalty) / PENALTY_DENOMINATOR;
        }
        
        totalStaked -= amount;
        stakedBalance[msg.sender] -= amount;
        
        uint256 amountAfterPenalty = amount - penalty;
        
        stakingToken.safeTransfer(msg.sender, amountAfterPenalty);
        
        // Send penalty to owner (could be burned instead)
        if (penalty > 0) {
            stakingToken.safeTransfer(owner(), penalty);
        }
        
        emit Withdrawn(msg.sender, amountAfterPenalty, penalty);
    }
    
    /**
     * @notice Claim accumulated rewards
     */
    function claimRewards() external nonReentrant updateReward(msg.sender) {
        uint256 reward = rewards[msg.sender];
        require(reward > 0, "No rewards");
        
        rewards[msg.sender] = 0;
        rewardToken.safeTransfer(msg.sender, reward);
        
        emit RewardsClaimed(msg.sender, reward);
    }
    
    /**
     * @notice Emergency withdraw (forfeit all rewards)
     * @dev Use when contract needs to be shut down
     */
    function emergencyWithdraw() external nonReentrant {
        uint256 amount = stakedBalance[msg.sender];
        require(amount > 0, "Nothing staked");
        
        // Don't update rewards (forfeit them)
        totalStaked -= amount;
        stakedBalance[msg.sender] = 0;
        rewards[msg.sender] = 0;
        
        stakingToken.safeTransfer(msg.sender, amount);
        
        emit EmergencyWithdraw(msg.sender, amount);
    }
    
    /**
     * @notice Update reward rate (owner only)
     * @param _rewardRate New reward rate (tokens per second)
     */
    function setRewardRate(uint256 _rewardRate) external onlyOwner updateReward(address(0)) {
        rewardRate = _rewardRate;
        emit RewardRateUpdated(_rewardRate);
    }
    
    /**
     * @notice Update early withdrawal penalty (owner only)
     * @param _penalty New penalty (basis points, max 50% = 5000)
     */
    function setEarlyWithdrawalPenalty(uint256 _penalty) external onlyOwner {
        require(_penalty <= 5000, "Penalty too high"); // Max 50%
        earlyWithdrawalPenalty = _penalty;
    }
    
    /**
     * @notice Emergency function to recover stuck tokens (owner only)
     * @param token Token to recover
     * @param amount Amount to recover
     * @dev Should only be used for tokens accidentally sent to contract
     */
    function recoverToken(address token, uint256 amount) external onlyOwner {
        require(token != address(stakingToken), "Cannot recover staking token");
        IERC20(token).safeTransfer(owner(), amount);
    }
    
    /**
     * @notice Get stake info for an account
     * @param account Address to check
     * @return staked Amount staked
     * @return earned Amount of rewards earned
     * @return canWithdrawWithoutPenalty Whether can withdraw without penalty
     */
    function getStakeInfo(address account) external view returns (
        uint256 staked,
        uint256 earned_,
        bool canWithdrawWithoutPenalty
    ) {
        staked = stakedBalance[account];
        earned_ = earned(account);
        canWithdrawWithoutPenalty = block.timestamp >= stakeTimestamp[account] + MIN_STAKE_DURATION;
    }
}
```

**Testing (abbreviated - similar structure to Token tests):**

```solidity
// test/Staking.t.sol
contract StakingTest is Test {
    MyToken stakingToken;
    MyToken rewardToken;
    Staking staking;
    
    address owner = address(1);
    address user1 = address(2);
    address user2 = address(3);
    
    function setUp() public {
        // Deploy tokens
        vm.prank(owner);
        stakingToken = new MyToken();
        
        vm.prank(owner);
        rewardToken = new MyToken();
        
        // Deploy staking
        vm.prank(owner);
        staking = new Staking(address(stakingToken), address(rewardToken));
        
        // Fund staking contract with rewards
        vm.prank(owner);
        rewardToken.transfer(address(staking), 1_000_000 * 10**18);
        
        // Give users staking tokens
        vm.prank(owner);
        stakingToken.transfer(user1, 10_000 * 10**18);
        
        vm.prank(owner);
        stakingToken.transfer(user2, 10_000 * 10**18);
    }
    
    function testStake() public {
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        assertEq(staking.stakedBalance(user1), 1000 * 10**18);
        assertEq(staking.totalStaked(), 1000 * 10**18);
    }
    
    function testEarnRewards() public {
        // User1 stakes
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        // Fast forward 1 day
        vm.warp(block.timestamp + 1 days);
        
        // Check earned rewards
        // 100 tokens/sec * 86400 sec = 8,640,000 tokens
        uint256 earned = staking.earned(user1);
        assertApproxEqRel(earned, 8_640_000 * 10**18, 0.01e18); // 1% tolerance
    }
    
    function testClaimRewards() public {
        // Stake and wait
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        vm.warp(block.timestamp + 1 days);
        
        // Claim
        vm.prank(user1);
        staking.claimRewards();
        
        // Check balance
        assertTrue(rewardToken.balanceOf(user1) > 0);
        assertEq(staking.rewards(user1), 0);
    }
    
    function testEarlyWithdrawalPenalty() public {
        // Stake
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        // Try to withdraw immediately (within 1 day)
        vm.warp(block.timestamp + 1 hours);
        
        uint256 balanceBefore = stakingToken.balanceOf(user1);
        
        vm.prank(user1);
        staking.withdraw(1000 * 10**18);
        
        uint256 balanceAfter = stakingToken.balanceOf(user1);
        
        // Should receive 90% (10% penalty)
        assertEq(balanceAfter - balanceBefore, 900 * 10**18);
    }
    
    function testWithdrawAfterMinDuration() public {
        // Stake
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        // Wait minimum duration
        vm.warp(block.timestamp + 1 days + 1);
        
        uint256 balanceBefore = stakingToken.balanceOf(user1);
        
        vm.prank(user1);
        staking.withdraw(1000 * 10**18);
        
        uint256 balanceAfter = stakingToken.balanceOf(user1);
        
        // Should receive 100% (no penalty)
        assertEq(balanceAfter - balanceBefore, 1000 * 10**18);
    }
    
    function testMultipleStakers() public {
        // User1 stakes
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        // Fast forward
        vm.warp(block.timestamp + 1 days);
        
        // User2 stakes (50% of pool)
        vm.startPrank(user2);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        // Fast forward another day
        vm.warp(block.timestamp + 1 days);
        
        // User1 should have ~1.5 days worth (1 full day + 0.5 day at 50% share)
        // User2 should have ~0.5 days worth (at 50% share)
        uint256 user1Earned = staking.earned(user1);
        uint256 user2Earned = staking.earned(user2);
        
        assertTrue(user1Earned > user2Earned);
    }
    
    function testEmergencyWithdraw() public {
        // Stake
        vm.startPrank(user1);
        stakingToken.approve(address(staking), 1000 * 10**18);
        staking.stake(1000 * 10**18);
        vm.stopPrank();
        
        // Earn some rewards
        vm.warp(block.timestamp + 1 days);
        
        uint256 earnedBefore = staking.earned(user1);
        assertTrue(earnedBefore > 0);
        
        // Emergency withdraw
        vm.prank(user1);
        staking.emergencyWithdraw();
        
        // Should get staked tokens back but lose rewards
        assertEq(staking.stakedBalance(user1), 0);
        assertEq(staking.rewards(user1), 0);
        assertEq(stakingToken.balanceOf(user1), 10_000 * 10**18);
    }
}
```

**Tasks:**

1. Implement staking contract
2. Write comprehensive tests
3. Test reward calculation accuracy
4. Test penalty mechanics
5. Run gas report
6. Security audit (reentrancy, overflow, etc.)

**Testing:**
- All tests pass
- Fuzz tests pass
- Reward calculation accurate to 0.1%
- Coverage >95%

### Phase 6D: Contract 3 - Simple DEX (AMM)

**Requirements:**
- Constant product AMM (x * y = k)
- Swap tokens
- Add/remove liquidity
- LP tokens
- 0.3% fee
- Slippage protection

**Implementation:**

```solidity
// src/DEX.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title LPToken
 * @notice ERC20 token representing liquidity provider shares
 */
contract LPToken is ERC20 {
    address public immutable dex;
    
    constructor(string memory name, string memory symbol) ERC20(name, symbol) {
        dex = msg.sender;
    }
    
    function mint(address to, uint256 amount) external {
        require(msg.sender == dex, "Only DEX can mint");
        _mint(to, amount);
    }
    
    function burn(address from, uint256 amount) external {
        require(msg.sender == dex, "Only DEX can burn");
        _burn(from, amount);
    }
}

/**
 * @title SimpleDEX
 * @notice Automated Market Maker using constant product formula
 * @dev x * y = k, where x and y are token reserves
 */
contract SimpleDEX is ReentrancyGuard {
    using SafeERC20 for IERC20;
    
    /// @notice First token in the pair
    IERC20 public immutable token0;
    
    /// @notice Second token in the pair
    IERC20 public immutable token1;
    
    /// @notice LP token for this pair
    LPToken public immutable lpToken;
    
    /// @notice Reserve of token0
    uint256 public reserve0;
    
    /// @notice Reserve of token1
    uint256 public reserve1;
    
    /// @notice Trading fee (30 = 0.3%)
    uint256 private constant FEE = 3;
    uint256 private constant FEE_DENOMINATOR = 1000;
    
    /// @notice Minimum liquidity locked forever
    uint256 private constant MINIMUM_LIQUIDITY = 1000;
    
    event Swap(
        address indexed user,
        address tokenIn,
        uint256 amountIn,
        uint256 amountOut
    );
    
    event LiquidityAdded(
        address indexed provider,
        uint256 amount0,
        uint256 amount1,
        uint256 lpAmount
    );
    
    event LiquidityRemoved(
        address indexed provider,
        uint256 amount0,
        uint256 amount1,
        uint256 lpAmount
    );
    
    /**
     * @notice Constructor
     * @param _token0 Address of first token
     * @param _token1 Address of second token
     */
    constructor(address _token0, address _token1) {
        require(_token0 != _token1, "Identical tokens");
        require(_token0 != address(0) && _token1 != address(0), "Zero address");
        
        token0 = IERC20(_token0);
        token1 = IERC20(_token1);
        
        lpToken = new LPToken(
            "DEX LP Token",
            "DEX-LP"
        );
    }
    
    /**
     * @notice Add liquidity to the pool
     * @param amount0 Amount of token0 to add
     * @param amount1 Amount of token1 to add
     * @return lpAmount Amount of LP tokens minted
     * @dev For first liquidity, amounts can be anything. After that, must maintain price ratio.
     */
    function addLiquidity(uint256 amount0, uint256 amount1)
        external
        nonReentrant
        returns (uint256 lpAmount)
    {
        require(amount0 > 0 && amount1 > 0, "Amounts must be > 0");
        
        // First liquidity provision
        if (reserve0 == 0 && reserve1 == 0) {
            lpAmount = sqrt(amount0 * amount1);
            require(lpAmount > MINIMUM_LIQUIDITY, "Insufficient liquidity");
            
            // Lock minimum liquidity forever
            lpToken.mint(address(1), MINIMUM_LIQUIDITY);
            lpAmount -= MINIMUM_LIQUIDITY;
        } else {
            // Calculate optimal amounts to maintain price ratio
            uint256 amount1Optimal = (amount0 * reserve1) / reserve0;
            
            if (amount1Optimal <= amount1) {
                amount1 = amount1Optimal;
            } else {
                amount0 = (amount1 * reserve0) / reserve1;
            }
            
            // LP tokens proportional to share of pool
            lpAmount = min(
                (amount0 * lpToken.totalSupply()) / reserve0,
                (amount1 * lpToken.totalSupply()) / reserve1
            );
        }
        
        require(lpAmount > 0, "LP amount too small");
        
        // Transfer tokens
        token0.safeTransferFrom(msg.sender, address(this), amount0);
        token1.safeTransferFrom(msg.sender, address(this), amount1);
        
        // Update reserves
        reserve0 += amount0;
        reserve1 += amount1;
        
        // Mint LP tokens
        lpToken.mint(msg.sender, lpAmount);
        
        emit LiquidityAdded(msg.sender, amount0, amount1, lpAmount);
    }
    
    /**
     * @notice Remove liquidity from the pool
     * @param lpAmount Amount of LP tokens to burn
     * @return amount0 Amount of token0 returned
     * @return amount1 Amount of token1 returned
     */
    function removeLiquidity(uint256 lpAmount)
        external
        nonReentrant
        returns (uint256 amount0, uint256 amount1)
    {
        require(lpAmount > 0, "Amount must be > 0");
        
        uint256 totalSupply = lpToken.totalSupply();
        
        // Calculate amounts to return (proportional to share)
        amount0 = (lpAmount * reserve0) / totalSupply;
        amount1 = (lpAmount * reserve1) / totalSupply;
        
        require(amount0 > 0 && amount1 > 0, "Insufficient liquidity burned");
        
        // Burn LP tokens
        lpToken.burn(msg.sender, lpAmount);
        
        // Update reserves
        reserve0 -= amount0;
        reserve1 -= amount1;
        
        // Transfer tokens
        token0.safeTransfer(msg.sender, amount0);
        token1.safeTransfer(msg.sender, amount1);
        
        emit LiquidityRemoved(msg.sender, amount0, amount1, lpAmount);
    }
    
    /**
     * @notice Swap tokens
     * @param tokenIn Address of token to swap
     * @param amountIn Amount of tokenIn to swap
     * @param minAmountOut Minimum amount of tokenOut (slippage protection)
     * @return amountOut Amount of tokenOut received
     * @dev Uses constant product formula with 0.3% fee
     */
    function swap(
        address tokenIn,
        uint256 amountIn,
        uint256 minAmountOut
    ) external nonReentrant returns (uint256 amountOut) {
        require(amountIn > 0, "Amount must be > 0");
        require(
            tokenIn == address(token0) || tokenIn == address(token1),
            "Invalid token"
        );
        
        bool isToken0 = tokenIn == address(token0);
        
        (IERC20 tIn, IERC20 tOut, uint256 resIn, uint256 resOut) = isToken0
            ? (token0, token1, reserve0, reserve1)
            : (token1, token0, reserve1, reserve0);
        
        // Transfer in
        tIn.safeTransferFrom(msg.sender, address(this), amountIn);
        
        // Calculate output with fee
        // amountOut = (amountIn * 997 * resOut) / (resIn * 1000 + amountIn * 997)
        uint256 amountInWithFee = amountIn * (FEE_DENOMINATOR - FEE);
        amountOut = (amountInWithFee * resOut) / (resIn * FEE_DENOMINATOR + amountInWithFee);
        
        require(amountOut >= minAmountOut, "Slippage exceeded");
        require(amountOut < resOut, "Insufficient liquidity");
        
        // Update reserves
        if (isToken0) {
            reserve0 += amountIn;
            reserve1 -= amountOut;
        } else {
            reserve1 += amountIn;
            reserve0 -= amountOut;
        }
        
        // Transfer out
        tOut.safeTransfer(msg.sender, amountOut);
        
        emit Swap(msg.sender, tokenIn, amountIn, amountOut);
    }
    
    /**
     * @notice Get quote for swap
     * @param tokenIn Address of token to swap
     * @param amountIn Amount of tokenIn
     * @return amountOut Amount of tokenOut that would be received
     */
    function getAmountOut(address tokenIn, uint256 amountIn)
        public
        view
        returns (uint256 amountOut)
    {
        require(amountIn > 0, "Amount must be > 0");
        require(
            tokenIn == address(token0) || tokenIn == address(token1),
            "Invalid token"
        );
        
        bool isToken0 = tokenIn == address(token0);
        (uint256 resIn, uint256 resOut) = isToken0
            ? (reserve0, reserve1)
            : (reserve1, reserve0);
        
        uint256 amountInWithFee = amountIn * (FEE_DENOMINATOR - FEE);
        amountOut = (amountInWithFee * resOut) / (resIn * FEE_DENOMINATOR + amountInWithFee);
    }
    
    /**
     * @notice Get current price of token0 in terms of token1
     * @return price Price as token1 per token0 (scaled by 1e18)
     */
    function getPrice0() external view returns (uint256 price) {
        require(reserve0 > 0, "No liquidity");
        price = (reserve1 * 1e18) / reserve0;
    }
    
    /**
     * @notice Get current price of token1 in terms of token0
     * @return price Price as token0 per token1 (scaled by 1e18)
     */
    function getPrice1() external view returns (uint256 price) {
        require(reserve1 > 0, "No liquidity");
        price = (reserve0 * 1e18) / reserve1;
    }
    
    /**
     * @notice Square root function (Babylonian method)
     * @param y Input
     * @return z Square root of y
     */
    function sqrt(uint256 y) private pure returns (uint256 z) {
        if (y > 3) {
            z = y;
            uint256 x = y / 2 + 1;
            while (x < z) {
                z = x;
                x = (y / x + x) / 2;
            }
        } else if (y != 0) {
            z = 1;
        }
    }
    
    /**
     * @notice Minimum of two numbers
     */
    function min(uint256 a, uint256 b) private pure returns (uint256) {
        return a < b ? a : b;
    }
}
```

**Testing (abbreviated):**

```solidity
// test/DEX.t.sol
contract DEXTest is Test {
    MyToken token0;
    MyToken token1;
    SimpleDEX dex;
    
    address user1 = address(1);
    address user2 = address(2);
    
    function setUp() public {
        // Deploy tokens
        token0 = new MyToken();
        token1 = new MyToken();
        
        // Deploy DEX
        dex = new SimpleDEX(address(token0), address(token1));
        
        // Give users tokens
        token0.transfer(user1, 100_000 * 10**18);
        token1.transfer(user1, 100_000 * 10**18);
        
        token0.transfer(user2, 10_000 * 10**18);
        token1.transfer(user2, 10_000 * 10**18);
    }
    
    function testAddInitialLiquidity() public {
        vm.startPrank(user1);
        token0.approve(address(dex), 1000 * 10**18);
        token1.approve(address(dex), 1000 * 10**18);
        
        uint256 lpAmount = dex.addLiquidity(1000 * 10**18, 1000 * 10**18);
        vm.stopPrank();
        
        assertTrue(lpAmount > 0);
        assertEq(dex.reserve0(), 1000 * 10**18);
        assertEq(dex.reserve1(), 1000 * 10**18);
    }
    
    function testSwap() public {
        // Add liquidity
        vm.startPrank(user1);
        token0.approve(address(dex), 1000 * 10**18);
        token1.approve(address(dex), 1000 * 10**18);
        dex.addLiquidity(1000 * 10**18, 1000 * 10**18);
        vm.stopPrank();
        
        // Swap
        vm.startPrank(user2);
        token0.approve(address(dex), 100 * 10**18);
        
        uint256 balanceBefore = token1.balanceOf(user2);
        dex.swap(address(token0), 100 * 10**18, 0);
        uint256 balanceAfter = token1.balanceOf(user2);
        
        vm.stopPrank();
        
        uint256 received = balanceAfter - balanceBefore;
        assertTrue(received > 0);
        assertTrue(received < 100 * 10**18); // Should be less due to slippage + fee
    }
    
    function testPriceImpact() public {
        // Add liquidity (1:1 ratio)
        vm.startPrank(user1);
        token0.approve(address(dex), 1000 * 10**18);
        token1.approve(address(dex), 1000 * 10**18);
        dex.addLiquidity(1000 * 10**18, 1000 * 10**18);
        vm.stopPrank();
        
        // Small swap should have minimal impact
        uint256 amountOut = dex.getAmountOut(address(token0), 10 * 10**18);
        assertApproxEqRel(amountOut, 9.97 * 10**18, 0.01e18); // ~0.3% fee
        
        // Large swap should have significant impact
        uint256 amountOutLarge = dex.getAmountOut(address(token0), 500 * 10**18);
        assertTrue(amountOutLarge < 500 * 10**18 * 997 / 1000); // More than just fee
    }
    
    function testRemoveLiquidity() public {
        // Add liquidity
        vm.startPrank(user1);
        token0.approve(address(dex), 1000 * 10**18);
        token1.approve(address(dex), 1000 * 10**18);
        uint256 lpAmount = dex.addLiquidity(1000 * 10**18, 1000 * 10**18);
        vm.stopPrank();
        
        // Remove half
        vm.startPrank(user1);
        (uint256 amount0, uint256 amount1) = dex.removeLiquidity(lpAmount / 2);
        vm.stopPrank();
        
        assertApproxEqRel(amount0, 500 * 10**18, 0.01e18);
        assertApproxEqRel(amount1, 500 * 10**18, 0.01e18);
    }
    
    function testSlippageProtection() public {
        // Add liquidity
        vm.startPrank(user1);
        token0.approve(address(dex), 1000 * 10**18);
        token1.approve(address(dex), 1000 * 10**18);
        dex.addLiquidity(1000 * 10**18, 1000 * 10**18);
        vm.stopPrank();
        
        // Try swap with unrealistic minAmountOut
        vm.startPrank(user2);
        token0.approve(address(dex), 100 * 10**18);
        
        vm.expectRevert("Slippage exceeded");
        dex.swap(address(token0), 100 * 10**18, 100 * 10**18); // Expect 1:1, won't get it
        
        vm.stopPrank();
    }
    
    function testConstantProduct() public {
        // Add liquidity
        vm.startPrank(user1);
        token0.approve(address(dex), 1000 * 10**18);
        token1.approve(address(dex), 1000 * 10**18);
        dex.addLiquidity(1000 * 10**18, 1000 * 10**18);
        vm.stopPrank();
        
        uint256 kBefore = dex.reserve0() * dex.reserve1();
        
        // Swap
        vm.startPrank(user2);
        token0.approve(address(dex), 100 * 10**18);
        dex.swap(address(token0), 100 * 10**18, 0);
        vm.stopPrank();
        
        uint256 kAfter = dex.reserve0() * dex.reserve1();
        
        // k should increase slightly due to fees
        assertTrue(kAfter >= kBefore);
    }
}
```

**Tasks:**

1. Implement DEX contract
2. Write comprehensive tests
3. Test edge cases (price impacts, slippage, etc.)
4. Gas optimization
5. Security audit

**Testing:**
- All tests pass
- Fuzz tests for edge cases
- Math verified (constant product)
- Coverage >95%

### Phase 6E: Deployment & Documentation

**Tasks:**

1. **Deploy all contracts to local testnet**
   ```bash
   forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast
   ```

2. **Save contract addresses**
   ```
   Token: 0x...
   Staking: 0x...
   DEX: 0x...
   ```

3. **Verify contracts work**
   - Mint tokens
   - Stake tokens
   - Add liquidity to DEX
   - Perform swaps

4. **Write comprehensive Natspec**
   - Document all functions
   - Explain formulas
   - Note security considerations

5. **Generate documentation**
   ```bash
   forge doc
   ```

6. **Gas optimization document**
   - Document all optimizations
   - Before/after gas reports
   - Trade-offs made

7. **Security audit document**
   - Tools used (slither, aderyn)
   - Findings
   - Mitigations

**Deliverables:**
- ✅ 3 production-grade contracts
- ✅ >95% test coverage
- ✅ Fuzz tests passing
- ✅ Deployed to local testnet
- ✅ Contract addresses documented
- ✅ Gas optimization documented
- ✅ Security audit performed
- ✅ Natspec complete

---

## **PHASE 7: dAPP FRONTENDS**

### Goal
Routes in Next.js app (`/dapps/*`) to interact with deployed contracts

### Phase 7A: Contract Integration Setup

**Tasks:**

1. **Export contract ABIs**
   ```bash
   cd contracts
   forge inspect Token abi > ../frontend/src/lib/abis/Token.json
   forge inspect Staking abi > ../frontend/src/lib/abis/Staking.json
   forge inspect DEX abi > ../frontend/src/lib/abis/DEX.json
   ```

2. **Create contract config**
   ```typescript
   // frontend/src/lib/contracts.ts
   import TokenABI from './abis/Token.json';
   import StakingABI from './abis/Staking.json';
   import DEXABI from './abis/DEX.json';
   
   export const contracts = {
     token: {
       address: '0x...' as const,
       abi: TokenABI,
     },
     staking: {
       address: '0x...' as const,
       abi: StakingABI,
     },
     dex: {
       address: '0x...' as const,
       abi: DEXABI,
     },
   };
   ```

3. **Create wagmi hooks**
   ```typescript
   // frontend/src/hooks/useContracts.ts
   import { useReadContract, useWriteContract } from 'wagmi';
   import { contracts } from '@/lib/contracts';
   
   export function useToken() {
     // Read: balance, allowance, etc.
     // Write: transfer, approve, etc.
   }
   
   export function useStaking() {
     // Read: staked, earned, etc.
     // Write: stake, withdraw, claim
   }
   
   export function useDEX() {
     // Read: reserves, price, getAmountOut
     // Write: swap, addLiquidity, removeLiquidity
   }
   ```

**Testing:**
- Can read contract data
- Hooks return expected values

### Phase 7B: Token Dashboard

**Route:** `/dapps/token`

**Features:**
- View token info (name, symbol, supply)
- View your balance
- Transfer tokens
- Approve spender
- Burn tokens (if you have balance)
- Mint tokens (if you're owner)

**Implementation:**

```typescript
// frontend/src/app/dapps/token/page.tsx
'use client';

import { useState } from 'react';
import { useAccount, useReadContract, useWriteContract } from 'wagmi';
import { parseEther, formatEther } from 'viem';
import { contracts } from '@/lib/contracts';

export default function TokenDashboard() {
  const { address } = useAccount();
  const [transferTo, setTransferTo] = useState('');
  const [transferAmount, setTransferAmount] = useState('');
  
  // Read token data
  const { data: name } = useReadContract({
    ...contracts.token,
    functionName: 'name',
  });
  
  const { data: symbol } = useReadContract({
    ...contracts.token,
    functionName: 'symbol',
  });
  
  const { data: totalSupply } = useReadContract({
    ...contracts.token,
    functionName: 'totalSupply',
  });
  
  const { data: balance, refetch: refetchBalance } = useReadContract({
    ...contracts.token,
    functionName: 'balanceOf',
    args: [address!],
    query: { enabled: !!address },
  });
  
  // Write: transfer
  const { writeContract: transfer, isPending: isTransferring } = useWriteContract();
  
  const handleTransfer = async () => {
    if (!transferTo || !transferAmount) return;
    
    transfer({
      ...contracts.token,
      functionName: 'transfer',
      args: [transferTo as `0x${string}`, parseEther(transferAmount)],
    }, {
      onSuccess: () => {
        refetchBalance();
        setTransferTo('');
        setTransferAmount('');
      },
    });
  };
  
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-8">Token Dashboard</h1>
      
      {/* Token Info */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="text-sm text-gray-500 mb-1">Token Name</div>
          <div className="text-2xl font-semibold">{name as string}</div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="text-sm text-gray-500 mb-1">Symbol</div>
          <div className="text-2xl font-semibold">{symbol as string}</div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="text-sm text-gray-500 mb-1">Total Supply</div>
          <div className="text-2xl font-semibold">
            {totalSupply ? formatEther(totalSupply as bigint) : '0'} {symbol as string}
          </div>
        </div>
      </div>
      
      {/* Your Balance */}
      <div className="bg-white p-6 rounded-lg shadow mb-8">
        <h2 className="text-2xl font-semibold mb-4">Your Balance</h2>
        <div className="text-4xl font-mono">
          {balance ? formatEther(balance as bigint) : '0'} {symbol as string}
        </div>
      </div>
      
      {/* Transfer */}
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-2xl font-semibold mb-4">Transfer Tokens</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">Recipient Address</label>
            <input
              type="text"
              value={transferTo}
              onChange={(e) => setTransferTo(e.target.value)}
              placeholder="0x..."
              className="w-full px-4 py-2 border rounded-lg font-mono text-sm"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2">Amount</label>
            <input
              type="number"
              value={transferAmount}
              onChange={(e) => setTransferAmount(e.target.value)}
              placeholder="0.0"
              className="w-full px-4 py-2 border rounded-lg"
              step="0.01"
            />
          </div>
          <button
            onClick={handleTransfer}
            disabled={isTransferring || !transferTo || !transferAmount}
            className="w-full bg-blue-600 text-white py-3 rounded-lg disabled:bg-gray-400 hover:bg-blue-700"
          >
            {isTransferring ? 'Transferring...' : 'Transfer'}
          </button>
        </div>
      </div>
    </div>
  );
}
```

**Tasks:**

1. Implement token dashboard UI
2. Add approve functionality
3. Add burn functionality
4. Add mint (if owner)
5. Show transaction status
6. Link to explorer for txs

**Testing:**
- Can view token info
- Can transfer tokens
- Transaction appears in wallet
- Balance updates after transfer

### Phase 7C: Staking dApp

**Route:** `/dapps/staking`

**Features:**
- View staked amount
- View earned rewards
- Stake tokens
- Withdraw tokens
- Claim rewards
- APY calculator

**Implementation (similar pattern to token dashboard)**

**Key components:**
- Approve token before staking
- Real-time rewards counter
- Warning about early withdrawal penalty

**Testing:**
- Can stake tokens
- Rewards update in real-time
- Can withdraw after duration
- Penalty applied if early

### Phase 7D: DEX dApp

**Route:** `/dapps/dex`

**Features:**
- Swap tokens
- Add liquidity
- Remove liquidity
- View reserves & price
- Price impact calculator
- LP token balance

**Implementation:**

Two tabs:
1. **Swap** - Swap token0 <-> token1
2. **Liquidity** - Add/remove liquidity

**Key features:**
- Price impact warning
- Slippage settings
- Gas estimation
- Transaction summary

**Testing:**
- Can swap tokens
- Can add liquidity
- Can remove liquidity
- Price impact shown correctly

### Phase 7E: dApp Home

**Route:** `/dapps`

**Simple landing page:**
- Cards for each dApp
- Brief description
- Link to each

**Testing:**
- Links work
- Layout responsive

### Phase 7F: Polish & Testing

**Tasks:**

1. **Consistent styling across all dApps**
2. **Loading states for all contract calls**
3. **Error handling (reject, insufficient funds, etc.)**
4. **Success notifications**
5. **Transaction history (link to explorer)**
6. **Mobile responsive**

**Deliverables:**
- ✅ Functional dApp UIs for all contracts
- ✅ Users can interact via MetaMask
- ✅ All contract functions accessible
- ✅ Good UX (loading, errors, success states)
- ✅ Deployed to Vercel

---

## **PHASE 8: TESTING, CI/CD & POLISH**

### Goal
Professional-quality codebase with testing, automation, and documentation

### Phase 8A: Backend Testing

**Tasks:**

1. **Ensure >80% coverage**
   ```bash
   cd backend
   go test -cover ./...
   ```

2. **Add integration tests**
   - Full API test suite
   - End-to-end indexer test

3. **Add table-driven tests**
   ```go
   func TestBlockHandler(t *testing.T) {
       tests := []struct{
           name string
           input string
           want int
           wantErr bool
       }{
           {"valid block", "123", 200, false},
           {"invalid block", "abc", 400, true},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // test logic
           })
       }
   }
   ```

4. **Add benchmark tests**
   ```go
   func BenchmarkInsertBlock(b *testing.B) {
       for i := 0; i < b.N; i++ {
           db.InsertBlock(block)
       }
   }
   ```

**Testing:**
- All tests pass
- Coverage >80%
- Benchmarks show acceptable performance

### Phase 8B: Smart Contract Security

**Tasks:**

1. **Run Slither**
   ```bash
   cd contracts
   slither .
   ```

2. **Run Aderyn**
   ```bash
   aderyn .
   ```

3. **Fix critical/high findings**

4. **Document accepted risks**
   ```markdown
   # Security Audit Results
   
   ## Tools Used
   - Slither v0.10.0
   - Aderyn v0.1.0
   
   ## Findings
   
   ### Critical: 0
   ### High: 0
   ### Medium: 2
   1. ...
   2. ...
   
   ### Accepted Risks
   - Low: Unused variable in test contract (informational)
   ```

**Testing:**
- No critical/high vulnerabilities
- Medium/low documented

### Phase 8C: GitHub Actions (CI/CD)

**Tasks:**

1. **Backend tests workflow**
   ```yaml
   # .github/workflows/backend-tests.yml
   name: Backend Tests
   
   on:
     push:
       paths:
         - 'backend/**'
     pull_request:
       paths:
         - 'backend/**'
   
   jobs:
     test:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - uses: actions/setup-go@v5
           with:
             go-version: '1.21'
         
         - name: Run tests
           working-directory: ./backend
           run: |
             go test -v -cover ./...
             
         - name: Check coverage
           working-directory: ./backend
           run: |
             go test -coverprofile=coverage.out ./...
             go tool cover -func=coverage.out
   ```

2. **Contract tests workflow**
   ```yaml
   # .github/workflows/contract-tests.yml
   name: Contract Tests
   
   on:
     push:
       paths:
         - 'contracts/**'
     pull_request:
       paths:
         - 'contracts/**'
   
   jobs:
     test:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
           with:
             submodules: recursive
             
         - name: Install Foundry
           uses: foundry-rs/foundry-toolchain@v1
         
         - name: Run tests
           working-directory: ./contracts
           run: |
             forge test -vvv
             forge coverage
   ```

3. **Frontend build workflow**
   ```yaml
   # .github/workflows/frontend-build.yml
   name: Frontend Build
   
   on:
     push:
       paths:
         - 'frontend/**'
     pull_request:
       paths:
         - 'frontend/**'
   
   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - uses: actions/setup-node@v4
           with:
             node-version: '20'
             
         - name: Install dependencies
           working-directory: ./frontend
           run: npm ci
           
         - name: Build
           working-directory: ./frontend
           run: npm run build
   ```

**Testing:**
- Push to GitHub
- Verify workflows run
- All checks pass

### Phase 8D: Documentation

**Tasks:**

1. **Root README.md**
   - Project overview
   - Quick start guide
   - Architecture diagram
   - Component links
   - Setup instructions
   - Contributing guide

2. **docs/ARCHITECTURE.md**
   - System design
   - Data flow
   - Component interactions
   - Database schema
   - API design

3. **docs/DECISIONS.md (ADR)**
   ```markdown
   # Architecture Decision Records
   
   ## ADR-001: Clique PoA Consensus
   
   **Status:** Accepted
   
   **Context:**
   Need consensus mechanism for private testnet.
   
   **Decision:**
   Use Clique PoA with 2 signer nodes.
   
   **Consequences:**
   - ✅ Low resource usage
   - ✅ Fast block times
   - ✅ Industry standard for dev networks
   - ❌ Centralized (acceptable for dev)
   
   **Alternatives Considered:**
   - PoW: Wasteful, slow
   - PoS: Complex, overkill
   
   ---
   
   ## ADR-002: Indexed Block Explorer
   
   **Status:** Accepted
   
   **Context:**
   Need fast queries for address history, token transfers.
   
   **Decision:**
   Index blockchain data into PostgreSQL.
   
   **Consequences:**
   - ✅ Fast complex queries
   - ✅ Production-realistic pattern
   - ❌ More complex than direct RPC
   - ❌ Need to handle reorgs
   
   **Alternatives Considered:**
   - Direct RPC: Too slow for complex queries
   
   ---
   
   ## ADR-003: Single Frontend Application
   
   **Status:** Accepted
   
   **Context:**
   Need UIs for explorer, wallet, and dApps.
   
   **Decision:**
   Single Next.js app with route groups.
   
   **Consequences:**
   - ✅ Shared components
   - ✅ Single deployment
   - ✅ Easier to maintain
   - ✅ Better UX (one site, no context switching)
   
   **Alternatives Considered:**
   - 3 separate apps: More complex, duplicate code
   
   [Continue for all major decisions...]
   ```

4. **docs/SECURITY.md**
   - Wallet security warnings
   - Smart contract security
   - API security
   - Production recommendations

5. **docs/API.md**
   - All endpoints documented
   - Request/response examples
   - Error codes
   - Rate limits

6. **docs/DEPLOYMENT.md**
   - Deploy frontend (Vercel)
   - Deploy backend (optional VPS)
   - Environment variables
   - Troubleshooting

7. **Component READMEs**
   - blockchain/README.md
   - backend/README.md
   - frontend/README.md
   - contracts/README.md

**Testing:**
- All docs accurate
- Links work
- Code examples correct

### Phase 8E: Video Demo

**Tasks:**

1. **Write script** (5-10 minutes)
   - Introduction (30s)
   - Show architecture diagram (1min)
   - Demo testnet running (1min)
   - Demo block explorer (2min)
   - Demo wallet (2min)
   - Demo smart contracts + dApps (2min)
   - Switch networks (30s)
   - Conclusion (30s)

2. **Record demo**
   - Use OBS or Loom
   - Clear audio
   - Screen at 1920x1080

3. **Edit**
   - Trim mistakes
   - Add captions (important)
   - Add intro/outro

4. **Upload to YouTube**
   - Public or unlisted
   - Add to README

**Testing:**
- Video flows well
- Audio clear
- Demonstrates key features

### Phase 8F: Blog Post

**Tasks:**

1. **Write technical blog post** (1500-2000 words)
   - Why I built this
   - Architecture overview
   - Key technical challenges
   - Interesting solutions
   - What I learned
   - Code snippets
   - Diagrams

2. **Publish**
   - Dev.to
   - Medium
   - Personal blog
   - LinkedIn

3. **Share**
   - Twitter
   - LinkedIn
   - Reddit (r/ethereum, r/ethdev)

**Testing:**
- Post readable
- Code snippets correct
- Links work

### Phase 8G: Deployment

**Tasks:**

1. **Deploy frontend to Vercel**
   - Connect GitHub repo
   - Configure env vars
   - Deploy
   - Test deployed version

2. **(Optional) Deploy backend**
   - Railway, Render, or DigitalOcean
   - Setup PostgreSQL
   - Configure env vars
   - Run migrations
   - Test endpoints

3. **Update README with live links**

**Testing:**
- Frontend loads
- API works (if deployed)
- Can connect MetaMask to testnet

**Deliverables:**
- ✅ All tests passing
- ✅ >80% backend coverage
- ✅ >95% contract coverage
- ✅ CI/CD workflows working
- ✅ Comprehensive documentation
- ✅ Architecture Decision Records
- ✅ Video demo published
- ✅ Blog post published
- ✅ Frontend deployed
- ✅ GitHub repo public and polished

---

## VIII. FUTURE ENHANCEMENTS (Phase 9+)

**These are optional but make it even more impressive:**

### 1. Rate Limiting (Redis)
- Add Redis to docker-compose
- Implement rate limiting middleware
- Track requests per IP/API key

### 2. Caching Layer
- Cache frequently accessed data
- Reduce database load
- Faster API responses

### 3. Load Balancer (Nginx)
- Multiple API instances
- Reverse proxy
- SSL termination

### 4. Monitoring (Prometheus + Grafana)
- Metrics collection
- Dashboards for:
  - Sync status
  - API performance
  - Error rates

### 5. Advanced Contract Features
- Governance token
- Timelock for admin actions
- Upgradeability (proxy pattern)

### 6. Multi-Chain Deployment
- Deploy contracts to Sepolia
- Deploy contracts to Polygon Mumbai
- Show dApps work on all chains

### 7. Advanced dApp Features
- Chart.js for price charts (DEX)
- Historical APY charts (Staking)
- Transaction simulation
- Gas estimation improvements

### 8. Security Enhancements
- API authentication (JWT)
- Admin dashboard with auth
- DDoS protection
- Input sanitization

---

## IX. SUCCESS CRITERIA

**You'll know you've succeeded when:**

### Technical ✅
- [ ] Private network produces blocks reliably
- [ ] Indexer syncs without errors for 24+ hours
- [ ] API responds <100ms for simple queries
- [ ] Explorer shows real-time updates
- [ ] Wallet can send transactions on multiple chains
- [ ] Smart contracts have >95% test coverage
- [ ] All dApps interact with contracts successfully
- [ ] No memory leaks in backend (tested with profiler)
- [ ] Frontend performs well (Lighthouse score >90)

### Quality ✅
- [ ] >80% backend test coverage
- [ ] >95% contract test coverage
- [ ] All CI workflows passing
- [ ] No critical/high security findings
- [ ] Code formatted consistently
- [ ] Comprehensive documentation
- [ ] Architecture decisions documented
- [ ] Security considerations documented

### Portfolio ✅
- [ ] GitHub repo has professional README
- [ ] Video demo is clear and impressive
- [ ] Blog post demonstrates deep understanding
- [ ] Code is clean, commented, and organized
- [ ] Can explain every architecture decision
- [ ] Deployed and accessible
- [ ] Non-technical person can run locally

### Interview-Ready ✅
- [ ] Can walk through code confidently
- [ ] Can explain trade-offs made
- [ ] Can discuss alternative approaches
- [ ] Can identify areas for improvement
- [ ] Demonstrates blockchain fundamentals
- [ ] Demonstrates full-stack engineering
- [ ] Demonstrates security awareness
- [ ] Shows testing best practices

---

## X. FINAL NOTES

### This Blueprint is Your Guide

**Use it as:**
- Reference when stuck
- Checklist for completeness
- Interview prep (explain decisions)
- Template for documentation

**Don't:**
- Follow it blindly
- Skip testing
- Rush through security
- Ignore best practices

### Key Reminders

1. **Security First** - Especially wallet implementation
2. **Test Everything** - >80% backend, >95% contracts
3. **Document Decisions** - ADRs are interview gold
4. **No Timeline Pressure** - Quality > Speed
5. **Ask Questions** - When stuck, research or ask

### What Makes This Special

Most candidates have:
- Followed tutorials
- Deployed a simple contract
- Maybe used Hardhat

You'll have:
- Built ENTIRE ecosystem from scratch
- Production-architecture code
- Deep understanding of every layer
- Multi-chain support
- Comprehensive testing
- Security considerations documented
- Deployed and running
- Video demo
- Blog post explaining it all

**This is 10x more impressive than a tutorial project.**
