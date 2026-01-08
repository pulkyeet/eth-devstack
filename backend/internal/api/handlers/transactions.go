package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
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

	txs, err := h.db.GetTransactions(c.Context(), int64(chainID), limit, offset)
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch transactions", err.Error())
	}

	total, _ := h.db.CountTransactions(c.Context(), int64(chainID))
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	cID := int64(chainID)
	return responses.Success(c, fiber.Map{
		"transactions": txs,
		"pagination": responses.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
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

	cID := int64(chainID)
	return responses.Success(c, tx, &cID)
}