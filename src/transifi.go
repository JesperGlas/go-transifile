package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
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
		var data []byte = []byte(*filePtr)
		send(host, &data)
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
	fmt.Printf("Handling request from: %s..\n", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	// close connection
	defer conn.Close()

	// write response
	time := time.Now().Format(time.ANSIC)
	resStr := fmt.Sprintf("Received %d bytes at: %v", n, time)
	conn.Write([]byte(resStr))
}

func send(host string, data *[]byte) {
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

	wn, err := conn.Write([]byte(*data))
	if err != nil {
		log.Fatal("Write data faild: ", err)
	}
	fmt.Printf("Wrote %d bytes..\n", wn)

	res := make([]byte, 1024)
	rn, err := conn.Read(res)
	if err != nil {
		log.Fatal("Read data failed: ", err)
	}
	fmt.Printf("Read %d bytes: %s\n", rn, string(res))
}
