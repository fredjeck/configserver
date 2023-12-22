package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fredjeck/configserver/internal/repository"
	"strings"
	"time"

	"github.com/fredjeck/configserver/internal/encryption"
)

// JSONWebTokenPayload represents the Payload part of a JWT
type JSONWebTokenPayload struct {
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	Audience  []string `json:"aud"`
	NotBefore int64    `json:"nbf"`
	IssuedAt  int64    `json:"iat"`
	Expires   int64    `json:"exp"`
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

const TokenValidityDays = 1

// NewJSONWebToken creates a new empty JSON Web Token
func NewJSONWebToken() *JSONWebToken {
	return &JSONWebToken{
		Header: &JSONWebTokenHeader{Alg: "HS256", Type: "JWT"},
		Payload: &JSONWebTokenPayload{
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Expires:   time.Now().AddDate(0, 0, TokenValidityDays).Unix(),
		},
	}
}

// Pack generates the token by marshalling the content to JSON, formatting the output into b64UrlEncoded strings
// and by appending the token signature
func (jwt *JSONWebToken) Pack(secret *encryption.HmacSha256Secret) string {
	tk := jwt.token()
	hash := encryption.HmacSha256Hash([]byte(tk), secret)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s.%s", tk, b64Hash)
}

func (jwt *JSONWebToken) token() string {
	header, _ := jsonb64UrlEncode(jwt.Header)
	body, _ := jsonb64UrlEncode(jwt.Payload)
	return fmt.Sprintf("%s.%s", header, body)
}

func (jwt *JSONWebToken) ClientId() string {
	return jwt.Payload.Subject
}

func (jwt *JSONWebToken) IsAllowedRepository(mgr *repository.Manager, repository string) bool {
	found := false
	for _, aud := range jwt.Payload.Audience {
		if strings.EqualFold(aud, repository) {
			found = true
			break
		}
	}
	return found && mgr.IsClientAllowed(repository, jwt.Payload.Subject)
}

// ParseJwt parses a jwt token string, validates its signature and expiration and returns the token
func ParseJwt(token string, secret *encryption.HmacSha256Secret) (*JSONWebToken, error) {
	jwt, err := Unpack(token, secret)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if jwt.Payload.NotBefore > now || now > jwt.Payload.Expires {
		return nil, errors.New("token is expired")
	}

	return jwt, nil
}

// Unpack unmarshall a b64 encoded JWT to a JSONWebToken object
func Unpack(token string, secret *encryption.HmacSha256Secret) (*JSONWebToken, error) {
	if err := VerifySignature(token, secret); err != nil {
		return nil, err
	}
	components := strings.Split(token, ".")

	headerStr, err := base64.RawURLEncoding.DecodeString(components[0])
	if err != nil {
		return nil, err
	}

	tk := NewJSONWebToken()
	err = json.Unmarshal(headerStr, tk.Header)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal jwt header: %w", err)
	}

	bodyStr, err := base64.RawURLEncoding.DecodeString(components[1])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyStr, tk.Payload)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal jwt payload: %w", err)
	}

	return tk, nil
}

// VerifySignature validates a token has been signed with the provided key
// Simplistic approach, does not verify the alg and enforces HMAC
func VerifySignature(token string, secret *encryption.HmacSha256Secret) error {
	components := strings.Split(token, ".")
	if len(components) != 3 {
		return fmt.Errorf("malformed jwt token - expecting three components but found %d parts", len(components))
	}

	tk := components[0] + "." + components[1]
	signature := components[2]

	hash := encryption.HmacSha256Hash([]byte(tk), secret)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash)

	if strings.Compare(b64Hash, signature) != 0 {
		return fmt.Errorf("invalid token signature")
	}

	return nil
}

func jsonb64UrlEncode(e interface{}) (string, error) {
	jsonStr, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(jsonStr), nil
}
