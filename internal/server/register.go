package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fredjeck/configserver/internal/config"
	"github.com/google/uuid"
)

// RegisterClientResponse represents the API's output
type RegisterClientResponse struct {
	ClientId     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// handleClientRegistration responds to client registration requests
func handleClientRegistration(c *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientId := r.URL.Query().Get("client_id")
		if len(clientId) == 0 {
			uid, _ := uuid.NewV7()
			clientId = uid.String()
		}

		clientSecret, expires := generateClientSecret(clientId, c.Server.SecretExpiryDays, c.Server.PassPhrase)

		jsonStr, err := json.Marshal(&RegisterClientResponse{
			clientId, clientSecret, expires,
		})
		if err != nil {
			HttpInternalServerError(w, err.Error())
			return
		}
		Ok(w, jsonStr, "application/json;charset=utf-8")
	}
}
