package encrypt

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	key := NewEncryptionKey()
	text := []byte("this is a sample text!")
	s, err := Encrypt(text, key)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s", s)

	dec, derr := Decrypt(s, key)
	if derr != nil {
		t.Error(derr)
	}
	assert.EqualValues(t, text, dec)
}

func TestEncryptionKeyStorage(t *testing.T) {
	keyfile, _ := os.CreateTemp("", "keyfile")
	key := NewEncryptionKey()
	_ = StoreEncryptionKey(key, keyfile.Name())
	retrievedKey, err := ReadEncryptionKey(keyfile.Name(), false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, key, retrievedKey)
}
