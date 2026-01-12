package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
)

type TransactionHandler struct {
	db *database.DB
}

func NewTransactionHandler(db *database.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) GetTransactions(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	
	// NEW: Support filtering by block_number
	blockNumber := c.QueryInt("block_number", 0)
	
	var txs interface{}
	var err error
	
	if blockNumber > 0 {
		// Get transactions for specific block
		txs, err = h.db.GetTransactionsByBlock(c.Context(), int64(chainID), int64(blockNumber))
	} else {
		// Get all transactions (paginated)
		txs, err = h.db.GetTransactions(c.Context(), int64(chainID), limit, offset)
	}
	
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch transactions", err.Error())
	}

	cID := int64(chainID)
	
	if blockNumber > 0 {
		// No pagination for block-specific queries
		return responses.Success(c, fiber.Map{
			"transactions": txs,
		}, &cID)
	}
	
	// With pagination
	count, _ := h.db.CountTransactions(c.Context(), int64(chainID))
	totalPages := (int(count) + limit - 1) / limit

	return responses.Success(c, fiber.Map{
		"transactions": txs,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       count,
			"total_pages": totalPages,
		},
	}, &cID)
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	hash := c.Params("hash")

	tx, err := h.db.GetTransactionByHash(c.Context(), int64(chainID), hash)
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch transaction", err.Error())
	}
	if tx == nil {
		return responses.Error(c, 404, "RESOURCE_NOT_FOUND", "Transaction not found", nil)
	}

	// Get logs for this transaction
	logs, _ := h.db.GetLogsByTransaction(c.Context(), int64(chainID), hash)

	cID := int64(chainID)
	return responses.Success(c, fiber.Map{
		"hash":             tx.Hash,
		"block_number":     tx.BlockNumber,
		"block_hash":       tx.BlockHash,
		"timestamp":        tx.Timestamp,
		"from_address":     tx.FromAddress,
		"to_address":       tx.ToAddress,
		"value":            tx.Value,
		"gas":              tx.Gas,
		"gas_price":        tx.GasPrice,
		"gas_used":         tx.GasUsed,
		"effective_gas_price": tx.EffectiveGasPrice,
		"nonce":            tx.Nonce,
		"transaction_index": tx.TransactionIndex,
		"transaction_type": tx.TransactionType,
		"input":            tx.Input,
		"status":           tx.Status,
		"contract_address": tx.ContractAddress,
		"logs":             logs,
	}, &cID)
}