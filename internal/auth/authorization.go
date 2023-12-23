package auth

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/fredjeck/configserver/internal/encryption"
	"github.com/fredjeck/configserver/internal/repository"
	"net/http"
	"slices"
	"strings"
	"time"
)

const ClientSecretSeparatorChar = "|"            // Char used to separate client secret component
const ClientSecretComponents = 3                 // Number of components used in client secrets
const ClientSecretValidity = time.Hour * 24 * 30 // 30 Days

// GenerateClientSecret creates a client secret for the provided client id using given encryption key
// This is purely experimental - goal is to generate a self-contained secret which can easily be validated
// and which does need to be stored locally
func GenerateClientSecret(clientId string, vault *encryption.KeyVault) (string, error) {
	secret, err := vault.Encrypt([]byte(clientId))
	if err != nil {
		return "", err
	}

	h, err := vault.Hash(secret)
	if err != nil {
		return "", err
	}

	//return b64.StdEncoding.EncodeToString(secret), nil
	return b64.StdEncoding.EncodeToString(h), nil
}

// ValidateClientSecret ensure the clientId and clientSecret pairs are matching
func ValidateClientSecret(clientId string, clientSecret string, vault *encryption.KeyVault) bool {
	secret, err := GenerateClientSecret(clientId, vault)
	if err != nil {
		return false
	}

	return secret == clientSecret
}

type AuthorizationKind string

var AuthorizationKindBasic = AuthorizationKind("basic")
var AuthorizationKindBearer = AuthorizationKind("bearer")
var AuthorizationKindNone = AuthorizationKind("none")

type Authorization interface {
	IsAllowedRepository(mgr *repository.Manager, repository string) bool
	ClientId() string
}

func FromRequest(r *http.Request, key *encryption.Keystore, allowedMethods ...AuthorizationKind) (Authorization, error) {
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
		return performJWTAuthorization(authComponents[1], key)
	case AuthorizationKindBasic:
		return performBasicAuthorization(authComponents[1], key)
	default:
		return nil, fmt.Errorf("'%s' unsupported authorization method", authKind)
	}
}

func performJWTAuthorization(credentials string, key *encryption.Keystore) (Authorization, error) {
	err := VerifySignature(credentials, key.HmacSha256Secret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	jwt, err := ParseJwt(credentials, key.HmacSha256Secret)
	if err != nil {
		return nil, errors.New("json token cannot be parsed")
	}

	return jwt, nil
}

func performBasicAuthorization(credentials string, key *encryption.Keystore) (Authorization, error) {
	basicAuth, err := b64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return nil, errors.New("invalid authorization header")
	}

	loginPwd := strings.Split(string(basicAuth), ":")
	if len(loginPwd) != 2 {
		return nil, errors.New("invalid authorization header")
	}

	if !ValidateClientSecret(loginPwd[0], loginPwd[1], key.Aes256Key) {
		return nil, errors.New("unauthorized")
	}
	return &BasicAuth{clientId: loginPwd[0]}, nil
}

type BasicAuth struct {
	clientId string
}

func (a *BasicAuth) IsAllowedRepository(mgr *repository.Manager, repository string) bool {
	return mgr.IsClientAllowed(repository, a.clientId)
}

func (a *BasicAuth) ClientId() string {
	return a.clientId
}

type NoneAuth struct {
	clientId string
}

func (a *NoneAuth) IsAllowedRepository(_ *repository.Manager, _ string) bool {
	return true
}

func (a *NoneAuth) ClientId() string {
	return "None"
}
