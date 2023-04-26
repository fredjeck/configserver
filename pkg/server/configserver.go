package server

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"path"
	"time"

	"github.com/fredjeck/configserver/pkg/cache"
	"github.com/fredjeck/configserver/pkg/config"
	"github.com/fredjeck/configserver/pkg/repo"
	"go.uber.org/zap"
)

// GitUrlPrefix URL prefix from which git repository accesses are served
const GitUrlPrefix string = "/git"

type ConfigServer struct {
	configuration *config.Config
	key           *[32]byte
	repositories  *repo.RepositoryManager
	logger        *zap.Logger
	cache         *cache.MemoryCache
}

type ConfigServerError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func New(configuration *config.Config, key *[32]byte, logger *zap.Logger) *ConfigServer {
	return &ConfigServer{
		configuration: configuration,
		key:           key,
		repositories:  repo.NewManager(configuration, logger),
		cache:         cache.NewMemoryCache(time.Duration(configuration.CacheEvictorIntervalSeconds), logger),
		logger:        logger,
	}
}

// Start starts the server
// - Enables the repository manager to pull changes from configured repositories
// - Start serving hosted repositories request
// - Start serving api requests
func (server *ConfigServer) Start() {

	err := server.repositories.Checkout()
	if err != nil {
		server.logger.Sugar().Fatal("error starting configserver, cannot checkout repositories:", err.Error())
		return
	}

	router := http.NewServeMux()
	middleware := server.createGitMiddleWare()
	loggingMiddleware := RequestLoggingMiddleware(server.logger)

	ui := http.FileServer(http.Dir(path.Join(server.configuration.Home, "static")))

	router.HandleFunc("/api/encrypt", server.encryptValue)
	router.HandleFunc("/api/stats", server.statistics)
	router.HandleFunc("/api/repositories", server.listRepositories)
	router.HandleFunc("/api/register", server.registerClient)
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/", ui)

	server.logger.Sugar().Info("Now listening on %s", server.configuration.ListenOn)
	err = http.ListenAndServe(server.configuration.ListenOn, loggingMiddleware(middleware(router)))
	if err != nil {
		server.logger.Sugar().Fatal("error starting configserver:", err.Error())
		return
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// Writes the Git Middleware response
func (server *ConfigServer) writeResponse(status int, content []byte, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, _ = w.Write(content)
}

func (server *ConfigServer) writeError(status int, w http.ResponseWriter, message string) {
	w.WriteHeader(status)
	serverError := &ConfigServerError{
		Status:  status,
		Message: message,
	}
	j, err := json.Marshal(serverError)
	if err != nil {
		server.logger.Sugar().Error(err)
	} else {
		_, _ = w.Write(j)
	}
}

func (server *ConfigServer) writeErrorF(status int, w http.ResponseWriter, message string, params ...interface{}) {
	server.writeError(status, w, fmt.Sprintf(message, params...))
}
