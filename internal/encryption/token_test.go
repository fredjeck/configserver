package encryption

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenSubstitution(t *testing.T) {
	kp, err := NewKeyVault()
	assert.NoError(t, err)

	t1, _ := kp.CreateToken([]byte("value 1"))
	t2, _ := kp.CreateToken([]byte("value 2"))
	t3, _ := kp.CreateToken([]byte("value 3"))
	t4, _ := kp.CreateToken([]byte("value 4"))

	text := fmt.Sprintf("p1='%s';p2='%s';p3='%s';p4='%s';", t1, t2, t3, t4)

	clearText, err := SubstituteTokens([]byte(text), kp)
	assert.NoError(t, err)
	assert.Equal(t, []byte("p1='value 1';p2='value 2';p3='value 3';p4='value 4';"), clearText)
}

func TestTokenize(t *testing.T) {
	kp, err := NewKeyVault()
	assert.NoError(t, err)

	text := "p1='{enc:value1}';p2='{enc:value2}';"

	tokenized, err := Tokenize([]byte(text), kp)
	assert.NoError(t, err)
	assert.NotEqual(t, text, tokenized)
	assert.NotContains(t, "value1", string(tokenized))
	assert.NotContains(t, "value2", string(tokenized))
}
