package server

import (
	"io"
	"net/http"

	"github.com/fredjeck/configserver/internal/encryption"
)

func (server *ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	value, err := io.ReadAll(req.Body)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot parse the request body")
		return
	}

	token, err := encryption.NewEncryptedToken(value, server.keystore.Aes256Key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.writeResponse(http.StatusOK, []byte(token), w)
}
