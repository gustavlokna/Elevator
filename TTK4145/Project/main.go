package main

import (
	. "Project/dataenums"
	"Project/elevatordriver"
	"Project/hwelevio"
	"Project/orderassigner"
	"flag"
	//"time"
)

func main() {
	nodeID := parseArgs()
	print("hei: ", nodeID)

	//INITILIZE DRIVER
	hwelevio.Init(Addr)

	// Ensure that hwelevio.Init() has completed successfully before continuing
	print("Initialization of hwelevio completed.")

	var (
		assignerToElevatorChannel = make(chan [NFloors][NButtons]bool, 10)
		elevatorToAssignerChannel = make(chan Elevator, 10)
		elevatorLifelineChannel   = make(chan bool)
	)
	print("hei jeg starter go")
	go elevatordriver.ElevatorDriver(
		assignerToElevatorChannel,
		elevatorToAssignerChannel,
		elevatorLifelineChannel,
		nodeID,
	)

	go orderassigner.OrderAssigner(
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
