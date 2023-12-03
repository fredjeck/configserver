// Package encryption based on cryptopasta - basic cryptography examples
// Written in 2015 by George Tankersley <george.tankersley@gmail.com>
//
// Slightly modified by FredJeck to add support for HmacSha256 and various other things
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
)

// Aes256Key is an Alias for AES256 operations required a secret key
type Aes256Key struct {
	Key []byte
}

var ErrMalformedCipherText = errors.New("malformed cipher text")
var ErrInvalidToken = errors.New("invalid token")
var ErrCannotDecryptTokenContent = errors.New("unable to decrypt the token's content")

// NewAes256Key generates a random 256-bit key
func NewAes256Key() (*Aes256Key, error) {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}
	return &Aes256Key{Key: key[:]}, nil
}

// AesEncrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func AesEncrypt(plaintext []byte, key *Aes256Key) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key.Key[:])
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
func AesDecrypt(ciphertext []byte, key *Aes256Key) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key.Key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, ErrMalformedCipherText
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

// NewEncryptedToken encrypts the provided text into a substitution token
// Substition tokens are
func NewEncryptedToken(plaintext []byte, key *Aes256Key) (string, error) {
	enc, err := AesEncrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{enc:%s}", b64.StdEncoding.EncodeToString(enc)), nil
}

// Regex used to extract the value from a substitution token
var reTokenContent = regexp.MustCompile(`\{enc:(.*?)}`)

// DecryptToken extracts the payload from the provided substition token and decrypts its value using the given key
func DecryptToken(token string, key *Aes256Key) ([]byte, error) {
	if len(token) == 0 {
		return nil, ErrInvalidToken
	}

	match := reTokenContent.FindStringSubmatch(token)
	if len(match) != 2 {
		return nil, ErrInvalidToken
	}

	decoded, err := b64.StdEncoding.DecodeString(match[1])
	if err != nil {
		return nil, ErrCannotDecryptTokenContent
	}

	value, derr := AesDecrypt(decoded, key)
	if derr != nil {
		return nil, ErrCannotDecryptTokenContent
	}

	return value, nil
}
