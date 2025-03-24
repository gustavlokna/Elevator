package network

import (
	. "Project/config"
	. "Project/dataenums"
)

func resetHallCalls() [NFloors][NButtons]ButtonState {
	var hallOrderList [NFloors][NButtons]ButtonState
	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			hallOrderList[floor][btn] = Initial
		}
	}
	return hallOrderList
}

func initializeElevatorList() [NElevators]HRAElevState {
	var list [NElevators]HRAElevState
	for elevator := 0; elevator < NElevators; elevator++ {
		list[elevator] = HRAElevState{
			Behaviour:   "idle",
			Floor:       0,
			Direction:   "stop",
			CabRequests: make([]bool, NFloors),
		}
	}
	return list
}

func allAcknowledged(ackMap [NElevators]bool, aliveList [NElevators]bool, id int) bool {
	for elevator := 0; elevator < NElevators; elevator++ {
		if id == elevator {
			continue
		}
		if aliveList[elevator] && !ackMap[elevator] {
			return false
		}
	}
	return true
}
