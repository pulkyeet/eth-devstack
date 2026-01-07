-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Chains
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

-- Blocks
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
    base_fee_per_gas NUMERIC(78, 0),
    tx_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, block_number),
    UNIQUE(chain_id, hash)
);

CREATE INDEX idx_blocks_chain_number ON blocks(chain_id, block_number DESC);
CREATE INDEX idx_blocks_chain_timestamp ON blocks(chain_id, timestamp DESC);
CREATE INDEX idx_blocks_chain_miner ON blocks(chain_id, miner);
CREATE INDEX idx_blocks_hash ON blocks(hash);

-- Transactions
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    hash VARCHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    block_hash VARCHAR(66) NOT NULL,
    transaction_index INT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    value NUMERIC(78, 0) NOT NULL DEFAULT 0,
    gas BIGINT NOT NULL,
    gas_price NUMERIC(78, 0),
    max_fee_per_gas NUMERIC(78, 0),
    max_priority_fee_per_gas NUMERIC(78, 0),
    input TEXT,
    nonce BIGINT NOT NULL,
    transaction_type INT DEFAULT 0,
    status INT,
    gas_used BIGINT,
    cumulative_gas_used BIGINT,
    effective_gas_price NUMERIC(78, 0),
    contract_address VARCHAR(42),
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

-- Transaction Logs
CREATE TABLE transaction_logs (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    transaction_hash VARCHAR(66) NOT NULL,
    log_index INT NOT NULL,
    address VARCHAR(42) NOT NULL,
    data TEXT,
    topic0 VARCHAR(66),
    topic1 VARCHAR(66),
    topic2 VARCHAR(66),
    topic3 VARCHAR(66),
    block_number BIGINT NOT NULL,
    block_hash VARCHAR(66) NOT NULL,
    transaction_index INT NOT NULL,
    removed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(chain_id, transaction_hash, log_index)
);

CREATE INDEX idx_logs_chain_address ON transaction_logs(chain_id, address);
CREATE INDEX idx_logs_chain_topic0 ON transaction_logs(chain_id, topic0);
CREATE INDEX idx_logs_chain_block ON transaction_logs(chain_id, block_number DESC);
CREATE INDEX idx_logs_topics ON transaction_logs(chain_id, topic0, topic1) WHERE topic1 IS NOT NULL;

-- Addresses
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

-- Sync Status
CREATE TABLE sync_status (
    chain_id BIGINT PRIMARY KEY REFERENCES chains(chain_id) ON DELETE CASCADE,
    last_synced_block BIGINT NOT NULL DEFAULT 0,
    latest_block BIGINT NOT NULL DEFAULT 0,
    is_syncing BOOLEAN DEFAULT false,
    sync_rate FLOAT,
    last_sync_time TIMESTAMP,
    error_count INT DEFAULT 0,
    last_error TEXT,
    last_error_time TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sync_status_syncing ON sync_status(is_syncing) WHERE is_syncing = true;

-- Reorgs
CREATE TABLE reorgs (
    id SERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    old_block_number BIGINT NOT NULL,
    old_block_hash VARCHAR(66) NOT NULL,
    new_block_hash VARCHAR(66) NOT NULL,
    depth INT NOT NULL,
    detected_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reorgs_chain ON reorgs(chain_id, detected_at DESC);

-- Initial chain data
INSERT INTO chains (chain_id, name, short_name, native_symbol, rpc_endpoint, block_time_seconds, is_testnet, is_active) VALUES
(1337, 'Local Testnet', 'local', 'ETH', 'http://localhost:8545', 15, true, true);

-- Triggers
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