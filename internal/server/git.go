package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/fredjeck/configserver/internal/repository"
	"github.com/fredjeck/configserver/internal/utils"
)

// handleGitRepositoryAccess matches requests with git repositories and returns the request files
func handleGitRepositoryAccess(mgr *repository.Manager, c *configuration.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Context().Value(ctxClientID{}).(string)
		requestID := r.Context().Value(ctxRequestID{}).(string)
		repo := r.PathValue("repository")
		path := r.PathValue("path")

		content, err := mgr.Get(repo, path, clientID)
		if err != nil {
			if errors.Is(err, repository.ErrRepositoryNotFound) {
				HTTPNotFound(w, r, "repository '%s' was not found on this server", repo)
			} else if errors.Is(err, repository.ErrClientNotAllowed) {
				HTTPUnauthorized(w, r, "client '%s' is not allowed to access this repository", clientID)
			} else {
				HTTPInternalServerError(w, r, "%w", []interface{}{err}...)
			}
			return
		}

		clear, err := utils.Detokenize(string(content[:]), c.Server.PassPhrase)
		if err != nil {
			slog.Error("An error occured while detokenizing the requested file", "error", err, HTTPRequestID, requestID)
			HTTPInternalServerError(w, r, "An error occured while detokenizing the requested file")
		}

		Ok(w, []byte(clear), "text/plain")
	}
}
