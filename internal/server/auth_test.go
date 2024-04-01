package server

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/stretchr/testify/assert"
)

var AuthTestConfiguration = config.DefaultConfiguration

func TestInvalidAuthorizationScheme(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {}
	mdw := authenticatedOnly(AuthTestConfiguration)

	req := httptest.NewRequest(http.MethodGet, "https://www.your-domain.com/git", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer:%s", "token"))

	w := httptest.NewRecorder()
	mdw(http.HandlerFunc(next)).ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestMalformedBasicAuth(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {}
	mdw := authenticatedOnly(AuthTestConfiguration)

	req := httptest.NewRequest(http.MethodGet, "https://www.your-domain.com/git", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", "token"))

	w := httptest.NewRecorder()
	mdw(http.HandlerFunc(next)).ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestInvalidClientSecret(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {}
	mdw := authenticatedOnly(AuthTestConfiguration)

	token := b64.StdEncoding.EncodeToString([]byte("a:b:c"))

	req := httptest.NewRequest(http.MethodGet, "https://www.your-domain.com/git", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", token))

	w := httptest.NewRecorder()
	mdw(http.HandlerFunc(next)).ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestValidClientSecret(t *testing.T) {
	id := "AClientId"
	secret, _ := generateClientSecret(id, 360, AuthTestConfiguration.Server.PassPhrase)
	token := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", id, secret)))

	next := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, id, r.Context().Value(ctxClientID{}))
	}
	mdw := authenticatedOnly(AuthTestConfiguration)

	req := httptest.NewRequest(http.MethodGet, "https://www.your-domain.com/git", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", token))

	w := httptest.NewRecorder()
	mdw(http.HandlerFunc(next)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
