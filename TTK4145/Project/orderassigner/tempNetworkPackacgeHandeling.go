package orderassigner
import (
	. "Project/dataenums"
)

func mergeHRA(localHRA HRAInput, externalHRA HRAInput, localElevatorName string, incomingElevatorName string) HRAInput {
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