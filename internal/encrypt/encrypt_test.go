package encrypt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAesDecrypt(t *testing.T) {
	key, _ := NewAes256Key()
	text := []byte("this is a sample text!")
	s, err := AesEncrypt(text, key)
	assert.NoError(t, err, "Encryption should not raise an error")

	dec, derr := AesDecrypt(s, key)
	assert.NoError(t, derr, "Decryption should not raise an error")
	assert.EqualValues(t, text, dec)
}

func TestEncryptionKeyStorage(t *testing.T) {
	keyfile, _ := os.CreateTemp("", "keyfile")
	defer os.Remove(keyfile.Name())

	key, _ := NewAes256Key()
	_ = StoreEncryptionKey(key, keyfile.Name())

	retrievedKey, err := ReadEncryptionKey(keyfile.Name(), false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, key, retrievedKey)
}

func TestTokenEncryption(t *testing.T) {
	text := []byte("this is a sample text!")
	key, _ := NewAes256Key()
	token, err := NewEncryptedToken(text, key)
	assert.NoError(t, err, "Token cannot be encrypted")

	dec, derr := DecryptToken(token, key)
	assert.NoError(t, derr, "Token cannot be decrypted")
	assert.Equal(t, text, dec)
}

func TestInvalidTokenDecryption(t *testing.T) {
	key, _ := NewAes256Key()
	token := "{e:abc}"

	_, derr := DecryptToken(token, key)
	assert.Error(t, derr)
}

func TestZeroLengthTokenDecryption(t *testing.T) {
	key, _ := NewAes256Key()
	token := ""

	_, derr := DecryptToken(token, key)
	assert.Error(t, derr)
}

func TestInvalidPayloadTokenDecryption(t *testing.T) {
	key, _ := NewAes256Key()
	token := "{enc:ABCDEFG}"

	_, derr := DecryptToken(token, key)
	assert.Error(t, derr)
}

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
