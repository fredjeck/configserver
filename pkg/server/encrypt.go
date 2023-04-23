package server

import (
	"github.com/fredjeck/configserver/pkg/encrypt"
	"io"
	"net/http"
)

func (server *ConfigServer) encryptValue(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	w.Header().Add("Content-Type", "application/json")

	value, err := io.ReadAll(req.Body)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot parse the request body")
		return
	}

	token, err := encrypt.NewEncryptedToken(value, server.key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.writeResponse(http.StatusOK, []byte(token), w)
}
