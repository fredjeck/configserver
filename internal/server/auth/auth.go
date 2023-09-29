// Package auth provides a simple client secret generation mechanism to identify remote clients
// This package is temporary only and shall be replaced by using somthing more straightforward like JWT
package auth

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fredjeck/configserver/internal/encrypt"
)

// ClientSpecElements holds the number of elements stored in a ClientSpec
const keySpecItemsCount = 2

// ClientSpecSeparator represents the separator used when serializing ClientSpec
const keySpecItemSeparatorChar = ":"

// ClientKeySpec represents the unencrypted state of a client key/token
type ClientKeySpec struct {
	Repositories []string // The list of repositories this client is allowed to access
	ClientID     string   // A unique client identifier
}

// GenerateSecret generates a client secret out of a client key specification.
// The provided key is used to encrypt information embedded into the client secret it is therefore important to use the same key throughout the whole execution context
// Client Secret Specification :
// Values are separated by colons (":')
// [0] - Bound repositories name
// [1] - Client ID
func (spec *ClientKeySpec) GenerateSecret(key *[32]byte) (string, error) {
	secret, err := encrypt.Encrypt([]byte(fmt.Sprintf("%s%s%s", strings.Join(spec.Repositories, "|"), keySpecItemSeparatorChar, spec.ClientID)), key)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secret), nil
}

// CanAccessRepository returns true whenever the provided repository name is included in the ClientKeySpec repositories
func (spec *ClientKeySpec) CanAccessRepository(repository string) bool {
	for _, r := range spec.Repositories {
		if strings.EqualFold(r, repository) {
			return true
		}
	}
	return false
}

// NewClientKeySpec generates a new client specification.
func NewClientKeySpec(clientID string, repositories []string) *ClientKeySpec {
	return &ClientKeySpec{Repositories: repositories, ClientID: clientID}
}

// ClientKeySpecFromSecret unmarshalls a ClientSpec out of a client secret
func ClientKeySpecFromSecret(clientSecret string, key *[32]byte) (*ClientKeySpec, error) {
	bytes, err := b64.StdEncoding.DecodeString(clientSecret)
	if err != nil {
		return nil, err
	}

	secret, err := encrypt.Decrypt(bytes, key)
	if err != nil {
		return nil, err
	}

	elements := strings.Split(string(secret), keySpecItemSeparatorChar)
	if len(elements) != keySpecItemsCount {
		return nil, ErrorMalformedClientSecret
	}

	return NewClientKeySpec(elements[1], strings.Split(elements[0], "|")), nil
}

var (
	// ErrorAuthRequired is thrown whenever auth header is missing from an inbound http request
	ErrorAuthRequired = errors.New("authentication required")
	// ErrorMissingCredentials is thrown when the auth header is malformed or when credentials are missing
	ErrorMissingCredentials = errors.New("missing credentials")
	// ErrorUnauthorized is thrown when a client is not authorized on a specific repository
	ErrorUnauthorized = errors.New("repository unauthorized")
	// ErrorMalformedClientSecret is thrown when the provided client secret cannot be decoded
	ErrorMalformedClientSecret = errors.New("malformed client secret")
)

// ClientKeySpecFromBasicAuth ensures basic auth is enabled on the inbound request and validates the ClientID and Client Secret
func ClientKeySpecFromBasicAuth(r http.Request, key *[32]byte) (*ClientKeySpec, error) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		return nil, ErrorAuthRequired
	}

	auth, err := b64.StdEncoding.DecodeString(strings.ReplaceAll(authorization, "Basic ", ""))
	if err != nil {
		return nil, ErrorMissingCredentials
	}

	credentials := strings.Split(string(auth), ":")
	if len(credentials) != 2 {
		return nil, ErrorMissingCredentials
	}

	spec, err := ClientKeySpecFromSecret(credentials[1], key)
	if err != nil {
		return nil, ErrorMalformedClientSecret
	}

	if credentials[0] == spec.ClientID {
		return spec, nil
	}

	return nil, ErrorUnauthorized
}
