package jwt

import (
	"testing"

	"github.com/fredjeck/configserver/internal/encrypt"
	"github.com/stretchr/testify/assert"
)

func TestJSONWebToken(t *testing.T) {
	token := NewJSONWebToken()

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	secret, _ := encrypt.NewHmacSha256Secret()
	jwt := token.Pack(secret)
	assert.Greater(t, len(jwt), 20)
}

func TestVerifySignature(t *testing.T) {
	token := NewJSONWebToken()

	token.Payload.Subject = "test"
	token.Payload.Issuer = "ConfigServer"

	secret, _ := encrypt.NewHmacSha256Secret()
	jwt := token.Pack(secret)
	err := VerifySignature(jwt, secret)
	assert.NoError(t, err)
}
