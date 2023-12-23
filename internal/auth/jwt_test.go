package auth

import (
	"testing"

	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/stretchr/testify/assert"
)

func TestJSONWebToken(t *testing.T) {
	token := NewJSONWebToken()
	vault, err := encryption.NewKeyVault()
	assert.NoError(t, err)

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	jwt := token.Pack(vault)
	assert.Greater(t, len(jwt), 20)
}

func TestVerifySignature(t *testing.T) {
	token := NewJSONWebToken()
	vault, err := encryption.NewKeyVault()
	assert.NoError(t, err)

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	jwt := token.Pack(vault)
	err = VerifySignature(jwt, vault)
	assert.NoError(t, err)
}

func TestUnpack(t *testing.T) {
	token := NewJSONWebToken()
	vault, err := encryption.NewKeyVault()
	assert.NoError(t, err)

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	jwt := token.Pack(vault)
	tk, err := Unpack(jwt, vault)
	assert.NoError(t, err)

	assert.Equal(t, "test", tk.Payload.Subject)
	assert.Equal(t, "ConfigServer", tk.Payload.Issuer)
}
