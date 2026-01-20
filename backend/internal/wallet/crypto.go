package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	Argon2Time = 1
	Argon2Memory = 64 * 1024
	Argon2Threads = 4
	Argon2KeyLen = 32

	GCMNonceSize = 12
	GCMTagSize = 16
	SaltSize = 32
)

var (
	ErrInvalidPassword = errors.New("Invalid Password")
	ErrDecryptionFailed = errors.New("Decryption Failed")
	ErrInvalidCiphertext = errors.New("Invalid Ciphertext")
)

type EncryptedData struct {
	Ciphertext []byte
	Nonce []byte
	Tag []byte
	Salt []byte
}

func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		Argon2KeyLen,
	)
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err!=nil {
		return nil, err
	}
	return salt, nil
}

func Encrypt(plaintext []byte, password string) (*EncryptedData, error) {
	salt, err := GenerateSalt()
	if err!=nil {
		return nil, err
	}

	key := DeriveKey(password, salt)
	defer clearBytes(key)

	block, err := aes.NewCipher(key)
	if err!=nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err!=nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err!=nil {
		return nil, err
	}

	ciphertextAndTag := gcm.Seal(nil, nonce, plaintext, nil)

	ciphertext := ciphertextAndTag[:len(ciphertextAndTag)-GCMTagSize]
	tag := ciphertextAndTag[len(ciphertextAndTag)-GCMTagSize:]

	return &EncryptedData{
		Ciphertext: ciphertext,
		Nonce: nonce,
		Tag: tag,
		Salt: salt,
	}, nil
}

func Decrypt(encrypted *EncryptedData, password string) ([]byte, error) {
	key := DeriveKey(password, encrypted.Salt)
	defer clearBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertextAndTag := append(encrypted.Ciphertext, encrypted.Tag...)

	plaintext, err := gcm.Open(nil, encrypted.Nonce, ciphertextAndTag, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

func VerifyPassword(password string, encrypted *EncryptedData, expectedPlaintext []byte) bool {
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		return false
	}

	defer clearBytes(decrypted)

	return subtle.ConstantTimeCompare(decrypted, expectedPlaintext) == 1
}

func clearBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func (e *EncryptedData) ToHex() (ciphertext, nonce, tag, salt string) {
	return hex.EncodeToString(e.Ciphertext), hex.EncodeToString(e.Nonce), hex.EncodeToString(e.Tag), hex.EncodeToString(e.Salt)
}

func FromHex(ciphertext, nonce, tag, salt string) (*EncryptedData, error) {
	c, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	n, err := hex.DecodeString(nonce)
	if err != nil {
		return nil, err
	}
	t, err := hex.DecodeString(tag)
	if err != nil {
		return nil, err
	}
	s, err := hex.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	return &EncryptedData{
		Ciphertext: c,
		Nonce: n,
		Tag: t,
		Salt: s,
	}, nil
}