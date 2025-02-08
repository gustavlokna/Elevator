package orderassigner

import (
	. "Project/dataenums"
	"time"
	//"Project/elevatordriver"
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
		
		//PayloadFromNetwork PayloadFromNetworkToAssigner
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
	time.Sleep(3 *time.Second)
	for {
		select {
		case btnEvent := <-drv_buttons:
			//Note make cylick counter own module and put this there ? 
			fmt.Println("button pressed")
			// TODO do not overwrite this is fixed when we get the fromNetworkChannel working
			PayloadFromassignerToNetwork = buttonPressed(PayloadFromassignerToNetwork, 
				nodeID, btnEvent)
			//PrintHRAInput(hraInput)
			toNetworkChannel <- PayloadFromassignerToNetwork
		
		case payload := <-payloadFromElevator:
			/*
			hraInput = handlePayloadFromElevator(hraInput, payload.Elevator, nodeID)
			hraInput = orderComplete(hraInput, nodeID, payload.CompletedOrders)
			fmt.Println("elevator was changed")
			*/
			PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
				PayloadFromassignerToNetwork, nodeID)
			//print("hallo")
			toNetworkChannel <- PayloadFromassignerToNetwork
		
		case PayloadFromNetwork := <-fromNetworkChannel:
			//TODO why this. 
			
			PayloadFromassignerToNetwork = handlePayloadFromNetwork(PayloadFromassignerToNetwork, 
				PayloadFromNetwork, myID)
			
			newOrderChannel <- assignOrders(PayloadFromNetwork,myID)
			
			fromAsstoLight <- PayloadFromassignerToNetwork.HallRequests
		
		}
	}
	
}
