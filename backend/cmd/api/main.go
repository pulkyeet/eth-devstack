package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pulkyeet/eth-devstack/backend/internal/config"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/utils"
	"github.com/pulkyeet/eth-devstack/backend/internal/api"
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

	server := api.NewServer(db, logger, cfg.Server.Port)

	go func() {
		if err := server.Start(); err != nil {
			sugar.Fatalw("Server error", "error", err)
		}
	}()

	sugar.Info("API Server started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down server...")
	if err := server.Shutdown(); err != nil {
		sugar.Errorw("Error shutting down", "error", err)
	}
	sugar.Info("Server stopped")
}
