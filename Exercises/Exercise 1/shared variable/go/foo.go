// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	//"time"
)

func numberserver(inchan chan int, dechan chan int, i_chan chan int) {
	var i = 0
	for {
		select {
		case <-inchan:
			i++
		case <-dechan:
			i--
		case i_chan <- i:
		}
	}
}

func incrementing(inchan chan int, quit chan int) {
	//TODO: increment i 1000000 times
	for j := 0; j < 1000000; j++ {
		inchan <- 1
	}
	quit <- 1
}

func decrementing(dechan chan int, quit chan int) {
	//TODO: decrement i 1000000 times
	for j := 0; j < 1000000; j++ {
		dechan <- 1
	}
	quit <- 1
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	// GOMAXPROCS determines the max number of simultaneous threads. If set to 1, the program will run sequentially.
	runtime.GOMAXPROCS(3)
	inchan := make(chan int)
	dechan := make(chan int)
	i_chan := make(chan int)
	quit := make(chan int)

	// TODO: Spawn both functions as goroutines
	go incrementing(inchan, quit)
	go decrementing(dechan, quit)
	go numberserver(inchan, dechan, i_chan)

	<-quit
	<-quit

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	Println("The magic number is:", <-i_chan)
}


