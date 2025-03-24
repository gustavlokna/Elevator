package main

import (
	"Project/assigner"
	. "Project/config"
	. "Project/dataenums"
	"Project/driver"
	"Project/hwelevio"
	"Project/lights"
	"Project/network"
	"flag"
	"strconv"
)

func main() {

	nodeID := parseArgs()

	hwelevio.Init(Addr)

	var (
		newOrders      = make(chan [NFloors][NButtons]bool, ChannelBufferSize)
		driverEvents   = make(chan FromDriverToAssigner, ChannelBufferSize)
		worldview      = make(chan FromAssignerToNetwork, ChannelBufferSize)
		stateBroadcast = make(chan FromNetworkToAssigner, ChannelBufferSize)
		localLights    = make(chan FromDriverToLight, ChannelBufferSize)
		sharedLights   = make(chan [NFloors][NButtons]ButtonState, ChannelBufferSize)
	)

	go assigner.Assigner(
		newOrders,
		driverEvents,
		worldview,
		stateBroadcast,
		sharedLights,
		nodeID,
	)

	go driver.Driver(
		newOrders,
		driverEvents,
		localLights,
	)

	go network.Network(
		worldview,
		stateBroadcast,
		nodeID,
	)

	go lights.Lights(
		sharedLights,
		localLights,
	)
	select {}
}

func parseArgs() string {
	var nodeID int
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.Parse()
	return strconv.Itoa(nodeID)
}
