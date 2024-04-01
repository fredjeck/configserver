package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/google/uuid"
)

// RegisterClientResponse represents the API's output
type RegisterClientResponse struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// handleClientRegistration responds to client registration requests
func handleClientRegistration(c *configuration.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		if len(clientID) == 0 {
			uid, _ := uuid.NewV7()
			clientID = uid.String()
		}

		clientSecret, expires := generateClientSecret(clientID, c.Server.SecretExpiryDays, c.Server.PassPhrase)

		jsonStr, err := json.Marshal(&RegisterClientResponse{
			clientID, clientSecret, expires,
		})
		if err != nil {
			HTTPInternalServerError(w, err.Error())
			return
		}
		Ok(w, jsonStr, "application/json;charset=utf-8")
	}
}
