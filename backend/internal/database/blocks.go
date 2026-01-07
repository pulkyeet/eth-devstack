package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulkyeet/eth-devstack/backend/internal/models"
)

func (db *DB) InsertBlock(ctx context.Context, block *models.Block) error {
	query := `
		INSERT INTO blocks (
			chain_id, block_number, hash, parent_hash, nonce, sha3_uncles,
			miner, state_root, transactions_root, receipts_root,
			difficulty, total_difficulty, size, gas_limit, gas_used,
			timestamp, extra_data, mix_hash, base_fee_per_gas, tx_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (chain_id, block_number) DO UPDATE SET
			hash = EXCLUDED.hash,
			parent_hash = EXCLUDED.parent_hash,
			miner = EXCLUDED.miner,
			gas_used = EXCLUDED.gas_used,
			timestamp = EXCLUDED.timestamp,
			tx_count = EXCLUDED.tx_count
		RETURNING id
	`

	err := db.conn.QueryRowContext(ctx, query,
		block.ChainID, block.BlockNumber, block.Hash, block.ParentHash,
		block.Nonce, block.Sha3Uncles, block.Miner, block.StateRoot,
		block.TransactionsRoot, block.ReceiptsRoot, block.Difficulty,
		block.TotalDifficulty, block.Size, block.GasLimit, block.GasUsed,
		block.Timestamp, block.ExtraData, block.MixHash, block.BaseFeePerGas,
		block.TxCount,
	).Scan(&block.ID)

	if err != nil {
		return fmt.Errorf("failed to insert block: %w", err)
	}

	return nil
}

func (db *DB) GetBlockByNumber(ctx context.Context, chainID, blockNumber int64) (*models.Block, error) {
	query := `
		SELECT id, chain_id, block_number, hash, parent_hash, nonce, sha3_uncles,
			   miner, state_root, transactions_root, receipts_root,
			   difficulty, total_difficulty, size, gas_limit, gas_used,
			   timestamp, extra_data, mix_hash, base_fee_per_gas, tx_count, created_at
		FROM blocks
		WHERE chain_id = $1 AND block_number = $2
	`

	block := &models.Block{}
	err := db.conn.QueryRowContext(ctx, query, chainID, blockNumber).Scan(
		&block.ID, &block.ChainID, &block.BlockNumber, &block.Hash,
		&block.ParentHash, &block.Nonce, &block.Sha3Uncles, &block.Miner,
		&block.StateRoot, &block.TransactionsRoot, &block.ReceiptsRoot,
		&block.Difficulty, &block.TotalDifficulty, &block.Size,
		&block.GasLimit, &block.GasUsed, &block.Timestamp, &block.ExtraData,
		&block.MixHash, &block.BaseFeePerGas, &block.TxCount, &block.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}

	return block, nil
}

func (db *DB) GetBlockByHash(ctx context.Context, chainID int64, hash string) (*models.Block, error) {
	query := `
		SELECT id, chain_id, block_number, hash, parent_hash, nonce, sha3_uncles,
			   miner, state_root, transactions_root, receipts_root,
			   difficulty, total_difficulty, size, gas_limit, gas_used,
			   timestamp, extra_data, mix_hash, base_fee_per_gas, tx_count, created_at
		FROM blocks
		WHERE chain_id = $1 AND hash = $2
	`

	block := &models.Block{}
	err := db.conn.QueryRowContext(ctx, query, chainID, hash).Scan(
		&block.ID, &block.ChainID, &block.BlockNumber, &block.Hash,
		&block.ParentHash, &block.Nonce, &block.Sha3Uncles, &block.Miner,
		&block.StateRoot, &block.TransactionsRoot, &block.ReceiptsRoot,
		&block.Difficulty, &block.TotalDifficulty, &block.Size,
		&block.GasLimit, &block.GasUsed, &block.Timestamp, &block.ExtraData,
		&block.MixHash, &block.BaseFeePerGas, &block.TxCount, &block.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}

	return block, nil
}

func (db *DB) GetLatestBlock(ctx context.Context, chainID int64) (*models.Block, error) {
	query := `
		SELECT id, chain_id, block_number, hash, parent_hash, nonce, sha3_uncles,
			   miner, state_root, transactions_root, receipts_root,
			   difficulty, total_difficulty, size, gas_limit, gas_used,
			   timestamp, extra_data, mix_hash, base_fee_per_gas, tx_count, created_at
		FROM blocks
		WHERE chain_id = $1
		ORDER BY block_number DESC
		LIMIT 1
	`

	block := &models.Block{}
	err := db.conn.QueryRowContext(ctx, query, chainID).Scan(
		&block.ID, &block.ChainID, &block.BlockNumber, &block.Hash,
		&block.ParentHash, &block.Nonce, &block.Sha3Uncles, &block.Miner,
		&block.StateRoot, &block.TransactionsRoot, &block.ReceiptsRoot,
		&block.Difficulty, &block.TotalDifficulty, &block.Size,
		&block.GasLimit, &block.GasUsed, &block.Timestamp, &block.ExtraData,
		&block.MixHash, &block.BaseFeePerGas, &block.TxCount, &block.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	return block, nil
}

func (db *DB) GetBlocks(ctx context.Context, chainID int64, limit, offset int) ([]*models.Block, error) {
	query := `
		SELECT id, chain_id, block_number, hash, parent_hash, nonce, sha3_uncles,
			   miner, state_root, transactions_root, receipts_root,
			   difficulty, total_difficulty, size, gas_limit, gas_used,
			   timestamp, extra_data, mix_hash, base_fee_per_gas, tx_count, created_at
		FROM blocks
		WHERE chain_id = $1
		ORDER BY block_number DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.QueryContext(ctx, query, chainID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocks: %w", err)
	}
	defer rows.Close()

	var blocks []*models.Block
	for rows.Next() {
		block := &models.Block{}
		err := rows.Scan(
			&block.ID, &block.ChainID, &block.BlockNumber, &block.Hash,
			&block.ParentHash, &block.Nonce, &block.Sha3Uncles, &block.Miner,
			&block.StateRoot, &block.TransactionsRoot, &block.ReceiptsRoot,
			&block.Difficulty, &block.TotalDifficulty, &block.Size,
			&block.GasLimit, &block.GasUsed, &block.Timestamp, &block.ExtraData,
			&block.MixHash, &block.BaseFeePerGas, &block.TxCount, &block.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan block: %w", err)
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (db *DB) DeleteBlocksFromHeight(ctx context.Context, chainID, blockNumber int64) error {
	query := `DELETE FROM blocks WHERE chain_id = $1 AND block_number >= $2`
	_, err := db.conn.ExecContext(ctx, query, chainID, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to delete blocks: %w", err)
	}
	return nil
}
