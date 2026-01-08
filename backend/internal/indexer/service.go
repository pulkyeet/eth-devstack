package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"
	"github.com/pulkyeet/eth-devstack/backend/internal/blockchain"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"go.uber.org/zap"
)

type Service struct {
	db *database.DB
	chainManager *blockchain.ChainManager
	blockProcessor *BlockProcessor
	logger *zap.SugaredLogger
	stopChan chan struct{}
	batchSize int
}

func NewService(db *database.DB, chainManager *blockchain.ChainManager, logger *zap.Logger) *Service {
	return &Service{
		db:             db,
		chainManager:   chainManager,
		blockProcessor: NewBlockProcessor(db, logger),
		logger:         logger.Sugar(),
		stopChan:       make(chan struct{}),
		batchSize:      100, // Process 100 blocks at a time
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.logger.Info("Starting indexer service")

	chains := s.chainManager.GetActiveChains()
	if len(chains) == 0 {
		return fmt.Errorf("No active chains to index")
	}

	for _, chainConfig := range chains {
		go s.indexChain(ctx, chainConfig.ChainID)
	}

	<-s.stopChan
	s.logger.Info("Indexer service stopped")
	return nil
}

func (s *Service) Stop() {
	close(s.stopChan)
}

func (s *Service) indexChain(ctx context.Context, chainID int64) {
	logger := s.logger.With("chain_id", chainID)
	logger.Info("Starting chain indexer")

	client, err := s.chainManager.GetClient(chainID)
	if err!=nil {
		logger.Errorw("Failed to get chain client", "error", err)
	}

	txProcessor := NewTxProcessor(s.db, client, s.logger.Desugar().Sugar().Desugar())

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled, stoppping indexer")
			return
		case <-s.stopChan:
			logger.Info("Stop signal received. Stopping indexer")
			return
		case <-ticker.C:
			if err := s.syncChain(ctx, client, txProcessor, chainID); err!=nil {
				logger.Errorw("Sync error", "error", err)
			}
		}
	}
}

func (s *Service) syncChain(ctx context.Context, client *blockchain.ChainClient, txProcessor *TxProcessor, chainID int64) error {
	latestChainBlock, err := client.GetLatestBlockNumber(ctx)
	if err!=nil{
		return fmt.Errorf("Failed to get latest block number: %w", err)
	}
	latestDBBlock, err := s.db.GetLatestBlock(ctx, chainID)
	if err!=nil {
		return fmt.Errorf("Failed to get latest DB block: %w", err)
	}

	var startBlock int64 
	if latestDBBlock == nil {
		startBlock = 0
	} else {
		startBlock = latestDBBlock.BlockNumber + 1
	}

	blocksToSync := int64(latestChainBlock) - startBlock
	if blocksToSync <= 0 {
		return nil
	}
	
	endBlock := startBlock + int64(s.batchSize)
	if endBlock > int64(latestChainBlock) {
		endBlock = int64(latestChainBlock)
	}

	s.logger.Infow("Syncing blocks", "chain_id", chainID, "from", startBlock, "to", endBlock, "total_behind", blocksToSync)

	for blockNum := startBlock; blockNum <= endBlock; blockNum++ {
		if err := s.processBlock(ctx, client, txProcessor, blockNum, chainID); err!=nil {
			return fmt.Errorf("Failed to process block %d: %w", blockNum, err)
		}
	}
	s.logger.Infow("Sync batch complete", "chain_id", chainID, "synced_from", startBlock, "synced_to", endBlock)
	return nil
}

func (s *Service) processBlock(ctx context.Context, client *blockchain.ChainClient, txProcessor *TxProcessor, blockNum int64, chainID int64) error {
	block, err := client.GetBlockByNumber(ctx, big.NewInt(blockNum))
	if err!=nil {
		return fmt.Errorf("Failed to get block: %w", err)
	}

	if block == nil {
		return fmt.Errorf("Block %d not found", blockNum)
	}

	if err := s.blockProcessor.ProcessBlock(ctx, block, chainID); err!=nil {
		return err
	}

	blockTime := time.Unix(int64(block.Time()), 0)
	for txIndex, tx := range block.Transactions() {
		if err := txProcessor.ProcessTransaction(ctx, tx, block.NumberU64(), block.Hash().Hex(), txIndex, blockTime, chainID); err!=nil {
			s.logger.Warnw("Failed to process transaction", "block", blockNum, "tx_hash", tx.Hash().Hex(), "error", err)
		}
	}
	return nil
}