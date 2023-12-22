package server

import (
	"encoding/json"
	"github.com/fredjeck/configserver/internal/auth"
	"log/slog"
	"net/http"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

// authorize exposes an oauth2 compatible authorization endpoints.
// authorize only supports the client_credentials grant type and expect to find credentials within the authorization header
// Example request :
// POST http://localhost:8080/oauth2/authorize HTTP/1.1
// content-type: application/x-www-form-urlencoded
// Authorization: Basic basic_auth_b64
//
// grant_type=client_crendentials&scope=repo1%2Frepo2
func (server *ConfigServer) authorize(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		die(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	authorization, err := auth.FromRequest(req, server.keystore, auth.AuthorizationKindBasic)
	if err != nil {
		dieErr(w, req, http.StatusUnauthorized, "authorization failed", err)
		return
	}

	_ = req.ParseForm()
	if "client_credentials" != req.Form.Get("grant_type") {
		die(w, http.StatusBadRequest, "unsupported grant type, only client_credentials is supported")
		return
	}

	scopeStr := req.Form.Get("scope")
	if len(scopeStr) == 0 {
		die(w, http.StatusBadRequest, "'scope' is required")
		return
	}

	scopes := strings.Split(scopeStr, " ")
	if len(scopes) == 0 {
		die(w, http.StatusBadRequest, "at least one scope is required")
	}

	slog.Info("successful login", "client_id", authorization.ClientId())

	var allowedScopes []string
	for _, scope := range scopes {
		if server.repository.IsClientAllowed(scope, authorization.ClientId()) {
			allowedScopes = append(allowedScopes, scope)
		}
	}

	response := &AccessTokenResponse{TokenType: "bearer"}
	response.TokenType = "bearer"

	token := auth.NewJSONWebToken()
	token.Payload.Audience = allowedScopes
	token.Payload.Issuer = "ConfigServer"
	token.Payload.Subject = authorization.ClientId()

	response.AccessToken = token.Pack(server.keystore.HmacSha256Secret)
	response.ExpiresIn = token.Payload.Expires - token.Payload.IssuedAt
	response.Scope = strings.Join(allowedScopes, " ")

	values, err := json.Marshal(response)
	if err != nil {
		die(w, http.StatusInternalServerError, "unable to marshall token response")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	server.writeResponse(http.StatusOK, values, w)
}
