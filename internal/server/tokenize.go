package server

import (
	"github.com/fredjeck/configserver/internal/config"
	"github.com/fredjeck/configserver/internal/utils"
	"io"
	"net/http"
	"strings"
)

// Handles the clients file tokenization requests
func handleFileTokenization(c *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if len(contentType) == 0 || !strings.HasPrefix(contentType, "text") {
			HttpUnsupportedMediaType(w, "Unsupported content type '%s' only text/* is supported", contentType)
			return
		}

		value, err := io.ReadAll(r.Body)
		if err != nil {
			HttpInternalServerError(w, "Cannot parse request body")
			return
		}

		tokenized, err := utils.Tokenize(string(value), c.Server.PassPhrase)
		if err != nil {
			HttpInternalServerError(w, "An error occured while tokenizing the content")
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		Ok(w, []byte(tokenized))
	}
}
