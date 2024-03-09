package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// RegisterClientResponse represents the API's output
type RegisterClientResponse struct {
	ClientId     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func handleClientRegistration(c *Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientId := r.URL.Query().Get("client_id")
		if len(clientId) == 0 {
			uid, _ := uuid.NewV7()
			clientId = uid.String()
		}

		clientSecret, expires := Generate(clientId, c.SecretExpiryDays, c.PassPhrase)

		jsonStr, err := json.Marshal(&RegisterClientResponse{
			clientId, clientSecret, expires,
		})
		if err != nil {
			HttpInternalServerError(w, err.Error())
			return
		}

		Ok(w, jsonStr)
	}
}
