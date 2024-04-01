package server

import (
	"net/http"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/fredjeck/configserver/internal/repository"
)

func addRoutes(mux *http.ServeMux, c *configuration.Configuration, m *repository.Manager) {
	mux.HandleFunc("GET /api/register", handleClientRegistration(c))
	mux.HandleFunc("POST /api/tokenize", handleFileTokenization(c))
	mux.HandleFunc("GET /stats", handleStatistics(m))
	requireAuth := authenticatedOnly(c)
	mux.Handle("GET /git/{repository}/{path...}", requireAuth(http.HandlerFunc(handleGitRepositoryAccess(m, c))))
}
