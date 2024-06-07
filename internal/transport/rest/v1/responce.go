package v1

import (
	"github.com/gofiber/fiber"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *fiber.Ctx, statusCode int, message string) {
	c.Status(statusCode).JSON(errorResponse{Message: message})
}

func newSuccessResponse(c *fiber.Ctx, statusCode int, message string) {
	c.Status(statusCode).JSON(statusResponse{Status: message})
}
