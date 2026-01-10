package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
)

type StatsHandler struct {
	db *database.DB
}

func NewStatsHandler(db *database.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

func (h *StatsHandler) GetStats(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)

	stats, err := h.db.GetNetworkStats(c.Context(), int64(chainID))
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch stats", err.Error())
	}

	cID := int64(chainID)
	return responses.Success(c, stats, &cID)
}