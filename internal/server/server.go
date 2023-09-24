// Package server needs a better description
package server

import (
	"log/slog"
	"net/http"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ConfigServer represents the server itself
type ConfigServer struct {
	configuration *config.Configuration
}

// New creates a new instance of ConfigServer using the provided configuration
func New(configuration *config.Configuration) *ConfigServer {
	return &ConfigServer{configuration: configuration}
}

// Start starts the server
// - Enables the repository manager to pull changes from configured repositories
// - Start serving hosted repositories request
// - Start serving api requests
func (server *ConfigServer) Start() {

	router := http.NewServeMux()
	loggingMiddleware := middleware.RequestLoggingMiddleware()

	// router.HandleFunc("/api/encrypt", server.encryptValue)
	// router.HandleFunc("/api/stats", server.statistics)
	// router.HandleFunc("/api/repositories", server.listRepositories)
	// router.HandleFunc("/api/register", server.registerClient)
	router.Handle("/metrics", promhttp.Handler())

	slog.Info("Now listening on", slog.String("port", ":8080"))
	err := http.ListenAndServe(":8080", loggingMiddleware(router))
	if err != nil {
		slog.Error("error starting configserver:", err.Error())
		panic("")
	}
}
