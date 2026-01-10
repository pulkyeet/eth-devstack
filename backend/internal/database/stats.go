package database

import (
	"context"
	"fmt"
)

type NetworkStats struct {
	LatestBlock int64 `json:"latest_block"`
	TotalTransactions int64 `json:"total_transactions"`
	TotalAddresses int64 `json:"total_addresses"`
	TotalTokens int64 `json:"total_tokens"`
	AvgBlockTime string `json:"avg_block_time"`
	TPS24h string `json:"tps_24h"`
}

func (db *DB) GetNetworkStats(ctx context.Context, chainID int64) (*NetworkStats, error) {
	stats := &NetworkStats{}

	//Get latest block
	err := db.conn.QueryRowContext(ctx, `SELECT COALESCE(MAX(block_number), 0) FROM blocks WHERE chain_id = $1`, chainID).Scan(&stats.LatestBlock)
	if err!=nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	// Get total transactions
	err = db.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE chain_id = $1`, chainID).Scan(&stats.TotalTransactions)
	if err!=nil {
		return nil, err
	}

	//Get total addresses
	err = db.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM addresses WHERE chain_id = $1`, chainID).Scan(&stats.TotalAddresses)
	if err!=nil {
		return nil, err
	}

	//Get total tokens
	err = db.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM tokens WHERE chain_id = $1`, chainID).Scan(&stats.TotalTokens)
	if err!=nil {
		return nil, err
	}

	err = db.conn.QueryRowContext(ctx, `
		SELECT COALESCE(
			EXTRACT(EPOCH FROM (MAX(timestamp) - MIN(timestamp))) / NULLIF(COUNT(*) - 1, 0),
			0
		)::text || 's'
		FROM (
			SELECT timestamp FROM blocks
			WHERE chain_id = $1
			ORDER BY block_number DESC
			LIMIT 100
		) recent
	`, chainID).Scan(&stats.AvgBlockTime)
	if err != nil {
		stats.AvgBlockTime = "0s"
	}

	// TPS last 24h
	err = db.conn.QueryRowContext(ctx, `
		SELECT COALESCE(
			COUNT(*)::float / NULLIF(EXTRACT(EPOCH FROM (MAX(timestamp) - MIN(timestamp))), 0),
			0
		)::text
		FROM transactions
		WHERE chain_id = $1 AND timestamp > NOW() - INTERVAL '24 hours'
	`, chainID).Scan(&stats.TPS24h)
	if err != nil {
		stats.TPS24h = "0"
	}

	return stats, nil
}