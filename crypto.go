package main

import (
	"crypto"
	cryptorand "crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/chacha20poly1305"
)

// chacha20poly1305 package
// https://pkg.go.dev/golang.org/x/crypto@v0.16.0/chacha20poly1305

func Decrypt(data []byte, key string, associated_data string) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(hash_256(key))
	if err != nil {
		return nil, err
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := data[:aead.NonceSize()], data[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(associated_data))
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func Encrypt(data []byte, key string, associated_data string) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(hash_256(key))
	if err != nil {
		return nil, err
	}

	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(data)+aead.Overhead())
	if _, err := cryptorand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encrypted := aead.Seal(nonce, nonce, data, []byte(associated_data))
	return encrypted, nil
}

func hash_256(s string) []byte {
	h := crypto.SHA3_256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func Hash(s string) string {
	return hex.EncodeToString(hash_256(s))
}
