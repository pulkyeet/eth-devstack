package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pulkyeet/eth-devstack/backend/internal/blockchain"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/models"
	"go.uber.org/zap"
)

type TxProcessor struct {
	db     *database.DB
	client *blockchain.ChainClient
	logger *zap.SugaredLogger
}

func NewTxProcessor(db *database.DB, client *blockchain.ChainClient, logger *zap.Logger) *TxProcessor {
	return &TxProcessor{
		db:     db,
		client: client,
		logger: logger.Sugar(),
	}
}

func (tp *TxProcessor) ProcessBlockTransactions(ctx context.Context, block *types.Block, chainID int64) error {
	blockTime := time.Unix(int64(block.Time()), 0)

	for i, tx := range block.Transactions() {
		// Get transaction receipt - FIXED: use GetTransactionReceipt with string
		receipt, err := tp.client.GetTransactionReceipt(ctx, tx.Hash().Hex())
		if err != nil {
			tp.logger.Warnw("Failed to get receipt", "tx_hash", tx.Hash().Hex(), "error", err)
			continue
		}

		// Build transaction model
		txModel := &models.Transaction{
			ChainID:          chainID,
			Hash:             tx.Hash().Hex(),
			BlockNumber:      block.Number().Int64(),
			BlockHash:        block.Hash().Hex(),
			TransactionIndex: i,
			Nonce:            int64(tx.Nonce()),
			TransactionType:  int(tx.Type()),
			Timestamp:        blockTime,
		}

		// From address
		msg, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			tp.logger.Warnw("Failed to get sender", "tx_hash", tx.Hash().Hex(), "error", err)
			continue
		}
		txModel.FromAddress = msg.Hex()

		// To address
		if tx.To() != nil {
			toAddr := tx.To().Hex()
			txModel.ToAddress = &toAddr
		}

		// Value
		txModel.Value = tx.Value().String()

		// Gas - FIXED: GasUsed is *int64
		txModel.Gas = int64(tx.Gas())
		gasUsed := int64(receipt.GasUsed)
		txModel.GasUsed = &gasUsed
		cumulativeGas := int64(receipt.CumulativeGasUsed)
		txModel.CumulativeGasUsed = &cumulativeGas

		// Gas price
		if tx.GasPrice() != nil {
			txModel.GasPrice = toStringPtr(tx.GasPrice().String())
		}
		if tx.GasFeeCap() != nil {
			txModel.MaxFeePerGas = toStringPtr(tx.GasFeeCap().String())
		}
		if tx.GasTipCap() != nil {
			txModel.MaxPriorityFeePerGas = toStringPtr(tx.GasTipCap().String())
		}

		// Effective gas price
		if receipt.EffectiveGasPrice != nil {
			txModel.EffectiveGasPrice = toStringPtr(receipt.EffectiveGasPrice.String())
		}

		// Input data
		if len(tx.Data()) > 0 {
			data := fmt.Sprintf("0x%x", tx.Data())
			txModel.Input = &data
		}

		// Status
		status := int(receipt.Status)
		txModel.Status = &status

		// Contract address
		if receipt.ContractAddress != (common.Address{}) {
			contractAddr := receipt.ContractAddress.Hex()
			txModel.ContractAddress = &contractAddr
		}

		// Logs bloom
		if len(receipt.Bloom) > 0 {
			bloom := fmt.Sprintf("0x%x", receipt.Bloom)
			txModel.LogsBloom = &bloom
		}

		// Insert transaction
		if err := tp.db.InsertTransaction(ctx, txModel); err != nil {
			tp.logger.Errorw("Failed to insert transaction", "tx_hash", tx.Hash().Hex(), "error", err)
			continue
		}

		// Update address balances
		if err := tp.updateAddressBalances(ctx, txModel, chainID, block.Number()); err != nil {
			tp.logger.Errorw("Failed to update address balances", "tx_hash", tx.Hash().Hex(), "error", err)
		}
	}

	return nil
}

func (tp *TxProcessor) updateAddressBalances(ctx context.Context, tx *models.Transaction, chainID int64, blockNum *big.Int) error {
	// Update FROM address - FIXED: GetBalance takes string, Balance is int64
	fromBalance, err := tp.client.GetBalance(ctx, tx.FromAddress, blockNum)
	if err != nil {
		return fmt.Errorf("failed to get balance for from address: %w", err)
	}

	fromAddr := &models.Address{
		ChainID:        chainID,
		Address:        tx.FromAddress,
		Balance:        fromBalance.String(), // FIXED: convert to int64
		Nonce:          0,
		FirstSeenBlock: &tx.BlockNumber,
		LastSeenBlock:  &tx.BlockNumber,
		FirstSeenAt:    &tx.Timestamp,
		LastSeenAt:     &tx.Timestamp,
	}

	if err := tp.db.UpsertAddress(ctx, fromAddr); err != nil {
		return fmt.Errorf("failed to upsert from address: %w", err)
	}

	// Update TO address
	if tx.ToAddress != nil {
		toBalance, err := tp.client.GetBalance(ctx, *tx.ToAddress, blockNum)
		if err != nil {
			return fmt.Errorf("failed to get balance for to address: %w", err)
		}

		toAddr := &models.Address{
			ChainID:        chainID,
			Address:        *tx.ToAddress,
			Balance:        toBalance.String(), // FIXED
			Nonce:          0,
			FirstSeenBlock: &tx.BlockNumber,
			LastSeenBlock:  &tx.BlockNumber,
			FirstSeenAt:    &tx.Timestamp,
			LastSeenAt:     &tx.Timestamp,
		}

		if err := tp.db.UpsertAddress(ctx, toAddr); err != nil {
			return fmt.Errorf("failed to upsert to address: %w", err)
		}
	}

	// Update contract address
	if tx.ContractAddress != nil {
		contractBalance, err := tp.client.GetBalance(ctx, *tx.ContractAddress, blockNum)
		if err != nil {
			return fmt.Errorf("failed to get balance for contract address: %w", err)
		}

		contractAddr := &models.Address{
			ChainID:         chainID,
			Address:         *tx.ContractAddress,
			Balance:         contractBalance.String(), // FIXED
			IsContract:      true,
			ContractCreator: &tx.FromAddress,
			CreationTxHash:  &tx.Hash,
			FirstSeenBlock:  &tx.BlockNumber,
			LastSeenBlock:   &tx.BlockNumber,
			FirstSeenAt:     &tx.Timestamp,
			LastSeenAt:      &tx.Timestamp,
		}

		if err := tp.db.UpsertAddress(ctx, contractAddr); err != nil {
			return fmt.Errorf("failed to upsert contract address: %w", err)
		}
	}

	return nil
}

