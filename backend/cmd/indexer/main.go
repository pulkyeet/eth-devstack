package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pulkyeet/eth-devstack/backend/internal/blockchain"
	"github.com/pulkyeet/eth-devstack/backend/internal/config"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/indexer"
	"github.com/pulkyeet/eth-devstack/backend/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	logger, err := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatal("Failed to initialise logger", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("Starting Ethereum Indexer")

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

	indexerService := indexer.NewService(db, chainManager, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		sugar.Info("Shutdown signal received")
		indexerService.Stop()
		cancel()
	}()

	if err := indexerService.Start(ctx); err != nil {
		sugar.Fatalw("Indexed failed", "error", err)
	}

	sugar.Info("Indexer shutdown complete")
}
