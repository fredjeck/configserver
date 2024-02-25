package scratchpad

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateClientSecret(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret := Generate(clientId, 360, passPhrase)
	assert.True(t, Validate(clientId, secret, passPhrase, false))
}

func TestValidateClientSecretWithWrongKey(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret := Generate(clientId, 360, passPhrase)
	assert.False(t, Validate(clientId, secret, "Incorrect key", false))
}

func TestValidateClientSecretWithWrongClientId(t *testing.T) {
	clientId := "sample-client"
	passPhrase := "magic passphrase"

	secret := Generate(clientId, 360, passPhrase)
	assert.False(t, Validate("wrong-client", secret, passPhrase, false))
}
