package transfer

import (
	"log"
	"net"
)

func Advertise() {
	address := net.IPv4(255, 255, 255, 255)
	port := 49505
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   address,
		Port: port,
	})
	if err != nil {
		log.Fatal("Could not dial UDP: ", err.Error())
	}
	log.Printf("Successfully dialed %s\n", socket.RemoteAddr().String())
	defer socket.Close()

	greeting := []byte("Hello from client!")
	socket.Write(greeting)
}

func FindSender() {
	address := net.IPv4(255, 255, 255, 255)
	port := 49505
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   address,
		Port: port,
	})
	if err != nil {
		log.Fatal("Could not listen to UDP: ", err.Error())
	}
	log.Printf("Successfully established connections: \n")
	defer socket.Close()

	data := make([]byte, 1024)
	read, remote, err := socket.ReadFromUDP(data)
	if err != nil {
		log.Fatal("Could not read from UDP: ", err.Error())
	}
	log.Printf("Read %d bytes from address %s\n", read, remote.IP.String())
}

func handleHandshake(conn *net.Conn) {
	log.Println("Got a connection!")
}
