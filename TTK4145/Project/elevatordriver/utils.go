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

func shouldClearImmediately(e Elevator, btnEvent ButtonEvent) bool {
	btn_floor  := btnEvent.Floor
	btn_type := btnEvent.Button
	switch e.Config.ClearRequestVariant {
	case CRVAll:
		return e.CurrentFloor == btn_floor

	case CRVInDirn:
		return e.CurrentFloor == btn_floor &&
			((e.Dirn == DirUp && btn_type == BHallUp) ||
				(e.Dirn == DirDown && btn_type == BHallDown) ||
				e.Dirn == DirStop ||
				btn_type == BCab)
	default:
		return false
	}
}