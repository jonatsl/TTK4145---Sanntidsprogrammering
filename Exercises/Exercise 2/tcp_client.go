package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	port := "10.100.23.129:34933"

	// Resolve the string address to a TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connect to the address with tcp
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	//wbyte := make([]byte, 1024)

	for {
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))

		// Send a message to the server
		_, err = conn.Write(append([]byte("Hello TCP Server\n"), 0))
		fmt.Println("send...")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("etter")
		time.Sleep(1 * time.Second)
	}

	// 	// Read from the connection untill a new line is send
	// 	data, err := bufio.NewReader(conn).ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	// Print the data read from the connection to the terminal
	// 	fmt.Print("> ", string(data))
}
