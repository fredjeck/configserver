package server

import (
	"encoding/json"
	"github.com/fredjeck/configserver/internal/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var RefactorTestConfiguration = config.DefaultConfiguration

var ClientId = "SampleClientId"

func TestRegisterClientId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register?client_id="+ClientId, nil)
	w := httptest.NewRecorder()
	f := handleClientRegistration(RefactorTestConfiguration)
	f(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRegisterPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register?client_id="+ClientId, nil)
	w := httptest.NewRecorder()
	f := handleClientRegistration(RefactorTestConfiguration)
	f(w, req)

	res := w.Result()
	data, _ := io.ReadAll(res.Body)

	m := &RegisterClientResponse{}
	_ = json.Unmarshal(data, &m)

	assert.Equal(t, m.ClientId, ClientId)
	assert.True(t, validateClientSecret(ClientId, m.ClientSecret, RefactorTestConfiguration.Server.PassPhrase, true))
}

func TestGenerateClientId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)
	w := httptest.NewRecorder()
	f := handleClientRegistration(RefactorTestConfiguration)
	f(w, req)

	res := w.Result()
	data, _ := io.ReadAll(res.Body)

	m := &RegisterClientResponse{}
	_ = json.Unmarshal(data, &m)

	assert.NotNil(t, m.ClientId, ClientId)
	assert.Len(t, m.ClientId, 36)
}

func TestRegistrationExpiry(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register?client_id="+ClientId, nil)
	w := httptest.NewRecorder()
	f := handleClientRegistration(RefactorTestConfiguration)
	f(w, req)

	res := w.Result()
	data, _ := io.ReadAll(res.Body)

	m := &RegisterClientResponse{}
	_ = json.Unmarshal(data, &m)

	shouldExpire := time.Now().Add(time.Hour * 24 * time.Duration(RefactorTestConfiguration.Server.SecretExpiryDays))

	assert.True(t, time.Now().Before(m.ExpiresAt))
	assert.Equal(t, shouldExpire.Truncate(24*time.Hour), m.ExpiresAt.Truncate(24*time.Hour))
}
