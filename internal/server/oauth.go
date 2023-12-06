package server

import (
	b64 "encoding/base64"
	"github.com/fredjeck/configserver/internal/auth"
	"net/http"
	"strings"
)

func (server *ConfigServer) authorize(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		die(w, http.StatusMethodNotAllowed, "only POST allowed")
		return
	}

	authorization := req.Header.Get("Authorization")
	if len(authorization) == 0 {
		die(w, http.StatusUnauthorized, "authorization required")
		return
	}

	basicAuth, err := b64.StdEncoding.DecodeString(strings.ReplaceAll(authorization, "Basic ", ""))
	if err != nil {
		die(w, http.StatusUnauthorized, "missing credentials")
		return
	}

	credentials := strings.Split(string(basicAuth), ":")
	if len(credentials) != 2 {
		die(w, http.StatusUnauthorized, "missing credentials")
		return
	}

	_ = req.ParseForm()
	if "token" != req.Form.Get("response_type") {
		die(w, http.StatusBadRequest, "unsupported response type, only token is currently supported")
		return
	}

	scope := req.Form.Get("response_type")
	if len(scope) == 0 {
		die(w, http.StatusBadRequest, "missing property : 'scope' is required")
		return
	}

	clientId := req.Form.Get("client_id")
	if len(clientId) == 0 {
		die(w, http.StatusBadRequest, "missing property : 'client_id' is required")
		return
	}
	if clientId != credentials[0] {
		die(w, http.StatusUnauthorized, "'client_id' property differs from the authorization header")
		return
	}

	if !auth.ValidateClientSecret(credentials[0], credentials[1], server.keystore.Aes256Key) {
		die(w, http.StatusUnauthorized, "unknown user")
		return
	}

	//scopes := strings.Split(scope, " ")

}
