package wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

const (
	EthereumCoinType = 60
	DefaultAccount = 0
)

var (
	ErrInvalidMnemonic = errors.New("Inalid mnemonic phrase")
	ErrInvalidPath = errors.New("Invalid derivatino path")
)

type HDWallet struct {
	mnemonic string
	seed []byte
}

func NewHDWallet(mnemonic string) (*HDWallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	seed := bip39.NewSeed(mnemonic, "")

	return &HDWallet{
		mnemonic: mnemonic,
		seed: seed,
	}, nil
}

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128) // 128 bits = 12 words
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func GenerateMnemonic24() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err!=nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err!=nil {
		return "", err
	}

	return mnemonic, nil
}

func (w *HDWallet) DeriveAccount(index uint32) (*Account, error) {
	path := DerivationPath(index)
	return w.DeriveAccountFromPath(path)
}

func (w *HDWallet) DeriveAccountFromPath(path accounts.DerivationPath) (*Account, error) {
	masterKey, err := deriveMasterKey(w.seed)
	if err!=nil {
		return nil, err
	}

	defer clearPrivateKey(masterKey)

	privateKey, err := derivePrivateKey(masterKey, path)
	if err!=nil {
		return nil, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	return &Account{
		PrivateKey: privateKey,
		PublicKey: publicKey,
		Address: address,
		Path: path.String(), 
	}, nil
}

func DerivationPath(index uint32) accounts.DerivationPath {
	return accounts.DerivationPath{
		0x80000000 + 44,
		0x80000000 + EthereumCoinType,
		0x80000000 + DefaultAccount,
		0,
		index,
	}
}

func ParseDerivationPath(path string) (accounts.DerivationPath, error) {
	 return accounts.ParseDerivationPath(path)
}

type Account struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey *ecdsa.PublicKey
	Address common.Address
	Path string
}

func (a *Account) PrivateKeyBytes() []byte {
	return crypto.FromECDSA(a.PrivateKey)
}

func (a *Account) PrivateKeyHex() string {
	return common.Bytes2Hex(a.PrivateKeyBytes())
}

func (a *Account) Clear() {
	clearPrivateKey(a.PrivateKey)
}

func ImportPrivateKey(privateKeyHex string) (*Account, error) {
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err!=nil {
		return nil, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	return &Account{
		PrivateKey: privateKey,
		PublicKey: publicKey,
		Address: address,
		Path: "imported",
	}, nil
}

func deriveMasterKey(seed []byte) (*ecdsa.PrivateKey, error) {
	return deriveKeyFromSeed(seed)
}

func deriveKeyFromSeed(seed []byte) (*ecdsa.PrivateKey, error) {
	hash := crypto.Keccak256(seed)

	privateKey, err := crypto.ToECDSA(hash)
	if err!=nil {
		return nil, err
	}

	return privateKey, nil
}

func derivePrivateKey(masterKey *ecdsa.PrivateKey, path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	masterBytes := crypto.FromECDSA(masterKey)
	pathBytes := []byte(path.String())

	combined := append(masterBytes, pathBytes...)
	hash := crypto.Keccak256(combined)

	childKey, err := crypto.ToECDSA(hash)
	if err!=nil {
		return nil, err
	}

	return childKey, nil
}

func clearPrivateKey(key *ecdsa.PrivateKey) {
	if key == nil || key.D == nil {
		return
	}
	key.D.SetInt64(0)
}

func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func MnemonicToSeed(mnemonic, passphrase string) []byte {
	return bip39.NewSeed(mnemonic, passphrase)
}