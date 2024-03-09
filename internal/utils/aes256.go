package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

// deriveKeyFromPassPhrase derives an AES-256 compatible key from the provided password
func deriveKeyFromPassPhrase(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}

// AesEncrypt uses AES256GCP to encrypt the provided plainText string using the given passPhrase
func AesEncrypt(plainText string, passPhrase string) []byte {
	secretKey := deriveKeyFromPassPhrase(passPhrase)

	aesCipher, err := aes.NewCipher(secretKey[:])
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		panic(err)
	}

	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	return gcm.Seal(nonce, nonce, []byte(plainText), nil)
}

// AesDecrypt attempts to decrypt the provided bytes with the given passPhrase provided that it has been
// encrypted with AES256-GCM encryption
func AesDecrypt(cipherText []byte, passPhrase string) (string, error) {
	secretKey := deriveKeyFromPassPhrase(passPhrase)

	aesCipher, err := aes.NewCipher(secretKey[:])
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		panic(err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
