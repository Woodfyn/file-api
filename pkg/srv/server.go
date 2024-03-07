package server

import (
	"crypto/tls"
	"time"

	"github.com/gofiber/fiber"
)

func Run(port string, app *fiber.App) error {
	settings := fiber.Settings{
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	*app.Settings = settings

	return app.Listen(":"+port, &tls.Config{})
}

func Shutdown(app *fiber.App) error {
	return app.Shutdown()
}
