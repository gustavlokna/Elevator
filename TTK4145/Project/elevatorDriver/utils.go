package elevatorDriver

import (
	. "Project/dataenums"
)

func initelevator() Elevator {
	return Elevator{
		CurrentFloor:     -1,
		Direction:        MDDown,
		CurrentBehaviour: Moving,
		ActiveStatus:     true,
	}
}

func directionToButton(direction MotorDirection) Button {
	switch direction {
	case MDUp:
		return BHallUp
	case MDDown:
		return BHallDown
	default:
		panic("invalid Direction in directionToBtn ")
	}
}

func buttonToDirection(elevator Elevator) MotorDirection {
	switch {
	case elevator.Requests[elevator.CurrentFloor][BHallUp]:
		return MDUp
	case elevator.Requests[elevator.CurrentFloor][BHallDown]:
		return MDDown
	default:
		return MDStop
	}
}

func setMotorOppositeDirection(elevator Elevator) MotorDirection {
	switch elevator.Direction {
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
