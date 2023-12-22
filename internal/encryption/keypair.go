package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path"
)

type KeyPair struct {
	PublicKey  rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewKeyPair() (*KeyPair, error) {
	kp := &KeyPair{}

	private, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	kp.PrivateKey = private
	kp.PublicKey = private.PublicKey

	return kp, nil
}

func FromLocation(location string) (*KeyPair, error) {
	bytes, err := os.ReadFile(path.Join(location, "id_rsa"))
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("invalid id_rsa private key")
	}
	return nil, nil
}

func (kp *KeyPair) StoreToLocation(location string) error {
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
