package orderassigner

import (
	. "Project/dataenums"
)

func buttonPressed(payload PayloadFromassignerToNetwork, ElevatorName string,
	btnEvent ButtonEvent) PayloadFromassignerToNetwork {
	switch btnEvent.Button {
	case BHallUp:
		if payload.Orders[btnEvent.Floor][BHallUp] != OrderComplete {
			payload.Orders[btnEvent.Floor][BHallUp] = ButtonPressed
		}
	case BHallDown:
		if payload.Orders[btnEvent.Floor][BHallDown] != OrderComplete {
			payload.Orders[btnEvent.Floor][BHallDown] = ButtonPressed
		}
	case BCab:
		payload.Orders[btnEvent.Floor][BCab] = OrderAssigned
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
					payload.Orders[floor][BHallUp] = OrderComplete
				case BHallDown:
					payload.Orders[floor][BHallDown] = OrderComplete
				case BCab:
					payload.Orders[floor][BCab] = OrderComplete
				}
			}
		}
	}
	return payload
}
