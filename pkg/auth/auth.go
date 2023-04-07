package auth

import (
	b64 "encoding/base64"
	"github.com/fredjeck/configserver/pkg/encrypt"
)

// NewClientSecret generates a new client secret.
// The provided key is used to encrypt information embedded into the client secret it is therefore important to use the same key throughout the whole execution context
// Client Secret Specification :
// Values are separated by colons (":')
// [0] - Bound repository name
// [1] - Client ID
func NewClientSecret(clientId string, repository string, key *[32]byte) (string, error) {
	secret, err := encrypt.Encrypt([]byte(repository+":"+clientId), key)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secret), nil
}
