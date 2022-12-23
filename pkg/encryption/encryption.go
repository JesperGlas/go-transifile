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
		log.Println("[encryption/EncryptData] Could not create cipher: ")
		log.Fatal(err.Error())
	}

	// generate GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("[encryption/EncryptData] Could not generate GCM: ")
		log.Fatal(err.Error())
	}

	// generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("[EncryptData] Could not generate nonce: ")
		log.Fatal(err.Error())
	}

	// encrypt data
	return gcm.Seal(nonce, nonce, *plainData, nil)
}

func DecryptData(passphrase string, cipherData *[]byte) []byte {
	key := genSHA256Key(passphrase)

	// generate cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		log.Println("[encryption/DecryptData] Could not create cipher: ")
		log.Fatal(err.Error())
	}

	// generate GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("[encryption/DecryptData] Could not generate GCM: ")
		log.Fatal(err.Error())
	}

	// generate nonce
	nonce := (*cipherData)[:gcm.NonceSize()]
	content := (*cipherData)[gcm.NonceSize():]
	plainData, err := gcm.Open(nil, nonce, content, nil)
	if err != nil {
		log.Println("[encrption/DecrytData] Could not decipher data: ")
		log.Fatal(err.Error())
	}
	return plainData
}

// ### PRIVATE FUNCTIONS ###

func genSHA256Key(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}
