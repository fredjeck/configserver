// Package server contains all the server related functions of the ConfigServer app
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/repository"
)

// ConfigServer is a standalone server which aims to securely serve git repositories via http
type ConfigServer struct {
	Configuration *config.Configuration
}

// NewConfigServer initializes a new ConfigServer instance
func NewConfigServer(c *config.Configuration) *ConfigServer {
	return &ConfigServer{c}
}

// Start is where all the magic happens
func (c *ConfigServer) Start() {

	manager, mgrErr := repository.NewManager(c.Configuration.Repositories)
	if mgrErr != nil {
		slog.Error("error starting the repository manager, aborting:", "error", mgrErr)
		os.Exit(1)
	}
	manager.Start()

	mux := http.NewServeMux()
	addRoutes(mux, c.Configuration, manager)
	logger := requestLogger()
	slog.Info(fmt.Sprintf("ConfigServer started and listening on %s", c.Configuration.ListenOn))
	err := http.ListenAndServe(c.Configuration.ListenOn, logger(mux))
	if err != nil {
		slog.Error("error starting configserver:", "error", err)
		os.Exit(1)
	}
}
