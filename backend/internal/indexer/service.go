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

type Service struct {
	db             *database.DB
	chainManager   *blockchain.ChainManager
	blockProcessor *BlockProcessor
	logger         *zap.SugaredLogger
	stopChan       chan struct{}
	batchSize      int
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
	if err != nil {
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
			if err := s.syncChain(ctx, client, txProcessor, chainID); err != nil {
				logger.Errorw("Sync error", "error", err)
			}
		}
	}
}

func (s *Service) syncChain(ctx context.Context, client *blockchain.ChainClient, txProcessor *TxProcessor, chainID int64) error {
	latestChainBlock, err := client.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get latest block number: %w", err)
	}
	latestDBBlock, err := s.db.GetLatestBlock(ctx, chainID)
	if err != nil {
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
		if err := s.processBlock(ctx, client, txProcessor, blockNum, chainID); err != nil {
			return fmt.Errorf("Failed to process block %d: %w", blockNum, err)
		}
	}
	s.logger.Infow("Sync batch complete", "chain_id", chainID, "synced_from", startBlock, "synced_to", endBlock)
	return nil
}

func (s *Service) processBlock(ctx context.Context, client *blockchain.ChainClient, txProcessor *TxProcessor, blockNum int64, chainID int64) error {
	block, err := client.GetBlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return fmt.Errorf("Failed to get block: %w", err)
	}

	if block == nil {
		return fmt.Errorf("Block %d not found", blockNum)
	}

	if err := s.blockProcessor.ProcessBlock(ctx, block, chainID); err != nil {
		return err
	}

	blockTime := time.Unix(int64(block.Time()), 0)
	for txIndex, tx := range block.Transactions() {
		if err := txProcessor.ProcessTransaction(ctx, tx, block.NumberU64(), block.Hash().Hex(), txIndex, blockTime, chainID); err != nil {
			s.logger.Warnw("Failed to process transaction", "block", blockNum, "tx_hash", tx.Hash().Hex(), "error", err)
			continue
		}

		// Get receipt for log processing
		receipt, err := client.GetTransactionReceipt(ctx, tx.Hash().Hex())
		if err == nil && receipt != nil {
			// Process logs
			s.processLogs(ctx, receipt, tx, blockTime, chainID)
		}

		// Update addresses
		s.updateAddresses(ctx, tx, blockNum, blockTime, chainID)
	}
	return nil
}

func (s *Service) processLogs(ctx context.Context, receipt *types.Receipt, tx *types.Transaction, blockTime time.Time, chainID int64) error {
	for _, log := range receipt.Logs {
		logModel := &models.TransactionLog{
			ChainID:          chainID, // Changed from s.chainID
			TransactionHash:  tx.Hash().Hex(),
			LogIndex:         int(log.Index),
			Address:          log.Address.Hex(),
			Data:             toStringPtr(fmt.Sprintf("0x%x", log.Data)),
			BlockNumber:      int64(log.BlockNumber),
			BlockHash:        log.BlockHash.Hex(),
			TransactionIndex: int(log.TxIndex),
			Removed:          log.Removed,
		}

		if len(log.Topics) > 0 {
			logModel.Topic0 = toStringPtr(log.Topics[0].Hex())
		}
		if len(log.Topics) > 1 {
			logModel.Topic1 = toStringPtr(log.Topics[1].Hex())
		}
		if len(log.Topics) > 2 {
			logModel.Topic2 = toStringPtr(log.Topics[2].Hex())
		}
		if len(log.Topics) > 3 {
			logModel.Topic3 = toStringPtr(log.Topics[3].Hex())
		}

		if err := s.db.InsertLog(ctx, logModel); err != nil {
			s.logger.Warnw("Failed to insert log", "error", err)
		}

		// Detect ERC20 Transfer events
		if len(log.Topics) == 3 && log.Topics[0].Hex() == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
			s.processERC20Transfer(ctx, log, tx, blockTime, chainID)
		}
	}
	return nil
}

func (s *Service) processERC20Transfer(ctx context.Context, log *types.Log, tx *types.Transaction, blockTime time.Time, chainID int64) {
	// Ensure token exists
	token := &models.Token{
		ChainID: chainID, // Changed from s.chainID
		Address: log.Address.Hex(),
		Type:    "ERC20",
	}
	s.db.UpsertToken(ctx, token)

	// Parse transfer
	from := common.HexToAddress(log.Topics[1].Hex())
	to := common.HexToAddress(log.Topics[2].Hex())
	value := new(big.Int).SetBytes(log.Data)

	transfer := &models.TokenTransfer{
		ChainID:         chainID, // Changed from s.chainID
		TransactionHash: tx.Hash().Hex(),
		LogIndex:        int(log.Index),
		TokenAddress:    log.Address.Hex(),
		FromAddress:     from.Hex(),
		ToAddress:       to.Hex(),
		Value:           toStringPtr(value.String()),
		BlockNumber:     int64(log.BlockNumber),
		Timestamp:       blockTime,
	}
	s.db.InsertTokenTransfer(ctx, transfer)

	// Update balances (simplified - real impl would call balanceOf)
	if from.Hex() != "0x0000000000000000000000000000000000000000" {
		s.db.UpsertTokenBalance(ctx, &models.TokenBalance{
			ChainID:       chainID, // Changed from s.chainID
			TokenAddress:  log.Address.Hex(),
			HolderAddress: from.Hex(),
			Balance:       "0", // Placeholder
		})
	}
	if to.Hex() != "0x0000000000000000000000000000000000000000" {
		s.db.UpsertTokenBalance(ctx, &models.TokenBalance{
			ChainID:       chainID, // Changed from s.chainID
			TokenAddress:  log.Address.Hex(),
			HolderAddress: to.Hex(),
			Balance:       value.String(),
		})
	}
}

func (s *Service) updateAddresses(ctx context.Context, tx *types.Transaction, blockNum int64, blockTime time.Time, chainID int64) error {
	signer := types.LatestSignerForChainID(big.NewInt(chainID))
	from, _ := types.Sender(signer, tx)

	// Update from address
	fromAddr := &models.Address{
		ChainID:        chainID, // Changed from s.chainID
		Address:        from.Hex(),
		Balance:        0,
		Nonce:          int64(tx.Nonce()),
		FirstSeenBlock: &blockNum,
		LastSeenBlock:  &blockNum,
		FirstSeenAt:    &blockTime,
		LastSeenAt:     &blockTime,
	}
	s.db.UpsertAddress(ctx, fromAddr)

	// Update to address
	if tx.To() != nil {
		toAddr := &models.Address{
			ChainID:        chainID, // Changed from s.chainID
			Address:        tx.To().Hex(),
			Balance:        0,
			Nonce:          0,
			FirstSeenBlock: &blockNum,
			LastSeenBlock:  &blockNum,
			FirstSeenAt:    &blockTime,
			LastSeenAt:     &blockTime,
		}
		s.db.UpsertAddress(ctx, toAddr)
	}

	return nil
}
