package encryption

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"log/slog"
)

// KeyFileKindIndentifier represents the keyfile kind and is used when storing the keys
type KeyFileKindIndentifier string

const (
	// KindAes represents the AES keyfile header tag
	KindAes KeyFileKindIndentifier = "AES"
	// KindSha represents the SHA keyfile header tag
	KindSha KeyFileKindIndentifier = "SHA"
	// AesKeyFileName represents the name of the AES Keyfile name
	AesKeyFileName = "id_aes"
	// ShaKeyFileName represents the name of the SHA Keyfile name
	ShaKeyFileName = "id_sha"
	// WarnKeyfileGenerated is a warning message indicating a keyfile was not found and generated automatically
	WarnKeyfileGenerated = "A key file was not found. A default key file has been created for you. If the key storage location is not persistent it is highly recommended that you either provide a persistent storage or your provide your own keys. Please refer to the user manual."
	// ArgKeyFilePath is used as context param for keyfile path
	ArgKeyFilePath = "keyfile_path"
)

// Keystore is a utility struct storing the various keys and secret used by ConfigServer
// This is a temporary solution which needs to be improved - me don't like it
type Keystore struct {
	Aes256Key        *Aes256Key
	HmacSha256Secret *HmacSha256Secret
}

// LoadKeyStoreFromPath loads the keystore from the provided path
// Expects the path to exist and to contain the following files id_aes and id_sha
// When not found the files will be created, no error will be triggered - as this is not recommended a warning will be logged
func LoadKeyStoreFromPath(path string) (*Keystore, error) {
	store := &Keystore{}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	// TODO Refactor and remove dupliucation maybe using a keygen function type
	aesKeyPath := filepath.Join(path, AesKeyFileName)
	if _, err := os.Stat(aesKeyPath); errors.Is(err, os.ErrNotExist) {
		slog.Warn(WarnKeyfileGenerated, ArgKeyFilePath, aesKeyPath)
		if store.Aes256Key, err = NewAes256Key(); err != nil {
			return nil, err
		}
		if err := StoreKeyToPath(store.Aes256Key.Key, KindAes, aesKeyPath); err != nil {
			return nil, err
		}
	}
	aes, err := LoadKeyFromPath(aesKeyPath)
	if err != nil {
		return nil, err
	}
	store.Aes256Key = &Aes256Key{Key: aes}
	slog.Info("aes256 keyfile loaded", ArgKeyFilePath, aesKeyPath)

	shaKeyPath := filepath.Join(path, ShaKeyFileName)
	if _, err := os.Stat(shaKeyPath); errors.Is(err, os.ErrNotExist) {
		slog.Warn(WarnKeyfileGenerated, ArgKeyFilePath, shaKeyPath)
		if store.HmacSha256Secret, err = NewHmacSha256Secret(); err != nil {
			return nil, err
		}
		if err := StoreKeyToPath(store.HmacSha256Secret.Key, KindSha, shaKeyPath); err != nil {
			return nil, err
		}
	}
	sha, err := LoadKeyFromPath(shaKeyPath)
	if err != nil {
		return nil, err
	}
	store.HmacSha256Secret = &HmacSha256Secret{Key: sha}
	slog.Info("hmacsha256 keyfile loaded", ArgKeyFilePath, shaKeyPath)

	return store, nil
}

// StoreKeyToPath stores the encryption key at the provided location
func StoreKeyToPath(key []byte, kindIdentifier KeyFileKindIndentifier, keyFilePath string) error {
	encoded := b64.StdEncoding.EncodeToString(key[:])
	kind := strings.ToUpper(string(kindIdentifier))
	content := fmt.Sprintf("-----BEGIN %s PRIVATE KEY-----\n%s\n-----END %s PRIVATE KEY-----", kind, encoded, kind)

	err := os.WriteFile(keyFilePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadKeyFromPath loads a key from the local storage
func LoadKeyFromPath(keyFilePath string) ([]byte, error) {
	contents, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}

	components := strings.Split(string(contents), "\n")
	if len(components) != 3 {
		return nil, fmt.Errorf("'%s' invalid keyfile format", keyFilePath)
	}

	decoded, err := b64.StdEncoding.DecodeString(components[1])
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
