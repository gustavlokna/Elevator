package network

import (
	. "Project/dataenums"
    "fmt"
)
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func resetHallCalls() [NFloors][NButtons]ButtonState {
	var hallOrderList [NFloors][NButtons]ButtonState
	for floor := 0; floor < NFloors; floor++ {
		for button := 0; button < NButtons; button++ {
			hallOrderList[floor][button] = Initial
		}
	}
	return hallOrderList
}

func initializeElevatorList() [NUM_ELEVATORS]HRAElevState {
	var list [NUM_ELEVATORS]HRAElevState
	for i := 0; i < NUM_ELEVATORS; i++ {
		list[i] = HRAElevState{
			Behavior:    "EBIdle",
			Floor:       0,
			Direction:   "MDStop",
			CabRequests: make([]bool, NFloors),
		}
	}
	return list
}

// Helper function to check if HRAInput is empty
// TODO ALL BELOW IS GARBAGE ? 
func isEmptyHRAInput(input HRAInput) bool {
	return len(input.HallRequests) == 0 && len(input.States) == 0
}
// Helper function to get a readable string for ButtonState
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
func printHallOrderList(hallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState) {
	for elevator := 0; elevator < NUM_ELEVATORS; elevator++ {
		fmt.Printf("Elevator %d:\n", elevator)
		fmt.Println("Hall Requests:")
		for floor := 0; floor < NFloors; floor++ {
			fmt.Printf("  Floor %d: [", floor)
			fmt.Printf("Up: %s, ", buttonStateToString(hallOrderList[elevator][floor][BHallUp]))
			fmt.Printf("Down: %s", buttonStateToString(hallOrderList[elevator][floor][BHallDown]))
			fmt.Println("]")
		}
		fmt.Println("States:")
		// Add any additional states or information related to the elevator if needed.
		fmt.Println()
	}
}

func printElevatorList(elevatorList [NUM_ELEVATORS]HRAElevState) {
	fmt.Println("Elevator List:")
	for i, state := range elevatorList {
		fmt.Printf("  Elevator %d:\n", i)
		fmt.Printf("    Behavior: %s\n", state.Behavior)
		fmt.Printf("    Floor: %d\n", state.Floor)
		fmt.Printf("    Direction: %s\n", state.Direction)
		fmt.Printf("    Cab Requests: %v\n", state.CabRequests)
	}
}
