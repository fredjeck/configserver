package auth

import (
	"testing"

	"github.com/fredjeck/configserver/internal/encryption"

	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	vault, err := encryption.NewKeyVault()
	assert.NoError(t, err)

	secret, err := GenerateClientSecret("clientid", vault)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
}

func TestValidateSecret(t *testing.T) {
	vault, err := encryption.NewKeyVault()
	assert.NoError(t, err)

	secret, err := GenerateClientSecret("clientid", vault)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)

	assert.True(t, ValidateClientSecret("clientid", secret, vault))
}
