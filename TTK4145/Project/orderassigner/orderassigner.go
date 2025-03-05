package orderassigner

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"fmt"
	"strconv"
	//"time"
)

func OrderAssigner(
	newOrderChannel chan<- [NFloors][NButtons]bool,
	payloadFromElevator <-chan PayloadFromElevator,
	toNetworkChannel chan<- PayloadFromassignerToNetwork,
	fromNetworkChannel <-chan PayloadFromNetworkToAssigner,
	// TODO: use this fromAsstoLight
	fromAsstoLight chan<- [NFloors][NButtons]ButtonState,
	nodeID string,
) {
	var (
		PayloadFromassignerToNetwork = initPayloadToNetwork()
		prevAssignedOrders           [NFloors][NButtons]bool
		drv_buttons                  = make(chan ButtonEvent)
	)
	// Convert nodeID to int
	myID, err := strconv.Atoi(nodeID)
	if err != nil {
		fmt.Printf("Invalid nodeID: %v\n", err)
		return
	}
	payload := <-payloadFromElevator
	PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
		PayloadFromassignerToNetwork, nodeID)
	//check if it creates error by sending to network here
	toNetworkChannel <- PayloadFromassignerToNetwork

	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btnEvent := <-drv_buttons:
			PayloadFromassignerToNetwork = buttonPressed(PayloadFromassignerToNetwork,
				nodeID, btnEvent)
			toNetworkChannel <- PayloadFromassignerToNetwork

		case payload := <-payloadFromElevator:
			PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
				PayloadFromassignerToNetwork, nodeID)

			toNetworkChannel <- PayloadFromassignerToNetwork

		case PayloadFromNetwork := <-fromNetworkChannel:

			PayloadFromassignerToNetwork = handlePayloadFromNetwork(PayloadFromassignerToNetwork,
				PayloadFromNetwork, myID)

			newOrders := assignOrders(PayloadFromNetwork, myID)

			if newOrders != prevAssignedOrders {
				newOrderChannel <- newOrders
				prevAssignedOrders = newOrders
			}
			fromAsstoLight <- updateLightStates(PayloadFromNetwork, myID)

		}
	}

}
