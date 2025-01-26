package main

import (
	. "Project/dataenums"
	"Project/elevatordriver"
	"Project/hwelevio"
	"Project/network"
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
		newOrderChannel    = make(chan [NFloors][NButtons]bool, 100)
		payloadFromElevator     = make(chan PayloadFromElevator, 100)
		toNetworkChannel   = make(chan PayloadFromassignerToNetwork, 100)
		fromNetworkChannel = make(chan PayloadFromNetworkToAssigner, 100)
		fromDriverToLight = make(chan PayloadFromDriver, 100)
		fromAsstoLight = make(chan [NFloors][NButtons]ButtonState, 100)
	)

	//todo set ip as id in main? 
	go elevatordriver.ElevatorDriver(
		newOrderChannel,
		payloadFromElevator,
		fromDriverToLight,
		nodeID,
	)

	go orderassigner.OrderAssigner(
		newOrderChannel,
		payloadFromElevator,
		toNetworkChannel,
		fromNetworkChannel,
		fromAsstoLight,
		nodeID,
	)

	go network.Network(
		toNetworkChannel,
		fromNetworkChannel,
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