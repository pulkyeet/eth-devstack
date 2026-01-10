package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulkyeet/eth-devstack/backend/internal/models"
)

func (db *DB) UpsertToken(ctx context.Context, token *models.Token) error {
	query := `INSERT INTO tokens (chain_id, address, type, name, symbol, decimals, total_supply) VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (chain_id, address) DO UPDATE SET
	name = COALESCE(EXCLUDED.name, tokens.name),
	symbol = COALESCE(EXCLUDED.symbol, tokens.symbol),
			decimals = COALESCE(EXCLUDED.decimals, tokens.decimals),
			total_supply = COALESCE(EXCLUDED.total_supply, tokens.total_supply),
			updated_at = NOW()
		RETURNING id`

	err := db.conn.QueryRowContext(ctx, query,
		token.ChainID, token.Address, token.Type, token.Name,
		token.Symbol, token.Decimals, token.TotalSupply,
	).Scan(&token.ID)
	return err
}

func (db *DB) InsertTokenTransfer(ctx context.Context, transfer *models.TokenTransfer) error {
	query := `
		INSERT INTO token_transfers (
			chain_id, transaction_hash, log_index, token_address,
			from_address, to_address, value, token_id, block_number, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (chain_id, transaction_hash, log_index) DO NOTHING
		RETURNING id
	`
	err := db.conn.QueryRowContext(ctx, query,
		transfer.ChainID, transfer.TransactionHash, transfer.LogIndex,
		transfer.TokenAddress, transfer.FromAddress, transfer.ToAddress,
		transfer.Value, transfer.TokenID, transfer.BlockNumber, transfer.Timestamp,
	).Scan(&transfer.ID)

	if err != nil && err.Error()!="no rows in result set" {
		return fmt.Errorf("failed to insert token transfer: %w", err)
	}
	return nil
}

func (db *DB) GetTokensByAddress(ctx context.Context, chainID int64, address string) ([]*models.Token, error) {
	query := `SELECT DISTINCT t.id, t.chain_id, t.address, t.type, t.name, t.symbol, t.decimals, t.total_supply, t.holder_count, t.transfer_count, t.created_at, t.updated_at
	FROM tokens t INNER JOIN token_balances ON t.chain_id = tb.chain_id AND t.address = tb.token_address
	WHERE t.chain_id = $1 AND tb.holder_address = $2 AND tb.balance!=0`

	rows, err := db.conn.QueryContext(ctx, query, chainID, address)
	if err!=nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*models.Token
	for rows.Next() {
		token := &models.Token{}
		err := rows.Scan(
			&token.ID, &token.ChainID, &token.Address, &token.Type,
			&token.Name, &token.Symbol, &token.Decimals, &token.TotalSupply,
			&token.HolderCount, &token.TransferCount, &token.CreatedAt, &token.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (db *DB) UpsertTokenBalance(ctx context.Context, balance *models.TokenBalance) error {
	query := `
		INSERT INTO token_balances (chain_id, token_address, holder_address, balance)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (chain_id, token_address, holder_address) DO UPDATE SET
			balance = EXCLUDED.balance,
			updated_at = NOW()
		RETURNING id
	`
	err := db.conn.QueryRowContext(ctx, query,
		balance.ChainID, balance.TokenAddress, balance.HolderAddress, balance.Balance,
	).Scan(&balance.ID)
	return err
}

func (db *DB) GetTokenBalance(ctx context.Context, chainID int64, tokenAddress, holderAddress string) (*models.TokenBalance, error) {
	query := `
		SELECT id, chain_id, token_address, holder_address, balance, updated_at
		FROM token_balances
		WHERE chain_id = $1 AND token_address = $2 AND holder_address = $3
	`
	balance := &models.TokenBalance{}
	err := db.conn.QueryRowContext(ctx, query, chainID, tokenAddress, holderAddress).Scan(
		&balance.ID, &balance.ChainID, &balance.TokenAddress,
		&balance.HolderAddress, &balance.Balance, &balance.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return balance, err
}