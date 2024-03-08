package scratchpad

import (
	"encoding/json"
	"io"
	"net/http"
)

func handleClientRegistration(c *Configuration) http.Handler {
	// RegisterClientRequest represents the API's input body payload
	type registerRequest struct {
		ClientID string `json:"client_id"`
	}

	// RegisterClientResponse represents the API's output
	type registerResponse struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		ExpiresAt    int64  `json:"expires_at"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var registerRequest registerRequest
			body, err := io.ReadAll(r.Body)
			if err != nil {
				//dieErr(w, req, http.StatusBadRequest, "unable to parse request body", err)
				return
			}
			err = json.Unmarshal(body, &registerRequest)
			if err != nil {
				//dieErr(w, req, http.StatusBadRequest, "unable to parse request body", err)
				return
			}

		},
	)
}
