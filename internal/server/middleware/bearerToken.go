package middleware

import (
	"net/http"
	"strings"

	"github.com/fredjeck/configserver/internal/auth"
	"github.com/fredjeck/configserver/internal/encryption"
)

// BearerTokenMiddleware validates the provided bearer token signature is valid
func BearerTokenMiddleware(secret encryption.HmacSha256Secret) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path[0:4] == "/api" {
				next.ServeHTTP(w, r)
				return
			}

			authorization, ok := r.Header["Authorization"]
			if ok && len(authorization) == 1 {
				authstr := strings.ToLower(authorization[0])
				if strings.Contains(authstr, "bearer") {
					token := strings.Replace(authstr, "bearer ", "", -1)
					err := auth.VerifySignature(token, secret)
					if err != nil {
						http.Error(w, "Not authorized", http.StatusUnauthorized)
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
