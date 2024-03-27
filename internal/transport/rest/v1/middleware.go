package v1

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp"
)

func (h *Handler) loggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) {
		slog.Info("request", "method", c.Method(), "URI", c.OriginalURL())
		c.Next()
	}
}

func (h *Handler) authMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) {
		token, err := getTokenFromRequest(&c.Fasthttp.Request)
		if err != nil {
			newErrorResponse(c, fiber.StatusUnauthorized, err.Error())
			return
		}

		if _, err = h.auth.Parse(token); err != nil {
			newErrorResponse(c, fiber.StatusUnauthorized, err.Error())
			return
		}

		if !h.auth.IsTokenExpired(token) {
			newErrorResponse(c, fiber.StatusUnauthorized, "token is Expired")
			return
		}

		c.Next()
	}
}

func getTokenFromRequest(r *fasthttp.Request) (string, error) {
	header := r.Header.Peek("Authorization")

	strHeader := string(header)
	if strHeader == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(strHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
