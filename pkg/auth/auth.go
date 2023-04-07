package auth

import (
	b64 "encoding/base64"
	"errors"
	"github.com/fredjeck/configserver/pkg/encrypt"
	"net/http"
	"strings"
)

type ClientSpec struct {
	Repository string
	ClientId   string
}

// ClientSpecElements holds the number of elements stored in a ClientSpec
const ClientSpecElements = 2

// ClientSpecSeparator represents the separator used when serializing ClientSpec
const ClientSpecSeparator = ":"

// ClientSecret generates a client secret out of a client specification.
// The provided key is used to encrypt information embedded into the client secret it is therefore important to use the same key throughout the whole execution context
// Client Secret Specification :
// Values are separated by colons (":')
// [0] - Bound repository name
// [1] - Client ID
func (spec ClientSpec) ClientSecret(key *[32]byte) (string, error) {
	secret, err := encrypt.Encrypt([]byte(spec.Repository+":"+spec.ClientId), key)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secret), nil
}

// NewClientSpec generates a new client specification.
func NewClientSpec(clientId string, repository string) *ClientSpec {
	return &ClientSpec{Repository: repository, ClientId: clientId}
}

// UnmarshalClientSpec unmarhaslls a ClientSpec out of a client secret
func UnmarshalClientSecret(clientSecret string, key *[32]byte) (*ClientSpec, error) {
	bytes, err := b64.StdEncoding.DecodeString(clientSecret)
	if err != nil {
		return nil, err
	}

	secret, decryptionErr := encrypt.Decrypt(bytes, key)
	if decryptionErr != nil {
		return nil, decryptionErr
	}

	elements := strings.Split(string(secret), ClientSpecSeparator)
	if len(elements) != ClientSpecElements {
		return nil, ErrMalformedClientSecret
	}

	return NewClientSpec(elements[1], elements[0]), nil
}

var (
	ErrAuthRequired          = errors.New("authentication required")
	ErrMissingCredentials    = errors.New("missing credentials")
	ErrUnauthorized          = errors.New("repository unauthorized")
	ErrMalformedClientSecret = errors.New("malformed client secret")
)

// FromBasicAuth ensures basic auth is enabled on the inbout request and validates the ClientID and Client Secret
func FromBasicAuth(r http.Request, key *[32]byte) (*ClientSpec, error) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		return nil, ErrAuthRequired
	}

	auth, err := b64.StdEncoding.DecodeString(strings.ReplaceAll(authorization, "Basic ", ""))
	if err != nil {
		return nil, ErrMissingCredentials
	}

	credentials := strings.Split(string(auth), ":")
	if len(credentials) != 2 {
		return nil, ErrMissingCredentials
	}

	spec, specErr := UnmarshalClientSecret(credentials[1], key)
	if specErr != nil {
		return nil, ErrMalformedClientSecret
	}

	if credentials[0] == spec.ClientId {
		return spec, nil
	}

	return nil, ErrUnauthorized
}
