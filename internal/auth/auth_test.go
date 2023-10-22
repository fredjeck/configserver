package auth

import (
	b64 "encoding/base64"
	"net/http"
	"testing"

	"github.com/fredjeck/configserver/internal/encryption"

	"github.com/stretchr/testify/assert"
)

func TestCreateGenerateSecret(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	spec := NewClientKeySpec("clientid", []string{"repo1", "repo2"})
	secret, err := spec.GenerateSecret(key)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
}

func TestUnmarshalGenerateSecret(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	spec := NewClientKeySpec("clientid", []string{"repo1", "repo2"})
	secret, _ := spec.GenerateSecret(key)

	unmarshalled, err := ClientKeySpecFromSecret(secret, key)
	assert.NoError(t, err)
	assert.NotNil(t, unmarshalled)
	assert.Equal(t, "clientid", unmarshalled.ClientID)
	assert.Contains(t, unmarshalled.Repositories, "repo1")
	assert.Contains(t, unmarshalled.Repositories, "repo2")
}

func TestClientKeySpecFromBasicAuth(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	spec := NewClientKeySpec("clientid", []string{"repo1", "repo2"})
	secret, _ := spec.GenerateSecret(key)
	auth := b64.StdEncoding.EncodeToString([]byte("clientid:" + secret))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Basic "+auth)

	unmarshalled, err := ClientKeySpecFromBasicAuth(*request, key)
	assert.NoError(t, err)
	assert.NotNil(t, unmarshalled)
	assert.Equal(t, "clientid", unmarshalled.ClientID)
	assert.Contains(t, unmarshalled.Repositories, "repo1")
	assert.Contains(t, unmarshalled.Repositories, "repo2")
}

func TestAuthRequired(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	request, _ := http.NewRequest("GET", "/", nil)

	unmarshalled, err := ClientKeySpecFromBasicAuth(*request, key)
	assert.Nil(t, unmarshalled)
	assert.ErrorIs(t, err, ErrorAuthRequired)
}

func TestCorruptedAuth(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Basic randomstringwithoutcolon")

	unmarshalled, err := ClientKeySpecFromBasicAuth(*request, key)
	assert.Nil(t, unmarshalled)
	assert.ErrorIs(t, err, ErrorMissingCredentials)
}

func TestUnauthorized(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	spec := NewClientKeySpec("clientid", []string{"repo1", "repo2"})
	secret, _ := spec.GenerateSecret(key)
	auth := b64.StdEncoding.EncodeToString([]byte("wrongclientid:" + secret))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Basic "+auth)

	unmarshalled, err := ClientKeySpecFromBasicAuth(*request, key)
	assert.Nil(t, unmarshalled)
	assert.ErrorIs(t, err, ErrorUnauthorized)
}
