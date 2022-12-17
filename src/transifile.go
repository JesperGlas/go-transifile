package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Using TransiFi version 0.1.0..")

	// Parse flags
	filePtr := flag.String("file", "", "The name of the file (Required).")
	batchSizePtr := flag.Int("batch", 1024, "Transfer batch size (Default 1024)")

	flag.Parse()

	fmt.Printf("File: %s\n", *filePtr)
	fmt.Printf("Batch: %d\n", *batchSizePtr)
}
