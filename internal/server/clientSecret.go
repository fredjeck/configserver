package server

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/fredjeck/configserver/internal/utils"
	"log/slog"
	"strings"
	"time"
)

const (
	ClientSecretSeparatorChar = "|" // Char used to separate client secret component
	ClientSecretComponents    = 2   // Number of components used in client secrets
)

// generateClientSecret creates a new client secret which will be valid for the given number of days
// The generated client secret is bound to the provided client id
func generateClientSecret(clientId string, expiresInDays int, passPhrase string) (string, time.Time) {
	validity := time.Hour * 24 * time.Duration(expiresInDays)
	expires := time.Now().Add(validity)
	idStr := fmt.Sprintf("%s%s%s", expires.Format(time.RFC3339), ClientSecretSeparatorChar, clientId)
	return b64.StdEncoding.EncodeToString(utils.AesEncrypt(idStr, passPhrase)), expires
}

// validateClientSecret checks the provided client secret is valid and bound to the provided clientId
// if enforceValidity is true and if the secret expired validateClientSecret will consider the seccret as invalid
func validateClientSecret(clientId string, clientSecret string, passPhrase string, enforceValidity bool) bool {
	bytes, err := b64.StdEncoding.DecodeString(clientSecret)
	if err != nil {
		return false
	}

	secret, err := utils.AesDecrypt(bytes, passPhrase)
	if err != nil {
		return false
	}

	elements := strings.Split(string(secret), ClientSecretSeparatorChar)
	if len(elements) != ClientSecretComponents {
		return false
	}

	expiresAt, err := time.Parse(time.RFC3339, elements[0])
	if err == nil && time.Now().After(expiresAt) {
		if enforceValidity {
			return false
		}
		slog.Warn("client secret is expired, consider regenerating it", "client_id", clientId, "time_generated", expiresAt)
	}

	return elements[1] == clientId
}
