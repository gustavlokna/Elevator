package orderassigner

import (
	. "Project/dataenums"
	//"Project/elevatordriver"
	"Project/hwelevio"
	"fmt"
	//"time"
)

func OrderAssigner(
	newOrderChannel chan<- [NFloors][NButtons]bool,
	payloadFromElevator <-chan PayloadFromElevator,
	toNetworkChannel chan<- HRAInput,
	fromNetworkChannel <-chan Message,
	nodeID string,
) {
	var (
		hraInput = InitialiseHRAInput()
	)
	payload := <-payloadFromElevator
	hraInput = handlePayloadFromElevator(hraInput, payload.Elevator, nodeID)
	//check if it creates error by sending to network here
	toNetworkChannel <- hraInput

	drv_buttons := make(chan ButtonEvent)
	go hwelevio.PollButtons(drv_buttons)
	print("PENIS")
	
	for {
		select {
		case btnEvent := <-drv_buttons:
			//Note make cylick counter own module and put this there ? 
			fmt.Println("button pressed")

			hraInput = buttonPressed(hraInput, nodeID, btnEvent)
			//PrintHRAInput(hraInput)
			toNetworkChannel <- hraInput
		/*
		case payload := <-payloadFromElevator:
			hraInput = handlePayloadFromElevator(hraInput, payload.Elevator, nodeID)
			hraInput = orderComplete(hraInput, nodeID, payload.CompletedOrders)
			fmt.Println("elevator was changed")
			toNetworkChannel <- hraInput
		*/
		case  <-fromNetworkChannel:
			//TODO DEFINIE incomingmsg as the from NetworkChannel
			//aInput = mergeHRA(hraInput, incomingmsg.Payload, incomingmsg.SenderId)
			//newOrderChannel <- assignOrders(hraInput, nodeID)
			fmt.Println("nye meldinger incomming")
		
		}
	}
	
}
