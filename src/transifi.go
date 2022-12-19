package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func loadFile(fileName string) ([]byte, error) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("[Load File Error] Could not determine executable directory: ", err.Error())
	}
	path := filepath.Dir(exe) + "/" + fileName
	fmt.Printf("Reading file from: %s\n", path)
	return os.ReadFile(path)
}

func writeFile(fileName string, data *[]byte) error {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("[Write File Error] Could not determine executable directory: ", err.Error())
	}
	path := filepath.Dir(exe) + "/" + fileName
	fmt.Printf("Writing to: %s\n", path)
	return os.WriteFile(path, *data, 0644)
}

func main() {
	fileName := "test.txt"
	fmt.Println("TransiFile")

	// Load file in to byte-array
	data, err := loadFile(fileName)
	if err != nil {
		log.Fatal("[Main Error] Could not load data from file: ", err.Error())
	}
	fmt.Printf("Loaded %d bytes from file %s\n", len(data), fileName)

	// print incoming data
	fmt.Printf("Plain data: %s\n", data)

	// encrypt, print and output to file
	cipher := encryptData(&data)
	fmt.Printf("Encrypted data: %s\n", cipher)
	err = writeFile("out_"+fileName, &cipher)
	if err != nil {
		log.Fatal("[Main Error] Could not output cipher: ", err.Error())
	}

	plain := decryptData(&cipher)
	fmt.Printf("Decrypted data: %s\n", plain)

}
