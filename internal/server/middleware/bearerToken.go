package middleware

import (
	"net/http"
	"strings"

	"github.com/fredjeck/configserver/internal/encrypt"
	"github.com/fredjeck/configserver/internal/jwt"
)

// BearerTokenMiddleware validates the provided bearer token signature is valid
func BearerTokenMiddleware(secret encrypt.HmacSha256Secret) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth, ok := r.Header["Authorization"]
			if ok && len(auth) == 1 {
				authstr := strings.ToLower(auth[0])
				if strings.Contains(authstr, "bearer") {
					token := strings.Replace(authstr, "bearer ", "", -1)
					err := jwt.VerifySignature(token, secret)
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
