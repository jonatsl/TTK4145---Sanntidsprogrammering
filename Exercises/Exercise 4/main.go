package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var addr string = "localhost:8080"

func primary(counter int) {
	//Gjør det primary skal.
	fmt.Printf("Primary: ")
	// i := 0
	var msg string = ""

	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	for i := 0; i < 10; i++ {
		// fmt.Println(i)
		// hei, _ := net.Dial("udp", addr)
		// hei.Write([]byte(string(i)))
		// i += 1
		// time.Sleep(3 * time.Second)
		conn, err := net.Dial("udp", addr)
		if err != nil {
			fmt.Println("The following error occured", err)
		} else {
			fmt.Println("The connection was established to", conn.RemoteAddr())
		}
		msg = strconv.Itoa(counter)
		fmt.Println("sender: ", msg)
		conn.Write([]byte(msg))
		counter += 1
		time.Sleep(time.Second)

	}

}
func backup() int {
	counter := 0

	udpConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		// Terminate the program with a non-zero exit code.
		os.Exit(1)
	}
	defer udpConn.Close()

	buffer := make([]byte, 1024)
	for {
		err = udpConn.SetReadDeadline(time.Now().Add(time.Millisecond * 2000))
		if err != nil {
			fmt.Println("Error:", err)
			// Terminate the program with a non-zero exit code.
			os.Exit(1)
		}
		number, address, err := udpConn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			// Terminate the program with a non-zero exit code.
			// os.Exit(1)
			return counter
		}

		fmt.Println(" backup server: ", address)
		counter, _ = strconv.Atoi(string(buffer[:number]))
		fmt.Println(counter)
		time.Sleep(time.Millisecond * 100)

	}

}

func main() {
	fmt.Println("Program started.")

	counter := backup()
	primary(counter)

	fmt.Println("Prekæs")
	time.Sleep(1 * time.Second)

}
