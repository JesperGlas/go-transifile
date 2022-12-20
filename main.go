package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JesperGlas/go-transifile/pkg/encryption"
	"github.com/JesperGlas/go-transifile/pkg/transfer"
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

func encryptFile(fileName string) []byte {
	// Load file in to byte-array
	data, err := loadFile(fileName)
	if err != nil {
		log.Fatal("[Encrypt File Error] Could not load data from file: ", err.Error())
	}
	fmt.Printf("Loaded %d bytes from file %s\n", len(data), fileName)

	return encryption.EncryptData(&data)
}

func decryptPayload(outFileName string, cipherData *[]byte) {
	cipher := encryption.DecryptData(cipherData)
	err := writeFile(outFileName, &cipher)
	if err != nil {
		log.Fatal("[Decrypt Payload Error] Could not output cipher: ", err.Error())
	}
}

func main() {
	fmt.Println("TransiFile")

	modePtr := flag.String("m", "", "Advertise or Send [a|s]")
	flag.Parse()

	if *modePtr == "a" {
		transfer.Advertise()
	} else if *modePtr == "s" {
		sender := transfer.FindSender()
		if sender == "" {
			log.Fatalln("[Main Error] Could not find sender!")
		}
	} else {
		log.Fatal("Flag must be specified! [-m=a | -m=s]")
	}
}
