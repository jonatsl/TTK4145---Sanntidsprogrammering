package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func receive(receiverport string) {

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", receiverport)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Read from UDP listener in endless loop
	for {
		buf := make([]byte, 1024)
		_, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("> ", string(buf[0:]))

		// Write back the message over UPD
		conn.WriteToUDP(buf[0:], addr)
	}
}

func sender() {
	send_to_port := "255.255.255.255:20023"

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", send_to_port) // os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)

	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send a message to the server
	for {
		_, err = conn.Write([]byte("Hallo min broder\n"))
		fmt.Println("send...")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	port := ":30000"
	go receive(port)

	for {
	}

}
