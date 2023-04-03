package server

import (
	"embed"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
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
		cache:         cache.NewMemoryCache(time.Duration(configuration.CacheEvicterIntervalSeconds), logger),
		logger:        logger,
	}
}

func (server ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {
	value, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ciphered, error := encrypt.Encrypt(value, server.key)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	base := b64.StdEncoding.EncodeToString(ciphered[:])
	w.Write([]byte(base))

}

func (server ConfigServer) Start() {

	server.repositories.Checkout()

	router := http.NewServeMux()
	middleware := server.gitRepoMiddleWare()
	routerWithMiddleware := middleware(router)

	serverRoot, err := fs.Sub(content, "resources")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/api/encrypt", server.encryptValue)
	router.Handle("/", http.FileServer(http.FS(serverRoot)))

	err = http.ListenAndServe(":8090", routerWithMiddleware)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
}

// TODO cleanup the mess
func (s ConfigServer) gitRepoMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.RequestURI, "/git") {
				// handle
				elements := strings.Split(r.RequestURI, "/")
				if len(elements) < 4 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				r := elements[2]
				path := strings.Join(elements[3:], string(os.PathSeparator))

				content, err := s.cache.Get(path)
				if errors.Is(err, cache.ErrKeyNotInCache) {
					content, err = s.repositories.Get(r, path)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					eviction := time.Now().Add(time.Duration(s.configuration.CacheEvicterIntervalSeconds) * time.Second)
					s.cache.Set(path, content, eviction)
					s.logger.Sugar().Debugf("'%s' : '%s' retrieved from filesystem (cached until %s)", r, path, eviction)
				} else {
					s.logger.Sugar().Debugf("'%s' : '%s' retrieved from memory cache", r, path)
				}

				w.WriteHeader(http.StatusOK)
				w.Write(content)
				return
			}

			// call next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
