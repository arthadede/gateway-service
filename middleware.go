package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

func SetupMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		log.WithFields(log.Fields{
			"method":      c.Method(),
			"path":        c.Path(),
			"remote_addr": c.IP(),
			"user_agent":  c.Get("User-Agent"),
			"duration_ms": time.Since(start).Milliseconds(),
		}).Info("Request processed")

		return err
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))
}
