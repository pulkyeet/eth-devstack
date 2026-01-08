package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/models"
	"go.uber.org/zap"
)

type BlockProcessor struct {
	db     *database.DB
	logger *zap.SugaredLogger
}

func NewBlockProcessor(db *database.DB, logger *zap.Logger) *BlockProcessor {
	return &BlockProcessor{
		db:     db,
		logger: logger.Sugar(),
	}
}

func (bp *BlockProcessor) ProcessBlock(ctx context.Context, block *types.Block, chainID int64) error {
	blockModel := &models.Block{
		ChainID:          chainID,
		BlockNumber:      block.Number().Int64(),
		Hash:             block.Hash().Hex(),
		ParentHash:       block.ParentHash().Hex(),
		Nonce:            toStringPtr(fmt.Sprintf("0x%x", block.Nonce())),
		Sha3Uncles:       toStringPtr(block.UncleHash().Hex()),
		Miner:            block.Coinbase().Hex(),
		StateRoot:        toStringPtr(block.Root().Hex()),
		TransactionsRoot: toStringPtr(block.TxHash().Hex()),
		ReceiptsRoot:     toStringPtr(block.ReceiptHash().Hex()),
		Difficulty:       toStringPtr(block.Difficulty().String()),
		TotalDifficulty:  nil,
		Size:             toInt64Ptr(int64(block.Size())),
		GasLimit:         int64(block.GasLimit()),
		GasUsed:          int64(block.GasUsed()),
		Timestamp:        time.Unix(int64(block.Time()), 0),
		ExtraData:        toStringPtr(fmt.Sprintf("0x%x", block.Extra())),
		MixHash:          toStringPtr(block.MixDigest().Hex()),
		BaseFeePerGas:    bigIntToStringPtr(block.BaseFee()),
		TxCount:          len(block.Transactions()),
	}

	if err := bp.db.InsertBlock(ctx, blockModel); err != nil {
		return fmt.Errorf("Failed to insert block %d: %w", block.Number().Int64(), err)
	}
	bp.logger.Debugw("Processed block", "block_number", block.Number().Int64(), "hash", block.Hash().Hex(), "tx_count", len(block.Transactions()))
	return nil
}

func toStringPtr(s string) *string {
	if s == " " || s == "0x" {
		return nil
	}
	return &s
}

func toInt64Ptr(i int64) *int64 {
	return &i
}

func bigIntToStringPtr(bi *big.Int) *string {
	if bi == nil {
		return nil
	}
	s := bi.String()
	return &s
}
