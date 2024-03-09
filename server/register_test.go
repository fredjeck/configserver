package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var RefactorTestConfiguration = &Configuration{
	PassPhrase:             "This is a passphrase used to protect yourself",
	ListenOn:               "127.0.0.1:4200",
	SecretExpiryDays:       60,
	ValidateSecretLifeSpan: true,
}

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
	assert.True(t, Validate(ClientId, m.ClientSecret, RefactorTestConfiguration.PassPhrase, true))
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

	shouldExpire := time.Now().Add(time.Hour * 24 * time.Duration(RefactorTestConfiguration.SecretExpiryDays))

	assert.True(t, time.Now().Before(m.ExpiresAt))
	assert.Equal(t, shouldExpire.Truncate(24*time.Hour), m.ExpiresAt.Truncate(24*time.Hour))
}
