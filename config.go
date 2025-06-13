package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// ServiceConfig holds the configuration for a microservice
type ServiceConfig struct {
	Name     string
	URL      string
	Prefixes []string
}

// Config holds the gateway configuration
type Config struct {
	Port     string
	Services []ServiceConfig
}

func init() {
	// Configure logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// loadConfig loads the gateway configuration
// In a real-world scenario, this would likely load from a file or environment variables
func loadConfig() Config {
	return Config{
		Port: "8080",
		Services: []ServiceConfig{
			{
				Name: "user-service",
				URL:  "http://localhost:8081",
				Prefixes: []string{
					"/api/users",
					"/api/auth",
				},
			},
			{
				Name: "product-service",
				URL:  "http://localhost:8082",
				Prefixes: []string{
					"/api/products",
					"/api/categories",
				},
			},
			{
				Name: "order-service",
				URL:  "http://localhost:8083",
				Prefixes: []string{
					"/api/orders",
					"/api/payments",
				},
			},
		},
	}
}
