// cryptopasta - basic cryptography examples
//
// Written in 2015 by George Tankersley <george.tankersley@gmail.com>
//
// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along
// with this software. If not, see // <http://creativecommons.org/publicdomain/zero/1.0/>.

// Package encrypt symmetric authenticated encryption using 256-bit AES-GCM with a random nonce.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

// ReadEncryptionKey reads the encryption key stored at the provided location.
// If createIfMissing is set to true, this function will attempt to create a new key if the file cannot be found
func ReadEncryptionKey(keyFilePath string, createIfMissing bool) (*[32]byte, error) {
	key := [32]byte{}

	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		if createIfMissing {
			err := StoreEncryptionKey(NewEncryptionKey(), keyFilePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	base64, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}

	decoded, err := b64.StdEncoding.DecodeString(string(base64))
	if err != nil {
		return nil, err
	}

	copy(key[:], decoded)
	return &key, nil
}

// StoreEncryptionKey stores the encryption key at the provided location
// Encryption keys are stored base 64 encoded
func StoreEncryptionKey(key *[32]byte, keyFilePath string) error {
	encoded := b64.StdEncoding.EncodeToString(key[:])
	err := os.WriteFile(keyFilePath, []byte(encoded), 0644)
	if err != nil {
		return err
	}
	return nil
}

// NewEncryptedToken encrypts the provided text into a substitution token
// Substition tokens are
func NewEncryptedToken(plaintext []byte, key *[32]byte) (string, error) {
	enc, err := Encrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{enc:%s}", b64.StdEncoding.EncodeToString(enc)), nil
}

// Regex used to extract the value from a substitution token
var reTokenContent = regexp.MustCompile(`\{enc:(.*?)\}`)

// DecryptToken extracts the payload from the provided substition token and decrypts its value using the given key
func DecryptToken(token string, key *[32]byte) ([]byte, error) {
	if len(token) == 0 {
		return nil, errors.New("Invalid token size 0")
	}

	match := reTokenContent.FindStringSubmatch(token)
	if len(match) != 2 {
		return nil, errors.New("Invalid token")
	}

	decoded, err := b64.StdEncoding.DecodeString(match[1])
	if err != nil {
		return nil, errors.New("Token cannot be decoded")
	}

	value, derr := Decrypt(decoded, key)
	if derr != nil {
		return nil, errors.New("Unable to decrypt the token's content")
	}

	return value, nil
}
