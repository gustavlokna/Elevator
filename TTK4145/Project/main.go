package main

import (
	. "Project/dataenums"
	"Project/elevatordriver"
	"Project/hwelevio"
	"Project/lights"
	"Project/network"
	"Project/orderassigner"
	"flag"
	"strconv"
)

func main() {

	nodeID := parseArgs()

	hwelevio.Init(Addr)

	var (
		newOrderChannel     = make(chan [NFloors][NButtons]bool, 100)
		payloadFromElevator = make(chan PayloadFromElevator, 100)
		toNetworkChannel    = make(chan PayloadFromassignerToNetwork, 100)
		fromNetworkChannel  = make(chan PayloadFromNetworkToAssigner, 100)
		fromDriverToLight   = make(chan PayloadFromDriver, 100)
		fromAsstoLight      = make(chan [NFloors][NButtons]ButtonState, 100)
	)

	go elevatordriver.ElevatorDriver(
		newOrderChannel,
		payloadFromElevator,
		fromDriverToLight,
		//nodeID,
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

	go lights.LightsHandler(
		fromAsstoLight,
		fromDriverToLight,
	)
	// TODO is the select needed ?
	select {}
}

func parseArgs() string {
	var nodeID int
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.Parse()
	return strconv.Itoa(nodeID)
}
