package network

import (
	. "Project/config"
	. "Project/dataenums"
)

func cyclicCounter(
	orders [NElevators][NFloors][NButtons]ButtonState,
	myID int,
) [NElevators][NFloors][NButtons]ButtonState {

	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			origState := orders[myID][floor][btn]
			myState := origState

			if myState == Initial {
				for elevator := 0; elevator < NElevators; elevator++ {
					if elevator != myID && orders[elevator][floor][btn] != Initial {
						myState = orders[elevator][floor][btn]
						break
					}
				}
				orders[myID][floor][btn] = myState
				continue
			}

			var peers []ButtonState
			for elevator := 0; elevator < NElevators; elevator++ {
				if elevator != myID && orders[elevator][floor][btn] != Initial {
					peers = append(peers, orders[elevator][floor][btn])
				}
			}

			// Attempt a valid transition.
			switch origState {
			case Inactive:
				if allIn(peers, Inactive, ButtonPressed) && anyIs(peers, ButtonPressed) {
					myState = ButtonPressed
				}
			case ButtonPressed:
				if allIn(peers, ButtonPressed, OrderAssigned) {
					myState = OrderAssigned
				}
			case OrderAssigned:
				if allIn(peers, OrderAssigned, OrderComplete) && anyIs(peers, OrderComplete) {
					myState = OrderComplete
				}
			case OrderComplete:
				if allIn(peers, OrderComplete, Inactive) {
					myState = Inactive
				}
			}

			// If no valid transition occurred, check for an illegal combination.
			if myState == origState {
				switch origState {
				case Inactive:
					if !(allIn(peers, Inactive, ButtonPressed) || allIn(peers, Inactive, OrderComplete)) {
						myState = Initial
					}
				case ButtonPressed:
					if !(allIn(peers, ButtonPressed, OrderAssigned) || allIn(peers, Inactive, ButtonPressed)) {
						myState = Initial
					}
				case OrderAssigned:
					if !(allIn(peers, OrderAssigned, OrderComplete) || allIn(peers, OrderAssigned, ButtonPressed)) {
						myState = Initial
					}
				case OrderComplete:
					if !(allIn(peers, OrderComplete, Inactive) || allIn(peers, OrderAssigned, OrderComplete)) {
						myState = Initial
					}
				}
			}

			orders[myID][floor][btn] = myState
		}
	}
	return orders
}

func allIn(peers []ButtonState, optA, optB ButtonState) bool {
	for _, p := range peers {
		if p != optA && p != optB {
			return false
		}
	}
	return true
}

func anyIs(peers []ButtonState, target ButtonState) bool {
	for _, p := range peers {
		if p == target {
			return true
		}
	}
	return false
}
