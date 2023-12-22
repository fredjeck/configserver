package auth

import (
	"testing"

	"github.com/fredjeck/configserver/internal/encryption"

	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	secret, err := GenerateClientSecret("clientid", key)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
}

func TestValidateSecret(t *testing.T) {
	key, _ := encryption.NewAes256Key()
	secret, err := GenerateClientSecret("clientid", key)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)

	assert.True(t, ValidateClientSecret("clientid", secret, key))
}
