// Package server contains all the HTTP function of the ConfigServer app
package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fredjeck/configserver/internal/repository"
	"log/slog"
	"net/http"
	"os"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/fredjeck/configserver/internal/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ConfigServer is a simple HTTP server allowing to securely expose the files provided by the underlying git repositories configuration.
type ConfigServer struct {
	configuration *config.Configuration
	keystore      *encryption.Keystore
	repository    *repository.Manager
}

// New creates a new instance of ConfigServer using the supplioed configuration
func New(configuration *config.Configuration, keystore *encryption.Keystore, repository *repository.Manager) *ConfigServer {
	return &ConfigServer{configuration: configuration, keystore: keystore, repository: repository}
}

// Start starts the server
// - Enables the repository manager to pull changes from configured repositories
// - Start serving hosted repositories request
// - Start serving api requests
// - Adds support for request logging and prometheus metrics
func (server *ConfigServer) Start() {

	secret, _ := encryption.NewHmacSha256Secret()

	slog.Info("Secret", "secret", base64.StdEncoding.EncodeToString(secret.Key))

	router := http.NewServeMux()
	loggingMiddleware := middleware.RequestLoggingMiddleware()
	//bearerTokenMiddleware := middleware.BearerTokenMiddleware(secret)
	gitMiddleware := server.GitRepoMiddleware()

	// router.HandleFunc("/api/stats", server.statistics)
	// router.HandleFunc("/api/repositories", server.listRepositories)
	router.HandleFunc("/api/encrypt", server.encryptValue)
	router.HandleFunc("/api/register", server.registerClient)
	router.HandleFunc("/api/register/jwt", server.registerClientJwt)
	router.HandleFunc("/api/keygen/aes", server.GenAes256)
	router.HandleFunc("/api/keygen/hmac", server.GenHmacSha256)
	router.Handle("/metrics", promhttp.Handler())

	// TODO change bearerToken middleware so that it is explicitely wrapped around the pieces which need auth
	// https://drstearns.github.io/tutorials/gomiddleware/

	slog.Info(fmt.Sprintf("Now istening on %s", server.configuration.Server.ListenOn))
	err := http.ListenAndServe(server.configuration.Server.ListenOn, loggingMiddleware(gitMiddleware(router)))
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
