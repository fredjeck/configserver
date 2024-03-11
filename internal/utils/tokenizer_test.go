package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateDecryptToken(t *testing.T) {
	pass := "This is a sample passphrase"
	content := "token content"
	token := CreateToken(content, pass)
	decrypted, err := DecryptToken(token, pass)

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(token, "{enc:"))
	assert.True(t, strings.HasSuffix(token, "}"))
	assert.NotContains(t, content, token)
	assert.Equal(t, content, decrypted)
}

func TestTokenize(t *testing.T) {
	pass := "This is a sample passphrase"

	text := "p1='{enc:value1}';p2='{enc:value2}';"

	tokenized, err := Tokenize(text, pass)
	assert.NoError(t, err)
	assert.NotEqual(t, text, tokenized)
	assert.NotContains(t, "value1", tokenized)
	assert.NotContains(t, "value2", tokenized)
}

func TestTokenSubstitution(t *testing.T) {
	pass := "This is a sample passphrase"

	text := fmt.Sprintf("p1='%s';p2='%s';p3='%s';p4='%s';",
		CreateToken("value 1", pass),
		CreateToken("value 2", pass),
		CreateToken("value 3", pass),
		CreateToken("value 4", pass),
	)

	clearText, err := Detokenize(text, pass)
	assert.NoError(t, err)
	assert.Equal(t, "p1='value 1';p2='value 2';p3='value 3';p4='value 4';", clearText)
}
