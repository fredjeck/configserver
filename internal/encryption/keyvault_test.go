package encryption

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestKeypairStorage(t *testing.T) {
	dir, err := os.MkdirTemp("", "configserver_tests")
	assert.NoError(t, err)
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir) // clean up

	kp, err := NewKeyVault()
	assert.NoError(t, err)

	err = kp.SaveTo(dir)
	assert.NoError(t, err)
}

func TestLoadAndCreate(t *testing.T) {
	dir, err := os.MkdirTemp("", "configserver_tests")
	assert.NoError(t, err)
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir) // clean up

	kp, err := LoadKeyVault(dir, true)
	assert.NoError(t, err)
	assert.NotNil(t, kp)
}

func TestSaveAndLoad(t *testing.T) {
	dir, err := os.MkdirTemp("", "configserver_tests")
	assert.NoError(t, err)
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir) // clean up

	kp, err := NewKeyVault()
	assert.NoError(t, err)

	err = kp.SaveTo(dir)
	assert.NoError(t, err)

	loaded, err := LoadKeyVault(dir, true)
	assert.NoError(t, err)
	assert.NotNil(t, kp)

	assert.Equal(t, kp.PrivateKey, loaded.PrivateKey)
}

func TestEncryptDecrypt(t *testing.T) {
	kp, err := NewKeyVault()
	assert.NoError(t, err)

	message := "We can't change the world with only ideas in our minds. We need conviction in our hearts."
	encrypted, err := kp.Encrypt([]byte(message))
	assert.NoError(t, err)

	decrypted, err := kp.Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, []byte(message), decrypted)
}

func TestSignature(t *testing.T) {
	kp, err := NewKeyVault()
	assert.NoError(t, err)

	message := "We can't change the world with only ideas in our minds. We need conviction in our hearts."
	signature, err := kp.Sign([]byte(message))
	assert.NoError(t, err)

	assert.NoError(t, kp.Verify([]byte(message), signature))
}

func TestTokens(t *testing.T) {
	kp, err := NewKeyVault()
	assert.NoError(t, err)

	value := "abcdefg123456"
	token, err := kp.CreateToken([]byte(value))
	assert.NoError(t, err)

	decoded, err := kp.DecryptToken(token)
	assert.NoError(t, err)
	assert.Equal(t, value, string(decoded))
}
