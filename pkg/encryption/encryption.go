package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
)

// ### PUBLIC FUNCTIONS ###

func EncryptData(passphrase string, plainData *[]byte) []byte {
	key := genSHA256Key(passphrase)

	// generate cipher block
	block, err := aes.NewCipher(key[:])
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

func DecryptData(passphrase string, cipherData *[]byte) []byte {
	key := genSHA256Key(passphrase)

	// generate cipher block
	block, err := aes.NewCipher(key[:])
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

// ### PRIVATE FUNCTIONS ###

func genSHA256Key(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}
