package orderassigner
import (
	. "Project/dataenums"
)




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
