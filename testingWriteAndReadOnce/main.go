package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)


func main() {

	initializeCabCalls := checkAndLoadCabCalls()

	println("Initial cab calls: ")
	fmt.Println(initializeCabCalls)
	println()

	processPairCh := make(chan []bool)

	println("Sending updated cab calls to channel and writing to file")
	updatedCabCalls := []bool{true, false, true, false}
	println()

	writeToLocalBackup(updatedCabCalls)

	time.Sleep(2*time.Second)

	processPairCh <- updatedCabCalls

	println("")
	fmt.Println(checkAndLoadCabCalls())

	select {}
}




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



func checkAndLoadCabCalls()  []bool{
	if _, err := os.Stat("localBackup.txt"); os.IsNotExist(err) {
		return []bool{false, false, false, false}
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
			} else if char == 'f'{
				cabCalls = append(cabCalls, false)
			}
		}
		return cabCalls
	}
}