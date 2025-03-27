package assigner

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/hwelevio"
)

func Assigner(
	newOrders chan<- [NFloors][NButtons]bool,
	driverEvents <-chan FromDriverToAssigner,
	worldview chan<- FromAssignerToNetwork,
	stateBroadcast <-chan FromNetworkToAssigner,
	sharedLights chan<- [NFloors][NButtons]ButtonState,
	nodeID int) {
	var (
		drv_buttons        = make(chan ButtonEvent)
		prevAssignedOrders [NFloors][NButtons]bool
	)
	elevatorState := <-driverEvents
	globalWorldview := <-stateBroadcast
	localWorldview := initLocalWorldview(elevatorState,
		globalWorldview, nodeID)
	worldview <- localWorldview

	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btnEvent := <-drv_buttons:
			localWorldview = handleButtonPressed(localWorldview,
				nodeID, btnEvent)
			worldview <- localWorldview

		case elevatorState := <-driverEvents:
			localWorldview = syncElevatorState(elevatorState,
				localWorldview, nodeID)

			worldview <- localWorldview

		case globalWorldview := <-stateBroadcast:
			localWorldview = mergeNetworkHallOrders(localWorldview,
				globalWorldview, nodeID)

			localOrders := assignOrders(globalWorldview, nodeID)
			if localOrders != prevAssignedOrders {
				newOrders <- localOrders
				prevAssignedOrders = localOrders
			}
			sharedLights <- updateLightStates(globalWorldview, nodeID)
		}
	}

}
