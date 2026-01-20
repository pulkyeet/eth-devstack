package wallet

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic()
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	// Should be 12 words
	words := strings.Split(mnemonic, " ")
	assert.Len(t, words, 12)

	// Should be valid
	assert.True(t, ValidateMnemonic(mnemonic))
}

func TestGenerateMnemonic24(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	// Should be 24 words
	words := strings.Split(mnemonic, " ")
	assert.Len(t, words, 24)

	// Should be valid
	assert.True(t, ValidateMnemonic(mnemonic))
}

func TestNewHDWallet(t *testing.T) {
	mnemonic, _ := GenerateMnemonic()

	wallet, err := NewHDWallet(mnemonic)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	assert.Equal(t, mnemonic, wallet.mnemonic)
	assert.NotEmpty(t, wallet.seed)
}

func TestNewHDWalletInvalidMnemonic(t *testing.T) {
	invalidMnemonic := "invalid mnemonic words that dont make sense"

	wallet, err := NewHDWallet(invalidMnemonic)
	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Equal(t, ErrInvalidMnemonic, err)
}

func TestDeriveAccount(t *testing.T) {
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)

	// Derive first account
	account, err := wallet.DeriveAccount(0)
	require.NoError(t, err)
	require.NotNil(t, account)

	// Should have valid components
	assert.NotNil(t, account.PrivateKey)
	assert.NotNil(t, account.PublicKey)
	assert.NotEqual(t, common.Address{}, account.Address)
	assert.Contains(t, account.Path, "44")
	assert.Contains(t, account.Path, "60")

	// Private key should be 32 bytes
	assert.Len(t, account.PrivateKeyBytes(), 32)

	// Address should be 20 bytes (42 chars with 0x)
	assert.Len(t, account.Address.Hex(), 42)
}

func TestDeriveMultipleAccounts(t *testing.T) {
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)

	// Derive first 5 accounts
	accounts := make([]*Account, 5)
	for i := 0; i < 5; i++ {
		account, err := wallet.DeriveAccount(uint32(i))
		require.NoError(t, err)
		accounts[i] = account
	}

	// All accounts should be different
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 5; j++ {
			assert.NotEqual(t, accounts[i].Address, accounts[j].Address)
			assert.NotEqual(t, accounts[i].PrivateKeyHex(), accounts[j].PrivateKeyHex())
		}
	}
}

func TestSameMnemonicSameAccounts(t *testing.T) {
	mnemonic, _ := GenerateMnemonic()

	// Create two wallets from same mnemonic
	wallet1, _ := NewHDWallet(mnemonic)
	wallet2, _ := NewHDWallet(mnemonic)

	// Derive same account from both
	account1, _ := wallet1.DeriveAccount(0)
	account2, _ := wallet2.DeriveAccount(0)

	// Should be identical
	assert.Equal(t, account1.Address, account2.Address)
	assert.Equal(t, account1.PrivateKeyHex(), account2.PrivateKeyHex())
}

func TestImportPrivateKey(t *testing.T) {
	// Generate a wallet to get a valid private key
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)
	originalAccount, _ := wallet.DeriveAccount(0)

	// Export private key
	privateKeyHex := originalAccount.PrivateKeyHex()

	// Import it
	importedAccount, err := ImportPrivateKey(privateKeyHex)
	require.NoError(t, err)

	// Should have same address
	assert.Equal(t, originalAccount.Address, importedAccount.Address)
	assert.Equal(t, "imported", importedAccount.Path)
}

func TestImportPrivateKeyWithPrefix(t *testing.T) {
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)
	account, _ := wallet.DeriveAccount(0)

	// Add 0x prefix
	privateKeyHex := "0x" + account.PrivateKeyHex()

	// Should still import correctly
	imported, err := ImportPrivateKey(privateKeyHex)
	require.NoError(t, err)
	assert.Equal(t, account.Address, imported.Address)
}

func TestDerivationPath(t *testing.T) {
	path := DerivationPath(0)

	// Should be: m/44'/60'/0'/0/0
	expected := "m/44'/60'/0'/0/0"
	assert.Equal(t, expected, path.String())

	path5 := DerivationPath(5)
	expected5 := "m/44'/60'/0'/0/5"
	assert.Equal(t, expected5, path5.String())
}

func TestParseDerivationPath(t *testing.T) {
	pathStr := "m/44'/60'/0'/0/0"

	path, err := ParseDerivationPath(pathStr)
	require.NoError(t, err)
	assert.Equal(t, pathStr, path.String())
}

func TestAccountClear(t *testing.T) {
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)
	account, _ := wallet.DeriveAccount(0)

	// Private key should exist
	require.NotNil(t, account.PrivateKey)
	require.NotNil(t, account.PrivateKey.D)

	// Clear it
	account.Clear()

	// D should be zero
	assert.Equal(t, int64(0), account.PrivateKey.D.Int64())
}

func TestValidateMnemonic(t *testing.T) {
	validMnemonic, _ := GenerateMnemonic()
	assert.True(t, ValidateMnemonic(validMnemonic))

	invalidMnemonic := "invalid words that are not real"
	assert.False(t, ValidateMnemonic(invalidMnemonic))
}

func BenchmarkGenerateMnemonic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateMnemonic()
	}
}

func BenchmarkDeriveAccount(b *testing.B) {
	mnemonic, _ := GenerateMnemonic()
	wallet, _ := NewHDWallet(mnemonic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = wallet.DeriveAccount(0)
	}
}