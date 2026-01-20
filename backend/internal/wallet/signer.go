package wallet

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient funds")
	ErrNonceTooLow = errors.New("Nonce too low")
	ErrInvalidChainID = errors.New("Invalid ChainID")
)

type Signer struct {
	client *ethclient.Client
	chainID *big.Int
	nonceLock sync.Mutex
	nonceCache map[common.Address]uint64
}

func NewSigner(client *ethclient.Client, chainID *big.Int) *Signer {
	return &Signer{
		client: client,
		chainID: chainID,
		nonceCache: make(map[common.Address]uint64),
	}
}

type TxParams struct {
	To common.Address
	Value *big.Int
	GasLimit uint64
	Data []byte
	GasPrice *big.Int
	MaxFeePerGas *big.Int
	MaxPriorityFeePerGas *big.Int
}

func (s *Signer) SignTransaction(ctx context.Context, account *Account, params *TxParams) (*types.Transaction, error) {
	nonce, err := s.GetNonce(ctx, account.Address)
	if err!=nil {
		return nil, fmt.Errorf("Get nonce error: %w", err)
	}

	var tx *types.Transaction
	if params.MaxFeePerGas !=nil && params.MaxPriorityFeePerGas != nil {
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID: s.chainID,
			Nonce: nonce,
			GasTipCap: params.MaxPriorityFeePerGas,
			GasFeeCap: params.MaxFeePerGas,
			Gas: params.GasLimit,
			To: &params.To,
			Value: params.Value,
			Data: params.Data,
		})
	} else {
		gasPrice := params.GasPrice
		if gasPrice == nil {
			gasPrice, err = s.client.SuggestGasPrice(ctx)
			if err!=nil {
				return nil, fmt.Errorf("Suggest gas price error: %w", err)
			}
		}

		tx = types.NewTx(&types.LegacyTx{
			Nonce: nonce,
			GasPrice: gasPrice,
			Gas: params.GasLimit,
			To: &params.To,
			Value: params.Value,
			Data: params.Data,
		})
	}

	signer := types.NewLondonSigner(s.chainID)
	signedTx, err := types.SignTx(tx, signer, account.PrivateKey)
	if err!=nil {
		return nil, fmt.Errorf("Sign transaction error: %w", err)
	}

	s.incrementNonce(account.Address)
	return signedTx, nil
}

func (s *Signer) SendTransaction(ctx context.Context, account *Account, params *TxParams) (common.Hash, error) {
	signedTx, err := s.SignTransaction(ctx, account, params)
	if err!=nil {
		return common.Hash{}, err
	}

	err = s.client.SendTransaction(ctx, signedTx)
	if err !=nil {
		if isNonceTooLowError(err) {
			s.RefreshNonce(ctx, account.Address)
			return common.Hash{}, fmt.Errorf("Send transaction error: %w", err)
		}
	}
	return signedTx.Hash(), nil
} 

func (s *Signer) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
	s.nonceLock.Lock()
	defer s.nonceLock.Unlock()

	if nonce, exists := s.nonceCache[address]; exists {
		return nonce, nil
	}

	nonce, err := s.client.PendingNonceAt(ctx, address)
	if err!=nil {
		return 0, fmt.Errorf("Get nonce error: %w", err)
	}

	s.nonceCache[address] = nonce
	return nonce, nil
}

func (s *Signer) RefreshNonce(ctx context.Context, address common.Address) error {
	s.nonceLock.Lock()
	defer s.nonceLock.Unlock()

	nonce, err := s.client.PendingNonceAt(ctx, address)
	if err!=nil {
		return err
	}

	s.nonceCache[address] = nonce
	return nil
}

func (s *Signer) incrementNonce(address common.Address) {
	s.nonceLock.Lock()
	defer s.nonceLock.Unlock()

	if nonce, exists := s.nonceCache[address]; exists {
		s.nonceCache[address] = nonce + 1
	}
}

func (s *Signer) EstimateGas(ctx context.Context, from, to common.Address, value *big.Int, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From: from,
		To: &to,
		Value: value,
		Data: data,
	}

	gas, err := s.client.EstimateGas(ctx, msg)
	if err!=nil {
		return 0, err
	}
	return gas + (gas/10), nil
}

func (s *Signer) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return s.client.BalanceAt(ctx, address, nil)
}

func isNonceTooLowError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "nonce too low") || strings.Contains(s, "already known")
}