package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	config := loadConfig()

	// Create router
	router := mux.NewRouter()

	// Setup routes
	setupGatewayRoutes(router, config)

	// Apply middleware
	handler := setupMiddleware(router)

	// Configure server
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithFields(log.Fields{
			"port": config.Port,
		}).Info("Starting gateway server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Could not start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	server.Shutdown(ctx)
	log.Info("Server gracefully stopped")
}
