package main

import (
	"crypto"
	"encoding/hex"
)

func Decrypt(data []byte, key string) ([]byte, error) {
	return data, nil
}

func Encrypt(data []byte, key string) ([]byte, error) {
	return data, nil
}

func Hash(s string) string {
	h := crypto.SHA3_256.New()
	h.Write([]byte(s))
	r := h.Sum(nil)
	return hex.EncodeToString(r)
}
