package server

import (
	"github.com/fredjeck/configserver/internal/auth"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"net/http"
	"strings"
)

var (
	hitCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "git_hit_count",
		Help: "Total number of files retrieved",
	}, []string{
		// Repository from where the files where retrieved
		"repository",
		// File path
		"path",
	})
)

// GitRepoMiddleware validates the provided bearer token signature is valid
func (server *ConfigServer) GitRepoMiddleware() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path
			if path[0:4] != "/git" {
				next.ServeHTTP(w, r)
				return
			}

			authorization, err := auth.FromRequest(r, server.vault, server.authorization...)
			if err != nil {
				dieErr(w, r, http.StatusUnauthorized, "authorization failed", err)
				return
			}

			if r.Method != http.MethodGet {
				dief(w, http.StatusBadRequest, "'%s' unsupported http verb", r.Method)
				return
			}

			path = path[5:]
			idx := strings.Index(path, "/")
			if idx == -1 {
				dief(w, http.StatusBadRequest, "'%s' malformed url", path)
				return
			}

			repo := path[0:idx]
			filePath := path[idx+1:]

			authorized := authorization.IsAllowedRepository(server.repository, repo)

			if !authorized {
				die(w, http.StatusUnauthorized, "repository access is not allowed")
				return
			}

			content, err := server.repository.Get(repo, filePath)
			if err != nil {
				slog.Error("file or repository not found", "repository", repo, "file_path", filePath)
				dief(w, http.StatusNotFound, "'%s' was not found on this server", filePath)
				return
			}

			clearText, err := encryption.SubstituteTokens(content, server.vault)
			if err != nil {
				dief(w, http.StatusNotFound, "'%s' : unable to decrypt file", filePath)
				return
			}

			w.Header().Add("Content-Type", "text/plain")
			_, err = w.Write(clearText)

			if err != nil {
				slog.Error("an error occured while sending back repository file", "url_path", path, "error", err)
			}

			hitCount.With(prometheus.Labels{"repository": repo, "path": filePath}).Inc()
			return
		}
		return http.HandlerFunc(fn)
	}
}
