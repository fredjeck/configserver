// Package server contains all the HTTP function of the ConfigServer app
package server

import (
	"encoding/json"
	"fmt"
	"github.com/fredjeck/configserver/internal/auth"
	"github.com/fredjeck/configserver/internal/repository"
	"github.com/fredjeck/configserver/internal/server/middleware"
	"log/slog"
	"net/http"
	"os"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ConfigServer is a simple HTTP server allowing to securely expose the files provided by the underlying git repositories configuration.
type ConfigServer struct {
	configuration *config.Configuration
	repository    *repository.Manager
	authorization []auth.AuthorizationKind
	vault         *encryption.KeyVault
}

// New creates a new instance of ConfigServer using the supplioed configuration
func New(configuration *config.Configuration, repository *repository.Manager, vault *encryption.KeyVault) *ConfigServer {
	srv := &ConfigServer{configuration: configuration, repository: repository, authorization: []auth.AuthorizationKind{}, vault: vault}
	for _, akind := range configuration.Authorization {
		srv.authorization = append(srv.authorization, auth.AuthorizationKind(akind))
	}
	return srv
}

// Start starts the server
// - Enables the repository manager to pull changes from configured repositories
// - Start serving hosted repositories request
// - Start serving api requests
// - Adds support for request logging and prometheus metrics
func (server *ConfigServer) Start() {

	router := http.NewServeMux()
	loggingMiddleware := middleware.RequestLoggingMiddleware()
	gitMiddleware := middleware.GitMiddleware(server.vault, server.repository)
	authmdw := middleware.AuthMiddleware(server.authorization, server.vault, "/git")

	// router.HandleFunc("/api/stats", server.statistics)

	// Encrypts a value to a substitution token
	router.HandleFunc("/api/encrypt", server.encryptValue)

	// Tokenizes the provided data
	router.HandleFunc("/api/tokenize", server.tokenizeText)

	// Obtain a new ClientID
	router.HandleFunc("/api/register", server.generateClientSecret)

	// OAuth2 Authorization endpoint
	router.HandleFunc("/oauth2/authorize", server.authorize)

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	slog.Info(fmt.Sprintf("Now istening on %s", server.configuration.Server.ListenOn))
	err := http.ListenAndServe(server.configuration.Server.ListenOn, loggingMiddleware(authmdw(gitMiddleware(router))))
	if err != nil {
		slog.Error("error starting configserver:", "error", err)
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
		slog.Error("cannot convert error to Json", "error", err.Error)
	} else {
		_, _ = w.Write(j)
	}
}

// Writes the Git Middleware response
func (server *ConfigServer) writeResponse(status int, content []byte, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, _ = w.Write(content)
}
