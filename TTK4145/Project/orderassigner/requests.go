package orderassigner

import (
	. "Project/dataenums"
)

func buttonPressed(payload PayloadFromassignerToNetwork, ElevatorName string,
	btnEvent ButtonEvent) PayloadFromassignerToNetwork {
// Note did something like this ? 
// 	if requests.ShouldClearImmediately(elevator, btnFloor, btn) && (elevator.CurrentBehaviour == elev.EBDoorOpen) {
// er.Start(elevator.Config.DoorOpenDurationS)
// Send dir to driver ? 
// else : 
switch btnEvent.Button {
case BHallUp:
	payload.HallRequests[btnEvent.Floor][BHallUp] = ButtonPressed
case BHallDown:
	payload.HallRequests[btnEvent.Floor][BHallDown] = ButtonPressed
	
case BCab:
	// For Cab button press
	print("CAB BUTTON PRESSED")
	payload.States[ElevatorName].CabRequests[btnEvent.Floor] = true
}
return payload
}



func orderComplete(payload PayloadFromassignerToNetwork, elevatorName string,
	completedOrders [NFloors][NButtons]bool) PayloadFromassignerToNetwork {
	for floor := 0; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if completedOrders[floor][btn] {
				switch btn {
					case BHallUp:
						payload.HallRequests[floor][BHallUp] = OrderComplete
					case BHallDown:
						payload.HallRequests[floor][BHallDown] = OrderComplete
					case BCab:
						payload.States[elevatorName].CabRequests[floor] = false
					}
				}
			}
		}
	return payload
	}