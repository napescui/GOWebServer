package http

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yoruakio/gowebserver/config"
	"github.com/yoruakio/gowebserver/logger"
)

func Initialize() *fiber.App {
	logger.Info("Initializing HTTP Server")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             4 * 1024 * 1024, // 4MB of body limit
		IdleTimeout:           10 * time.Second,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
	})

	// logger middleware
	app.Use(func(c *fiber.Ctx) error {
		if config.GetConfig().Logger {
			logger.Infof("[%s] %s %s => %d", c.IP(), c.Method(), c.Path(), c.Response().StatusCode())
		}
		return c.Next()
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
