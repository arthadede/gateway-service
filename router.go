package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func setupGatewayRoutes(router *mux.Router, config Config) {
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	for _, service := range config.Services {
		targetURL, err := url.Parse(service.URL)
		if err != nil {
			log.WithFields(log.Fields{
				"service": service.Name,
				"url":     service.URL,
				"error":   err.Error(),
			}).Fatal("Failed to parse service URL")
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Header.Set("X-Gateway", "true")
			req.URL.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.Host = targetURL.Host

			log.WithFields(log.Fields{
				"service":      service.Name,
				"method":       req.Method,
				"path":         req.URL.Path,
				"forwarded_to": targetURL.String(),
			}).Info("Forwarding request")
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.WithFields(log.Fields{
				"service": service.Name,
				"method":  r.Method,
				"path":    r.URL.Path,
				"error":   err.Error(),
			}).Error("Proxy error")

			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Bad Gateway",
				"message": "The service is currently unavailable",
			})
		}

		for _, prefix := range service.Prefixes {
			router.PathPrefix(prefix).Handler(proxy)
			log.WithFields(log.Fields{
				"service": service.Name,
				"prefix":  prefix,
				"target":  service.URL,
			}).Info("Registered route")
		}
	}

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Warn("No route matched")

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Not Found",
			"message": "The requested resource does not exist",
		})
	})
}
