package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fredjeck/configserver/internal/encrypt"
)

// JSONWebTokenPayload represents the Payload part of a JWT
type JSONWebTokenPayload struct {
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	Audience  []string `json:"aud"`
	NotBefore string   `json:"nbf"`
}

// JSONWebTokenHeader represents the Payload part of a JWT
type JSONWebTokenHeader struct {
	Alg  string `json:"alg"`
	Type string `json:"typ"`
}

// JSONWebToken is a basic implementation of the JSON Web Token rfc7519
type JSONWebToken struct {
	Header  *JSONWebTokenHeader
	Payload *JSONWebTokenPayload
}

// NewJSONWebToken creates a new empty JSON Web Token
func NewJSONWebToken() *JSONWebToken {
	return &JSONWebToken{Header: &JSONWebTokenHeader{Alg: "HS256", Type: "JWT"}, Payload: &JSONWebTokenPayload{}}
}

// Pack generates the token by marshalling the content to JSON, formatting the output into b64UrlEncoded strings and by appending the token signature
func (jwt *JSONWebToken) Pack(secret encrypt.HmacSha256Secret) string {
	tk := jwt.token()
	hash := encrypt.HmacSha256Hash([]byte(tk), secret)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s.%s", tk, b64Hash)
}

// VerifySignature validates a token has been signed with the provided key
func VerifySignature(token string, secret encrypt.HmacSha256Secret) error {
	components := strings.Split(token, ".")
	if len(components) != 3 {
		return fmt.Errorf("malformed jwt token - expecting three components only %d parts found", len(components))
	}

	tk := components[0] + "." + components[1]
	signature := components[2]

	hash := encrypt.HmacSha256Hash([]byte(tk), secret)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash)

	if strings.Compare(b64Hash, signature) != 0 {
		return fmt.Errorf("invalid token signature")
	}

	return nil
}

func (jwt *JSONWebToken) token() string {
	header, _ := jsonb64UrlEncode(jwt.Header)
	body, _ := jsonb64UrlEncode(jwt.Payload)
	return fmt.Sprintf("%s.%s", header, body)
}

func jsonb64UrlEncode(e interface{}) (string, error) {
	json, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(json), nil
}
