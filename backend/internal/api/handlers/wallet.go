package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"github.com/pulkyeet/eth-devstack/backend/internal/wallet"
)

type WalletHandler struct {
	service *wallet.Service
}

func NewWalletHandler(service *wallet.Service) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

// POST /api/v1/wallet/create
func (h *WalletHandler) CreateWallet(c *fiber.Ctx) error {
	var req wallet.CreateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
	}

	// Validate
	if req.Password == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_PASSWORD", "Password is required", nil)
	}
	if len(req.Password) < 12 {
		return responses.Error(c, fiber.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 12 characters", nil)
	}

	resp, err := h.service.CreateWallet(c.Context(), req)
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, "CREATION_FAILED", err.Error(), nil)
	}

	c.Status(fiber.StatusCreated)
	return responses.Success(c, resp, nil)
}

// POST /api/v1/wallet/import
func (h *WalletHandler) ImportWallet(c *fiber.Ctx) error {
	var req wallet.ImportWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
	}

	// Validate
	if req.Password == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_PASSWORD", "Password is required", nil)
	}
	if len(req.Password) < 12 {
		return responses.Error(c, fiber.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 12 characters", nil)
	}
	if req.Method != "mnemonic" && req.Method != "private_key" {
		return responses.Error(c, fiber.StatusBadRequest, "INVALID_METHOD", "Method must be 'mnemonic' or 'private_key'", nil)
	}

	w, err := h.service.ImportWallet(c.Context(), req)
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, "IMPORT_FAILED", err.Error(), nil)
	}

	// Don't return encrypted keys in response
	c.Status(fiber.StatusCreated)
	return responses.Success(c, fiber.Map{
		"id":      w.ID,
		"address": w.Address,
		"name":    w.Name,
	}, nil)
}

// GET /api/v1/wallet/:id
func (h *WalletHandler) GetWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_ID", "Wallet ID is required", nil)
	}

	w, err := h.service.GetWallet(c.Context(), id)
	if err != nil {
		return responses.Error(c, fiber.StatusNotFound, "WALLET_NOT_FOUND", err.Error(), nil)
	}

	// Don't return encrypted keys
	return responses.Success(c, fiber.Map{
		"id":              w.ID,
		"address":         w.Address,
		"name":            w.Name,
		"derivation_path": w.DerivationPath,
		"is_imported":     w.IsImported,
		"created_at":      w.CreatedAt,
	}, nil)
}

// GET /api/v1/wallet/:id/balance
func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_ID", "Wallet ID is required", nil)
	}

	// Check if should update from chain
	updateFromChain := c.QueryBool("update", false)

	balances, err := h.service.GetBalances(c.Context(), id, updateFromChain)
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, "BALANCE_FETCH_FAILED", err.Error(), nil)
	}

	return responses.Success(c, fiber.Map{
		"wallet_id": id,
		"balances":  balances,
	}, nil)
}

// POST /api/v1/wallet/send
func (h *WalletHandler) SendTransaction(c *fiber.Ctx) error {
	var req wallet.SendTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
	}

	// Validate
	if req.WalletID == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_WALLET_ID", "Wallet ID is required", nil)
	}
	if req.To == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_TO", "Recipient address is required", nil)
	}
	if req.Value == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_VALUE", "Value is required", nil)
	}
	if req.Password == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_PASSWORD", "Password is required", nil)
	}

	resp, err := h.service.SendTransaction(c.Context(), req)
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, "TRANSACTION_FAILED", err.Error(), nil)
	}

	return responses.Success(c, resp, nil)
}

// POST /api/v1/wallet/sign
func (h *WalletHandler) SignMessage(c *fiber.Ctx) error {
	var req wallet.SignMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
	}

	// Validate
	if req.WalletID == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_WALLET_ID", "Wallet ID is required", nil)
	}
	if req.Message == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_MESSAGE", "Message is required", nil)
	}
	if req.Password == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_PASSWORD", "Password is required", nil)
	}

	resp, err := h.service.SignMessage(c.Context(), req)
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, "SIGNING_FAILED", err.Error(), nil)
	}

	return responses.Success(c, resp, nil)
}

// GET /api/v1/wallet/:id/transactions
func (h *WalletHandler) GetTransactions(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return responses.Error(c, fiber.StatusBadRequest, "MISSING_ID", "Wallet ID is required", nil)
	}

	// Optional chain filter
	var chainID *int64
	if c.Query("chain_id") != "" {
		cid := c.QueryInt("chain_id", 0)
		if cid > 0 {
			cidInt64 := int64(cid)
			chainID = &cidInt64
		}
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	if limit > 100 {
		limit = 100
	}

	txs, err := h.service.GetTransactions(c.Context(), id, chainID, limit, offset)
	if err != nil {
		return responses.Error(c, fiber.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
	}

	return responses.Success(c, fiber.Map{
		"wallet_id":    id,
		"transactions": txs,
	}, nil)
}