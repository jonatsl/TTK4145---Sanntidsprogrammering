package main

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"Heis/fsm"
	"Heis/timer"

	"Heis/network/bcast"
	"Heis/network/netfuncs"
	"Heis/network/peers"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//_________________________________________________________________________________________________

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

//_________________________________________________________________________________________________

func main() {

	//_________________________________________________________________________________________________

	//heisann, din gamle ørn36!
	processPairCh := make(chan elevator.Elevator)
	counter := backup()
	go primary(counter, processPairCh)

	_numFloors := elevio.NumFloors
	//_numButtons := elevio.NumButtons
	elevio.Init("localhost:15657", _numFloors)

	//Initialiserer en heisstruct
	//elev := elevator.InitElev()
	// elevUpdateBtnAndOrdersCh := make(chan elevator.Elevator)
	newOrderCh := make(chan map[string]elevator.Elevator)
	elevUpdateRealtimeCh := make(chan elevator.Elevator)
	doorTimerCh := make(chan bool)
	timedOut := make(chan int)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//networkchannels:
	//Peers
	eleviId := netfuncs.InitNet()
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15623, eleviId, peerTxEnable)
	go peers.Receiver(15623, peerUpdateCh)
	//broadcast
	ElevatorTx := make(chan elevator.Elevator)
	ElevatorRx := make(chan elevator.Elevator)
	go bcast.Transmitter(16523, ElevatorTx)
	go bcast.Receiver(16523, ElevatorRx)

	//go netfuncs.Bcast_message(ElevatorTx, elev, eleviId)
	//send cost func result to network
	mapOfElevsTx := make(chan map[string]elevator.Elevator)
	mapOfElevsRx := make(chan map[string]elevator.Elevator)
	go bcast.Transmitter(16524, mapOfElevsTx)
	go bcast.Receiver(16524, mapOfElevsRx)

	//Master slave
	isMaster := true

	fmt.Printf("Started!\n")
	go fsm.ButtonsAndRequests(eleviId, isMaster, elevUpdateRealtimeCh, drv_buttons, ElevatorTx, mapOfElevsTx, ElevatorRx, mapOfElevsRx, newOrderCh, processPairCh)
	go fsm.FloorObstrStop(isMaster, eleviId, elevUpdateRealtimeCh, drv_floors, drv_stop, drv_obstr, ElevatorTx, newOrderCh, doorTimerCh, timedOut)

	//go fsm.FSM(drv_buttons, drv_floors, drv_stop, drv_obstr, eleviId, isMaster, ElevatorTx, CostTx, masterOrders, ElevatorRx, CostRx)
	go netfuncs.Network_FSM(peerUpdateCh)
	go timer.Timer(doorTimerCh, timedOut)
	select {}
}

var addr string = "localhost:8080" // Process pair address

func primary(counter int, processPairCh chan elevator.Elevator) {
	//Gjør det primary skal.
	println("Primary: ")

	var msg string = ""

	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	for {
		select {
		case hallCalls := <-processPairCh:
			fmt.Println(hallCalls)
			println("printer noe bare for å se, men vil helst sende hall calls til backup")
		default:
		}
		conn, err := net.Dial("udp", addr)
		if err != nil {
			println("The following error occured", err)
		} else {
			fmt.Println("The connection was established to", conn.RemoteAddr())
		}
		msg = strconv.Itoa(counter)
		println("sender: ", msg)
		conn.Write([]byte(msg))
		counter += 1
		time.Sleep(1 * time.Second)
	}
}

func backup() int {
	counter := 0

	udpConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		println("Error:", err)
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
		println("number:", counter)
		time.Sleep(time.Millisecond * 100)

	}

}
