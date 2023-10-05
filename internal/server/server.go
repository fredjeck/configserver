// Package server contains all the HTTP function of the ConfigServer app
package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ConfigServer is a simple HTTP server allowing to securely expose the files provided by the underlying git repositories configuration.
type ConfigServer struct {
	configuration *config.Configuration
	key           *[32]byte
}

// New creates a new instance of ConfigServer using the supplioed configuration
func New(configuration *config.Configuration, key *[32]byte) *ConfigServer {
	return &ConfigServer{configuration: configuration, key: key}
}

// Start starts the server
// - Enables the repository manager to pull changes from configured repositories
// - Start serving hosted repositories request
// - Start serving api requests
// - Adds support for request logging and prometheus metrics
func (server *ConfigServer) Start() {

	router := http.NewServeMux()
	loggingMiddleware := middleware.RequestLoggingMiddleware()

	// router.HandleFunc("/api/stats", server.statistics)
	// router.HandleFunc("/api/repositories", server.listRepositories)
	router.HandleFunc("/api/encrypt", server.encryptValue)
	router.HandleFunc("/api/register", server.registerClient)
	router.Handle("/metrics", promhttp.Handler())

	slog.Info(fmt.Sprintf("Now istening on %s", server.configuration.Server.ListenOn))
	err := http.ListenAndServe(server.configuration.Server.ListenOn, loggingMiddleware(router))
	if err != nil {
		slog.Error("error starting configserver:", err.Error())
		os.Exit(1)
	}
}

// APIError represents an erroneous outcome from an API endpoint
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (server *ConfigServer) writeErrorF(status int, w http.ResponseWriter, message string, params ...interface{}) {
	server.writeError(status, w, fmt.Sprintf(message, params...))
}

func (server *ConfigServer) writeError(status int, w http.ResponseWriter, message string) {
	w.WriteHeader(status)
	serverError := &APIError{
		Status:  status,
		Message: message,
	}
	j, err := json.Marshal(serverError)
	if err != nil {
		slog.Error("cannot convert error to Json: %s", err.Error)
	} else {
		_, _ = w.Write(j)
	}
}

// Writes the Git Middleware response
func (server *ConfigServer) writeResponse(status int, content []byte, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, _ = w.Write(content)
}
