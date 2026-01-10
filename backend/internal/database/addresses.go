package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulkyeet/eth-devstack/backend/internal/models"
)

func (db *DB) GetAddress(ctx context.Context, chainID int64, address string) (*models.Address, error) {
	query := `SELECT id, chain_id, address, balance, nonce, is_contract, contract_creator,
			   creation_tx_hash, code_hash, tx_count, first_seen_block, last_seen_block,
			   first_seen_at, last_seen_at, created_at, updated_at
			   FROM addresses WHERE chain_id = $1 AND address = $2`
	
	addr := &models.Address{}
	err := db.conn.QueryRowContext(ctx, query, chainID, address).Scan(
		&addr.ID, &addr.ChainID, &addr.Address, &addr.Balance, &addr.Nonce,
		&addr.IsContract, &addr.ContractCreator, &addr.CreationTxHash,
		&addr.CodeHash, &addr.TxCount, &addr.FirstSeenBlock, &addr.LastSeenBlock,
		&addr.FirstSeenAt, &addr.LastSeenAt, &addr.CreatedAt, &addr.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err!=nil {
		return nil, fmt.Errorf("Failed to get address: %w", err)
	}
	return addr, nil
}

func (db *DB) UpsertAddress(ctx context.Context, addr *models.Address) error {
	query := `
		INSERT INTO addresses (
			chain_id, address, balance, nonce, is_contract,
			contract_creator, creation_tx_hash, code_hash,
			first_seen_block, last_seen_block, first_seen_at, last_seen_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (chain_id, address) DO UPDATE SET
			balance = EXCLUDED.balance,
			nonce = EXCLUDED.nonce,
			last_seen_block = EXCLUDED.last_seen_block,
			last_seen_at = EXCLUDED.last_seen_at,
			tx_count = addresses.tx_count + 1,
			updated_at = NOW()
		RETURNING id
	`
	err := db.conn.QueryRowContext(ctx, query,
		addr.ChainID, addr.Address, addr.Balance, addr.Nonce, addr.IsContract,
		addr.ContractCreator, addr.CreationTxHash, addr.CodeHash,
		addr.FirstSeenBlock, addr.LastSeenBlock, addr.FirstSeenAt, addr.LastSeenAt,
	).Scan(&addr.ID)
	return err
}

func (db *DB) IncrementAddressTxCount(ctx context.Context, chainID int64, address string) error {
	query := `
		UPDATE addresses
		SET tx_count = tx_count + 1, updated_at = NOW()
		WHERE chain_id = $1 AND address = $2
	`
	_, err := db.conn.ExecContext(ctx, query, chainID, address)
	return err
}