package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Handle flags
	filePtr := flag.String("f", "", "Relative path to file")
	peerPts := flag.String("p", "", "Peer address (IP:PORT)")
	flag.Parse()

	// Verify flags
	if *peerPts == "" {
		log.Fatal(*peerPts + " is not a valid peer.. Expecting IP:PORT")
	}
	host := string(*peerPts)

	// init receiver
	if *filePtr == "" {
		fmt.Println("Waiting for file transfers from " + *peerPts)
		initReciever(host)
	} else {
		fmt.Println("Trying to send file " + *filePtr + " to peer " + *peerPts)
		send(host, *filePtr)
	}
}

func initReciever(host string) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("Listening on " + host + "..")
	// Close listener
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// log incoming request
	fmt.Printf("Incomming file from: %s..\n", conn.RemoteAddr().String())

	// read file description sent by sender
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	// close connection
	defer conn.Close()

	// print file desc
	fileDesc := strings.Split(string(buf[:n]), ":")
	fileName := fileDesc[0]
	fileSize, _ := strconv.Atoi(fileDesc[1])
	fmt.Printf("File: %s, Size: %d bytes\n", fileName, fileSize)

	// promt for incomming file
	var choice string
	fmt.Printf("Accept file? [y/N] ")
	fmt.Scanf("%s", &choice)

	// response
	response := "1"
	if choice == "y" {
		response = "0"
	}

	// write response
	conn.Write([]byte(response))

	// get file
	if response == "0" {
		var chunkSize int64 = 1024
		var i int64 = 0
		fileBuf := make([]byte, fileSize)
		for i < int64(fileSize) {
			fmt.Printf("Recieving bytes (%d-%d/%d)..\n", i, i+chunkSize, fileSize)
			conn.Read(fileBuf)
			i += chunkSize
		}
		fmt.Println("Revieved file!")
	}
}

func send(host string, fileName string) {
	receiver, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatal("Could not resolve address: "+host, err)
	}

	conn, err := net.DialTCP("tcp", nil, receiver)
	if err != nil {
		log.Fatal("TCP Dial failed: ", err)
	}
	// close connection
	defer conn.Close()

	sendFile(fileName, conn)
}

func sendFile(fileName string, conn net.Conn) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Could not open file %s..\n", fileName)
		log.Fatal(err)
	}
	// def close file
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatal("Error loading file stats: ", err)
	}
	fileSize := stats.Size()

	descStr := fmt.Sprintf("%s:%d", fileName, fileSize)
	_, err = conn.Write([]byte(descStr))
	if err != nil {
		log.Fatal("Failed to write transfer description to connection: ", err)
	}
	fmt.Println("Waiting on reciever..")

	res := make([]byte, 1)
	_, err = conn.Read(res)
	if err != nil {
		log.Fatal("Could not read response: ", err)
	}
	fmt.Println(string(res))

	// send file if accepted
	if string(res) == "0" {
		var fileBuf []byte = make([]byte, fileSize)
		bn, err := file.Read(fileBuf)
		if err != nil {
			log.Fatal("Could not load file to buffer: ", err)
		}
		fmt.Printf("Loaded %d bytes to filebuffer..\n", bn)

		var chunkSize int64 = 1024
		var i int64 = 0
		var end int64 = 0
		for i < fileSize {
			end = i + chunkSize
			if end >= fileSize {
				end = fileSize
			}
			conn.Write(fileBuf[i:end])
			fmt.Printf("Sending bytes (%d-%d/%d)..\n", i, end, fileSize)
			i += chunkSize
		}
		fmt.Println("File transfer complete!")
	}
}
