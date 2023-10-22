package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHmacSha256HashLength(t *testing.T) {
	secret, _ := NewHmacSha256Secret()
	data := "Wingardium Leviosa"

	hash := HmacSha256Hash([]byte(data), secret)
	assert.Len(t, hash, 32)
}

func TestHmacSha256Hash(t *testing.T) {
	secret, _ := NewHmacSha256Secret()
	data := "Wingardium Leviosa"

	hash1 := HmacSha256Hash([]byte(data), secret)
	hash2 := HmacSha256Hash([]byte(data), secret)
	assert.Equal(t, hash1, hash2)
}
