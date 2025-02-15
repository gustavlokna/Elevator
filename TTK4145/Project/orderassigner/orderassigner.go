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
	fromAsstoLight chan <- [NFloors][NButtons]ButtonState,
	nodeID string,
) {
	var (
		PayloadFromassignerToNetwork = InitialisePayloadFromassignerToNetwork()
		prevAssignedOrders           [NFloors][NButtons]bool // Track previous assigned orders
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

	drv_buttons := make(chan ButtonEvent)
	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btnEvent := <-drv_buttons:
			fmt.Println("button pressed")
			// TODO do not overwrite this is fixed when we get the fromNetworkChannel working
			PayloadFromassignerToNetwork = buttonPressed(PayloadFromassignerToNetwork, 
				nodeID, btnEvent)
			//PrintHRAInput(hraInput)
			toNetworkChannel <- PayloadFromassignerToNetwork
		
		case payload := <-payloadFromElevator:
			PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
				PayloadFromassignerToNetwork, nodeID)
			//PrintPayloadFromassignerToNetwork(PayloadFromassignerToNetwork)
			toNetworkChannel <- PayloadFromassignerToNetwork
		
		case PayloadFromNetwork := <-fromNetworkChannel:
			// TOOD CANNOT REMOVE YET: NEED FUNC, BUT I WANT TO 
			// THIS CAN HAPPEN IF THERE IS MISSHAP IN THE ORDER THINGS OCCUR 
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