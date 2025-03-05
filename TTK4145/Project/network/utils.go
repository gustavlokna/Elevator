package network

import (
	. "Project/dataenums"
	"fmt"
	"sync"
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

var printLock sync.Mutex

func diffByOne(a, b []bool) bool {
	d := 0
	for i := range a {
		if a[i] != b[i] {
			d++
		}
	}
	return d == 1
}

func printLists(list1, list2 []bool) {
	printLock.Lock()
	defer printLock.Unlock()

	fmt.Println("Index\tList1\tList2")
	for i := 0; i < len(list1); i++ {
		fmt.Printf("%d\t%t\t%t\n", i, list1[i], list2[i])
	}
}
