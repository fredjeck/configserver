package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateClientSecret(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, expires := Generate(clientId, 360, passPhrase)
	assert.True(t, Validate(clientId, secret, passPhrase, false))
	assert.NotNil(t, expires)
	assert.True(t, time.Now().Before(expires))
}

func TestValidateClientSecretExpired(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, _ := Generate(clientId, -2, passPhrase)
	assert.False(t, Validate(clientId, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongKey(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, _ := Generate(clientId, 360, passPhrase)
	assert.False(t, Validate(clientId, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongClientId(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret, _ := Generate(clientId, 360, passPhrase)
	assert.False(t, Validate("wrong-client", secret, passPhrase, false))
}
