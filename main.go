package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := LoadConfig()

	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	SetupGatewayRoutes(app, config)

	SetupMiddleware(app)

	go func() {
		log.WithFields(log.Fields{
			"port": config.Port,
		}).Info("Starting gateway server")

		if err := app.Listen(":" + config.Port); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Could not start server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err := app.ShutdownWithTimeout(15 * time.Second); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error during server shutdown")
	}

	log.Info("Server gracefully stopped")
}
