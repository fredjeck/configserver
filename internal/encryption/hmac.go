package encryption

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// HmacSha256Secret is an alias for Hmac256 secret key
type HmacSha256Secret *[64]byte

// NewHmacSha256Secret generates a random 512-bit secret
func NewHmacSha256Secret() (HmacSha256Secret, error) {
	key := [64]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// HmacSha256Hash generates an HmacSha256 Hash from the provided data and secret
func HmacSha256Hash(data []byte, secret HmacSha256Secret) []byte {
	hmac := hmac.New(sha256.New, secret[:])
	hmac.Write(data)
	return hmac.Sum(nil)
}