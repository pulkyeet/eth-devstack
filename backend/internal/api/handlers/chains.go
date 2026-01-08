package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
)

type ChainHandler struct {
	db *database.DB
}

func NewChainHandler(db *database.DB) *ChainHandler {
	return &ChainHandler{db: db}
}

func (h *ChainHandler) GetChains(c *fiber.Ctx) error {
	chains, err := h.db.GetChains(c.Context())
	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch chains", err.Error())
	}

	return responses.Success(c, fiber.Map{"chains": chains}, nil)
}

func (h *ChainHandler) GetHealth(c *fiber.Ctx) error {
	err := h.db.Ping(c.Context())
	status := "healthy"
	dbStatus := "connected"
	if err != nil {
		status = "unhealthy"
		dbStatus = "disconnected"
	}

	return responses.Success(c, fiber.Map{
		"status":   status,
		"database": dbStatus,
	}, nil)
}