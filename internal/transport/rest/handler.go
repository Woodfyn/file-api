package rest

import (
	"context"
	"mime/multipart"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type authService interface {
	SignUp(ctx context.Context, input core.SignUpRequest) error
	SignIn(ctx context.Context, input core.SignInRequest) (*core.TokenResp, error)
	Refresh(ctx context.Context, refreshToken string) (*core.TokenResp, error)
	Parse(token string) (string, error)
	IsTokenExpired(token string) bool
}

type fileService interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader, file multipart.File, userId string) error
	GetFiles(ctx context.Context, userId string) ([]*core.GetAllFilesResp, error)
}

type Handler struct {
	fileService fileService
	authService authService
}

func NewHandler(fileService fileService, authService authService) *Handler {
	return &Handler{
		fileService: fileService,
		authService: authService,
	}
}

func (h *Handler) Init() *fiber.App {
	r := fiber.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080, http://127.0.0.1:8080",
		AllowMethods:     "GET, POST, PATCH, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	h.loggingMiddleware()

	api := r.Group("/api")
	{
		h.initUserRouter(api)
		h.initFileRouter(api)
	}

	return r
}
