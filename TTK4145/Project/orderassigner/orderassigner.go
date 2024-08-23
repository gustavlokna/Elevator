package orderassigner

import (
	. "Project/dataenums"
	//"Project/elevatordriver"
	"Project/hwelevio"
	//"time"
)

func OrderAssigner(
	newOrderChannel chan<- [NFloors][NButtons]bool,
	newStateChanel <-chan Elevator,
	orderDoneChannel <-chan [NFloors][NButtons]bool,
	toNetworkChannel chan<- HRAInput,
	fromNetworkChannel <-chan HRAInput,
	nodeID string,
) {
	var (
		hraInput       = InitialiseHRAInput()
		onlyNodeOnline = true
	)

	elevator := <-newStateChanel
	hraInput = addElevatorToHRA(hraInput, elevator, nodeID)

	drv_buttons := make(chan ButtonEvent)
	go hwelevio.PollButtons(drv_buttons)

	for {
		select {

		case btnEvent := <-drv_buttons:
			print("button pressed")
			if !buttonAlreadyActive(hraInput, nodeID, btnEvent) {
				print("new order")
				hraInput = ButtonPressed(hraInput, nodeID, btnEvent)
				//newOrderChannel <- AssignOrders(hraInput)
				// newOrderChannel <- buttonPressed(btnEvent)
			}

		case completedOrders := <-orderDoneChannel:
			hraInput = OrderComplete(hraInput, nodeID, completedOrders)
			// Optionally, update other systems with the updated hraInput
			// newOrderChannel <- AssignOrders(hraInput)

		case elev := <-newStateChanel:
			hraInput = addElevatorToHRA(hraInput, elev, nodeID)
			//newOrderChannel <- AssignOrders(hraInput)
			print("elevator was changed")
			toNetworkChannel <- hraInput
		case hraInput = <-fromNetworkChannel:
			print("nye meldinger incomming")

		}
		

		if onlyNodeOnline {
			//print("assigns")
			newOrderChannel <- assignOrders(hraInput, nodeID)
		}
	}
}
