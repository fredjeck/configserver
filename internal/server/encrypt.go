package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fredjeck/configserver/internal/encryption"
)

type EncryptResponse struct {
	Token string `json:"token"`
}

// encryptValue generates encrypted substitution tokens
func (server *ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {

	value, err := io.ReadAll(req.Body)
	if err != nil {
		dieErr(w, http.StatusBadRequest, "cannot parse request body", err)
		return
	}

	token, err := encryption.NewEncryptedToken(value, server.keystore.Aes256Key)
	if err != nil {
		dieErr(w, http.StatusInternalServerError, "unable to encrypt the provided value", err)
		return
	}

	response := &EncryptResponse{
		token,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		dieErr(w, http.StatusInternalServerError, "unable to marshall token response", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	server.writeResponse(http.StatusOK, responseJson, w)
}
