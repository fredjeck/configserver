package server

import (
	"fmt"
	"github.com/fredjeck/configserver/internal/auth"
	"log/slog"
	"net/http"
	"strings"
)

func (server *ConfigServer) extractToken(w http.ResponseWriter, r *http.Request) (*auth.JSONWebToken, bool) {
	authorization, ok := r.Header["Authorization"]
	if !ok {
		die(w, http.StatusBadRequest, "missing authorization header")
		return nil, false
	}

	authStr := strings.ToLower(authorization[0])
	if !strings.Contains(authStr, "bearer") {
		die(w, http.StatusBadRequest, "only bearer authorization is supported")
		return nil, false
	}

	token := strings.Replace(authStr, "bearer ", "", -1)
	err := auth.VerifySignature(token, server.keystore.HmacSha256Secret)
	if err != nil {
		die(w, http.StatusUnauthorized, "not authorized")
		return nil, false
	}

	jwt, err := auth.ParseJwt(token, server.keystore.HmacSha256Secret)
	if err != nil {
		die(w, http.StatusUnauthorized, "invalid token")
		return nil, false
	}

	return jwt, true
}

// GitRepoMiddleware validates the provided bearer token signature is valid
func (server *ConfigServer) GitRepoMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path
			if path[0:4] != "/git" {
				next.ServeHTTP(w, r)
				return
			}

			_, ok := server.extractToken(w, r)
			if !ok {
				return
			}

			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(fmt.Sprintf("'%s' unsupported http verb", r.Method)))
				return
			}

			path = path[5:]
			idx := strings.Index(path, "/")
			if idx == -1 {
				slog.Error("malformed url", "path_url", path)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(fmt.Sprintf("'%s' malformed url", path)))
				return
			}

			repo := path[0:idx]
			filePath := path[idx+1:]
			content, err := server.repository.Get(repo, filePath)
			if err != nil {
				slog.Error("file or repository not found", "repository", repo, "file_path", filePath)
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(fmt.Sprintf("repository '%s' - file '%s' file or repository not found", repo, filePath)))
				return
			}

			w.Header().Add("Content-Type", "text/plain")
			_, err = w.Write(content)

			if err != nil {
				slog.Error("an error occured while sending back repository file", "url_path", path, "error", err)
			}
			return
		}
		return http.HandlerFunc(fn)
	}
}
