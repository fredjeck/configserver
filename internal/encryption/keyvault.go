package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path"
)

// KeyVault is a wrapper around an rsa Private key
type KeyVault struct {
	PrivateKey *rsa.PrivateKey
}

// Generates a new keyvault - creates a new rsa Private key
func NewKeyVault() (*KeyVault, error) {
	private, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	return &KeyVault{private}, nil
}

// LoadKeyVault Searches the provided location for an id_rsa pem file.
// if the file is not found and if create is set to true, LoadKeyVault will generate a new keyvault and store it
func LoadKeyVault(location string, create bool) (*KeyVault, error) {
	prvKey := path.Join(location, "id_rsa")

	if _, err := os.Stat(prvKey); err == nil {
		return decodePem(prvKey)
	} else if errors.Is(err, os.ErrNotExist) && create {
		// No file found and creation is requested
		vault, err := NewKeyVault()
		if err != nil {
			return nil, err
		}
		err = vault.SaveTo(location)
		if err != nil {
			return nil, err
		}
		return vault, nil
	} else {
		// Well, something wrong is going on
		return nil, err
	}
}

// decodePem unmarshalls an RSA key from a PEM file
// Only unencrypted keys are supported
func decodePem(location string) (*KeyVault, error) {
	bytes, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(bytes)
	if pemBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("RSA private key is of the wrong type :%s", pemBlock.Type)
	}
	pemBytes := pemBlock.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(pemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(pemBytes); err != nil { // note this returns type `interface{}`
			return nil, fmt.Errorf("unable to parse RSA private key: %w", err)
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unable to parse RSA private key: %w", err)
	}
	return &KeyVault{privateKey}, nil
}

// SaveTo saves the KeyVault to the specified location, generating the id_rsa and id_rsa.pub keyfiles
func (kp *KeyVault) SaveTo(location string) error {
	privateKeyFile, err := os.Create(path.Join(location, "id_rsa"))
	defer func(privateKeyFile *os.File) {
		_ = privateKeyFile.Close()
	}(privateKeyFile)
	if err != nil {
		return err
	}

	publicKeyFile, err := os.Create(path.Join(location, "id_rsa.pub"))
	defer func(privateKeyFile *os.File) {
		_ = privateKeyFile.Close()
	}(privateKeyFile)
	if err != nil {
		return err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(kp.PrivateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&kp.PrivateKey.PublicKey)
	if err != nil {
		return err
	}

	publicKeyPEM := &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return err
	}
	return nil
}

// Encrypt encrypts the provided message with this KeyVault key
func (kp *KeyVault) Encrypt(message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		&kp.PrivateKey.PublicKey,
		message,
		nil)
}

// Decrypt decrypts the provided message with this KeyVault private key
func (kp *KeyVault) Decrypt(encryptedBytes []byte) ([]byte, error) {
	return kp.PrivateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA512})
}

// hash creates a message digest for its further signature
func hash(message []byte) ([]byte, error) {
	hasher := sha256.New()
	if _, err := hasher.Write(message); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// Sign signs the message
func (kp *KeyVault) Sign(message []byte) ([]byte, error) {
	messageDigest, err := hash(message)
	if err != nil {
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, kp.PrivateKey, crypto.SHA256, messageDigest)
}

// Verify verifies the provided signature is correct, if so return error is nil
func (kp *KeyVault) Verify(message []byte, signature []byte) error {
	messageDigest, err := hash(message)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(&kp.PrivateKey.PublicKey, crypto.SHA256, messageDigest, signature)
}
