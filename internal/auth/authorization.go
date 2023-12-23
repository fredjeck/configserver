package auth

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/fredjeck/configserver/internal/repository"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"
)

const ClientSecretSeparatorChar = "|"            // Char used to separate client secret component
const ClientSecretComponents = 2                 // Number of components used in client secrets
const ClientSecretValidity = time.Hour * 24 * 30 // 30 Days

// GenerateClientSecret creates a client secret for the provided client id using given encryption key
// This is purely experimental - goal is to generate a self-contained secret which can easily be validated
// and which does need to be stored locally
func GenerateClientSecret(clientId string, vault *encryption.KeyVault) (string, error) {
	secret, err := vault.Encrypt([]byte(fmt.Sprintf("%s%s%s", time.Now().Format(time.RFC3339), ClientSecretSeparatorChar, clientId)))
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secret), nil
}

// ValidateClientSecret ensure the clientId and clientSecret pairs are matching
func ValidateClientSecret(clientId string, clientSecret string, vault *encryption.KeyVault) bool {
	bytes, err := b64.StdEncoding.DecodeString(clientSecret)
	if err != nil {
		return false
	}

	secret, err := vault.Decrypt(bytes)
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

type AuthorizationKind string

var AuthorizationKindBasic = AuthorizationKind("basic")
var AuthorizationKindBearer = AuthorizationKind("bearer")
var AuthorizationKindNone = AuthorizationKind("none")

// Authorization interface allows support for various authorization kinds
type Authorization interface {
	IsAllowedRepository(mgr *repository.Manager, repository string) bool
	ClientId() string
}

// FromRequest generates the Authorization from the provided http request
func FromRequest(r *http.Request, vault *encryption.KeyVault, allowedMethods ...AuthorizationKind) (Authorization, error) {
	authorization := r.Header.Get("Authorization")

	// No auth enabled, we skip
	if slices.Contains(allowedMethods, AuthorizationKindNone) {
		return &NoneAuth{}, nil
	}

	if len(authorization) == 0 {
		return nil, errors.New("missing authorization header")
	}

	authComponents := strings.Split(authorization, " ")
	if len(authComponents) != 2 {
		return nil, errors.New("invalid authorization header")
	}

	authKind := AuthorizationKind(strings.ToLower(authComponents[0]))

	if !slices.Contains(allowedMethods, authKind) {
		return nil, fmt.Errorf("'%s' unsupported authorization method", authKind)
	}

	switch authKind {
	case AuthorizationKindBearer:
		return performJWTAuthorization(authComponents[1], vault)
	case AuthorizationKindBasic:
		return performBasicAuthorization(authComponents[1], vault)
	default:
		return nil, fmt.Errorf("'%s' unsupported authorization method", authKind)
	}
}

func performJWTAuthorization(credentials string, vault *encryption.KeyVault) (Authorization, error) {
	err := VerifySignature(credentials, vault)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	jwt, err := ParseJwt(credentials, vault)
	if err != nil {
		return nil, errors.New("json token cannot be parsed")
	}

	return jwt, nil
}

func performBasicAuthorization(credentials string, vault *encryption.KeyVault) (Authorization, error) {
	basicAuth, err := b64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return nil, errors.New("invalid authorization header")
	}

	loginPwd := strings.Split(string(basicAuth), ":")
	if len(loginPwd) != 2 {
		return nil, errors.New("invalid authorization header")
	}

	if !ValidateClientSecret(loginPwd[0], loginPwd[1], vault) {
		return nil, errors.New("unauthorized")
	}
	return &BasicAuth{clientId: loginPwd[0]}, nil
}

// BasicAuth represent a basic client:password Authorization
type BasicAuth struct {
	clientId string
}

func (a *BasicAuth) IsAllowedRepository(mgr *repository.Manager, repository string) bool {
	return mgr.IsClientAllowed(repository, a.clientId)
}

func (a *BasicAuth) ClientId() string {
	return a.clientId
}

// NoneAuth allows for disabling authorization
type NoneAuth struct {
	clientId string
}

func (a *NoneAuth) IsAllowedRepository(_ *repository.Manager, _ string) bool {
	return true
}

func (a *NoneAuth) ClientId() string {
	return "None"
}
