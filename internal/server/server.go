package server

import (
	"fmt"
	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/repository"
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

	manager, mgr_err := repository.NewManager(c.Configuration.Repositories)
	if mgr_err != nil {
		slog.Error("error starting the repository manager, aborting:", "error", mgr_err)
		os.Exit(1)
	}
	manager.Start()

	mux := http.NewServeMux()
	addRoutes(mux, c.Configuration)
	logger := requestLogger()
	slog.Info(fmt.Sprintf("ConfigServer started and listening on %s", c.Configuration.ListenOn))
	err := http.ListenAndServe(c.Configuration.ListenOn, logger(mux))
	if err != nil {
		slog.Error("error starting configserver:", "error", err)
		os.Exit(1)
	}
}
