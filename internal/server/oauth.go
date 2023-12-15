package server

import (
	b64 "encoding/base64"
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

	authorization := req.Header.Get("Authorization")
	if len(authorization) == 0 {
		die(w, http.StatusUnauthorized, "missing authorization header")
		return
	}

	basicAuth, err := b64.StdEncoding.DecodeString(strings.ReplaceAll(authorization, "Basic ", ""))
	if err != nil {
		dieErr(w, http.StatusUnauthorized, "incorrect authorization header", err)
		return
	}

	credentials := strings.Split(string(basicAuth), ":")
	if len(credentials) != 2 {
		die(w, http.StatusUnauthorized, "incorrect authorization header")
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

	clientId := credentials[0]
	if !auth.ValidateClientSecret(clientId, credentials[1], server.keystore.Aes256Key) {
		slog.Warn("login rejected", "client_id", clientId)
		dief(w, http.StatusUnauthorized, "'%s' : unauthorized", clientId)
		return
	}

	slog.Info("successful login", "client_id", clientId)

	var allowedScopes []string
	for _, scope := range scopes {
		if repo, ok := server.repository.Repositories[scope]; ok {
			if repo.IsClientAllowed(clientId) {
				allowedScopes = append(allowedScopes, scope)
			}
		}
	}

	response := &AccessTokenResponse{TokenType: "bearer"}
	response.TokenType = "bearer"

	token := auth.NewJSONWebToken()
	token.Payload.Audience = allowedScopes
	token.Payload.Issuer = "ConfigServer"
	token.Payload.Subject = clientId

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

func (server *ConfigServer) extractToken(w http.ResponseWriter, r *http.Request) (*auth.JSONWebToken, bool) {
	authorization, ok := r.Header["Authorization"]
	if !ok {
		die(w, http.StatusBadRequest, "missing authorization header")
		return nil, false
	}

	authStr := strings.ToLower(authorization[0])
	if !strings.Contains(authStr, "bearer") {
		die(w, http.StatusBadRequest, "only bearer authorization is supported")
		return nil, false
	}

	token := strings.Replace(authStr, "bearer ", "", -1)
	err := auth.VerifySignature(token, server.keystore.HmacSha256Secret)
	if err != nil {
		dieErr(w, http.StatusUnauthorized, "not authorized", err)
		return nil, false
	}

	jwt, err := auth.ParseJwt(token, server.keystore.HmacSha256Secret)
	if err != nil {
		dieErr(w, http.StatusUnauthorized, "invalid token", err)
		return nil, false
	}

	return jwt, true
}
