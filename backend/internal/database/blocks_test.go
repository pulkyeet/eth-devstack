package database

import (
	"context"
	"testing"
	"time"
	"fmt"
	"github.com/pulkyeet/eth-devstack/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func getTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func setupTestDB(t *testing.T) *DB {
	// Use test database
	connStr := "postgresql://eth_user:eth_pass_dev_only@localhost:5432/ethereum_explorer_test?sslmode=disable"
	logger := getTestLogger()
	db, err := NewDB(connStr, 5, 2, logger)
	require.NoError(t, err)
	
	// Clean tables
	db.conn.Exec("TRUNCATE blocks, transactions, addresses, chains CASCADE")
	
	db.conn.Exec(`
		INSERT INTO chains (chain_id, name, short_name, native_symbol, rpc_endpoint, block_time_seconds, is_testnet, is_active)
		VALUES (1337, 'Local Testnet', 'local', 'ETH', 'http://localhost:8545', 15, true, true)
		ON CONFLICT (chain_id) DO NOTHING
	`)

	return db
}

func TestInsertAndGetBlock(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC()

	block := &models.Block{
		ChainID:     1337,
		BlockNumber: 1,
		Hash:        "0xabc123",
		ParentHash:  "0x000000",
		Miner:       "0xminer",
		GasLimit:    8000000,
		GasUsed:     100000,
		Timestamp:   now,
		TxCount:     5,
	}

	err := db.InsertBlock(ctx, block)
	require.NoError(t, err)
	assert.NotZero(t, block.ID)

	retrieved, err := db.GetBlockByNumber(ctx, 1337, 1)
	require.NoError(t, err)
	assert.Equal(t, block.Hash, retrieved.Hash)
	assert.Equal(t, block.BlockNumber, retrieved.BlockNumber)
}

func TestGetBlocks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Insert 10 blocks
	for i := 1; i <= 10; i++ {
		block := &models.Block{
			ChainID:     1337,
			BlockNumber: int64(i),
			Hash:        fmt.Sprintf("0x%d", i),
			ParentHash:  "0x000000",
			Miner:       "0xminer",
			GasLimit:    8000000,
			GasUsed:     100000,
			Timestamp:   time.Now(),
			TxCount:     0,
		}
		db.InsertBlock(ctx, block)
	}

	blocks, err := db.GetBlocks(ctx, 1337, 5, 0)
	require.NoError(t, err)
	assert.Len(t, blocks, 5)
	
	if len(blocks) > 0 {
		assert.Equal(t, int64(10), blocks[0].BlockNumber) // DESC order
	}
}

func TestCountBlocks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		block := &models.Block{
			ChainID:     1337,
			BlockNumber: int64(i),
			Hash:        fmt.Sprintf("0x%d", i),
			ParentHash:  "0x000000",
			Miner:       "0xminer",
			GasLimit:    8000000,
			GasUsed:     100000,
			Timestamp:   time.Now(),
			TxCount:     0,
		}
		db.InsertBlock(ctx, block)
	}

	count, err := db.CountBlocks(ctx, 1337)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}