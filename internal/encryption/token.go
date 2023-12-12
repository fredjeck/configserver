package encryption

import (
	"regexp"
	"strings"
)

// SubstituteTokens replaces all the encoded token by their clear text value
func SubstituteTokens(file []byte, key *Aes256Key) ([]byte, error) {
	text := string(file)

	rx := regexp.MustCompile("({enc:.*?})")
	matches := rx.FindAllString(text, -1)
	if matches == nil {
		return file, nil
	}

	for _, match := range matches {
		clearText, err := DecryptToken(match, key)
		if err != nil {
			continue
		}

		text = strings.Replace(text, match, string(clearText), -1)
	}

	return []byte(text), nil
}
