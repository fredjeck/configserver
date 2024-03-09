package server

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var TokenizeTestConfiguration = &Configuration{
	PassPhrase:             "This is a passphrase used to protect yourself",
	ListenOn:               "127.0.0.1:4200",
	SecretExpiryDays:       60,
	ValidateSecretLifeSpan: true,
}

func TestMissingContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/tokenize", nil)
	w := httptest.NewRecorder()
	f := handleFileTokenization(TokenizeTestConfiguration)
	f(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestInvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/tokenize", nil)
	req.Header.Add("Content-Type", "image/png")
	w := httptest.NewRecorder()
	f := handleFileTokenization(TokenizeTestConfiguration)
	f(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestTokenization(t *testing.T) {
	body := `
{
	"property":"value",
	"otherProperty":{
		"subproperty":"{enc:EncodeMeFirst}"
	},
	"lastProperty""{enc:EncodeMeLast}"
}
`
	req := httptest.NewRequest(http.MethodPost, "/api/tokenize", strings.NewReader(body))
	req.Header.Add("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	f := handleFileTokenization(TokenizeTestConfiguration)
	f(w, req)

	res := w.Result()
	data, _ := io.ReadAll(res.Body)
	tokenized := string(data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, "{enc:EncodeMeFirst}", tokenized)
	assert.NotContains(t, "{enc:EncodeMeLast}", tokenized)
}
