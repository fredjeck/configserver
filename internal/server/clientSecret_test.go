package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const clientID = "sample-client"
const passPhrase = "magic passphrase"

func TestGenerateClientSecret(t *testing.T) {
	secret, expires := generateClientSecret(clientID, 360, passPhrase)
	assert.True(t, validateClientSecret(clientID, secret, passPhrase, false))
	assert.NotNil(t, expires)
	assert.True(t, time.Now().Before(expires))
}

func TestValidateClientSecretExpired(t *testing.T) {
	secret, _ := generateClientSecret(clientID, -2, passPhrase)
	assert.False(t, validateClientSecret(clientID, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongKey(t *testing.T) {
	secret, _ := generateClientSecret(clientID, 360, passPhrase)
	assert.False(t, validateClientSecret(clientID, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongClientID(t *testing.T) {
	secret, _ := generateClientSecret(clientID, 360, passPhrase)
	assert.False(t, validateClientSecret("wrong-client", secret, passPhrase, false))
}
