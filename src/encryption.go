package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

const (
	KEYSTR string = "@McQfThWmZq4t7w!z%C*F-JaNdRgUkXn"
)

func encryptData(plainData *[]byte) []byte {
	key := []byte(KEYSTR)

	// generate cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("[Encryption Error] Could not create cipher: ", err.Error())
	}

	// generate GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("[Encryption Error] Could not generate GCM: ", err.Error())
	}

	// generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal("[Encryption Error] Could not generate nonce: ", err.Error())
	}

	// encrypt data
	return gcm.Seal(nonce, nonce, *plainData, nil)
}

func decryptData(cipherData *[]byte) []byte {
	key := []byte(KEYSTR)

	// generate cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("[Decryption Error] Could not create cipher: ", err.Error())
	}

	// generate GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("[Decryption Error] Could not generate GCM: ", err.Error())
	}

	// generate nonce
	nonce := (*cipherData)[:gcm.NonceSize()]
	content := (*cipherData)[gcm.NonceSize():]
	plainData, err := gcm.Open(nil, nonce, content, nil)
	if err != nil {
		log.Fatal("[Decrytion Error] Could not decipher data: ", err.Error())
	}
	return plainData
}
