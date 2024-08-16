package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yoruakio/gowebserver/logger"
)

func Initialize() *fiber.App {
	logger.Info("Initializing HTTP Server")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return app
}

func Start(app *fiber.App) {
	logger.Info("Starting HTTP Server")

	log.Fatal(app.ListenTLS(":443", "ssl/server.crt", "ssl/server.key"))
}
