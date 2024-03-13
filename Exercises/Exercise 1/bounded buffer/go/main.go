package main

import (
	"fmt"
	"time"
)

func producer(buffer chan int, size int) {

	for i := 0; i < size; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[producer]: pushing %d\n", i)
		// TODO: push real value to buffer
		buffer <- i

	}

}

func consumer(buffer chan int) {

	time.Sleep(1 * time.Second)
	for {
		i := 0 //TODO: get real value from buffer
		i = <-buffer
		fmt.Printf("[consumer]: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}

}

func main() {

	// TODO: make a bounded buffer
	size := 10
	boundbuf := make(chan int, size)

	go consumer(boundbuf)
	go producer(boundbuf, size)

	select {}
}
