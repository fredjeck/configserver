package server

import (
	"errors"
	"net/http"

	"github.com/fredjeck/configserver/internal/repository"
)

// handleGitRepositoryAccess matches requests with git repositories and returns the request files
func handleGitRepositoryAccess(mgr *repository.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value("clientId").(string)
		repo := r.PathValue("repository")
		path := r.PathValue("path")

		content, err := mgr.Get(repo, path, clientID)
		if err != nil {
			if errors.Is(err, repository.ErrRepositoryNotFound) {
				HttpNotFound(w, "repository '%s' was not found on this server", repo)
			} else if errors.Is(err, repository.ErrClientNotAllowed) {
				HttpUnauthorized(w, "client '%s' is not allowed to access this repository", clientID)
			} else {
				HttpInternalServerError(w, "%w", err)
			}
			return

		}

		Ok(w, content, "text/plain")
	}
}
