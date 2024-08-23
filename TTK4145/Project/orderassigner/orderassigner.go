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

	drv_buttons := make(chan ButtonEvent)
	go hwelevio.PollButtons(drv_buttons)

	for {
		select {
		case btnEvent := <-drv_buttons:
			fmt.Println("button pressed")
			hraInput = ButtonPressed(hraInput, nodeID, btnEvent)
			toNetworkChannel <- hraInput

		case payload := <-payloadFromElevator:
			hraInput = handlePayloadFromElevator(hraInput, payload.Elevator, nodeID)
			hraInput = orderComplete(hraInput, nodeID, payload.CompletedOrders)
			fmt.Println("elevator was changed")
			toNetworkChannel <- hraInput

		case incomingmsg := <-fromNetworkChannel:
			hraInput = mergeHRA(hraInput, incomingmsg.Payload, incomingmsg.SenderId)
			newOrderChannel <- assignOrders(hraInput, nodeID)
			fmt.Println("nye meldinger incomming")
		}
	}
}
