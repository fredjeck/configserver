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

var TestConfiguration = &Configuration{
	PassPhrase:       "This is a passphrase used to protect yourself",
	ListenOn:         "127.0.0.1:4200",
	SecretExpiryDays: 60,
}

var ClientId = "SampleClientId"

func TestRegisterClientId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register?client_id="+ClientId, nil)
	w := httptest.NewRecorder()
	handleClientRegistration(TestConfiguration).ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRegisterPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register?client_id="+ClientId, nil)
	w := httptest.NewRecorder()
	handleClientRegistration(TestConfiguration).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	m := map[string]interface{}{}
	_ = json.Unmarshal(data, &m)

	assert.Equal(t, m["client_id"], ClientId)
	assert.True(t, Validate(ClientId, m["client_secret"].(string), TestConfiguration.PassPhrase, true))
}

func TestGenerateClientId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)
	w := httptest.NewRecorder()
	handleClientRegistration(TestConfiguration).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	m := map[string]interface{}{}
	_ = json.Unmarshal(data, &m)

	assert.NotNil(t, m["client_id"], ClientId)
	assert.Len(t, m["client_id"].(string), 36)
}

func TestRegistrationExpiry(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/register?client_id="+ClientId, nil)
	w := httptest.NewRecorder()
	handleClientRegistration(TestConfiguration).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	m := map[string]interface{}{}
	_ = json.Unmarshal(data, &m)

	expires, _ := time.Parse(time.RFC3339, m["expires_at"].(string))
	shouldExpire := time.Now().Add(time.Hour * 24 * time.Duration(TestConfiguration.SecretExpiryDays))

	assert.True(t, time.Now().Before(expires))
	assert.Equal(t, shouldExpire.Truncate(24*time.Hour), expires.Truncate(24*time.Hour))
}
