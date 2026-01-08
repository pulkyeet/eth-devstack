package responses

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta"`
}

type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Meta struct {
	Timestamp string `json:"timestamp"`
	ChainID   *int64 `json:"chain_id,omitempty"`
	Version   string `json:"version"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(c *fiber.Ctx, data interface{}, chainID *int64) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			ChainID:   chainID,
			Version:   "1.0.0",
		},
	})
}

func Error(c *fiber.Ctx, statusCode int, code, message string, details interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   "1.0.0",
		},
	})
}