package srv

import (
	"time"

	"github.com/gofiber/fiber"
)

type Server struct {
	app *fiber.App
}

func NewServer(app *fiber.App) *Server {
	return &Server{app: app}
}

func (s *Server) Run(port string) error {
	settings := fiber.Settings{
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	*s.app.Settings = settings

	return s.app.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
