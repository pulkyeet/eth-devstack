package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChains(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	chains, err := db.GetChains(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, chains)

	// Should have local chain
	found := false
	for _, chain := range chains {
		if chain.ChainID == 1337 {
			found = true
		}
	}
	assert.True(t, found)
}