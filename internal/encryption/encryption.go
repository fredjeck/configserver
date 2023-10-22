// Package encryption based on cryptopasta - basic cryptography examples
// Written in 2015 by George Tankersley <george.tankersley@gmail.com>
//
// Slightly modified by FredJeck to add support for HmacSha256 and various other things
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Aes256Key is an Alias for AES256 operations required a secret key
type Aes256Key *[32]byte

// HmacSha256Secret is an alias for Hmac256 secret key
type HmacSha256Secret *[64]byte

// NewAes256Key generates a random 256-bit key
func NewAes256Key() (Aes256Key, error) {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// NewHmacSha256Secret generates a random 512-bit secret
func NewHmacSha256Secret() (HmacSha256Secret, error) {
	key := [64]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// HmacSha256Hash generates an HmacSha256 Hash from the provided data and secret
func HmacSha256Hash(data []byte, secret HmacSha256Secret) []byte {
	hmac := hmac.New(sha256.New, secret[:])
	hmac.Write(data)
	return hmac.Sum(nil)
}

// AesEncrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func AesEncrypt(plaintext []byte, key Aes256Key) (ciphertext []byte, err error) {
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

// AesDecrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func AesDecrypt(ciphertext []byte, key Aes256Key) (plaintext []byte, err error) {
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
func ReadEncryptionKey(keyFilePath string, createIfMissing bool) (Aes256Key, error) {
	key := [32]byte{}

	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		if createIfMissing {
			key, err := NewAes256Key()
			if err != nil {
				return nil, err
			}
			err = StoreEncryptionKey(key, keyFilePath)
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
func StoreEncryptionKey(key Aes256Key, keyFilePath string) error {
	encoded := b64.StdEncoding.EncodeToString(key[:])
	err := os.WriteFile(keyFilePath, []byte(encoded), 0644)
	if err != nil {
		return err
	}
	return nil
}

// NewEncryptedToken encrypts the provided text into a substitution token
// Substition tokens are
func NewEncryptedToken(plaintext []byte, key Aes256Key) (string, error) {
	enc, err := AesEncrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{enc:%s}", b64.StdEncoding.EncodeToString(enc)), nil
}

// Regex used to extract the value from a substitution token
var reTokenContent = regexp.MustCompile(`\{enc:(.*?)\}`)

// DecryptToken extracts the payload from the provided substition token and decrypts its value using the given key
func DecryptToken(token string, key Aes256Key) ([]byte, error) {
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

	value, derr := AesDecrypt(decoded, key)
	if derr != nil {
		return nil, errors.New("Unable to decrypt the token's content")
	}

	return value, nil
}
