package middleware

import (
	"context"
	"github.com/fredjeck/configserver/internal/auth"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/fredjeck/configserver/internal/server"
	"net/http"
	"strings"
)

func AuthMiddleware(supportedAuth []auth.AuthorizationKind, vault *encryption.KeyVault, paths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			requiresAuth := false
			for _, p := range paths {
				if strings.Contains(path, p) {
					requiresAuth = true
					break
				}
			}

			if !requiresAuth {
				next.ServeHTTP(w, r)
			}

			authorization, err := auth.FromRequest(r, vault, supportedAuth...)
			if err != nil {
				server.HttpNotAuthorized(w, "You are not authorized to access '%s' on this server", r.URL.Path)
				return
			}

			ctx := context.WithValue(r.Context(), "authorization", authorization)
			newReq := r.WithContext(ctx)

			next.ServeHTTP(w, newReq)
		}
		return http.HandlerFunc(fn)
	}
}
