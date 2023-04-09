package auth

import (
	b64 "encoding/base64"
	"github.com/fredjeck/configserver/pkg/encrypt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClientSecret(t *testing.T) {
	key := encrypt.NewEncryptionKey()
	spec := NewClientSpec("clientid", []string{"repo1", "repo2"})
	secret, err := spec.ClientSecret(key)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
}

func TestUnmarshalClientSecret(t *testing.T) {
	key := encrypt.NewEncryptionKey()
	spec := NewClientSpec("clientid", []string{"repo1", "repo2"})
	secret, _ := spec.ClientSecret(key)

	unmarshalled, err := UnmarshalClientSecret(secret, key)
	assert.NoError(t, err)
	assert.NotNil(t, unmarshalled)
	assert.Equal(t, "clientid", unmarshalled.ClientId)
	assert.Contains(t, unmarshalled.Repositories, "repo1")
	assert.Contains(t, unmarshalled.Repositories, "repo2")
}

func TestFromBasicAuth(t *testing.T) {
	key := encrypt.NewEncryptionKey()
	spec := NewClientSpec("clientid", []string{"repo1", "repo2"})
	secret, _ := spec.ClientSecret(key)
	auth := b64.StdEncoding.EncodeToString([]byte("clientid:" + secret))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Basic "+auth)

	unmarshalled, err := FromBasicAuth(*request, key)
	assert.NoError(t, err)
	assert.NotNil(t, unmarshalled)
	assert.Equal(t, "clientid", unmarshalled.ClientId)
	assert.Contains(t, unmarshalled.Repositories, "repo1")
	assert.Contains(t, unmarshalled.Repositories, "repo2")
}

func TestAuthRequired(t *testing.T) {
	key := encrypt.NewEncryptionKey()
	request, _ := http.NewRequest("GET", "/", nil)

	unmarshalled, err := FromBasicAuth(*request, key)
	assert.Nil(t, unmarshalled)
	assert.ErrorIs(t, err, ErrAuthRequired)
}

func TestCorruptedAuth(t *testing.T) {
	key := encrypt.NewEncryptionKey()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Basic randomstringwithoutcolon")

	unmarshalled, err := FromBasicAuth(*request, key)
	assert.Nil(t, unmarshalled)
	assert.ErrorIs(t, err, ErrMissingCredentials)
}

func TestUnauthorized(t *testing.T) {
	key := encrypt.NewEncryptionKey()
	spec := NewClientSpec("clientid", []string{"repo1", "repo2"})
	secret, _ := spec.ClientSecret(key)
	auth := b64.StdEncoding.EncodeToString([]byte("wrongclientid:" + secret))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Basic "+auth)

	unmarshalled, err := FromBasicAuth(*request, key)
	assert.Nil(t, unmarshalled)
	assert.ErrorIs(t, err, ErrUnauthorized)
}
