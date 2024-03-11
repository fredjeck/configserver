package utils

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrCannotDecryptTokenContent = errors.New("unable to decrypt the token's content")

// Regex used to extract the value from a substitution token
var reToken = regexp.MustCompile(`\{enc:(.*?)}`)

// DecryptToken extracts the payload from the provided substition token and decrypts its value using the given key
func DecryptToken(token string, passphrase string) (string, error) {
	if len(token) == 0 {
		return "", ErrInvalidToken
	}

	match := reToken.FindStringSubmatch(token)
	if len(match) != 2 {
		return "", ErrInvalidToken
	}

	decoded, err := b64.StdEncoding.DecodeString(match[1])
	if err != nil {
		return "", ErrCannotDecryptTokenContent
	}

	value, err := AesDecrypt(decoded, passphrase)
	if err != nil {
		return "", ErrCannotDecryptTokenContent
	}

	return value, nil
}

// CreateToken creates an encrypted substitution token from the given value
func CreateToken(text string, passphrase string) string {
	return fmt.Sprintf("{enc:%s}", b64.StdEncoding.EncodeToString(AesEncrypt(text, passphrase)))
}

// Detokenize replaces all the encoded token by their clear text value
func Detokenize(text string, passphrase string) (string, error) {
	rx := regexp.MustCompile("({enc:.*?})")
	matches := rx.FindAllString(text, -1)
	if matches == nil {
		return text, nil
	}

	for _, match := range matches {
		clearText, err := DecryptToken(match, passphrase)
		if err != nil {
			continue
		}

		text = strings.Replace(text, match, clearText, -1)
	}

	return text, nil
}

// Tokenize replaces a pre-tokenized file tokens with encrypted tokens
func Tokenize(text string, passphrase string) (string, error) {
	rx := regexp.MustCompile("({enc:.*?})")
	matches := rx.FindAllString(text, -1)
	if matches == nil {
		return text, nil
	}

	for _, match := range matches {
		val := match[5 : len(match)-1]
		text = strings.Replace(text, match, CreateToken(val, passphrase), -1)
	}

	return text, nil
}
