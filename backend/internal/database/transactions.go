package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulkyeet/eth-devstack/backend/internal/models"
)

func (db *DB) InsertTransaction(ctx context.Context, tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			chain_id, hash, block_number, block_hash, transaction_index,
			from_address, to_address, value, gas, gas_price,
			max_fee_per_gas, max_priority_fee_per_gas, input, nonce,
			transaction_type, status, gas_used, cumulative_gas_used,
			effective_gas_price, contract_address, logs_bloom, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
		ON CONFLICT (chain_id, hash) DO UPDATE SET
			status = EXCLUDED.status,
			gas_used = EXCLUDED.gas_used,
			cumulative_gas_used = EXCLUDED.cumulative_gas_used,
			effective_gas_price = EXCLUDED.effective_gas_price
		RETURNING id
	`

	err := db.conn.QueryRowContext(ctx, query,
		tx.ChainID, tx.Hash, tx.BlockNumber, tx.BlockHash, tx.TransactionIndex,
		tx.FromAddress, tx.ToAddress, tx.Value, tx.Gas, tx.GasPrice,
		tx.MaxFeePerGas, tx.MaxPriorityFeePerGas, tx.Input, tx.Nonce,
		tx.TransactionType, tx.Status, tx.GasUsed, tx.CumulativeGasUsed,
		tx.EffectiveGasPrice, tx.ContractAddress, tx.LogsBloom, tx.Timestamp,
	).Scan(&tx.ID)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}

func (db *DB) GetTransactionByHash(ctx context.Context, chainID int64, hash string) (*models.Transaction, error) {
	query := `
		SELECT id, chain_id, hash, block_number, block_hash, transaction_index,
			   from_address, to_address, value, gas, gas_price,
			   max_fee_per_gas, max_priority_fee_per_gas, input, nonce,
			   transaction_type, status, gas_used, cumulative_gas_used,
			   effective_gas_price, contract_address, logs_bloom, timestamp, created_at
		FROM transactions
		WHERE chain_id = $1 AND hash = $2
	`

	tx := &models.Transaction{}
	err := db.conn.QueryRowContext(ctx, query, chainID, hash).Scan(
		&tx.ID, &tx.ChainID, &tx.Hash, &tx.BlockNumber, &tx.BlockHash,
		&tx.TransactionIndex, &tx.FromAddress, &tx.ToAddress, &tx.Value,
		&tx.Gas, &tx.GasPrice, &tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas,
		&tx.Input, &tx.Nonce, &tx.TransactionType, &tx.Status, &tx.GasUsed,
		&tx.CumulativeGasUsed, &tx.EffectiveGasPrice, &tx.ContractAddress,
		&tx.LogsBloom, &tx.Timestamp, &tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return tx, nil
}

func (db *DB) GetTransactionsByBlock(ctx context.Context, chainID, blockNumber int64) ([]*models.Transaction, error) {
	query := `
		SELECT id, chain_id, hash, block_number, block_hash, transaction_index,
			   from_address, to_address, value, gas, gas_price,
			   max_fee_per_gas, max_priority_fee_per_gas, input, nonce,
			   transaction_type, status, gas_used, cumulative_gas_used,
			   effective_gas_price, contract_address, logs_bloom, timestamp, created_at
		FROM transactions
		WHERE chain_id = $1 AND block_number = $2
		ORDER BY transaction_index ASC
	`

	rows, err := db.conn.QueryContext(ctx, query, chainID, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var txs []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.ChainID, &tx.Hash, &tx.BlockNumber, &tx.BlockHash,
			&tx.TransactionIndex, &tx.FromAddress, &tx.ToAddress, &tx.Value,
			&tx.Gas, &tx.GasPrice, &tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas,
			&tx.Input, &tx.Nonce, &tx.TransactionType, &tx.Status, &tx.GasUsed,
			&tx.CumulativeGasUsed, &tx.EffectiveGasPrice, &tx.ContractAddress,
			&tx.LogsBloom, &tx.Timestamp, &tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

func (db *DB) DeleteTransactionsFromBlock(ctx context.Context, chainID, blockNumber int64) error {
	query := `DELETE FROM transactions WHERE chain_id = $1 AND block_number >= $2`
	_, err := db.conn.ExecContext(ctx, query, chainID, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to delete transactions: %w", err)
	}
	return nil
}

func (db *DB) GetTransactions(ctx context.Context, chainID int64, limit, offset int) ([]*models.Transaction, error) {
	query := `
		SELECT id, chain_id, hash, block_number, block_hash, transaction_index,
			   from_address, to_address, value, gas, gas_price,
			   max_fee_per_gas, max_priority_fee_per_gas, input, nonce,
			   transaction_type, status, gas_used, cumulative_gas_used,
			   effective_gas_price, contract_address, logs_bloom, timestamp, created_at
		FROM transactions
		WHERE chain_id = $1
		ORDER BY block_number DESC, transaction_index DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.QueryContext(ctx, query, chainID, limit, offset)
	if err!=nil {
		return nil, fmt.Errorf("Failed to get transactions: %w", err)
	}
	defer rows.Close()

	var txs []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.ChainID, &tx.Hash, &tx.BlockNumber, &tx.BlockHash,
			&tx.TransactionIndex, &tx.FromAddress, &tx.ToAddress, &tx.Value,
			&tx.Gas, &tx.GasPrice, &tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas,
			&tx.Input, &tx.Nonce, &tx.TransactionType, &tx.Status, &tx.GasUsed,
			&tx.CumulativeGasUsed, &tx.EffectiveGasPrice, &tx.ContractAddress,
			&tx.LogsBloom, &tx.Timestamp, &tx.CreatedAt,
		)
		if err!=nil {
			return nil, fmt.Errorf("Failed to scan transaction: %w", err)
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (db *DB) CountTransactions(ctx context.Context, chainID int64) (int64, error) {
	var count int64
	err := db.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE chain_id = $1`, chainID).Scan(&count)
	return count, err
}

func (db *DB) GetTransactionsByAddress(ctx context.Context, chainID int64, address string, limit, offset int) ([]*models.Transaction, error) {
	query := `
		SELECT id, chain_id, hash, block_number, block_hash, transaction_index,
			   from_address, to_address, value, gas, gas_price,
			   max_fee_per_gas, max_priority_fee_per_gas, input, nonce,
			   transaction_type, status, gas_used, cumulative_gas_used,
			   effective_gas_price, contract_address, logs_bloom, timestamp, created_at
		FROM transactions
		WHERE chain_id = $1 AND (from_address = $2 OR to_address = $2)
		ORDER BY block_number DESC, transaction_index DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := db.conn.QueryContext(ctx, query, chainID, address, limit, offset)
	if err!=nil {
		return nil, fmt.Errorf("Failed to get transactions: %w", err)
	}
	defer rows.Close()

	var txs []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.ChainID, &tx.Hash, &tx.BlockNumber, &tx.BlockHash,
			&tx.TransactionIndex, &tx.FromAddress, &tx.ToAddress, &tx.Value,
			&tx.Gas, &tx.GasPrice, &tx.MaxFeePerGas, &tx.MaxPriorityFeePerGas,
			&tx.Input, &tx.Nonce, &tx.TransactionType, &tx.Status, &tx.GasUsed,
			&tx.CumulativeGasUsed, &tx.EffectiveGasPrice, &tx.ContractAddress,
			&tx.LogsBloom, &tx.Timestamp, &tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		txs = append(txs, tx)
	}
	return txs, nil
}