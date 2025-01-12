package orderassigner

import (
	. "Project/dataenums"
)

func buttonPressed(hraInput HRAInput, ElevatorName string,
	btnEvent ButtonEvent) HRAInput {
// Note did something like this ? 
// 	if requests.ShouldClearImmediately(elevator, btnFloor, btn) && (elevator.CurrentBehaviour == elev.EBDoorOpen) {
// er.Start(elevator.Config.DoorOpenDurationS)
// Send dir to driver ? 
// else : 
	switch btnEvent.Button {
		case BHallUp:
			if hraInput.HallRequests[btnEvent.Floor][BHallUp] == Idle{
			hraInput.HallRequests[btnEvent.Floor][BHallUp] = ButtonPressed
			}
		case BHallDown:
			if hraInput.HallRequests[btnEvent.Floor][BHallDown]== Idle{
			hraInput.HallRequests[btnEvent.Floor][BHallDown] = ButtonPressed
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
						hraInput.HallRequests[floor][BHallUp] = OrderComplete
					case BHallDown:
						hraInput.HallRequests[floor][BHallDown] = OrderComplete
					case BCab:
						hraInput.States[elevatorName].CabRequests[floor] = false
					}
				}
			}
		}
	return hraInput
	}
/*
func buttonPressed(btnEvent ButtonEvent) [NFloors][NButtons]bool {
	var buttons [NFloors][NButtons]bool
	buttons[btnEvent.Floor][btnEvent.Button] = true
	return buttons
}
func buttonAlreadyActive(HRAInput HRAInput,elevatorName string,
	btnEvent ButtonEvent) bool {
	switch btnEvent.Button {
	case BHallUp:
		return HRAInput.HallRequests[btnEvent.Floor][BHallUp]
	case BHallDown:
		return HRAInput.HallRequests[btnEvent.Floor][BHallDown]
	case BCab:
		return HRAInput.States[elevatorName].CabRequests[btnEvent.Floor]
	}
	return false
}
*/