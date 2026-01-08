package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) fiber.Handler {
	sugar := logger.Sugar()

	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				sugar.Errorw("Panic recovered", "error", r, "method", c.Method(), "path", c.Path())
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "INTERNAL_SERVER_ERROR",
						"message": "Internal Server Error",
					},
				})
			}
		}()
		return c.Next()
	}
}
