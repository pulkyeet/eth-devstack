package database

import (
	"context"
	"testing"
	"time"

	"github.com/pulkyeet/eth-devstack/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC()

	tx := &models.Transaction{
		ChainID:          1337,
		Hash:             "0xtx123",
		BlockNumber:      1,
		BlockHash:        "0xblock123",
		TransactionIndex: 0,
		FromAddress:      "0xfrom",
		Value:            "1000000000000000000",
		Gas:              21000,
		Nonce:            0,
		TransactionType:  0,
		Timestamp:        now,
	}

	err := db.InsertTransaction(ctx, tx)
	require.NoError(t, err)
	assert.NotZero(t, tx.ID)
}

func TestGetTransactionByHash(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	tx := &models.Transaction{
		ChainID:          1337,
		Hash:             "0xtxhash",
		BlockNumber:      1,
		BlockHash:        "0xblock",
		TransactionIndex: 0,
		FromAddress:      "0xfrom",
		Value:            "1000",
		Gas:              21000,
		Nonce:            0,
		TransactionType:  0,
		Timestamp:        time.Now(),
	}
	db.InsertTransaction(ctx, tx)

	retrieved, err := db.GetTransactionByHash(ctx, 1337, "0xtxhash")
	require.NoError(t, err)
	assert.Equal(t, "0xtxhash", retrieved.Hash)
}