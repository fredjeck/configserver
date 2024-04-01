package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const passphrase = "This is a sample passphrase"

func TestCreateDecryptToken(t *testing.T) {
	content := "token content"
	token := CreateToken(content, passphrase)
	decrypted, err := DecryptToken(token, passphrase)

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(token, "{enc:"))
	assert.True(t, strings.HasSuffix(token, "}"))
	assert.NotContains(t, content, token)
	assert.Equal(t, content, decrypted)
}

func TestTokenize(t *testing.T) {
	text := "p1='{enc:value1}';p2='{enc:value2}';"

	tokenized, err := Tokenize(text, passphrase)
	assert.NoError(t, err)
	assert.NotEqual(t, text, tokenized)
	assert.NotContains(t, "value1", tokenized)
	assert.NotContains(t, "value2", tokenized)
}

func TestTokenSubstitution(t *testing.T) {
	text := fmt.Sprintf("p1='%s';p2='%s';p3='%s';p4='%s';",
		CreateToken("value 1", passphrase),
		CreateToken("value 2", passphrase),
		CreateToken("value 3", passphrase),
		CreateToken("value 4", passphrase),
	)

	clearText, err := Detokenize(text, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, "p1='value 1';p2='value 2';p3='value 3';p4='value 4';", clearText)
}
