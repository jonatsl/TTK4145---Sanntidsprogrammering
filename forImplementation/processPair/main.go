package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	go primary(backup())
}

var addr string = "localhost:8080"

func primary(counter int) {
	fmt.Printf("This is now a primary.")

	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	for {
		conn, err := net.Dial("udp", addr)
		if err != nil {
			fmt.Println("The following error occured", err)
		}
		conn.Write([]byte(""))
		time.Sleep(20*time.Millisecond)
	}
}

func backup() int {
	println("This is a process pair backup.")
	udpConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	buffer := make([]byte, 1024)
	for {
		err = udpConn.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		_, _, err := udpConn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return 0
		}
	}
}
