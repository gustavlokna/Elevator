package driver

import (
	. "Project/config"
	. "Project/dataenums"
)

func orderAtCurrentFloorInDirection(elevator Elevator) bool {
	switch elevator.Direction {
	case MDUp:
		return elevator.Requests[elevator.CurrentFloor][BHallUp]
	case MDDown:
		return elevator.Requests[elevator.CurrentFloor][BHallDown]
	default:
		return elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BHallDown]
	}
}

func orderAtCurrentFloorOppositeDirection(elevator Elevator) bool {
	switch elevator.Direction {
	case MDUp:
		return elevator.Requests[elevator.CurrentFloor][BHallDown]
	case MDDown:
		return elevator.Requests[elevator.CurrentFloor][BHallUp]
	default:
		return elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BHallDown]
	}
}

func orderCurrentDirection(elevator Elevator) bool {
	switch elevator.Direction {
	case MDUp:
		return requestsAbove(elevator)
	case MDDown:
		return requestsBelow(elevator)
	}
	return false
}

func orderOppositeDirection(elevator Elevator) bool {
	switch elevator.Direction {
	case MDUp:
		return requestsBelow(elevator)
	case MDDown:
		return requestsAbove(elevator)
	default:
		return requestsBelow(elevator) || requestsAbove(elevator)

	}
}

func requestsAbove(elevator Elevator) bool {
	for floor := elevator.CurrentFloor + 1; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if elevator.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(elevator Elevator) bool {
	for floor := 0; floor < elevator.CurrentFloor; floor++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if elevator.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}
