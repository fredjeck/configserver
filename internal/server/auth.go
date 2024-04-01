package server

import (
	"context"
	b64 "encoding/base64"
	"net/http"
	"strings"

	"github.com/fredjeck/configserver/internal/configuration"
)

const msgInvalidAuthHeader = "Invalid authorization header"

// AuthenticatedOnly is a middleware which ensures the requests contains a valid Basic authentication.
// If the authentication succeeds the request context is augmented with the clientId key.
func authenticatedOnly(c *configuration.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")

			authComponents := strings.Split(authorization, " ")
			if len(authComponents) != 2 {
				HTTPUnauthorized(w, r, msgInvalidAuthHeader)
				return
			}

			if strings.ToLower(authComponents[0]) != "basic" {
				HTTPUnauthorized(w, r, "Unsupported authorization scheme '%s'", authComponents[0])
				return
			}

			basicAuth, err := b64.StdEncoding.DecodeString(authComponents[1])
			if err != nil {
				HTTPUnauthorized(w, r, msgInvalidAuthHeader)
				return
			}

			loginPwd := strings.Split(string(basicAuth), ":")
			if len(loginPwd) != 2 {
				HTTPUnauthorized(w, r, msgInvalidAuthHeader)
				return
			}

			if !validateClientSecret(loginPwd[0], loginPwd[1], c.Server.PassPhrase, c.Server.ValidateSecretLifeSpan) {
				HTTPUnauthorized(w, r, "client '%s' is not allowed to access this repository", loginPwd[0])
				return
			}

			ctx := context.WithValue(r.Context(), ctxClientID{}, loginPwd[0])
			rWithCtx := r.WithContext(ctx)

			next.ServeHTTP(w, rWithCtx)
		}
		return http.HandlerFunc(fn)
	}
}
