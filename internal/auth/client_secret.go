package auth

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/fredjeck/configserver/internal/encryption"
	"log/slog"
	"strings"
	"time"
)

const ClientSecretSeparatorChar = "|"            // Char used to separate client secret component
const ClientSecretComponents = 3                 // Number of components used in client secrets
const ClientSecretValidity = time.Hour * 24 * 30 // 30 Days

// GenerateClientSecret creates a client secret for the provided client id using given encryption key
// This is purely experimental - goal is to generate a self-contained secret which can easily be validated
// and which does need to be stored locally
func GenerateClientSecret(clientId string, key *encryption.Aes256Key) (string, error) {
	secret, err := encryption.AesEncrypt([]byte(fmt.Sprintf("%s%s%s%s%s", time.Now().Format(time.RFC3339), ClientSecretSeparatorChar, clientId, ClientSecretSeparatorChar, encryption.RandomSequence(5))), key)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secret), nil
}

// ValidateClientSecret ensure the clientId and clientSecret pairs are matching
func ValidateClientSecret(clientId string, clientSecret string, key *encryption.Aes256Key) bool {
	bytes, err := b64.StdEncoding.DecodeString(clientSecret)
	if err != nil {
		return false
	}

	secret, err := encryption.AesDecrypt(bytes, key)
	if err != nil {
		return false
	}

	elements := strings.Split(string(secret), ClientSecretSeparatorChar)
	if len(elements) != ClientSecretComponents {
		return false
	}

	generatedAt, err := time.Parse(time.RFC3339, elements[0])
	if err == nil && generatedAt.Add(ClientSecretValidity).Before(time.Now()) {
		slog.Warn("client secret was generated more than 30 days ago consider regenerating it", "client_id", clientId, "time_generated", generatedAt)
	}

	return elements[1] == clientId
}
