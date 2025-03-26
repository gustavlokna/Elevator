package assigner

import (
	. "Project/config"
	. "Project/dataenums"
)

func initPayloadToNetwork(driverEvents FromDriverToAssigner, stateBroadcast FromNetworkToAssigner, nodeID int) FromAssignerToNetwork {

	worldview := FromAssignerToNetwork{
		HallRequests: [NFloors][NButtons]ButtonState{},
		States:       make(map[int]HRAElevState),
	}
	worldview.States[nodeID] = HRAElevState{
		Behaviour:   ebToString(driverEvents.Elevator.CurrentBehaviour),
		Floor:       driverEvents.Elevator.CurrentFloor,
		Direction:   elevDirnToString(driverEvents.Elevator.Dirn),
		CabRequests: stateBroadcast.ElevatorList[nodeID].CabRequests,
	}
	return worldview
}
func updateLightStates(stateBroadcast FromNetworkToAssigner,
	nodeID int) [NFloors][NButtons]ButtonState {

	updatedLights := stateBroadcast.HallOrderList[nodeID]
	for floor := 0; floor < NFloors; floor++ {
		if stateBroadcast.ElevatorList[nodeID].CabRequests[floor] {
			updatedLights[floor][BCab] = OrderAssigned
		}
	}
	return updatedLights
}

func handlePayloadFromNetwork(
	worldview FromAssignerToNetwork,
	stateBroadcast FromNetworkToAssigner,
	nodeID int,
) FromAssignerToNetwork {
	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			incomingState := stateBroadcast.HallOrderList[nodeID][floor][btn]
			localState := worldview.HallRequests[floor][btn]
			if localState != OrderComplete || incomingState != OrderAssigned {
				worldview.HallRequests[floor][btn] = incomingState
			}
		}
	}
	return worldview
}

func handlePayloadFromElevator(driverEvents FromDriverToAssigner,
	worldview FromAssignerToNetwork, nodeID int) FromAssignerToNetwork {

	worldview.States[nodeID] = HRAElevState{
		Behaviour:   ebToString(driverEvents.Elevator.CurrentBehaviour),
		Floor:       driverEvents.Elevator.CurrentFloor,
		Direction:   elevDirnToString(driverEvents.Elevator.Dirn),
		CabRequests: worldview.States[nodeID].CabRequests,
	}
	worldview.ActiveStatus = driverEvents.Elevator.ActiveStatus
	worldview = handleOrderComplete(worldview, nodeID, driverEvents.CompletedOrders)

	return worldview
}

func ebToString(behaviour ElevatorBehaviour) string {
	switch behaviour {
	case Idle:
		return "idle"
	case DoorOpen:
		return "doorOpen"
	case Moving:
		return "moving"
	default:
		return "Unknown"
	}
}

func elevDirnToString(direction MotorDirection) string {
	switch direction {
	case MDDown:
		return "down"
	case MDStop:
		return "stop"
	case MDUp:
		return "up"
	default:
		return "DirUnknown"
	}
}
