package main

import (
	"Project/elevatordriver"
	. "Project/dataenums"
	"Project/hwelevio"
	"flag"
	//"time"
)

func main() {
	nodeID := parseArgs()
	print("hei: ", nodeID)

	//INITILIZE DRIVER 
	hwelevio.Init(Addr)

	var (
		assignerToElevatorChannel = make(chan bool)
		elevatorToAssignerChannel = make(chan bool)
		elevatorLifelineChannel   = make(chan bool)
	)
	go elevatordriver.ElevatorDriver(
		assignerToElevatorChannel,
		elevatorToAssignerChannel,
		elevatorLifelineChannel,
		nodeID,
	)
	// Sleep for a while to allow the goroutine to print the message
		// Hold main function indefinitely
	select {}
	//time.Sleep(1 * time.Second)


}

func parseArgs() (nodeID int) {
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.Parse()
	return nodeID
}
