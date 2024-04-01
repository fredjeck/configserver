package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/stretchr/testify/assert"
)

var TokenizeTestConfiguration = configuration.DefaultConfiguration

const tokenizeURL = "/api/tokenize"

func TestMissingContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, tokenizeURL, nil)
	w := httptest.NewRecorder()
	f := handleFileTokenization(TokenizeTestConfiguration)
	f(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestInvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, tokenizeURL, nil)
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
	req := httptest.NewRequest(http.MethodPost, tokenizeURL, strings.NewReader(body))
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
