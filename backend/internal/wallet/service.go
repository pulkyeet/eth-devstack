package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/config"
	"go.uber.org/zap"
)

// ⚠️ SECURITY WARNING: This service stores encrypted private keys server-side.
// This is FOR DEVELOPMENT/LEARNING ONLY. Production wallets should NEVER
// store private keys on servers. Use client-side signing (MetaMask, WalletConnect).

type Service struct {
	db      *database.DB
	logger  *zap.SugaredLogger
	clients map[int64]*ethclient.Client
	config  *config.Config
}

func NewService(db *database.DB, logger *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		db:      db,
		logger:  logger.Sugar(),
		clients: make(map[int64]*ethclient.Client),
		config:  cfg,
	}
}

type CreateWalletRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateWalletResponse struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Mnemonic string `json:"mnemonic"`
	Warning  string `json:"warning"`
}

type ImportWalletRequest struct {
	Method     string `json:"method"` // "mnemonic" or "private_key"
	Mnemonic   string `json:"mnemonic,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	Password   string `json:"password"`
	Name       string `json:"name"`
}

type SendTransactionRequest struct {
	WalletID string  `json:"wallet_id"`
	ChainID  int64   `json:"chain_id"`
	To       string  `json:"to"`
	Value    string  `json:"value"` // in wei
	Password string  `json:"password"`
	GasLimit *uint64 `json:"gas_limit,omitempty"`
	GasPrice *string `json:"gas_price,omitempty"`
}

type SendTransactionResponse struct {
	TxHash string `json:"tx_hash"`
	Status string `json:"status"`
}

type SignMessageRequest struct {
	WalletID string `json:"wallet_id"`
	Message  string `json:"message"`
	Password string `json:"password"`
}

type SignMessageResponse struct {
	Message     string `json:"message"`
	Signature   string `json:"signature"`
	MessageHash string `json:"message_hash"`
}

// Helper: Encrypt private key and return hex-encoded components
func encryptPrivateKey(privateKeyBytes []byte, password string) (ciphertext, nonce, tag, salt string, err error) {
	encrypted, err := Encrypt(privateKeyBytes, password)
	if err != nil {
		return "", "", "", "", err
	}
	ciphertext, nonce, tag, salt = encrypted.ToHex()
	return
}

// Helper: Decrypt private key from hex-encoded components
func decryptPrivateKey(ciphertext, nonce, tag, salt, password string) ([]byte, error) {
	encrypted, err := FromHex(ciphertext, nonce, tag, salt)
	if err != nil {
		return nil, err
	}
	return Decrypt(encrypted, password)
}

// CreateWallet generates new HD wallet, encrypts, stores
func (s *Service) CreateWallet(ctx context.Context, req CreateWalletRequest) (*CreateWalletResponse, error) {
	// Validate password
	if len(req.Password) < 12 {
		return nil, fmt.Errorf("password must be at least 12 characters")
	}

	// Generate mnemonic
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Create HD wallet
	hdWallet, err := NewHDWallet(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create HD wallet: %w", err)
	}

	// Derive first account (index 0)
	account, err := hdWallet.DeriveAccount(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %w", err)
	}
	defer account.Clear()

	address := account.Address.Hex()

	// Encrypt private key
	privateKeyBytes := account.PrivateKeyBytes()
	ciphertext, nonce, tag, salt, err := encryptPrivateKey(privateKeyBytes, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	// Clear private key bytes from memory
	for i := range privateKeyBytes {
		privateKeyBytes[i] = 0
	}

	// Store in database
	wallet := &database.Wallet{
		Address:             address,
		EncryptedPrivateKey: ciphertext,
		EncryptionNonce:     nonce,
		EncryptionTag:       tag,
		EncryptionSalt:      salt,
		Name:                req.Name,
		DerivationPath:      account.Path,
		IsImported:          false,
	}

	if err := s.db.CreateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to store wallet: %w", err)
	}

	s.logger.Infow("Wallet created", "wallet_id", wallet.ID, "address", address)

	return &CreateWalletResponse{
		ID:       wallet.ID,
		Address:  address,
		Name:     req.Name,
		Mnemonic: mnemonic,
		Warning:  "⚠️ SAVE YOUR MNEMONIC SECURELY - It cannot be recovered!",
	}, nil
}

// ImportWallet imports existing wallet from mnemonic or private key
func (s *Service) ImportWallet(ctx context.Context, req ImportWalletRequest) (*database.Wallet, error) {
	// Validate password
	if len(req.Password) < 12 {
		return nil, fmt.Errorf("password must be at least 12 characters")
	}

	var account *Account
	var err error

	switch req.Method {
	case "mnemonic":
		if req.Mnemonic == "" {
			return nil, fmt.Errorf("mnemonic required")
		}

		// Create HD wallet
		hdWallet, err := NewHDWallet(req.Mnemonic)
		if err != nil {
			return nil, fmt.Errorf("invalid mnemonic: %w", err)
		}

		// Derive first account
		account, err = hdWallet.DeriveAccount(0)
		if err != nil {
			return nil, fmt.Errorf("failed to derive account: %w", err)
		}

	case "private_key":
		if req.PrivateKey == "" {
			return nil, fmt.Errorf("private_key required")
		}

		account, err = ImportPrivateKey(req.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}

	default:
		return nil, fmt.Errorf("method must be 'mnemonic' or 'private_key'")
	}

	defer account.Clear()

	address := account.Address.Hex()

	// Check if wallet already exists
	existing, _ := s.db.GetWalletByAddress(ctx, address)
	if existing != nil {
		return nil, fmt.Errorf("wallet with this address already exists")
	}

	// Encrypt private key
	privateKeyBytes := account.PrivateKeyBytes()
	ciphertext, nonce, tag, salt, err := encryptPrivateKey(privateKeyBytes, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	// Clear private key bytes
	for i := range privateKeyBytes {
		privateKeyBytes[i] = 0
	}

	// Store in database
	wallet := &database.Wallet{
		Address:             address,
		EncryptedPrivateKey: ciphertext,
		EncryptionNonce:     nonce,
		EncryptionTag:       tag,
		EncryptionSalt:      salt,
		Name:                req.Name,
		DerivationPath:      account.Path,
		IsImported:          req.Method == "private_key",
	}

	if err := s.db.CreateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to store wallet: %w", err)
	}

	s.logger.Infow("Wallet imported", "wallet_id", wallet.ID, "address", address, "method", req.Method)

	return wallet, nil
}

// GetWallet retrieves wallet by ID
func (s *Service) GetWallet(ctx context.Context, walletID string) (*database.Wallet, error) {
	return s.db.GetWalletByID(ctx, walletID)
}

// GetBalances retrieves balances across all chains for a wallet
func (s *Service) GetBalances(ctx context.Context, walletID string, updateFromChain bool) ([]database.WalletBalance, error) {
	wallet, err := s.db.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if updateFromChain {
		// Fetch fresh balances from blockchain
		chains, err := s.getActiveChains(ctx)
		if err != nil {
			return nil, err
		}

		for _, chain := range chains {
			client, err := s.getClient(chain.ChainID, chain.RPCEndpoint)
			if err != nil {
				s.logger.Warnw("Failed to connect to chain", "chain_id", chain.ChainID, "error", err)
				continue
			}

			balance, err := client.BalanceAt(ctx, common.HexToAddress(wallet.Address), nil)
			if err != nil {
				s.logger.Warnw("Failed to fetch balance", "chain_id", chain.ChainID, "error", err)
				continue
			}

			// Update in database
			if err := s.db.UpsertWalletBalance(ctx, walletID, chain.ChainID, balance.String()); err != nil {
				s.logger.Warnw("Failed to update balance", "chain_id", chain.ChainID, "error", err)
			}
		}
	}

	// Return from database
	return s.db.GetWalletBalances(ctx, walletID)
}

// SendTransaction signs and broadcasts transaction
func (s *Service) SendTransaction(ctx context.Context, req SendTransactionRequest) (*SendTransactionResponse, error) {
	// Get wallet
	wallet, err := s.db.GetWalletByID(ctx, req.WalletID)
	if err != nil {
		return nil, err
	}

	// Decrypt private key
	privateKeyBytes, err := decryptPrivateKey(
		wallet.EncryptedPrivateKey,
		wallet.EncryptionNonce,
		wallet.EncryptionTag,
		wallet.EncryptionSalt,
		req.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key (wrong password?): %w", err)
	}
	defer func() {
		for i := range privateKeyBytes {
			privateKeyBytes[i] = 0
		}
	}()

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create Account struct for signer
	account := &Account{
		PrivateKey: privateKey,
		Address:    common.HexToAddress(wallet.Address),
		Path:       wallet.DerivationPath,
	}
	defer account.Clear()

	// Get chain RPC
	chain, err := s.getChain(ctx, req.ChainID)
	if err != nil {
		return nil, err
	}

	client, err := s.getClient(chain.ChainID, chain.RPCEndpoint)
	if err != nil {
		return nil, err
	}

	// Create signer
	chainIDBig := big.NewInt(req.ChainID)
	signer := NewSigner(client, chainIDBig)

	// Parse value
	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid value")
	}

	to := common.HexToAddress(req.To)

	// Estimate gas if not provided
	var gasLimit uint64
	if req.GasLimit != nil {
		gasLimit = *req.GasLimit
	} else {
		gasLimit, err = signer.EstimateGas(ctx, account.Address, to, value, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", err)
		}
	}

	// Get gas price if not provided
	var gasPrice *big.Int
	if req.GasPrice != nil {
		gasPrice, ok = new(big.Int).SetString(*req.GasPrice, 10)
		if !ok {
			return nil, fmt.Errorf("invalid gas price")
		}
	} else {
		gasPrice, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
		}
	}

	// Build transaction params
	params := &TxParams{
		To:       to,
		Value:    value,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	}

	// Sign transaction
	signedTx, err := signer.SignTransaction(ctx, account, params)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Broadcast
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	txHash := signedTx.Hash().Hex()

	// Store in database
	dbTx := &database.WalletTransaction{
		WalletID:    req.WalletID,
		ChainID:     req.ChainID,
		TxHash:      txHash,
		Direction:   "sent",
		FromAddress: wallet.Address,
		ToAddress:   &req.To,
		Value:       value.String(),
		Status:      2, // pending
	}
	if err := s.db.CreateWalletTransaction(ctx, dbTx); err != nil {
		s.logger.Warnw("Failed to store transaction", "tx_hash", txHash, "error", err)
	}

	s.logger.Infow("Transaction sent", "wallet_id", req.WalletID, "tx_hash", txHash, "to", req.To)

	return &SendTransactionResponse{
		TxHash: txHash,
		Status: "pending",
	}, nil
}

// SignMessage signs arbitrary message (EIP-191)
func (s *Service) SignMessage(ctx context.Context, req SignMessageRequest) (*SignMessageResponse, error) {
	// Get wallet
	wallet, err := s.db.GetWalletByID(ctx, req.WalletID)
	if err != nil {
		return nil, err
	}

	// Decrypt private key
	privateKeyBytes, err := decryptPrivateKey(
		wallet.EncryptedPrivateKey,
		wallet.EncryptionNonce,
		wallet.EncryptionTag,
		wallet.EncryptionSalt,
		req.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key (wrong password?): %w", err)
	}
	defer func() {
		for i := range privateKeyBytes {
			privateKeyBytes[i] = 0
		}
	}()

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Sign message (EIP-191)
	messageHash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n" + fmt.Sprint(len(req.Message)) + req.Message))
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Adjust V for Ethereum (27/28 instead of 0/1)
	signature[64] += 27

	return &SignMessageResponse{
		Message:     req.Message,
		Signature:   common.Bytes2Hex(signature),
		MessageHash: messageHash.Hex(),
	}, nil
}

// GetTransactions retrieves wallet transaction history
func (s *Service) GetTransactions(ctx context.Context, walletID string, chainID *int64, limit, offset int) ([]database.WalletTransaction, error) {
	return s.db.GetWalletTransactions(ctx, walletID, chainID, limit, offset)
}

// Helper: Get active chains
func (s *Service) getActiveChains(ctx context.Context) ([]struct {
	ChainID     int64
	RPCEndpoint string
}, error) {
	query := `SELECT chain_id, rpc_endpoint FROM chains WHERE is_active = true`
	rows, err := s.db.GetConn().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chains []struct {
		ChainID     int64
		RPCEndpoint string
	}
	for rows.Next() {
		var c struct {
			ChainID     int64
			RPCEndpoint string
		}
		if err := rows.Scan(&c.ChainID, &c.RPCEndpoint); err != nil {
			return nil, err
		}
		chains = append(chains, c)
	}
	return chains, nil
}

// Helper: Get chain by ID
func (s *Service) getChain(ctx context.Context, chainID int64) (struct {
	ChainID     int64
	RPCEndpoint string
}, error) {
	var chain struct {
		ChainID     int64
		RPCEndpoint string
	}
	query := `SELECT chain_id, rpc_endpoint FROM chains WHERE chain_id = $1`
	err := s.db.GetConn().QueryRowContext(ctx, query, chainID).Scan(&chain.ChainID, &chain.RPCEndpoint)
	return chain, err
}

// Helper: Get or create ethclient for chain
func (s *Service) getClient(chainID int64, rpcURL string) (*ethclient.Client, error) {
	if client, ok := s.clients[chainID]; ok {
		return client, nil
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	s.clients[chainID] = client
	return client, nil
}
