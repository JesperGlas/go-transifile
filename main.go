package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JesperGlas/go-transifile/pkg/encryption"
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

func encryptFile(fileName string, passphrase string) []byte {
	// Load file in to byte-array
	data, err := loadFile(fileName)
	if err != nil {
		log.Fatal("[Encrypt File Error] Could not load data from file: ", err.Error())
	}
	fmt.Printf("Loaded %d bytes from file %s\n", len(data), fileName)

	return encryption.EncryptData(passphrase, &data)
}

func decryptPayload(passphrase string, outFileName string, cipherData *[]byte) {
	cipher := encryption.DecryptData(passphrase, cipherData)
	err := writeFile(outFileName, &cipher)
	if err != nil {
		log.Fatal("[Decrypt Payload Error] Could not output cipher: ", err.Error())
	}
}

func main() {
	fmt.Println("TransiFile")

	modePtr := flag.String("m", "", "Specify mode ([e]ncrypt | [d]ecrypt)")
	filePtr := flag.String("f", "", "File name")
	passphrasePtr := flag.String("p", "", "Password for encryption")
	flag.Parse()

	if *modePtr == "" {
		log.Fatal("Please specify mode! (-m [e]ncrypt | [d]ecrypt)")
	}
	if *filePtr == "" {
		log.Fatal("Please specify file (-f <filename>)")
	}
	if *passphrasePtr == "" || len(*passphrasePtr) < 10 {
		log.Fatal("Please provide a passphrase with atleast 10 chars (-p <passphrase>)")
	}
}
