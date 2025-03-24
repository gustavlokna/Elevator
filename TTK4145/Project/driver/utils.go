package driver

import (
	. "Project/config"
	. "Project/dataenums"
	"fmt"
)

func dirnToBtn(dirn MotorDirection) Button {
	switch dirn {
	case MDUp:
		return BHallUp
	case MDDown:
		return BHallDown
	default:
		panic("invalid dirn in dirToEnum ")
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
func ElevatorPrint(elevator Elevator) {
	fmt.Println("\n  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12s|\n"+
			"  |behav = %-12s|\n",
		elevator.CurrentFloor,
		ElevDirnToString(elevator.Dirn),
		EBToString(elevator.CurrentBehaviour),
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
				if elevator.Requests[f][btn] {
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
	case Idle:
		return "idle"
	case DoorOpen:
		return "doorOpen"
	case Moving:
		return "moving"
	default:
		return "Unknown"
	}
}
func ElevDirnToString(direction MotorDirection) string {
	switch direction {
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
