package http

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
	"os"
	"encoding/base64"

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
		if !config.EnableGeo {
			return c.Next()
		}

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

		if len(config.GeoLocation) > 0 && record != nil {
			allowed := false
			for _, loc := range config.GeoLocation {
				if record.Country.IsoCode == loc {
					allowed = true
					break
				}
			}
			if !allowed {
				return c.Status(fiber.StatusForbidden).SendString("IP is not in the allowed GeoLocation")
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

	app.Static("/cache", "./cache")

	app.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/cache") {
			logger.Info("Connection from: " + c.IP() + " | Downloading: " + c.Path())

			pathname := filepath.Join("./cache", c.Path())

			if _, err := os.Stat(pathname); os.IsNotExist(err) {
				return c.Redirect(
					fmt.Sprintf("https://ubistatic-a.akamaihd.net/%s%s", config.serverCdn, c.Path()),
					fiber.StatusMovedPermanently,
				)
			}

			file, err := os.Open(pathname)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString("error from loading")
			}
			defer file.Close()

			buffer, err := io.ReadAll(file)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString("error")
			}

			contentTypes := map[string]string{
				".ico":  "image/x-icon",
				".html": "text/html",
				".js":   "text/javascript",
				".json": "application/json",
				".css":  "text/css",
				".png":  "image/png",
				".jpg":  "image/jpeg",
				".wav":  "audio/wav",
				".mp3":  "audio/mpeg",
				".svg":  "image/svg+xml",
				".pdf":  "application/pdf",
				".doc":  "application/msword",
			}

			ext := filepath.Ext(c.Path())
			c.Set("Content-Type", contentTypes[ext])

			return c.Send(buffer)
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	meta := fmt.Sprintf("K10WA_%d", rand.Intn(9000)+1000)
	loginUrl := config.LoginUrl
	if loginUrl == "default" {
		loginUrl = config.Host
	}
	content := fmt.Sprintf(
		"server|%s\n"+
			"port|%s\n"+
			"type|1\n"+
			"# maint|Server is currently down for maintenance. We will be back soon!\n"+
			"loginurl|%s\n"+
			"meta|%s\n"+
			"RTENDMARKERBS1001",
		config.Host, config.Port, loginUrl, meta)

	app.Post("/growtopia/server_data.php", func(c *fiber.Ctx) error {
		if c.Get("User-Agent") == "" || !strings.Contains(c.Get("User-Agent"), "UbiServices_SDK") {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.SendString(content)
	})

	app.Post("/player/login/dashboard", func(c *fiber.Ctx) error {
		htmlFile, err := os.ReadFile("html/dashboard.html")
		if err != nil {
			logger.Error(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
		}
	
		html := string(htmlFile)
		html = strings.Replace(html, "{serverName}", config.ServerName, -1)
		html = strings.Replace(html, "{serverSupport}", config.ServerSupport, -1)
		html = strings.Replace(html, "{serverHost}", config.Host, -1)
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	app.Post("/player/growid/login/validate", func(c *fiber.Ctx) error {
		token := c.FormValue("_token")
		growid := c.FormValue("growId")
		password := c.FormValue("password")

		encoded := fmt.Sprintf("_token=%s&growId=%s&password=%s", token, growid, password)
		encoded = base64.StdEncoding.EncodeToString([]byte(encoded))
		
		data := fmt.Sprintf(`{"status":"success","message":"Account Validated","token":"%s","url":"","accountType":"growtopia"}`, encoded)

		return c.SendString(data)
	})

	app.Get("/player/validate", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<script>window.close();</script>")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
	})

	return app
}

func Start(app *fiber.App) {
	logger.Info("Starting HTTP Server")

	log.Fatal(app.ListenTLS(":443", "ssl/server.crt", "ssl/server.key"))
}
