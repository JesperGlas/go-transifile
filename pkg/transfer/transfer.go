package transfer

import (
	"log"
	"net"
	"sync"
)

func Advertise() {
	conn, err := net.Dial("tcp", "192.168.0.255:49505")
	if err != nil {
		log.Fatal("Could not advertise service: ", err.Error())
	}
	log.Printf("Found reviever at: %s\n", conn.RemoteAddr().String())
}

func FindSender() {
	host := "192.168.0.255:49505"
	listener, err := net.Listen("tcp", ":49505")
	if err != nil {
		log.Fatal("Could not listen for services: ", err.Error())
	}
	log.Printf("Listening for senders on %s\n", host)
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Could not accept connection: ", err.Error())
	}
	log.Printf("Accepther connection\n")
	defer conn.Close()

	var wg sync.WaitGroup
	go func(conn *net.Conn, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		handleHandshake(conn)
	}(&conn, &wg)
}

func handleHandshake(conn *net.Conn) {
	log.Println("Got a connection!")
}
