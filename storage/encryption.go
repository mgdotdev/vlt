package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

var nonce []byte = []byte("super_secret")

func byteKey(key string) []byte {
	sha := sha256.New()
	sha.Write([]byte(key))
	_key := []byte(hex.EncodeToString(sha.Sum(nil)))
	result := make([]byte, 32)
	for i := 0; i < len(result); i++ {
		result[i] = _key[2*i] ^ _key[2*i+1]
	}
	return result
}

func encrypt(text string, key string) string {
	_key := byteKey(key)
	plaintext := []byte(text)
	block, err := aes.NewCipher(_key)
	check(err)
	aesgcm, err := cipher.NewGCM(block)
	check(err)
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	result := hex.EncodeToString(ciphertext)
	return result
}

func decrypt(text string, key string) string {
	_key := byteKey(key)
	ciphertext, err := hex.DecodeString(text)
	check(err)
	block, err := aes.NewCipher(_key)
	check(err)
	aesgcm, err := cipher.NewGCM(block)
	check(err)
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	check(err)
	return string(plaintext)
}
