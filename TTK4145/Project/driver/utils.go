package driver

import (
	. "Project/config"
	. "Project/dataenums"
)

func dirnToBtn(dirn MotorDirection) Button {
	switch dirn {
	case MDUp:
		return BHallUp
	case MDDown:
		return BHallDown
	default:
		panic("invalid dirn in dirnToBtn ")
	}
}

func btnToDirn(elevator Elevator) MotorDirection {
	switch {
	case elevator.Requests[elevator.CurrentFloor][BHallUp]:
		return MDUp
	case elevator.Requests[elevator.CurrentFloor][BHallDown]:
		return MDDown
	default:
		return MDStop
	}
}

func setMotorOppositeDirn(elevator Elevator) MotorDirection {
	switch elevator.Dirn {
	case MDUp:
		return MDDown
	case MDDown:
		return MDUp
	default:
		switch {
		case requestsAbove(elevator) || elevator.Requests[elevator.CurrentFloor][BHallUp]:
			return MDUp
		case requestsBelow(elevator) || elevator.Requests[elevator.CurrentFloor][BHallDown]:
			return MDDown
		default:
			return MDStop
		}
	}
}

func orderAtCurrentFloorInDirn(elevator Elevator) bool {
	switch elevator.Dirn {
	case MDUp:
		return elevator.Requests[elevator.CurrentFloor][BHallUp]
	case MDDown:
		return elevator.Requests[elevator.CurrentFloor][BHallDown]
	default:
		return elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BHallDown]
	}
}

func orderAtCurrentFloorOppositeDirn(elevator Elevator) bool {
	switch elevator.Dirn {
	case MDUp:
		return elevator.Requests[elevator.CurrentFloor][BHallDown]
	case MDDown:
		return elevator.Requests[elevator.CurrentFloor][BHallUp]
	default:
		return elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BHallDown]
	}
}

func orderCurrentDirn(elevator Elevator) bool {
	switch elevator.Dirn {
	case MDUp:
		return requestsAbove(elevator)
	case MDDown:
		return requestsBelow(elevator)
	}
	return false
}

func orderOppositeDirn(elevator Elevator) bool {
	switch elevator.Dirn {
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
