
-- ============================================================================

-- WALLET TABLES (⚠️ DEV ONLY - See docs/SECURITY.md)

-- ============================================================================

-- Enable UUID extension if not already enabled

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Wallets table

CREATE TABLE wallets (

    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    address VARCHAR(42) UNIQUE NOT NULL,

    encrypted_private_key TEXT NOT NULL,

    encryption_nonce VARCHAR(32) NOT NULL,

    encryption_tag VARCHAR(32) NOT NULL,

    encryption_salt VARCHAR(64) NOT NULL,

    name VARCHAR(100),

    derivation_path VARCHAR(100),

    is_imported BOOLEAN DEFAULT false,

    created_at TIMESTAMP DEFAULT NOW(),

    updated_at TIMESTAMP DEFAULT NOW()

);

CREATE INDEX idx_wallets_address ON wallets(address);

CREATE INDEX idx_wallets_created ON wallets(created_at DESC);

-- Wallet balances per chain

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

-- Wallet transactions (cache of sent transactions)

CREATE TABLE wallet_transactions (

    id BIGSERIAL PRIMARY KEY,

    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,

    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,

    tx_hash VARCHAR(66) NOT NULL,

    direction VARCHAR(10) NOT NULL,

    from_address VARCHAR(42) NOT NULL,

    to_address VARCHAR(42),

    value NUMERIC(78, 0) NOT NULL,

    gas_used BIGINT,

    gas_price NUMERIC(78, 0),

    status INT NOT NULL,

    block_number BIGINT,

    timestamp TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW(),

    

    CHECK (direction IN ('sent', 'received'))

);

CREATE INDEX idx_wallet_tx_wallet_chain ON wallet_transactions(wallet_id, chain_id, timestamp DESC);

CREATE INDEX idx_wallet_tx_hash ON wallet_transactions(chain_id, tx_hash);

CREATE INDEX idx_wallet_tx_status ON wallet_transactions(wallet_id, status);

-- Comments

COMMENT ON TABLE wallets IS '⚠️ DEV ONLY: Encrypted wallet keys (NOT for production)';

COMMENT ON COLUMN wallets.encrypted_private_key IS 'AES-256-GCM encrypted private key';

COMMENT ON COLUMN wallets.encryption_nonce IS 'Unique nonce (IV) for AES-GCM';

COMMENT ON COLUMN wallets.encryption_tag IS 'Authentication tag for AES-GCM';

COMMENT ON COLUMN wallets.encryption_salt IS 'Salt for Argon2id key derivation';

