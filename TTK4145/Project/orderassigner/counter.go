package orderassigner
import (
	. "Project/dataenums"
)



func InitialiseHRAInput() HRAInput {
	hraInput := HRAInput{
		HallRequests: make([][2]bool, NFloors),
		States:       make(map[string]HRAElevState),
		CounterHallRequests : make([][2]int, NFloors),
	}
	return hraInput
}
func handlePayloadFromElevator(hraInput HRAInput, e Elevator,
	elevatorName string) HRAInput{
	behavior, direction, cabRequests := convertElevatorState(e)
	hraInput.States[elevatorName] = HRAElevState{
		Behavior:    behavior,
		Floor:       e.CurrentFloor,
		Direction:   direction,
		CabRequests: cabRequests,
	}
	return hraInput
}

func handlePayloadFromNetwork(localHRA HRAInput, Incoming HRAInput)HRAInput{
	return Incoming
}

func convertElevatorState(e Elevator) (string, string, []bool) {
	behavior := EBToString(e.CurrentBehaviour)
	direction := ElevDirToString(e.Dirn)

	// Convert cab requests
	cabRequests := make([]bool, NFloors)
	for f := 0; f < NFloors; f++ {
		cabRequests[f] = e.Requests[f][BCab]
	}
	return behavior, direction, cabRequests
}

func ButtonPressed(hraInput HRAInput, ElevatorName string,
					btnEvent ButtonEvent) HRAInput {
	switch btnEvent.Button {
	case BHallUp:
		if !hraInput.HallRequests[btnEvent.Floor][BHallUp]{
			hraInput.HallRequests[btnEvent.Floor][BHallUp] = true
			hraInput.CounterHallRequests[btnEvent.Floor][BHallUp]++
		}
	case BHallDown:
		if !hraInput.HallRequests[btnEvent.Floor][BHallDown]{
			hraInput.HallRequests[btnEvent.Floor][BHallDown] = true
			hraInput.CounterHallRequests[btnEvent.Floor][BHallDown]++
		}
	case BCab:
		print("CAB BUTTON PRESSED")
		hraInput.States[ElevatorName].CabRequests[btnEvent.Floor] = true
	}
	return hraInput
}



func orderComplete(hraInput HRAInput, elevatorName string,
	completedOrders [NFloors][NButtons]bool) HRAInput {
	for floor := 0; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if completedOrders[floor][btn] {
				switch btn {
				case BHallUp:
					hraInput.HallRequests[floor][BHallUp] = false
					hraInput.CounterHallRequests[floor][BHallUp]++
				case BHallDown:
					hraInput.HallRequests[floor][BHallDown] = false
					hraInput.CounterHallRequests[floor][BHallDown]++
				case BCab:
					hraInput.States[elevatorName].CabRequests[floor] = false
				}
			}
		}
	}
	return hraInput
}
