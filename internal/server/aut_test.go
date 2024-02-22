package server

import (
	"context"
	"fmt"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {}
	vault, _ := encryption.NewKeyVault()
	mdw := AuthMiddleware(nil, vault, "/git")

	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/git", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer:%s", "token"))
	req = req.WithContext(context.WithValue(req.Context(), "some-key", "123ABC"))

	res := httptest.NewRecorder()

	tim := mdw(http.HandlerFunc(next))
	tim.ServeHTTP(res, req)
	assert.Equal(t, res.Code, http.StatusUnauthorized)
}
