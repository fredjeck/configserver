package auth

import (
	"testing"

	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/stretchr/testify/assert"
)

func TestJSONWebToken(t *testing.T) {
	token := NewJSONWebToken()

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	secret, _ := encryption.NewHmacSha256Secret()
	jwt := token.Pack(secret)
	assert.Greater(t, len(jwt), 20)
}

func TestVerifySignature(t *testing.T) {
	token := NewJSONWebToken()

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	secret, _ := encryption.NewHmacSha256Secret()
	jwt := token.Pack(secret)
	err := VerifySignature(jwt, secret)
	assert.NoError(t, err)
}

func TestUnpack(t *testing.T) {
	token := NewJSONWebToken()

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	secret, _ := encryption.NewHmacSha256Secret()
	jwt := token.Pack(secret)
	tk, err := Unpack(jwt, secret)
	assert.NoError(t, err)

	assert.Equal(t, "test", tk.Payload.Subject)
	assert.Equal(t, "ConfigServer", tk.Payload.Issuer)
}
