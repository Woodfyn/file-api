package rest

import (
	"log/slog"
	"strings"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) loggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		slog.Info("request", "method", c.Method(), "URI", c.OriginalURL())

		return c.Next()
	}
}

func (h *Handler) authMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Request().Header.Peek("Authorization")

		strHeader := string(header)
		if strHeader == "" {
			newErrorResponse(fiber.StatusUnauthorized, core.ErrEmptyAuthHeader)
		}

		headerParts := strings.Split(strHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return newErrorResponse(fiber.StatusUnauthorized, core.ErrInvalidAuthHeader)
		}

		if len(headerParts[1]) == 0 {
			return newErrorResponse(fiber.StatusUnauthorized, core.ErrAccessTokenEmpty)
		}

		id, err := h.authService.Parse(headerParts[1])
		if err != nil {
			return newErrorResponse(fiber.StatusUnauthorized, err)
		}

		if !h.authService.IsTokenExpired(headerParts[1]) {
			return newErrorResponse(fiber.StatusUnauthorized, core.ErrAccessTokenIsExpired)
		}

		c.Locals("userId", id)

		return c.Next()
	}
}
