package main

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	log "github.com/sirupsen/logrus"
)

func SetupGatewayRoutes(app *fiber.App, config Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	for _, service := range config.Services {
		targetURL, err := url.Parse(service.URL)
		if err != nil {
			log.WithFields(log.Fields{
				"service": service.Name,
				"url":     service.URL,
				"error":   err.Error(),
			}).Fatal("Failed to parse service URL")
		}

		for _, prefix := range service.Prefixes {
			app.All(prefix+"*", func(c *fiber.Ctx) error {
				log.WithFields(log.Fields{
					"service":      service.Name,
					"method":       c.Method(),
					"path":         c.Path(),
					"forwarded_to": targetURL.String() + c.Path(),
				}).Info("Forwarding request")

				c.Request().Header.Set("X-Gateway", "true")

				url := targetURL.String() + c.Path()
				if err := proxy.Do(c, url); err != nil {
					log.WithFields(log.Fields{
						"service": service.Name,
						"method":  c.Method(),
						"path":    c.Path(),
						"error":   err.Error(),
					}).Error("Proxy error")

					return c.Status(fiber.StatusBadGateway).JSON(map[string]string{
						"error":   "Bad Gateway",
						"message": "The service is currently unavailable",
					})
				}

				return nil
			})

			log.WithFields(log.Fields{
				"service": service.Name,
				"prefix":  prefix,
				"target":  service.URL,
			}).Info("Registered route")
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		log.WithFields(log.Fields{
			"method": c.Method(),
			"path":   c.Path(),
		}).Warn("No route matched")

		return c.Status(fiber.StatusNotFound).JSON(map[string]string{
			"error":   "Not Found",
			"message": "The requested resource does not exist",
		})
	})
}
