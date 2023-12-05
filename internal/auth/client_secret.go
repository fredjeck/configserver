package auth

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/fredjeck/configserver/internal/encryption"
	"log/slog"
	"math/rand"
	"strings"
	"time"
)

const ClientSecretSeparatorChar = "|"
const ClientSecretComponents = 3

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenerateClientSecret(clientId string, key *encryption.Aes256Key) (string, error) {
	secret, err := encryption.AesEncrypt([]byte(fmt.Sprintf("%s%s%s%s%s", time.Now().Format(time.RFC3339), ClientSecretSeparatorChar, clientId, ClientSecretSeparatorChar, randSeq(5))), key)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secret), nil
}

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
	if err == nil {
		if generatedAt.Add(time.Hour * 24 * 30).Before(time.Now()) {
			slog.Warn("Client secret was generated more than 30 days ago consider regenerating it", "client_id", clientId, "time_generated", generatedAt)
		}
	}

	return elements[1] == clientId
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
