package encryption

import (
	"errors"
	"os"
	"path/filepath"
)

// AesKeyFileName represents the name of the AES Keyfile name
const AesKeyFileName = "id_aes"

// ShaKeyFileName represents the name of the SHA Keyfile name
const ShaKeyFileName = "id_sha"

// WarnKeyfileGenerated is a warning message indicating a keyfile was not found and generated automatically
const WarnKeyfileGenerated = ""

// Keystore is an utility struct storing the various keys and secret used by ConfigServer
// This is a temporary solution which needs to be improved - me dont like it
type Keystore struct {
	AesAes256Key     Aes256Key
	HmacSha256Secret HmacSha256Secret
}

// LoadKeyStoreFromPath loads the keystore from the provided path
// Expects the path to exist and to contain the following files id_aes and id_sha
// When not found the files will be created, no error will be triggered
func LoadKeyStoreFromPath(path string) (*Keystore, error) {
	store := &Keystore{}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	if _, err := os.Stat(filepath.Join(path, AesKeyFileName)); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
	}

	if _, err := os.Stat(filepath.Join(path, ShaKeyFileName)); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
	}

	return store, nil
}
