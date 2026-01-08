package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
)

type AddressHandler struct {
	db *database.DB
}

func NewAddressHandler(db *database.DB) *AddressHandler {
	return &AddressHandler{db: db}
}

func (h *AddressHandler) GetAddress(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	address := c.Params("address")

	addr, err := h.db.GetAddress(c.Context(), int64(chainID), address)
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch address", err.Error())
	}
	if addr == nil {
		return responses.Error(c, 404, "RESOURCE_NOT_FOUND", "Address not found", nil)
	}

	cID := int64(chainID)
	return responses.Success(c, addr, &cID)
}

func (h *AddressHandler) GetAddressTransactions(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	address := c.Params("address")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	txs, err := h.db.GetTransactionsByAddress(c.Context(), int64(chainID), address, limit, offset)
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch transactions", err.Error())
	}

	cID := int64(chainID)
	return responses.Success(c, fiber.Map{
		"address":      address,
		"transactions": txs,
	}, &cID)
}