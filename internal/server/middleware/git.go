package middleware

import (
	"net/http"
)

type GitMiddleware struct {
}

// GitRepoMiddleware validates the provided bearer token signature is valid
func GitRepoMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path[0:4] != "/git" {
				next.ServeHTTP(w, r)
				return
			}

			// Do something with the request

			return
		}
		return http.HandlerFunc(fn)
	}
}
