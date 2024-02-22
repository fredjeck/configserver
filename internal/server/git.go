package server

import (
	"github.com/fredjeck/configserver/internal/auth"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/fredjeck/configserver/internal/repository"
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

// GitMiddleware validates the provided bearer token signature is valid
func GitMiddleware(vault *encryption.KeyVault, repositories *repository.Manager) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path
			if path[0:4] != "/git" {
				next.ServeHTTP(w, r)
				return
			}

			context := r.Context()
			authorization := context.Value("authorization").(auth.Authorization)

			if r.Method != http.MethodGet {
				HttpBadRequest(w, "'%s' unsupported http verb", r.Method)
				return
			}

			path = path[5:]
			idx := strings.Index(path, "/")
			if idx == -1 {
				HttpBadRequest(w, "'%s' malformed url", path)
				return
			}

			repo := path[0:idx]
			filePath := path[idx+1:]

			authorized := authorization.IsAllowedRepository(repositories, repo)

			if !authorized {
				HttpNotAuthorized(w, "Access to repository '%s' is not allowed", repo)
				return
			}

			content, err := repositories.Get(repo, filePath)
			if err != nil {
				slog.Error("file or repository not found", "repository", repo, "file_path", filePath)
				HttpNotFound(w, "'%s' was not found on this server", filePath)
				return
			}

			clearText, err := encryption.SubstituteTokens(content, vault)
			if err != nil {
				HttpInternalServerError(w, "'%s' : unable to decrypt file", filePath)
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
