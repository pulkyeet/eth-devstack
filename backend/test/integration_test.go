package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080/api/v1"

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestChainsEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/chains")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestBlocksEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/blocks?chain_id=1337&limit=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestSearchEndpoint(t *testing.T) {
	// Search by block number
	resp, err := http.Get(baseURL + "/search?q=1&chain_id=1337")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestStatsEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/stats?chain_id=1337")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}