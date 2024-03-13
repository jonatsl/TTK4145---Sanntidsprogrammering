package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	// if len(os.Args) == 1 {
	// 	fmt.Println("Please provide host:port to connect to")
	// 	os.Exit(1)
	// }

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
	// // Read from the connection untill a new line is send
	// data, err := bufio.NewReader(conn).ReadString('\n')
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // Print the data read from the connection to the terminal
	// fmt.Print("> ", string(data))
}
