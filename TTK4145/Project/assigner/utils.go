package assigner

import (
	. "Project/dataenums"
	. "Project/config"
)

//TODO Change payload to worldview or something and stateUpdate also

func initPayloadToNetwork() FromAssignerToNetwork {
	worldview := FromAssignerToNetwork{
		HallRequests: [NFloors][NButtons]ButtonState{},
		States:       make(map[string]HRAElevState),
	}
	return worldview
}
func updateLightStates(stateBroadcast FromNetworkToAssigner,
	myID int) [NFloors][NButtons]ButtonState {

	updatedLights := stateBroadcast.HallOrderList[myID]
	for floor := 0; floor < NFloors; floor++ {
		if stateBroadcast.ElevatorList[myID].CabRequests[floor] {
			updatedLights[floor][BCab] = OrderAssigned
		}
	}
	return updatedLights
}

// TODO ALEX WRITE BETTER 
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

// TODO BETTER NAME THAN toAssigner REQUIRED
func handlePayloadFromElevator(driverEvents FromDriverToAssigner,
	worldview FromAssignerToNetwork, nodeID string) FromAssignerToNetwork {

	cabRequests := make([]bool, NFloors)
	for floor := 0; floor < NFloors; floor++ {
		cabRequests[floor] = driverEvents.Elevator.Requests[floor][BCab]
	}

	worldview.States[nodeID] = HRAElevState{
		Behaviour:   ebToString(driverEvents.Elevator.CurrentBehaviour),
		Floor:       driverEvents.Elevator.CurrentFloor,
		Direction:   elevDirnToString(driverEvents.Elevator.Dirn),
		CabRequests: cabRequests,
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
