package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type crypt struct {
	key string
}

func newCrypt(key string) crypt {
	return crypt{key}
}

var iv []byte = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func (c *crypt) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return "", fmt.Errorf("could net encrypt %s: %v", text, err)
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext), nil
}

func (c *crypt) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return "", fmt.Errorf("could not decrypt %s: %v", text, err)
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}
