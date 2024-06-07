package srv

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func NewServer(app *fiber.App) *Server {
	return &Server{app: app}
}

func (s *Server) Run(port string) error {
	return s.app.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
