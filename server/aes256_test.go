package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := deriveKeyFromPassPhrase("This long passphrase is too long and it will be shortened to derive a key")
	assert.Len(t, key, 32)
}

func TestAesEncryptDecrypt(t *testing.T) {
	txt := "This text will be encrypted and decrypted with a passphrase"
	passPhrase := "This is a really long passphrase"

	b := Encrypt(txt, passPhrase)
	d, _ := Decrypt(b, passPhrase)

	assert.Equal(t, txt, d)
}

func TestAesFailToDecrypt(t *testing.T) {
	txt := "This text will be encrypted and decrypted with a passphrase"

	b := Encrypt(txt, "passPhrase used for encryption")
	_, err := Decrypt(b, "passPhrase used for decryption")

	assert.Error(t, err)
}
