package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Configuration struct {
	PassPhrase             string `yaml:"pass_phrase"`
	ListenOn               string `yaml:"listen_on"`
	SecretExpiryDays       int    `yaml:"secret_expiry_days"`
	ValidateSecretLifeSpan bool   `yaml:"validate_secret_lifespan"`
}

type ConfigServer struct {
	Configuration *Configuration
}

func NewConfigServer(c *Configuration) *ConfigServer {
	return &ConfigServer{c}
}

func (c *ConfigServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/register", handleClientRegistration(c.Configuration))
	mux.HandleFunc("POST /api/tokenize", handleFileTokenization(c.Configuration))
	logger := requestLogger()
	slog.Info(fmt.Sprintf("ConfigServer started and listening on %s", c.Configuration.ListenOn))
	err := http.ListenAndServe(c.Configuration.ListenOn, logger(mux))
	if err != nil {
		slog.Error("error starting configserver:", "error", err)
		os.Exit(1)
	}
}
