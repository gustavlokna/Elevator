package main

import (
	"Project/assigner"
	. "Project/dataenums"
	"Project/driver"
	"Project/hwelevio"
	"Project/lights"
	"Project/network"
	"Project/config"
	"flag"
	"strconv"
)

func main() {

	nodeID := parseArgs()

	hwelevio.Init(config.Addr)

	var (
		newOrders       = make(chan [config.NFloors][config.NButtons]bool, config.BufferSize)
		driverEvents  = make(chan FromDriverToAssigner, config.BufferSize)
		worldview = make(chan FromAssignerToNetwork, config.BufferSize)
		stateBroadcast = make(chan FromNetworkToAssigner, config.BufferSize)
		localLights     = make(chan FromDriverToLight, config.BufferSize)
		sharedLights   = make(chan [config.NFloors][config.NButtons]ButtonState, config.BufferSize)
	)

	go assigner.Assigner(
		newOrders,
		driverEvents,
		worldview,
		stateBroadcast,
		sharedLights,
		nodeID,
	)

	go driver.ElevatorDriver(
		newOrders,
		driverEvents,
		localLights,
	)

	go network.Network(
		worldview,
		stateBroadcast,
		nodeID,
	)

	go lights.LightsHandler(
		sharedLights,
		localLights,
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
