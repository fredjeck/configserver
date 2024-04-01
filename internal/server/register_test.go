package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/stretchr/testify/assert"
)

var RefactorTestConfiguration = configuration.DefaultConfiguration

const registerClientID = "SampleClientId"
const registerURL = "/api/register?client_id=" + registerClientID

func TestRegisterClientId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, registerURL, nil)
	w := httptest.NewRecorder()
	f := handleClientRegistration(RefactorTestConfiguration)
	f(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRegisterPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, registerURL, nil)
	w := httptest.NewRecorder()
	f := handleClientRegistration(RefactorTestConfiguration)
	f(w, req)

	res := w.Result()
	data, _ := io.ReadAll(res.Body)

	m := &RegisterClientResponse{}
	_ = json.Unmarshal(data, &m)

	assert.Equal(t, m.ClientID, registerClientID)
	assert.True(t, validateClientSecret(registerClientID, m.ClientSecret, RefactorTestConfiguration.Server.PassPhrase, true))
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

	assert.NotNil(t, m.ClientID, registerClientID)
	assert.Len(t, m.ClientID, 36)
}

func TestRegistrationExpiry(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, registerURL, nil)
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
