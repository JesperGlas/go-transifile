package discovery

import (
	"bytes"
	"log"
	"net"
	"strconv"
)

const IDENTIFIER string = "B&E)H@McQfTjWmZq4t7w!z%C*F-JaNdR"

func Advertise() string {
	address := net.IPv4(192, 168, 0, 255)
	port := 49505
	broadcast := address.String() + ":" + strconv.Itoa(port)
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   address,
		Port: port,
	})
	if err != nil {
		log.Printf("[discovery/Advertise] Could not dial broadcast address: %s\n", broadcast)
		log.Fatal(err.Error())
	}
	log.Printf("Successfully broadcasted on %s\n", socket.RemoteAddr().String())
	socket.Write([]byte(IDENTIFIER))
	socket.Close()

	log.Println("Awaiting handshake credentials from receiver..")
	listen, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   address,
		Port: port,
	})
	if err != nil {
		log.Println("[discovery/Advertise] Could not listen for credentials!")
		log.Fatal(err.Error())
	}
	payload := make([]byte, len(IDENTIFIER))
	n, remote, err := listen.ReadFromUDP(payload)
	if err != nil {
		log.Println("[discovery/Advertise] Could not read UDP payload!")
		log.Fatal(err.Error())
	}
	reciever := remote.IP.String() + ":" + strconv.Itoa(remote.Port)
	log.Printf("Handshake attempt from %s (%d bytes)\n", reciever, n)
	listen.Close()

	if bytes.EqualFold(payload, []byte(IDENTIFIER)) {
		log.Println("Handshake verified!")
		return reciever
	}
	return ""
}

func FindSender() string {
	address := net.IPv4(192, 168, 0, 255)
	port := 49505
	broadcast := address.String() + ":" + strconv.Itoa(port)
	listen, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   address,
		Port: port,
	})
	if err != nil {
		log.Printf("[discovery/FindSender] Could not find sender at %s\n", broadcast)
		log.Fatal(err.Error())
	}
	log.Printf("Successfully established connections: \n")

	payload := make([]byte, len(IDENTIFIER))
	n, remote, err := listen.ReadFromUDP(payload)
	if err != nil {
		log.Println("[discovery/FindSender] Could not read UDP payload!")
		log.Fatal(err.Error())
	}
	sender := remote.IP.String() + ":" + strconv.Itoa(remote.Port)
	log.Printf("Connection attempt by %s (%d bytes)\n", sender, n)
	listen.Close()

	if bytes.EqualFold(payload, []byte(IDENTIFIER)) {
		log.Println("Connection verified! Sending handshake!")
		socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
			IP:   address,
			Port: port,
		})
		if err != nil {
			log.Println("[discovery/FindSender] Could not send credentials to sender!")
			log.Fatal(err.Error())
		}
		socket.Write([]byte(IDENTIFIER))
		socket.Close()
		return sender
	}
	return ""
}
