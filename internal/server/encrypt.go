package server

import (
	"io"
	"net/http"

	"github.com/fredjeck/configserver/internal/encryption"
)

// encryptValue generates encrypted substitution tokens
func (server *ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {

	value, err := io.ReadAll(req.Body)
	if err != nil {
		die(w, http.StatusBadRequest, "cannot parse request body")
		return
	}

	token, err := encryption.NewEncryptedToken(value, server.keystore.Aes256Key)
	if err != nil {
		die(w, http.StatusInternalServerError, "unable to encrypt the provided value")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	server.writeResponse(http.StatusOK, []byte(token), w)
}
