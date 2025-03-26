package assigner

import (
	. "Project/config"
	. "Project/dataenums"
)

func initLocalWorldview(elevatorState FromDriverToAssigner,
	globaWorldview FromNetworkToAssigner,
	nodeID int) FromAssignerToNetwork {
	localWorldview := FromAssignerToNetwork{
		HallRequests: [NFloors][NButtons]ButtonState{},
		States:       make(map[int]HRAElevState),
	}
	localWorldview.States[nodeID] = HRAElevState{
		Behaviour:   elevbehaviourToString(elevatorState.Elevator.CurrentBehaviour),
		Floor:       elevatorState.Elevator.CurrentFloor,
		Direction:   elevDirectionToString(elevatorState.Elevator.Direction),
		CabRequests: globaWorldview.ElevatorList[nodeID].CabRequests,
	}
	return localWorldview
}

func updateLightStates(globaWorldview FromNetworkToAssigner,
	nodeID int) [NFloors][NButtons]ButtonState {

	updatedLights := globaWorldview.HallOrderList[nodeID]
	for floor := 0; floor < NFloors; floor++ {
		if globaWorldview.ElevatorList[nodeID].CabRequests[floor] {
			updatedLights[floor][BCab] = OrderAssigned
		}
	}
	return updatedLights
}

func mergeNetworkHallOrders(
	localWorldview FromAssignerToNetwork,
	globaWorldview FromNetworkToAssigner,
	nodeID int) FromAssignerToNetwork {
	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			incommingOrder := globaWorldview.HallOrderList[nodeID][floor][btn]
			localOrder := localWorldview.HallRequests[floor][btn]
			if localOrder != OrderComplete || incommingOrder != OrderAssigned {
				localWorldview.HallRequests[floor][btn] = incommingOrder
			}
		}
	}
	return localWorldview
}

func syncElevatorState(elevatorState FromDriverToAssigner,
	localWorldview FromAssignerToNetwork,
	nodeID int) FromAssignerToNetwork {

	localWorldview.States[nodeID] = HRAElevState{
		Behaviour:   elevbehaviourToString(elevatorState.Elevator.CurrentBehaviour),
		Floor:       elevatorState.Elevator.CurrentFloor,
		Direction:   elevDirectionToString(elevatorState.Elevator.Direction),
		CabRequests: localWorldview.States[nodeID].CabRequests,
	}
	localWorldview.ActiveStatus = elevatorState.Elevator.ActiveStatus
	localWorldview = handleOrderComplete(localWorldview, nodeID, elevatorState.CompletedOrders)

	return localWorldview
}

func elevbehaviourToString(behaviour ElevatorBehaviour) string {
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

func elevDirectionToString(direction MotorDirection) string {
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
