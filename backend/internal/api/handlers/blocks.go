package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
)

type BlockHandler struct {
	db *database.DB
}

func NewBlockHandler(db *database.DB) *BlockHandler {
	return &BlockHandler{db: db}
}

func (h *BlockHandler) GetBlocks(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	blocks, err := h.db.GetBlocks(c.Context(), int64(chainID), limit, offset)
	if err !=nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch blocks", err.Error())
	}

	total, _ := h.db.CountBlocks(c.Context(), int64(chainID))
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	cID := int64(chainID)
	return responses.Success(c, fiber.Map{
		"blocks": blocks,
		"pagination": responses.PaginationMeta{
			Page: page,
			Limit: limit,
			Total: total,
			TotalPages: totalPages,
		},
	}, &cID)
}

func (h *BlockHandler) GetBlock(c *fiber.Ctx) error {
	chainID := c.QueryInt("chain_id", 1337)
	blockID := c.Params("id")

	var block interface{}
	var err error

	if num, parseErr := strconv.ParseInt(blockID, 10,64); parseErr!=nil {
		block, err = h.db.GetBlockByNumber(c.Context(), int64(chainID), num)
	} else {
		block, err = h.db.GetBlockByHash(c.Context(), int64(chainID), blockID)
	}

	if err != nil {
		return responses.Error(c, 500, "DATABASE_ERROR", "Failed to fetch block", err.Error())
	}

	if block == nil {
		return responses.Error(c, 404, "RESOURCE_NOT_FOUND", "Block not found", nil)
	}

	cID := int64(chainID)
	return responses.Success(c, block, &cID)
}