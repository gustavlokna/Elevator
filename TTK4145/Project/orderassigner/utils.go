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


func PrintPayloadFromElevator(payload PayloadFromElevator) {
    fmt.Println("PayloadFromElevator:")

    fmt.Println("Elevator:")
    fmt.Printf("  Floor: %d\n", payload.Elevator.CurrentFloor)
    fmt.Printf("  Direction: %s\n", ElevDirToString(payload.Elevator.Dirn))
    fmt.Printf("  Behavior: %s\n", EBToString(payload.Elevator.CurrentBehaviour))
    fmt.Printf("  Active Status: %t\n", payload.Elevator.ActiveSatus)

    fmt.Println("Requests:")
    for floor := NFloors - 1; floor >= 0; floor-- {
        fmt.Printf("  Floor %d: [Up: %t, Down: %t, Cab: %t]\n",
            floor,
            payload.Elevator.Requests[floor][BHallUp],
            payload.Elevator.Requests[floor][BHallDown],
            payload.Elevator.Requests[floor][BCab],
        )
    }

    fmt.Println("Completed Orders:")
    for floor := NFloors - 1; floor >= 0; floor-- {
        fmt.Printf("  Floor %d: [Up: %t, Down: %t, Cab: %t]\n",
            floor,
            payload.CompletedOrders[floor][BHallUp],
            payload.CompletedOrders[floor][BHallDown],
            payload.CompletedOrders[floor][BCab],
        )
    }
}




func PrintPayloadFromassignerToNetwork(payload PayloadFromassignerToNetwork) {
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

func PrintHraInput(hraInput HRAInput) {
	fmt.Println("Hall Requests:")
	for floor, requests := range hraInput.HallRequests {
		fmt.Printf("  Floor %d:\n", floor)
		fmt.Printf("    Up: %s\n", boolToString(requests[BHallUp]))
		fmt.Printf("    Down: %s\n", boolToString(requests[BHallDown]))
	}

	fmt.Println("States:")
	for elevator, state := range hraInput.States {
		fmt.Printf("  Elevator: %s\n", elevator)
		fmt.Printf("    Behavior: %s\n", state.Behavior)
		fmt.Printf("    Floor: %d\n", state.Floor)
		fmt.Printf("    Direction: %s\n", state.Direction)
		fmt.Printf("    Cab Requests: %v\n", state.CabRequests)
	}
}

// Helper: Convert bool to string
func boolToString(state bool) string {
	if state {
		return "True"
	}
	return "False"
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
