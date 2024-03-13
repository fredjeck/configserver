package server

import (
	"github.com/fredjeck/configserver/internal/config"
	"net/http"
)

func addRoutes(mux *http.ServeMux, c *config.Configuration) {
	mux.HandleFunc("GET /api/register", handleClientRegistration(c))
	mux.HandleFunc("POST /api/tokenize", handleFileTokenization(c))
	requireAuth := authenticatedOnly(c)
	mux.Handle("GET /git/{repository}/{path...}", requireAuth(http.HandlerFunc(handleGitRepositoryAccess())))
}
