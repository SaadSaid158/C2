package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Static encryption key, should be different per build
var encryptionKey = []byte{0x12, 0x3F, 0xA7, 0xD9, 0x5B, 0x77, 0x43, 0x88, 0x62, 0x91, 0xC3, 0xF4, 0xAE, 0xB5, 0x69, 0x1E}

func encryptString(plaintext string) string {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
}

func decryptString(encrypted string) string {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return ""
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return ""
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	nonce, ciphertext := data[:12], data[12:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return ""
	}

	return string(plaintext)
}
