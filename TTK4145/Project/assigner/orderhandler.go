package assigner

import (
	. "Project/config"
	. "Project/dataenums"
	"strconv"
)

func handleButtonPressed(worldview FromAssignerToNetwork, nodeID int,
	btnEvent ButtonEvent) FromAssignerToNetwork {
	switch btnEvent.Button {
	case BHallUp:
		if worldview.HallRequests[btnEvent.Floor][BHallUp] != OrderComplete {
			worldview.HallRequests[btnEvent.Floor][BHallUp] = ButtonPressed
		}
	case BHallDown:
		if worldview.HallRequests[btnEvent.Floor][BHallDown] != OrderComplete {
			worldview.HallRequests[btnEvent.Floor][BHallDown] = ButtonPressed
		}

	case BCab:
		worldview.States[strconv.Itoa(nodeID)].CabRequests[btnEvent.Floor] = true

	}
	return worldview
}

func handleOrderComplete(worldview FromAssignerToNetwork, nodeID int,
	completedOrders [NFloors][NButtons]bool) FromAssignerToNetwork {
	for floor := 0; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if completedOrders[floor][btn] {
				switch btn {
				case BHallUp:
					worldview.HallRequests[floor][BHallUp] = OrderComplete
				case BHallDown:
					worldview.HallRequests[floor][BHallDown] = OrderComplete
				case BCab:
					worldview.States[strconv.Itoa(nodeID)].CabRequests[floor] = false
				}
			}
		}
	}
	return worldview
}
