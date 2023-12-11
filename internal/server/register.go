package server

import (
	"encoding/json"
	"github.com/fredjeck/configserver/internal/auth"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
)

// RegisterClientRequest represents the API's input body payload
type RegisterClientRequest struct {
	ClientID     string   `json:"clientID"`
	Repositories []string `json:"repositories"`
}

// RegisterClientResponse represents the API's output
type RegisterClientResponse struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

// generateClientSecret is a handlerFunc which generates a new client secret for the provided client id
func (server *ConfigServer) generateClientSecret(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		die(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var registerRequest RegisterClientRequest
	body, err := io.ReadAll(req.Body)
	if err != nil {
		dieErr(w, http.StatusBadRequest, "unable to parse request body", err)
		return
	}
	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		dieErr(w, http.StatusBadRequest, "unable to parse request body", err)
		return
	}

	if len(registerRequest.ClientID) == 0 {
		registerRequest.ClientID = uuid.NewString()
	}

	clientSecret, err := auth.GenerateClientSecret(registerRequest.ClientID, server.keystore.Aes256Key)
	if err != nil {
		die(w, http.StatusInternalServerError, "failed to generate client secret")
		return
	}

	resp := &RegisterClientResponse{ClientID: registerRequest.ClientID, ClientSecret: clientSecret}
	values, err := json.Marshal(resp)
	if err != nil {
		dieErr(w, http.StatusInternalServerError, "an error occurred while generating the server response", err)
		return
	}

	slog.Debug("new client secret generated", "clientID", registerRequest.ClientID)
	w.Header().Add("Content-Type", "application/json")
	server.writeResponse(http.StatusOK, values, w)
}
