package transport

import (
	"github.com/Woodfyn/file-api/internal/transport/rest"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	rest *rest.Handler
}

func NewHandler(rest *rest.Handler) *Handler {
	return &Handler{
		rest: rest,
	}
}

func (h *Handler) Init() *fiber.App {
	httpRouter := h.rest.Init()

	return httpRouter
}
