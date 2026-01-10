package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"go.uber.org/zap"
)

type StreamHandler struct {
	db     *database.DB
	logger *zap.SugaredLogger
}

func NewStreamHandler(db *database.DB, logger *zap.Logger) *StreamHandler {
	return &StreamHandler{
		db:     db,
		logger: logger.Sugar(),
	}
}

func (h *StreamHandler) StreamBlocks(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("X-Accel-Buffering", "no")

	ctx := context.Background()
	lastBlockNum := int64(0)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		// Send initial connection message
		fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"streaming\"}\n\n")
		if err := w.Flush(); err != nil {
			h.logger.Warnw("Failed to flush initial message", "error", err)
			return
		}

		for {
			select {
			case <-ticker.C:
				block, err := h.db.GetLatestBlock(ctx, int64(chainID))
				if err != nil {
					h.logger.Warnw("Failed to get latest block", "error", err)
					continue
				}

				if block == nil {
					continue
				}

				// Only send if new block
				if block.BlockNumber > lastBlockNum {
					lastBlockNum = block.BlockNumber
					
					data := map[string]interface{}{
						"block_number": block.BlockNumber,
						"hash":         block.Hash,
						"timestamp":    block.Timestamp.Format(time.RFC3339),
						"tx_count":     block.TxCount,
						"gas_used":     block.GasUsed,
					}

					jsonData, err := json.Marshal(data)
					if err != nil {
						h.logger.Warnw("Failed to marshal block data", "error", err)
						continue
					}

					_, writeErr := fmt.Fprintf(w, "event: block\ndata: %s\n\n", string(jsonData))
					if writeErr != nil {
						h.logger.Info("Client disconnected")
						return
					}

					if err := w.Flush(); err != nil {
						h.logger.Info("Client disconnected during flush")
						return
					}
				}
			}
		}
	})

	return nil
}