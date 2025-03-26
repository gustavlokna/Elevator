package assigner

import (
	. "Project/config"
	. "Project/dataenums"
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
	nodeID int,
) {
	var (
		drv_buttons                  = make(chan ButtonEvent)
	)
	payloadFormDriver := <-driverEvents
	PayloadFromNetwork := <-stateBroadcast:
	PayloadFromassignerToNetwork := initPayloadToNetwork(payloadFormDriver,
		PayloadFromNetwork, nodeID)
	worldview <- PayloadFromassignerToNetwork

	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btnEvent := <-drv_buttons:
			PayloadFromassignerToNetwork = handleButtonPressed(PayloadFromassignerToNetwork,
				nodeID, btnEvent)
			worldview <- PayloadFromassignerToNetwork

		case payloadFormDriver := <-driverEvents:
			PayloadFromassignerToNetwork = handlePayloadFromElevator(payloadFormDriver,
				PayloadFromassignerToNetwork, nodeID)

			worldview <- PayloadFromassignerToNetwork

		case PayloadFromNetwork := <-stateBroadcast:
			PayloadFromassignerToNetwork = handlePayloadFromNetwork(PayloadFromassignerToNetwork,
				PayloadFromNetwork, nodeID)

			newOrders <- assignOrders(PayloadFromNetwork, nodeID)
			sharedLights <- updateLightStates(PayloadFromNetwork, nodeID)
		}
	}

}
