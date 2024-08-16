package http

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Middleware
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(compress.New())

	app.Use(func(c *fiber.Ctx) error {
		if config.GetConfig().Logger {
			logger.Infof("[%s] %s %s => %d", c.IP(), c.Method(), c.Path(), c.Response().StatusCode())
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	meta := fmt.Sprintf("K10WA_%d", rand.Intn(9000)+1000)
	config := config.GetConfig()
	content := fmt.Sprintf(
		"server|%s\n"+
			"port|%s\n"+
			"type|1\n"+
			"# maint|Server is currently down for maintenance. We will be back soon!\n"+
			"loginurl|%s\n"+
			"meta|%s\n"+
			"RTENDMARKERBS1001",
		config.Host, config.Port, config.LoginUrl, meta)

	app.Get("/growtopia/server_data.php", func(c *fiber.Ctx) error {
		return c.SendString(content)
	})

	return app
}

func Start(app *fiber.App) {
	logger.Info("Starting HTTP Server")

	log.Fatal(app.ListenTLS(":443", "ssl/server.crt", "ssl/server.key"))
}
