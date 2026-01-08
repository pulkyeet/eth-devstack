package database

import (
	"context"
	"fmt"
	"database/sql"
	"github.com/pulkyeet/eth-devstack/backend/internal/models"
)

func (db *DB) GetChains(ctx context.Context) ([]*models.Chain, error) {
	query := `
		SELECT chain_id, name, short_name, native_symbol, rpc_endpoint, ws_endpoint,
			   block_time_seconds, is_active, is_testnet, last_indexed_block, icon_url,
			   explorer_url, created_at, updated_at
		FROM chains
		ORDER BY chain_id
	`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get chains: %w", err)
	}
	defer rows.Close()

	var chains []*models.Chain
	for rows.Next() {
		chain := &models.Chain{}
		err := rows.Scan(
			&chain.ChainID, &chain.Name, &chain.ShortName, &chain.NativeSymbol,
			&chain.RPCEndpoint, &chain.WSEndpoint, &chain.BlockTimeSeconds,
			&chain.IsActive, &chain.IsTestnet, &chain.LastIndexedBlock,
			&chain.IconURL, &chain.ExplorerURL, &chain.CreatedAt, &chain.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chain: %w", err)
		}
		chains = append(chains, chain)
	}
	return chains, nil
}

func (db *DB) GetChain(ctx context.Context, chainID int64) (*models.Chain, error) {
	query := `
		SELECT chain_id, name, short_name, native_symbol, rpc_endpoint, ws_endpoint,
			block_time_seconds, is_active, is_testnet, last_indexed_block, icon_url,
			explorer_url, created_at, updated_at
		FROM chains
		WHERE chain_id = $1
	`
	var chain models.Chain
	err := db.conn.QueryRowContext(ctx, query, chainID).Scan(
		&chain.ChainID, &chain.Name, &chain.ShortName, &chain.NativeSymbol,
		&chain.RPCEndpoint, &chain.WSEndpoint, &chain.BlockTimeSeconds,
		&chain.IsActive, &chain.IsTestnet, &chain.LastIndexedBlock,
		&chain.IconURL, &chain.ExplorerURL, &chain.CreatedAt, &chain.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &chain, err
}