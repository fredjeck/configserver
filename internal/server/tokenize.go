package server

import (
	"io"
	"net/http"

	"github.com/fredjeck/configserver/internal/encryption"
)

// tokenizeText encrypts a pre-tokenized file
func (server *ConfigServer) tokenizeText(w http.ResponseWriter, req *http.Request) {

	value, err := io.ReadAll(req.Body)
	if err != nil {
		dieErr(w, req, http.StatusBadRequest, "cannot parse request body", err)
		return
	}

	tokenized, err := encryption.Tokenize(value, server.vault)
	if err != nil {
		dieErr(w, req, http.StatusInternalServerError, "unable to encrypt the provided file", err)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	server.writeResponse(http.StatusOK, tokenized, w)
}
