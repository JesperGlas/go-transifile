package transfer

import (
	"log"
	"net"
)

func SendData(data *[]byte, reciever string) (int, error) {
	// resolve address
	addr, err := net.ResolveTCPAddr("tcp", reciever)
	if err != nil {
		log.Println("[transfer/SendData] Could not resolve TCP address of reciever: ")
		log.Fatal(err.Error())
	}

	// dial reciever
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println("[transfer/SendData] Could not establish TCP connection: ")
		log.Fatal(err.Error())
	}

	// send data
	return conn.Write(*data)
}

func RecieveData(payload *[]byte, sender string) (int, error) {
	// resolve address
	addr, err := net.ResolveTCPAddr("tcp", sender)
	if err != nil {
		log.Println("[transfer/RecieveData] Could not resolve TCP address of sender: ")
		log.Fatal(err.Error())
	}

	// listen for sender
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Printf("[transfer/RecieveData] Could not listen to %s\n", addr)
		log.Fatal(err.Error())
	}
	defer listener.Close()

	// accept connection
	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Println("[transfer/RecieveData] Refused connection: ")
		log.Print(err.Error())
	}

	// read data
	return conn.Read(*payload)
}
