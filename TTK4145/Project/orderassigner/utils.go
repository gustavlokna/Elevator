package orderassigner
import (
	. "Project/dataenums"
	 "fmt"
)


//these are double up and also in utils for elevatordriver. 
//should be deleted from there when project is complete 
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






func PrintHRAInput(payload PayloadFromassignerToNetwork) {
    fmt.Println("HRAInput:")
    fmt.Println("Hall Requests:")
    for floor, requests := range payload.HallRequests {
        fmt.Printf("  Floor %d: [Up: %s, Down: %s]\n",
            floor,
            buttonStateToString(requests[BHallUp]),
            buttonStateToString(requests[BHallDown]),
        )
    }

    fmt.Println("States:")
    for elevator, state := range payload.States {
        fmt.Printf("  Elevator: %s\n", elevator)
        fmt.Printf("    Behavior: %s\n", state.Behavior)
        fmt.Printf("    Floor: %d\n", state.Floor)
        fmt.Printf("    Direction: %s\n", state.Direction)
        fmt.Printf("    Cab Requests: %v\n", state.CabRequests)
    }
}

func buttonStateToString(state ButtonState) string {
    switch state {
    case Idle:
        return "Idle"
    case ButtonPressed:
        return "ButtonPressed"
    case OrderAssigned:
        return "OrderAssigned"
    case OrderComplete:
        return "OrderComplete"
    default:
        return "Unknown"
    }
}
