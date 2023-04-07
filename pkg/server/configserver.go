package server

import (
	"embed"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/fredjeck/configserver/pkg/auth"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fredjeck/configserver/pkg/cache"
	"github.com/fredjeck/configserver/pkg/config"
	"github.com/fredjeck/configserver/pkg/encrypt"
	"github.com/fredjeck/configserver/pkg/repo"
	"go.uber.org/zap"
)

//go:embed resources
var content embed.FS

// GitUrlPrefix URL prefix from which git repository accesses are served
const GitUrlPrefix string = "/git"

type ConfigServer struct {
	configuration config.Config
	key           *[32]byte
	repositories  *repo.RepositoryManager
	logger        zap.Logger
	cache         *cache.MemoryCache
}

func New(configuration config.Config, key *[32]byte, logger zap.Logger) *ConfigServer {
	return &ConfigServer{
		configuration: configuration,
		key:           key,
		repositories:  repo.NewManager(configuration, logger),
		cache:         cache.NewMemoryCache(time.Duration(configuration.CacheEvictorIntervalSeconds), logger),
		logger:        logger,
	}
}

func (server *ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {
	value, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ciphered, err := encrypt.Encrypt(value, server.key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	base := b64.StdEncoding.EncodeToString(ciphered[:])
	_, err = w.Write([]byte(base))
	if err != nil {
		return
	}
}

// Start starts the configserver
// - Enables the repository manager to pull changes from configured repositories
// - Start serving hosted repositories request
// - Start serving value encryption requests
func (server *ConfigServer) Start() {

	err := server.repositories.Checkout()
	if err != nil {
		server.logger.Sugar().Fatal("error starting configserver, cannot checkout repositories:", err.Error())
		return
	}

	router := http.NewServeMux()
	middleware := server.createGitMiddleWare()

	serverRoot, err := fs.Sub(content, "resources")
	if err != nil {
		server.logger.Sugar().Fatal("error starting configserver, cannot find static resources:", err.Error())
		return
	}

	router.HandleFunc("/api/encrypt", server.encryptValue)
	router.Handle("/", http.FileServer(http.FS(serverRoot)))

	err = http.ListenAndServe(":8090", middleware(router))
	if err != nil {
		server.logger.Sugar().Fatal("error starting configserver:", err.Error())
		return
	}
}

// Writes the Git Middleware response
func (server *ConfigServer) writeResponse(status int, content []byte, started time.Time, w http.ResponseWriter, r http.Request) {
	w.WriteHeader(status)
	_, _ = w.Write(content)
	server.logDuration(started, r, status)
}

// Generates an http like log
func (server *ConfigServer) logDuration(start time.Time, r http.Request, status int) {
	elapsed := time.Since(start)
	message := fmt.Sprintf("%s %s %d %s", r.Method, r.RequestURI, status, elapsed)
	server.logger.Info(message, zap.String("request.method", r.Method), zap.String("request.uri", r.RequestURI), zap.Int("response.status", status), zap.Duration("request.elapsed", elapsed))
}

func (server *ConfigServer) processGitRepoRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// element should at least contain ["", "git", "repository name", "file name"]
	// the first empty element is caused by the leading slash
	elements := strings.Split(r.RequestURI, "/")
	if len(elements) < 4 {
		message := fmt.Sprintf("Invalid repository path '%s' expected format is '%s/repository name/optional folder/file", r.RequestURI, GitUrlPrefix)
		server.logger.Warn(message, zap.String("request.path", r.RequestURI))
		server.writeResponse(http.StatusBadRequest, []byte(message), start, w, *r)
		return
	}
	repository := elements[2]
	path := strings.Join(elements[3:], string(os.PathSeparator))

	spec, err := auth.FromBasicAuth(*r, server.key)
	if err != nil {
		if errors.Is(err, auth.ErrAuthRequired) {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"ConfigServer\"")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			server.writeResponse(http.StatusUnauthorized, []byte(err.Error()), start, w, *r)
		}
	}

	if spec.Repository != repository {
		server.writeResponse(http.StatusUnauthorized, []byte(err.Error()), start, w, *r)
	}

	content, err := server.cache.Get(path)
	if errors.Is(err, cache.ErrKeyNotInCache) {
		content, err = server.repositories.Get(repository, path)
		if err != nil {
			message := fmt.Sprintf("'%s' file not found", path)
			if errors.Is(err, repo.ErrRepositoryNotFound) {
				message = fmt.Sprintf("'%s' repository does not exist", repository)
			}
			if errors.Is(err, repo.ErrFileNotFound) {
				message = fmt.Sprintf("'%s' file does not exsists", path)
			}
			if errors.Is(err, repo.ErrInvalidPath) {
				message = fmt.Sprintf("'%s' path is not valid or contains unsupported characters", path)
			}

			server.logger.Warn(message, zap.String("request.path", r.RequestURI))
			server.writeResponse(http.StatusNotFound, []byte(message), start, w, *r)
			return
		}
		eviction := time.Now().Add(time.Duration(server.configuration.CacheStorageSeconds) * time.Second)
		server.cache.Set(path, content, eviction)
		server.logger.Sugar().Debugf("'%s' : '%s' retrieved from filesystem (cached until %s)", repository, path, eviction)
	} else {
		server.logger.Sugar().Debugf("'%s' : '%s' retrieved from memory cache", repository, path)
	}

	server.writeResponse(http.StatusOK, []byte(content), start, w, *r)
	return
}

// Creates a middleware which intercepts requests retrieving files from the served GIT repositories
// Expects the URL with the following format : GitUrlPrefix/repository name/optional folder(s)/file name
// Example : /git/repository/folder/file.yaml
func (server *ConfigServer) createGitMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			if r.RequestURI[0:4] == GitUrlPrefix && r.Method == http.MethodGet {
				server.processGitRepoRequest(w, r)
				return
			}
			// call next handler
			next.ServeHTTP(w, r)
			// Todo, find a solution to intercept the wrapped middlewares http statuses without interfering with ResponseWriter other interfaces
			server.logDuration(start, *r, -1)
		}
		return http.HandlerFunc(fn)
	}
}
