package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateClientSecret(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, expires := generateClientSecret(clientId, 360, passPhrase)
	assert.True(t, validateClientSecret(clientId, secret, passPhrase, false))
	assert.NotNil(t, expires)
	assert.True(t, time.Now().Before(expires))
}

func TestValidateClientSecretExpired(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, _ := generateClientSecret(clientId, -2, passPhrase)
	assert.False(t, validateClientSecret(clientId, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongKey(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, _ := generateClientSecret(clientId, 360, passPhrase)
	assert.False(t, validateClientSecret(clientId, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongClientId(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, _ := generateClientSecret(clientId, 360, passPhrase)
	assert.False(t, validateClientSecret("wrong-client", secret, passPhrase, false))
}
