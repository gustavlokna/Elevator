package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"fmt"
)


func ElevatorPrint(e Elevator) {
	fmt.Println("\n  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12s|\n"+
			"  |behav = %-12s|\n",
		e.CurrentFloor,
		hwelevio.ElevDirToString(e.Dirn),
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