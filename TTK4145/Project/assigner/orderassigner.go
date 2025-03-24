package assigner

import (
	. "Project/dataenums"
	. "Project/config"
	"Project/hwelevio"
	"fmt"
	"strconv"
)

func Assigner(
	newOrders chan<- [NFloors][NButtons]bool,
	driverEvents <-chan FromDriverToAssigner,
	worldview chan<- FromAssignerToNetwork,
	stateBroadcast <-chan FromNetworkToAssigner,
	sharedLights chan<- [NFloors][NButtons]ButtonState,
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
	
	payload := <-driverEvents
	PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
		PayloadFromassignerToNetwork, nodeID)

		worldview <- PayloadFromassignerToNetwork
	
	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btnEvent := <-drv_buttons:
			PayloadFromassignerToNetwork = handleButtonPressed(PayloadFromassignerToNetwork,
				nodeID, btnEvent)
				worldview <- PayloadFromassignerToNetwork

		case payload := <-driverEvents:
			PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
				PayloadFromassignerToNetwork, nodeID)

				worldview <- PayloadFromassignerToNetwork

		case PayloadFromNetwork := <-stateBroadcast:
			
			PayloadFromassignerToNetwork = handlePayloadFromNetwork(PayloadFromassignerToNetwork,
				PayloadFromNetwork, myID)
			
			localOrders := assignOrders(PayloadFromNetwork, myID)
			if localOrders != prevAssignedOrders {
				newOrders <- localOrders
				prevAssignedOrders = localOrders
			}
			sharedLights <- updateLightStates(PayloadFromNetwork, myID)
		}
	}

}
