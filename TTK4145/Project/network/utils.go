package network

import (
	. "Project/dataenums"
	. "Project/config"
	"fmt"
)

func resetHallCalls() [NFloors][NButtons]ButtonState {
	var hallOrderList [NFloors][NButtons]ButtonState
	for floor := 0; floor < NFloors; floor++ {
		for button := 0; button < NButtons; button++ {
			hallOrderList[floor][button] = Initial
		}
	}
	return hallOrderList
}

func initializeElevatorList() [NElevators]HRAElevState {
	var list [NElevators]HRAElevState
	for i := 0; i < NElevators; i++ {
		list[i] = HRAElevState{
			Behaviour:   "idle",
			Floor:       0,
			Direction:   "stop",
			CabRequests: make([]bool, NFloors),
		}
	}
	return list
}

// Helper function to get a readable string for ButtonState
func buttonStateToString(state ButtonState) string {
	switch state {
	case Inactive:
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
func printHallOrderList(hallOrderList [NElevators][NFloors][NButtons]ButtonState) {
	for elevator := 0; elevator < NElevators; elevator++ {
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

func printElevatorList(elevatorList [NElevators]HRAElevState) {
	fmt.Println("Elevator List:")
	for i, state := range elevatorList {
		fmt.Printf("  Elevator %d:\n", i)
		fmt.Printf("    Behavior: %s\n", state.Behaviour)
		fmt.Printf("    Floor: %d\n", state.Floor)
		fmt.Printf("    Direction: %s\n", state.Direction)
		fmt.Printf("    Cab Requests: %v\n", state.CabRequests)
	}
}
