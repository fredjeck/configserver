package encryption

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionKeyStorage(t *testing.T) {
	keyfile, _ := os.CreateTemp("", "keyfile")
	defer func(name string) {
		_ = os.Remove(name)
	}(keyfile.Name())

	key, _ := NewAes256Key()
	_ = StoreKeyToPath(key.Key, KindSha, keyfile.Name())

	retrievedKey, err := LoadKeyFromPath(keyfile.Name())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, key.Key, retrievedKey)
}
