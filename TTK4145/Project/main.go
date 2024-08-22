package main

import (
	. "Project/dataenums"
	"Project/elevatordriver"
	"Project/hwelevio"
	"Project/orderassigner"
	"flag"
	"strconv"
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
		newOrderChannel = make(chan [NFloors][NButtons]bool, 10)
		newStateChanel = make(chan Elevator, 10)
		orderDoneChannel   = make(chan [NFloors][NButtons]bool,10)
		toNetworkChannel = make(chan HRAInput, 10)
	)
	go elevatordriver.ElevatorDriver(
		newOrderChannel,
		newStateChanel,
		orderDoneChannel,
		nodeID,
	)

	go orderassigner.OrderAssigner(
		newOrderChannel,
		newStateChanel,
		orderDoneChannel,
		toNetworkChannel,
		nodeID,
	)
	// Sleep for a while to allow the goroutine to print the message
	// Hold main function indefinitely
	select {}
	//time.Sleep(1 * time.Second)

}


func parseArgs() string {
	var nodeID int
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.Parse()
	return strconv.Itoa(nodeID)
}