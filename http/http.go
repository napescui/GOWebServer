package http

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/oschwald/geoip2-golang"
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

	config := config.GetConfig()

	var db *geoip2.Reader
	var err error
	if db, err = geoip2.Open("GeoLite2-City.mmdb"); err != nil {
		logger.Error(err)
	}

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(limiter.New(limiter.Config{
		Max:        config.RateLimit,
		Expiration: time.Duration(config.RateLimitDuration) * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			if config.Logger {
				logger.Infof("IP %s is rate limited", c.IP())
			}
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests, please try again later.")
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		if db == nil {
			return c.Next()
		}

		ip := net.ParseIP(c.IP())
		record, err := db.City(ip)
		if err != nil {
			logger.Error(err)
		}

		if config.Logger {
			if record != nil {
				logger.Infof("IP: %s, Country: %s, City: %s", c.IP(), record.Country.Names["en"], record.City.Names["en"])
			} else {
				logger.Infof("IP: %s", c.IP())
			}
		}

		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		if config.Logger {
			logger.Infof("[%s] %s %s => %d", c.IP(), c.Method(), c.Path(), c.Response().StatusCode())
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	meta := fmt.Sprintf("K10WA_%d", rand.Intn(9000)+1000)
	content := fmt.Sprintf(
		"server|%s\n"+
			"port|%s\n"+
			"type|1\n"+
			"# maint|Server is currently down for maintenance. We will be back soon!\n"+
			"loginurl|%s\n"+
			"meta|%s\n"+
			"RTENDMARKERBS1001",
		config.Host, config.Port, config.LoginUrl, meta)

	app.Post("/growtopia/server_data.php", func(c *fiber.Ctx) error {
		if c.Get("User-Agent") == "" || !strings.Contains(c.Get("User-Agent"), "UbiServices_SDK") {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.SendString(content)
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})

	return app
}

func Start(app *fiber.App) {
	logger.Info("Starting HTTP Server")

	log.Fatal(app.ListenTLS(":443", "ssl/server.crt", "ssl/server.key"))
}
