package middleware

import (
	"github.com/fredjeck/configserver/internal/repository"
	"net/http"
	"strings"
)

type GitMiddleware struct {
}

// GitRepoMiddleware validates the provided bearer token signature is valid
func GitRepoMiddleware(repo *repository.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path
			if path[0:4] != "/git" {
				next.ServeHTTP(w, r)
				return
			}

			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			path = path[5:]
			idx := strings.Index(path, "/")
			if idx == -1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			content, err := repo.Get(path[0:idx], path[idx+1:])
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Write(content)

			return
		}
		return http.HandlerFunc(fn)
	}
}
