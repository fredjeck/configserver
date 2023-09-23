package encrypt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	key := NewEncryptionKey()
	text := []byte("this is a sample text!")
	s, err := Encrypt(text, key)
	assert.NoError(t, err, "Encryption should not raise an error")

	dec, derr := Decrypt(s, key)
	assert.NoError(t, derr, "Decryption should not raise an error")
	assert.EqualValues(t, text, dec)
}

func TestEncryptionKeyStorage(t *testing.T) {
	keyfile, _ := os.CreateTemp("", "keyfile")
	defer os.Remove(keyfile.Name())

	key := NewEncryptionKey()
	_ = StoreEncryptionKey(key, keyfile.Name())

	retrievedKey, err := ReadEncryptionKey(keyfile.Name(), false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, key, retrievedKey)
}

func TestTokenEncryption(t *testing.T) {
	text := []byte("this is a sample text!")
	key := NewEncryptionKey()
	token, err := NewEncryptedToken(text, key)
	assert.NoError(t, err, "Token cannot be encrypted")

	dec, derr := DecryptToken(token, key)
	assert.NoError(t, derr, "Token cannot be decrypted")
	assert.Equal(t, text, dec)
}
