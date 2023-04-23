package server

import (
	"encoding/json"
	"github.com/fredjeck/configserver/pkg/auth"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type RegisterClientRequest struct {
	ClientId     string   `json:"clientId"`
	Repositories []string `json:"repositories"`
}

type RegisterClientResponse struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (server *ConfigServer) registerClient(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

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

	if len(registerRequest.ClientId) == 0 {
		registerRequest.ClientId = uuid.NewString()
	}

	spec := auth.NewClientSpec(registerRequest.ClientId, registerRequest.Repositories)
	secret, err := spec.ClientSecret(server.key)
	if err != nil {
		server.writeError(http.StatusBadRequest, w, "Cannot generate client secret")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	resp := &RegisterClientResponse{ClientId: registerRequest.ClientId, ClientSecret: secret}
	values, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.writeResponse(http.StatusOK, values, w)
}
