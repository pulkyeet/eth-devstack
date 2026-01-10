-- ============================================================================
-- TOKEN DATA (ERC20/ERC721/ERC1155)
-- ============================================================================

CREATE TABLE tokens (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    address VARCHAR(42) NOT NULL,
    type VARCHAR(20) NOT NULL,
    name VARCHAR(255),
    symbol VARCHAR(50),
    decimals INT,
    total_supply NUMERIC(78, 0),
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
    value NUMERIC(78, 0),
    token_id NUMERIC(78, 0),
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

-- Add trigger for tokens updated_at
CREATE TRIGGER update_tokens_updated_at BEFORE UPDATE ON tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();