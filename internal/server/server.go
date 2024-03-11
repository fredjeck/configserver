package server

import (
	"fmt"
	"github.com/fredjeck/configserver/internal/config"
	"log/slog"
	"net/http"
	"os"
)

type ConfigServer struct {
	Configuration *config.Configuration
}

func NewConfigServer(c *config.Configuration) *ConfigServer {
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
