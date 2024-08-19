package main

import (
	"Project/elevatordriver"
	"flag"
	"time"
)

func main() {
	myID := parseArgs()
	print("hei: ", myID)

	var (
		assignerToElevatorChannel = make(chan bool)
		elevatorToAssignerChannel = make(chan bool)
		elevatorLifelineChannel   = make(chan bool)
	)
	go elevatordriver.ElevatorDriver(
		assignerToElevatorChannel,
		elevatorToAssignerChannel,
		elevatorLifelineChannel,
	)
	// Sleep for a while to allow the goroutine to print the message
	time.Sleep(1 * time.Second)


}

func parseArgs() (nodeID int) {
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.Parse()
	return nodeID
}
