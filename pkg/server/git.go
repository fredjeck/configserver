package server

import (
	"errors"
	"fmt"
	"github.com/fredjeck/configserver/pkg/auth"
	"github.com/fredjeck/configserver/pkg/cache"
	"github.com/fredjeck/configserver/pkg/repo"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"time"
)

func (server *ConfigServer) processGitRepoRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// element should at least contain ["", "git", "repository name", "file name"]
	// the first empty element is caused by the leading slash
	elements := strings.Split(r.RequestURI, "/")
	if len(elements) < 4 {
		message := fmt.Sprintf("Invalid repository path '%s' expected format is '%s/repository name/optional folder/file", r.RequestURI, GitUrlPrefix)
		zap.L().Warn(message, zap.String("request.path", r.RequestURI))
		server.writeResponse(http.StatusBadRequest, []byte(message), w)
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
			server.writeResponse(http.StatusUnauthorized, []byte(err.Error()), w)
		}
	}

	if !spec.CanAccessRepository(repository) {
		server.writeResponse(http.StatusUnauthorized, []byte(err.Error()), w)
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

			zap.L().Warn(message, zap.String("request.path", r.RequestURI))
			server.writeResponse(http.StatusNotFound, []byte(message), w)
			return
		}
		eviction := time.Now().Add(time.Duration(server.configuration.CacheStorageSeconds) * time.Second)
		server.cache.Set(path, content, eviction)
		zap.L().Sugar().Debugf("'%s' : '%s' retrieved from filesystem (cached until %s)", repository, path, eviction)
	} else {
		zap.L().Sugar().Debugf("'%s' : '%s' retrieved from memory cache", repository, path)
	}

	server.writeResponse(http.StatusOK, content, w)
	return
}

// Creates a middleware which intercepts requests retrieving files from the served GIT repositories
// Expects the URL with the following format : GitUrlPrefix/repository name/optional folder(s)/file name
// Example : /git/repository/folder/file.yaml
func (server *ConfigServer) createGitMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if len(r.RequestURI) >= 4 && r.RequestURI[0:4] == GitUrlPrefix && r.Method == http.MethodGet {
				server.processGitRepoRequest(w, r)
				return
			}
			// call next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
