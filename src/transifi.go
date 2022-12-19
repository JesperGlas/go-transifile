package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	KEYSTR string = "@McQfThWmZq4t7w!z%C*F-JaNdRgUkXn"
)

func loadFile(fileName string) ([]byte, error) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Could not determine executable directory: ", err.Error())
	}
	exeDir := filepath.Dir(exe)
	return os.ReadFile(exeDir + "/" + fileName)
}

func encryptData(plainData *[]byte) []byte {
	key := []byte(KEYSTR)

	// generate cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("Could not create cipher block: ", err.Error())
	}

	// generate gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("Could not generate GCM: ", err.Error())
	}

	// generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal("Could not generate nonce: ", err.Error())
	}

	// encrypt data
	return gcm.Seal(nonce, nonce, *plainData, nil)
}

func main() {
	fileName := "test.txt"
	fmt.Println("TransiFile")

	// Load file in to byte-array
	data, err := loadFile(fileName)
	if err != nil {
		log.Fatal("Could not load data from file: ", err.Error())
	}
	fmt.Printf("Loaded %d bytes from file %s\n", len(data), fileName)

	fmt.Printf("Plain data: %s\n", data)
	cipher := encryptData(&data)
	fmt.Printf("Encrypted data: %s\n", cipher)
}
