package driver

import (
	. "Project/dataenums"
	"fmt"
)

func setMotorDir(dir HWMotorDirection) Button {
	switch dir {
	case MDUp:
		return BHallUp
	case MDDown:
		return BHallDown
	default:
		panic("invalid direction")
	}
}

func setMotorOppositeDir(e Elevator) HWMotorDirection {
	switch e.Dirn {
	case MDUp:
		return MDDown
	case MDDown:
		return MDUp
	case MDStop:
		switch {
		case requestsAbove(e) || e.Requests[e.CurrentFloor][BHallUp]:
			return MDDown
		case requestsBelow(e) || e.Requests[e.CurrentFloor][BHallDown]:
			return MDUp
		default:
			return MDStop
		}
	default: 
		panic("invalid direction")
	}
}

func orderAtCurrentFloorInDir(e Elevator) bool {
	switch e.Dirn {
	case MDUp:
		return e.Requests[e.CurrentFloor][BHallUp] || e.Requests[e.CurrentFloor][BCab]
	case MDDown:
		return e.Requests[e.CurrentFloor][BHallDown] || e.Requests[e.CurrentFloor][BCab]
	default:
		return e.Requests[e.CurrentFloor][BHallUp] || e.Requests[e.CurrentFloor][BHallDown] || e.Requests[e.CurrentFloor][BCab]
	}
}

func orderAtCurrentFloorOppositeDir(e Elevator) bool {
	switch e.Dirn {
	case MDUp:
		return e.Requests[e.CurrentFloor][BHallDown] || e.Requests[e.CurrentFloor][BCab]
	case MDDown:
		return e.Requests[e.CurrentFloor][BHallUp] || e.Requests[e.CurrentFloor][BCab]
	default:
		return e.Requests[e.CurrentFloor][BHallUp] || e.Requests[e.CurrentFloor][BHallDown] || e.Requests[e.CurrentFloor][BCab] //?
	}
}

func orderInCurrentDir(e Elevator) bool {
	switch e.Dirn {
	case MDUp:
		return requestsAbove(e)
	case MDDown:
		return requestsBelow(e)
	}
	return false
}

func orderOppositeDir(e Elevator) bool {
	switch e.Dirn {
	case MDUp:
		return requestsBelow(e)
	case MDDown:
		return requestsAbove(e)
	default:
		return requestsBelow(e) || requestsAbove(e)

	}
}

func requestsAbove(elevator Elevator) bool {
	for f := elevator.CurrentFloor + 1; f < NFloors; f++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if elevator.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(elevator Elevator) bool {
	for f := 0; f < elevator.CurrentFloor; f++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if elevator.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

// TODO REMOVE
func ElevatorPrint(e Elevator) {
	fmt.Println("\n  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12s|\n"+
			"  |behav = %-12s|\n",
		e.CurrentFloor,
		ElevDirToString(e.Dirn),
		EBToString(e.CurrentBehaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := NFloors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := BHallUp; btn <= BCab; btn++ {
			if (f == NFloors-1 && btn == BHallUp) ||
				(f == 0 && btn == BHallDown) {
				fmt.Print("|     ")
			} else {
				if e.Requests[f][btn] {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func EBToString(behaviour ElevatorBehaviour) string {
	switch behaviour {
	case EBIdle:
		return "idle"
	case EBDoorOpen:
		return "doorOpen"
	case EBMoving:
		return "moving"
	default:
		return "Unknown"
	}
}
func ElevDirToString(d HWMotorDirection) string {
	switch d {
	case MDDown:
		return "down"
	case MDStop:
		return "stop"
	case MDUp:
		return "up"
	default:
		return "DirUnknown"
	}
}
