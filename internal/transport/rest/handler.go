package rest

import (
	"github.com/Woodfyn/file-api/internal/service"
	v1 "github.com/Woodfyn/file-api/internal/transport/rest/v1"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/gofiber/fiber"
)

type Handler struct {
	service *service.Service
	auth    auth.TokenManager
}

func NewHandler(service *service.Service, auth auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

func (h *Handler) Init() *fiber.App {
	r := fiber.New()

	r.Get("/ping", func(c *fiber.Ctx) {
		c.SendString("pong")
	})

	h.initApi(r)

	return r
}

func (h *Handler) initApi(r *fiber.App) {
	handlerV1 := v1.NewHandler(h.service, h.auth)
	api := r.Group("/api")
	{
		handlerV1.Init(api)
	}
}
