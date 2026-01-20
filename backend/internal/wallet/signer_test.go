package wallet

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: These tests require a running testnet
// Skip if RPC not available
func getTestClient(t *testing.T) *ethclient.Client {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		t.Skip("Testnet not available, skipping integration tests")
	}
	return client
}

func TestSignTransaction(t *testing.T) {
	client := getTestClient(t)
	defer client.Close()

	chainID := big.NewInt(1337)
	signer := NewSigner(client, chainID)

	// Create test account
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)
	account, _ := wallet.DeriveAccount(0)
	defer account.Clear()

	// Build transaction params
	params := &TxParams{
		To:       common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"),
		Value:    big.NewInt(1000000000000000000), // 1 ETH
		GasLimit: 21000,
		GasPrice: big.NewInt(20000000000), // 20 gwei
	}

	// Sign transaction
	signedTx, err := signer.SignTransaction(context.Background(), account, params)
	require.NoError(t, err)
	require.NotNil(t, signedTx)

	// Verify transaction fields
	assert.Equal(t, params.To, *signedTx.To())
	assert.Equal(t, params.Value, signedTx.Value())
	assert.Equal(t, params.GasLimit, signedTx.Gas())
	assert.NotNil(t, signedTx.Hash())
}

func TestNonceManagement(t *testing.T) {
	client := getTestClient(t)
	defer client.Close()

	chainID := big.NewInt(1337)
	signer := NewSigner(client, chainID)

	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")

	// First call should fetch from network
	nonce1, err := signer.GetNonce(context.Background(), address)
	require.NoError(t, err)

	// Second call should use cache (same value)
	nonce2, err := signer.GetNonce(context.Background(), address)
	require.NoError(t, err)
	assert.Equal(t, nonce1, nonce2)

	// Increment nonce
	signer.incrementNonce(address)

	// Should be incremented
	nonce3, err := signer.GetNonce(context.Background(), address)
	require.NoError(t, err)
	assert.Equal(t, nonce1+1, nonce3)

	// Refresh should fetch latest from network
	err = signer.RefreshNonce(context.Background(), address)
	require.NoError(t, err)
}

func TestEstimateGas(t *testing.T) {
	client := getTestClient(t)
	defer client.Close()

	chainID := big.NewInt(1337)
	signer := NewSigner(client, chainID)

	from := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	to := common.HexToAddress("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed")
	value := big.NewInt(1000000000000000000)

	gas, err := signer.EstimateGas(context.Background(), from, to, value, nil)
	require.NoError(t, err)

	// Simple transfer should be around 21000 + 10% buffer
	assert.GreaterOrEqual(t, gas, uint64(21000))
	assert.LessOrEqual(t, gas, uint64(25000))
}

func TestGetBalance(t *testing.T) {
	client := getTestClient(t)
	defer client.Close()

	chainID := big.NewInt(1337)
	signer := NewSigner(client, chainID)

	// Use one of your pre-funded testnet addresses
	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")

	balance, err := signer.GetBalance(context.Background(), address)
	require.NoError(t, err)
	assert.NotNil(t, balance)
}