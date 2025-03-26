package main

import (
	"Project/assigner"
	. "Project/config"
	. "Project/dataenums"
	"Project/elevatorDriver"
	"Project/hwelevio"
	"Project/lights"
	"Project/network"
	"flag"
)

func main() {

	nodeID := parseArgs()

	hwelevio.Init(Addr)

	var (
		newOrders      = make(chan [NFloors][NButtons]bool, ChannelBufferSize)
		driverEvents   = make(chan FromDriverToAssigner, ChannelBufferSize)
		worldview      = make(chan FromAssignerToNetwork, ChannelBufferSize)
		stateBroadcast = make(chan FromNetworkToAssigner, ChannelBufferSize)
		elevatorLights = make(chan FromDriverToLight, ChannelBufferSize)
		orderLights    = make(chan [NFloors][NButtons]ButtonState, ChannelBufferSize)
	)

	go assigner.Assigner(
		newOrders,
		driverEvents,
		worldview,
		stateBroadcast,
		orderLights,
		nodeID,
	)

	go elevatorDriver.ElevatorDriver(
		newOrders,
		driverEvents,
		elevatorLights,
	)

	go network.Network(
		worldview,
		stateBroadcast,
		nodeID,
	)

	go lights.Lights(
		orderLights,
		elevatorLights,
	)
	select {}
}

func parseArgs() int {
	var nodeID int
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.Parse()
	return nodeID
}
