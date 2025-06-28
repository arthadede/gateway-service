package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type ServiceConfig struct {
	Name     string
	URL      string
	Prefixes []string
}

type Config struct {
	Port     string
	Services []ServiceConfig
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file, using default environment variables")
	}

	// Configure logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func LoadConfig() Config {
	//authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	//if authServiceURL == "" {
	//	authServiceURL = "http://auth-service"
	//	log.Warn("AUTH_SERVICE_URL not set, using default: ", authServiceURL)
	//}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Warn("PORT not set, using default: ", port)
	}

	return Config{
		Port:     port,
		Services: []ServiceConfig{
			//{
			//	Name: "auth-service",
			//	URL:  authServiceURL,
			//	Prefixes: []string{
			//		"/v1/api/auth",
			//	},
			//},
		},
	}
}
