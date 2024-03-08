package server

import (
	"log/slog"
	"net/http"
	"os"
)

type Configuration struct {
	PassPhrase       string `yaml:"pass_phrase"`
	ListenOn         string `yaml:"listen_on"`
	SecretExpiryDays int    `yaml:"secret_expiry_days"`
}

type ConfigServer struct {
	Configuration *Configuration
}

func NewConfigServer(c *Configuration) *ConfigServer {
	return &ConfigServer{c}
}

func (c *ConfigServer) Start() {
	mux := http.NewServeMux()
	mux.Handle("GET /api/register", handleClientRegistration(c.Configuration))
	err := http.ListenAndServe(c.Configuration.ListenOn, mux)
	if err != nil {
		slog.Error("error starting configserver:", "error", err)
		os.Exit(1)
	}
}
