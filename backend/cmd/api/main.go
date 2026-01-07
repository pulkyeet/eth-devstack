package main

import (
	"log"

	"github.com/pulkyeet/ethereum-explorer/internal/config"
	"github.com/pulkyeet/ethereum-explorer/internal/utils"
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
		"db_host", cfg.Database.Host,
	)
	sugar.Info("API server would start here...")
}
