package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	println("NTR")
}

// Can be run from anywhere. Needs either cab calls, elev or similar as input
func writeToLocalBackup(cabCalls []bool) error {
	file, err := os.Create("localBackup.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	var boolStr strings.Builder
	for _, value := range cabCalls {
		if value {
			boolStr.WriteString("true")
		} else {
			boolStr.WriteString("false")
		}
	}

	_, err = file.WriteString(boolStr.String())
	if err != nil {
		return err
	}
	return nil
}

// Function to read cab calls from localBackup.txt, if the file exists. If not it returns an all false list of length
// Place this function in init_elev.go
// Intended to only be run _once_ inside the initElev() function
func checkAndLoadCabCalls() []bool {
	if _, err := os.Stat("localBackup.txt"); os.IsNotExist(err) {
		cabCalls := make([]bool, elevio.NumFloors)
		for i := range cabCalls {
			cabCalls[i] = false
		}
		return cabCalls
	} else {
		file, err := os.ReadFile("localBackup.txt")
		if err != nil {
			fmt.Println(err)
		}
		var cabCalls []bool
		fileContent := strings.TrimSpace(string(file))
		for _, char := range fileContent {
			if char == 't' {
				cabCalls = append(cabCalls, true)
			} else if char == 'f' {
				cabCalls = append(cabCalls, false)
			}
		}
		return cabCalls
	}
}

// Function to merge requests and initializeCabCalls
// Place this function in init_elev.go
// Put this function in init_elev() by setting Requests: initializeRequests
func initializeRequests() [elevio.NumFloors][elevio.NumButtons]bool {
	initializeCabCalls := checkAndLoadCabCalls()
	requests := [elevio.NumFloors][elevio.NumButtons]bool{}

	for floor := range requests {
		requests[floor][0] = false
		requests[floor][1] = false
		requests[floor][2] = initializeCabCalls[floor]
	}
	return requests
}

/*
func initializeRequests() [][]bool {
	initializeCabCalls := checkAndLoadCabCalls()
	requests := make([][]bool, elevio.NumFloors)

	for floor := range requests {
		requests[floor] = make([]bool, Elevio.numButtons)
		requests[floor][0] = false
		requests[floor][1] = false
		requests[floor][2] = initializeCabCalls[floor]
	}
	return requests
}
*/
