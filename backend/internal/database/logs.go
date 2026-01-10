package database

import (
	"context"
	"fmt"
	"github.com/pulkyeet/eth-devstack/backend/internal/models"
)

func (db *DB) InsertLog(ctx context.Context, log *models.TransactionLog) error {
	query := `
		INSERT INTO transaction_logs (
			chain_id, transaction_hash, log_index, address, data,
			topic0, topic1, topic2, topic3, block_number, block_hash,
			transaction_index, removed
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (chain_id, transaction_hash, log_index) DO NOTHING
		RETURNING id
	`
	err := db.conn.QueryRowContext(ctx, query,
		log.ChainID, log.TransactionHash, log.LogIndex, log.Address, log.Data,
		log.Topic0, log.Topic1, log.Topic2, log.Topic3, log.BlockNumber,
		log.BlockHash, log.TransactionIndex, log.Removed,
	).Scan(&log.ID)

	if err !=nil && err.Error() != "no rows in result set" {
		return fmt.Errorf("failed to insert log: %w", err)
	}
	return nil
}

func (db *DB) GetLogsByTransaction(ctx context.Context, chainID int64, txHash string) ([]*models.TransactionLog, error) {
	query := `
		SELECT id, chain_id, transaction_hash, log_index, address, data,
			   topic0, topic1, topic2, topic3, block_number, block_hash,
			   transaction_index, removed, created_at
		FROM transaction_logs
		WHERE chain_id = $1 AND transaction_hash = $2
		ORDER BY log_index ASC
	`
	rows, err := db.conn.QueryContext(ctx, query, chainID, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.TransactionLog
	for rows.Next() {
		log := &models.TransactionLog{}
		err := rows.Scan(
			&log.ID, &log.ChainID, &log.TransactionHash, &log.LogIndex,
			&log.Address, &log.Data, &log.Topic0, &log.Topic1, &log.Topic2,
			&log.Topic3, &log.BlockNumber, &log.BlockHash, &log.TransactionIndex,
			&log.Removed, &log.CreatedAt,
		)
		if err!=nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

