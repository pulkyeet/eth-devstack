package main

import (
	"context"
	"log"
	"time"

	"github.com/pulkyeet/eth-devstack/backend/internal/blockchain"
	"github.com/pulkyeet/eth-devstack/backend/internal/config"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	logger, err := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatal("Failed to initialise logger:", err)
	}
	defer logger.Sync()

	// Use sugared logger for easier key-value logging
	sugar := logger.Sugar()
	sugar.Infow("Starting Ethereum Explorer API",
		"port", cfg.Server.Port,
		"environment", cfg.Server.Environment,
	)

	db, err := database.NewDB(
		cfg.Database.ConnectionString(),
		cfg.Database.MaxConnections,
		cfg.Database.MaxIdleConns,
		logger,
	)
	if err != nil {
		sugar.Fatalw("Failed to initialise database", "error", err)
	}
	defer db.Close()

	chainManager, err := blockchain.NewChainManager(cfg.Chains.ConfigPath, logger)
	if err != nil {
		sugar.Fatalw("Failed to initialise chain manager", "error", err)
	}
	defer chainManager.Close()

	client, err := chainManager.GetClient(1337)
	if err != nil {
		sugar.Fatalw("Failed to get chain client", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blockNum, err := client.GetLatestBlockNumber(ctx)
	if err != nil {
		sugar.Fatalw("Failed to get latest block", "error", err)
	}

	sugar.Infow("Connected to blockchain", "chain_id", client.ChainID(), "latest_block", blockNum)

	latestBlock, err := db.GetLatestBlock(ctx, 1337)
	if err != nil {
		sugar.Fatalw("Failed to query database", "error", err)
	}

	if latestBlock != nil {
		sugar.Infow("Latest block in database", "block_numer", latestBlock.BlockNumber, "hash", latestBlock.Hash)
	} else {
		sugar.Info("No blocks in database yet")
	}
	sugar.Info("Phase 2D complete - database layer working!")
}
