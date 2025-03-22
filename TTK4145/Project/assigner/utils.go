package assigner

import (
	. "Project/dataenums"
)

func initPayloadToNetwork() FromAssignerToNetwork {
	payload := FromAssignerToNetwork{
		HallRequests: [NFloors][NButtons]ButtonState{},
		States:       make(map[string]HRAElevState),
	}
	return payload
}
func updateLightStates(payload FromNetworkToAssigner,
	myID int) [NFloors][NButtons]ButtonState {

	updatedLights := payload.HallOrderList[myID]
	for floor := 0; floor < NFloors; floor++ {
		if payload.ElevatorList[myID].CabRequests[floor] {
			updatedLights[floor][BCab] = OrderAssigned
		}
	}
	return updatedLights
}

func handlePayloadFromNetwork(
	payload FromAssignerToNetwork,
	netPayload FromNetworkToAssigner,
	nodeID int,
) FromAssignerToNetwork {
	for f := 0; f < NFloors; f++ {
		for b := 0; b < NButtons; b++ {
			incomingState := netPayload.HallOrderList[nodeID][f][b]
			localState := payload.HallRequests[f][b]

			if localState == OrderComplete && incomingState == OrderAssigned {
				// do nothing; stay OrderComplete.
			} else {
				payload.HallRequests[f][b] = incomingState
			}

		}
	}
	return payload
}

// TODO BETTER NAME THAN msg REQUIRED
func handlePayloadFromElevator(msg FromDriverToAssigner,
	toNetwork FromAssignerToNetwork, nodeID string) FromAssignerToNetwork {

	cabRequests := make([]bool, NFloors)
	for f := 0; f < NFloors; f++ {
		cabRequests[f] = msg.Elevator.Requests[f][BCab]
	}

	toNetwork.States[nodeID] = HRAElevState{
		Behaviour:   ebToString(msg.Elevator.CurrentBehaviour),
		Floor:       msg.Elevator.CurrentFloor,
		Direction:   elevDirToString(msg.Elevator.Dirn),
		CabRequests: cabRequests,
	}
	toNetwork.ActiveSatus = msg.Elevator.ActiveSatus
	toNetwork = orderComplete(toNetwork, nodeID, msg.CompletedOrders)

	return toNetwork
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

func elevDirToString(d HWMotorDirection) string {
	switch d {
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
