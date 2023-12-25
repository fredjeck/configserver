package encryption

import (
	"regexp"
	"strings"
)

// SubstituteTokens replaces all the encoded token by their clear text value
func SubstituteTokens(file []byte, vault *KeyVault) ([]byte, error) {
	text := string(file)

	rx := regexp.MustCompile("({enc:.*?})")
	matches := rx.FindAllString(text, -1)
	if matches == nil {
		return file, nil
	}

	for _, match := range matches {
		clearText, err := vault.DecryptToken(match)
		if err != nil {
			continue
		}

		text = strings.Replace(text, match, string(clearText), -1)
	}

	return []byte(text), nil
}

// Tokenize replaces a pre-tokenized file tokens with encrypted tokens
func Tokenize(file []byte, vault *KeyVault) ([]byte, error) {
	text := string(file)

	rx := regexp.MustCompile("({enc:.*?})")
	matches := rx.FindAllString(text, -1)
	if matches == nil {
		return file, nil
	}

	for _, match := range matches {
		val := match[5 : len(match)-1]
		token, err := vault.CreateToken([]byte(val))
		if err != nil {
			continue
		}

		text = strings.Replace(text, match, token, -1)
	}

	return []byte(text), nil
}
