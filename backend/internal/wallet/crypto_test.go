package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("this is a secret pvt key")
	password := "veryBadPassword7899#"

	encrypted, err := Encrypt(plaintext, password)
	require.NoError(t, err)
	require.NotNil(t, encrypted)

	assert.NotEmpty(t, encrypted.Ciphertext)
	assert.Len(t, encrypted.Nonce, GCMNonceSize)
	assert.Len(t, encrypted.Tag, GCMTagSize)
	assert.Len(t, encrypted.Salt, SaltSize)

	decrypted, err := Decrypt(encrypted, password)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptWrongPassword(t *testing.T) {
	plaintext := []byte("secret")
	password := "correct"

	encrypted, err := Encrypt(plaintext, password)
	require.NoError(t, err)

	_, err = Decrypt(encrypted, "wrong")
	assert.Error(t, err)
	assert.Equal(t, ErrDecryptionFailed, err)
}

func TestEncryptionUniqueness(t *testing.T) {
	plaintext := []byte("same plaintext")
	password := "same password"

	encrypted1, err := Encrypt(plaintext, password)
	require.NoError(t, err)

	encrypted2, err := Encrypt(plaintext, password)
	require.NoError(t, err)

	assert.NotEqual(t, encrypted1.Salt, encrypted2.Salt)
	assert.NotEqual(t, encrypted1.Nonce, encrypted2.Nonce)
	assert.NotEqual(t, encrypted1.Ciphertext, encrypted2.Ciphertext)

	decrypted1, _ := Decrypt(encrypted1, password)
	decrypted2, _ := Decrypt(encrypted2, password)

	assert.Equal(t, plaintext, decrypted1)
	assert.Equal(t, plaintext, decrypted2)
}

func TestToHexFromHex(t *testing.T) {
	plaintext := []byte("test")
	password := "badPassword4341!"

	encrypted, err := Encrypt(plaintext, password)
	require.NoError(t, err)

	cHex, nHex, tHex, sHex := encrypted.ToHex()

	reconstructed, err := FromHex(cHex, nHex, tHex, sHex)
	require.NoError(t, err)

	decrypted, er := Decrypt(reconstructed, password)
	require.NoError(t, er)
	assert.Equal(t, plaintext, decrypted)
}

func TestVerifyPassword(t *testing.T) {
	plaintext := []byte("secret key")
	password := "passwordansiodna"

	encrypted, err := Encrypt(plaintext, password)
	require.NoError(t, err)

	assert.True(t, VerifyPassword(password, encrypted, plaintext))

	assert.False(t, VerifyPassword("wrongpassword", encrypted, plaintext))

	assert.False(t, VerifyPassword(password, encrypted, []byte("wrong plaintext")))
}

func TestDeriveSamePasswordSameSalt(t *testing.T) {
	password := "asdansda"
	salt, _ := GenerateSalt()

	key1 := DeriveKey(password, salt)
	key2 := DeriveKey(password, salt)

	assert.Equal(t, key1, key2, "Same password + salt should product same key")
}

func TestDeriveKeryDifferentSalts(t *testing.T) {
	password := "asndoasda"
	salt1, _ := GenerateSalt()
	salt2, _ := GenerateSalt()

	key1 := DeriveKey(password, salt1)
	key2 := DeriveKey(password, salt2)

	assert.NotEqual(t, key1, key2, "Different salts should produce different keys")

}

func TestClearBytes(t *testing.T) {
	data := []byte("sensitive data")
	original := make([]byte, len(data))
	copy(original, data)

	clearBytes(data)

	for _, b := range data {
		assert.Equal(t, byte(0), b)
	}

	assert.NotEqual(t, original, data)
}

func BenchmarkEncrypt(b *testing.B) {
	plaintext := []byte("benchmark data for encryption")
	password := "benchmarkPassword"

	b.ResetTimer()

	for i:=0; i<b.N; i++ {
		_, _ = Encrypt(plaintext, password)
	}
}

func BenchmarkDecryption(b *testing.B) {
	plaintext := []byte("benchmark data for decryption")
	password := "asjbdniajsda"

	encrypted, _ := Encrypt(plaintext, password)

	b.ResetTimer()

	for i:=0; i<b.N; i++ {
		_, _ = Decrypt(encrypted, password)
	}

}

func BenchmarkDeriveKey(b *testing.B) {
	password := "aosbnjkbdas"
	salt, _ := GenerateSalt()

	b.ResetTimer()

	for i:=0; i<b.N; i++ {
		_ = DeriveKey(password, salt)
	}
}