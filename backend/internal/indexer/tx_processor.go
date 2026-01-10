package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pulkyeet/eth-devstack/backend/internal/blockchain"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/models"
	"go.uber.org/zap"
)

type TxProcessor struct {
	db *database.DB
	client *blockchain.ChainClient
	logger *zap.SugaredLogger
}

func NewTxProcessor(db *database.DB, client *blockchain.ChainClient, logger *zap.Logger) *TxProcessor {
	return & TxProcessor{
		db: db,
		client: client,
		logger: logger.Sugar(),
	}
}

func (tp *TxProcessor) ProcessTransaction(ctx context.Context, tx *types.Transaction, blockNumber uint64, blockHash string, txIndex int, blockTime time.Time, chainID int64) error {
	msg, err := types.Sender(types.LatestSignerForChainID(big.NewInt(chainID)), tx)
	if err != nil {
		return fmt.Errorf("Failed to get sender: %w", err)
	}

	receipt, err := tp.client.GetTransactionReceipt(ctx, tx.Hash().Hex())
	if err != nil {
		tp.logger.Warnw("Failed to get receipt", "tx_hash", tx.Hash().Hex(), "error", err)
		receipt = nil
	}
	
	txType := int(tx.Type())
	txModel := &models.Transaction{
		ChainID:          chainID,
		Hash:             tx.Hash().Hex(),
		BlockNumber:      int64(blockNumber),
		BlockHash:        blockHash,
		TransactionIndex: txIndex,
		FromAddress:      msg.Hex(),
		ToAddress:        nil,
		Value:            tx.Value().String(),
		Gas:              int64(tx.Gas()),
		GasPrice:         bigIntToStringPtr(tx.GasPrice()),
		Input:            toStringPtr(fmt.Sprintf("0x%x", tx.Data())),
		Nonce:            int64(tx.Nonce()),
		TransactionType:  txType,
		Timestamp:        blockTime,
	}
	
	if tx.To() != nil {
		to := tx.To().Hex()
		txModel.ToAddress = &to
	}

	if tx.Type() == types.DynamicFeeTxType {
		txModel.MaxFeePerGas = bigIntToStringPtr(tx.GasFeeCap())
		txModel.MaxPriorityFeePerGas = bigIntToStringPtr(tx.GasTipCap())
	}

	if receipt != nil {
		status := int(receipt.Status)
		txModel.Status = &status
		gasUsed := int64(receipt.GasUsed)
		txModel.GasUsed = &gasUsed
		cumulativeGasUsed := int64(receipt.CumulativeGasUsed)
		txModel.CumulativeGasUsed = &cumulativeGasUsed
		txModel.EffectiveGasPrice = bigIntToStringPtr(receipt.EffectiveGasPrice)
		
		if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
			contractAddr := receipt.ContractAddress.Hex()
			txModel.ContractAddress = &contractAddr
		}

		logsBloom := fmt.Sprintf("0x%x", receipt.Bloom[:])
		txModel.LogsBloom = &logsBloom
	}
	
	if err := tp.db.InsertTransaction(ctx, txModel); err != nil {
		return fmt.Errorf("Failed to insert transaction: %w", err)
	}
	
	return nil
}