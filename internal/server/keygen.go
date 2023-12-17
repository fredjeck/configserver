package server

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fredjeck/configserver/internal/encryption"
)

// KeyGenResponse represents the API's output
type KeyGenResponse struct {
	Kind string `json:"kind"`
	Key  string `json:"key"`
}

// GenAes256 generates a new AES key
func (server *ConfigServer) genAes256(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		die(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	key, err := encryption.NewAes256Key()

	resp := &KeyGenResponse{Kind: "AES256", Key: base64.StdEncoding.EncodeToString(key.Key)}
	responseJson, err := json.Marshal(resp)
	if err != nil {
		dieErr(w, req, http.StatusInternalServerError, "unable to generate key", err)
		return
	}

	slog.Debug("new HS256 secret generated", "secret", resp.Key)
	w.Header().Add("Content-Type", "application/json")
	server.writeResponse(http.StatusOK, responseJson, w)
}

// GenHmacSha256 generates a new Hmac key
func (server *ConfigServer) genHmacSha256(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		die(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	key, err := encryption.NewHmacSha256Secret()

	resp := &KeyGenResponse{Kind: "HS256", Key: base64.StdEncoding.EncodeToString(key.Key)}
	responseJson, err := json.Marshal(resp)
	if err != nil {
		dieErr(w, req, http.StatusInternalServerError, "unable to generate key", err)
		return
	}
	slog.Debug("new HS256 secret generated", "secret", resp.Key)

	w.Header().Add("Content-Type", "application/json")
	server.writeResponse(http.StatusOK, responseJson, w)
}
