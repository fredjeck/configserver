package encryption

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionKeyStorage(t *testing.T) {
	keyfile, _ := os.CreateTemp("", "keyfile")
	defer os.Remove(keyfile.Name())

	key, _ := NewAes256Key()
	_ = StoreKeyToPath(key[:], KindSha, keyfile.Name())

	retrievedKey, err := LoadKeyFromPath(keyfile.Name())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, key, Aes256Key(retrievedKey))
}
