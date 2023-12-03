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

func (server *ConfigServer) GenAes256(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key, err := encryption.NewAes256Key()

	w.Header().Add("Content-Type", "application/json")
	resp := &KeyGenResponse{Kind: "AES256", Key: base64.StdEncoding.EncodeToString(key.Key)}
	json, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("New HS256 secret generated", "secret", resp.Key)
	server.writeResponse(http.StatusOK, json, w)
}

func (server *ConfigServer) GenHmacSha256(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key, err := encryption.NewHmacSha256Secret()

	w.Header().Add("Content-Type", "application/json")
	resp := &KeyGenResponse{Kind: "HS256", Key: base64.StdEncoding.EncodeToString(key.Key)}
	json, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("New HS256 secret generated", "secret", resp.Key)
	server.writeResponse(http.StatusOK, json, w)
}
