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

// RegisterClientResponseJwt represent a JWT response web token
type RegisterClientResponseJwt struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// registerClient is a handlerFunc which generates a new client secret from the provided data.
func (server *ConfigServer) registerClient(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var registerRequest RegisterClientRequest
	body, err := io.ReadAll(req.Body)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot parse the request body")
		return
	}
	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot unmarshal the request body")
		return
	}

	if len(registerRequest.ClientID) == 0 {
		registerRequest.ClientID = uuid.NewString()
	}

	spec := auth.NewClientKeySpec(registerRequest.ClientID, registerRequest.Repositories)
	secret, err := spec.GenerateSecret(server.keystore.Aes256Key)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot generate client secret")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	resp := &RegisterClientResponse{ClientID: registerRequest.ClientID, ClientSecret: secret}
	values, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug("New client secret generated", "clientID", registerRequest.ClientID)
	server.writeResponse(http.StatusOK, values, w)
}

// registerClient is a handlerFunc which generates a new client secret from the provided data.
func (server *ConfigServer) registerClientJwt(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var registerRequest RegisterClientRequest
	body, err := io.ReadAll(req.Body)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot parse the request body")
		return
	}
	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot unmarshal the request body")
		return
	}

	if len(registerRequest.ClientID) == 0 {
		registerRequest.ClientID = uuid.NewString()
	}

	token := auth.NewJSONWebToken()
	token.Payload.Audience = registerRequest.Repositories
	token.Payload.Issuer = "ConfigServer"
	token.Payload.Subject = registerRequest.ClientID

	w.Header().Add("Content-Type", "application/json")
	resp := &RegisterClientResponseJwt{TokenType: "Bearer", AccessToken: token.Pack(server.keystore.HmacSha256Secret)}
	values, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug("New bearer token generated", "ClientId", registerRequest.ClientID)
	server.writeResponse(http.StatusOK, values, w)
}

// registerClient is a handlerFunc which generates a new client secret from the provided data.
func (server *ConfigServer) generateClientSecret(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var registerRequest RegisterClientRequest
	body, err := io.ReadAll(req.Body)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot parse the request body")
		return
	}
	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot unmarshal the request body")
		return
	}

	if len(registerRequest.ClientID) == 0 {
		registerRequest.ClientID = uuid.NewString()
	}

	clientSecret, err := auth.GenerateClientSecret(registerRequest.ClientID, server.keystore.Aes256Key)
	if err != nil {
		server.writeError(http.StatusInternalServerError, w, "Cannot generate client secret")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	resp := &RegisterClientResponse{ClientID: registerRequest.ClientID, ClientSecret: clientSecret}
	values, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug("New client secret generated", "clientID", registerRequest.ClientID)
	server.writeResponse(http.StatusOK, values, w)
}
