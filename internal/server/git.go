package server

import (
	"log/slog"
	"net/http"
)

// handleGitRepositoryAccess matches requests with git repositories and returns the request files
func handleGitRepositoryAccess() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.Context().Value("clientId").(string))
		slog.Info(r.PathValue("repository"))
		slog.Info(r.PathValue("path"))
	}
}
