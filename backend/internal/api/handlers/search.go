package handlers

import (
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
)

type SearchHandler struct {
	db *database.DB
}

func NewSearchHandler(db *database.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	query := c.Query("q")

	if query == "" {
		return responses.Error(c, 400, "INVALID_QUERY", "Search query required", nil)
	}

	if matched, _ := regexp.MatchString(`0x[a-fA-F0-9]{64}$`, query); matched {
		if block, _ := h.db.GetBlockByHash(c.Context(), int64(chainID), query); block != nil {
			cID := int64(chainID)
			return responses.Success(c, fiber.Map{"type": "block", "result": block}, &cID)
		}
		if tx, _ := h.db.GetTransactionByHash(c.Context(), int64(chainID), query); tx != nil {
			cID := int64(chainID)
			return responses.Success(c, fiber.Map{"type": "transaction", "result": tx}, &cID)
		}
	} else if matched, _ := regexp.MatchString(`^0x[a-fA-F0-9]{40}$`, query); matched {
		if addr, _ := h.db.GetAddress(c.Context(), int64(chainID), query); addr != nil {
			cID := int64(chainID)
			return responses.Success(c, fiber.Map{"type": "address", "result": addr}, &cID)
		}
	} else if num, err := strconv.ParseInt(query, 10, 64); err == nil {
		if block, _ := h.db.GetBlockByNumber(c.Context(), int64(chainID), num); block != nil {
			cID := int64(chainID)
			return responses.Success(c, fiber.Map{"type": "block", "result": block}, &cID)
		}
	}

	return responses.Error(c, 404, "RESOURCE_NOT_FOUND", "No results found", nil)
}
