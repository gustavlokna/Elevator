package orderassigner
/*
import (
	. "Project/dataenums"
)

// DISCLAIMER: This implementation is temporary and should be replaced
// by cyclic counter logic for better efficiency and reliability in the future.

/*
func mergeHRA(localHRA HRAInput, externalHRA HRAInput, incomingElevatorName string) HRAInput {
	// Handle incoming hall requests
	for f := 0; f < NFloors; f++ {
		for btn := BHallUp; btn <= BHallDown; btn++ {
			if externalHRA.CounterHallRequests[f][btn] > localHRA.CounterHallRequests[f][btn] {
				localHRA.CounterHallRequests[f][btn] = externalHRA.CounterHallRequests[f][btn]
				localHRA.HallRequests[f][btn] = externalHRA.HallRequests[f][btn]
			}
			if externalHRA.CounterHallRequests[f][btn] == localHRA.CounterHallRequests[f][btn] {
				if localHRA.HallRequests[f][btn] != externalHRA.HallRequests[f][btn] {
					localHRA.HallRequests[f][btn] = false
				}
			}
		}
	}

	//should be no need for this 
	// we add an (active variable to msg. if it is false we append id to inactive elev list)
	if _, exists := localHRA.States[incomingElevatorName]; exists {
		localHRA.States[incomingElevatorName] = externalHRA.States[incomingElevatorName]
	} else {
		RemoveElevatorsFromState(localHRA, []string{incomingElevatorName})
	}

	return localHRA
}

func RemoveElevatorsFromState(hraInput HRAInput, elevatorNames []string) {
	for _, name := range elevatorNames {
		delete(hraInput.States, name)
	}
}
*/