package v1

import (
	"github.com/Woodfyn/file-api/internal/service"
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

// routing with fiber
func (h *Handler) Init(api fiber.Router) {
	api.Use(h.loggingMiddleware())

	v1 := api.Group("/v1")
	{
		h.initUserRouter(v1)
		h.initFileRouter(v1)
	}
}
