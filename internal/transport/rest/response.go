package rest

import (
	"github.com/gofiber/fiber/v2"
)

type dataResponse struct {
	Data interface{} `json:"data"`
}

func newErrorResponse(statusCode int, err error) *fiber.Error {
	return fiber.NewError(statusCode, err.Error())
}

func newDataResponse(c *fiber.Ctx, statusCode int, data interface{}) {
	c.Status(statusCode).JSON(dataResponse{Data: data})
}
