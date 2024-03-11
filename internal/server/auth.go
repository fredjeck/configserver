package server

import (
	"context"
	b64 "encoding/base64"
	"github.com/fredjeck/configserver/internal/config"
	"net/http"
	"strings"
)

const msgInvalidAuthHeader = "Invalid authorization header"

// AuthenticatedOnly is a middleware which ensures the requests contains a valid Basic authentication.
// If the authentication succeeds the request context is augmented with the clientId key.
func authenticatedOnly(c *config.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")

			authComponents := strings.Split(authorization, " ")
			if len(authComponents) != 2 {
				HttpUnauthorized(w, msgInvalidAuthHeader)
				return
			}

			if strings.ToLower(authComponents[0]) != "basic" {
				HttpUnauthorized(w, "Unsupported authorization scheme '%s'", authComponents[0])
				return
			}

			basicAuth, err := b64.StdEncoding.DecodeString(authComponents[1])
			if err != nil {
				HttpUnauthorized(w, msgInvalidAuthHeader)
				return
			}

			loginPwd := strings.Split(string(basicAuth), ":")
			if len(loginPwd) != 2 {
				HttpUnauthorized(w, msgInvalidAuthHeader)
				return
			}

			if !validateClientSecret(loginPwd[0], loginPwd[1], c.Server.PassPhrase, c.Server.ValidateSecretLifeSpan) {
				HttpUnauthorized(w, "Unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), "clientId", loginPwd[0])
			rWithCtx := r.WithContext(ctx)

			next.ServeHTTP(w, rWithCtx)
		}
		return http.HandlerFunc(fn)
	}
}
