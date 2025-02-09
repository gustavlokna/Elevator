package network

import (
	. "Project/dataenums"
)

func cyclicCounter(
	orders [NUM_ELEVATORS][NFloors][NButtons]ButtonState,
	alive [NUM_ELEVATORS]bool,
	myID int,
) [NUM_ELEVATORS][NFloors][NButtons]ButtonState {

	for f := 0; f < NFloors; f++ {
		for b := 0; b < NButtons; b++ {
			myState := orders[myID][f][b]

			// If I'm Initial, copy the first alive peer’s state that isn't Initial.
			if myState == Initial {
				for e := 0; e < NUM_ELEVATORS; e++ {
					if e != myID && alive[e] && orders[e][f][b] != Initial {
						myState = orders[e][f][b]
						break
					}
				}
				orders[myID][f][b] = myState
				continue
			}

			// Gather other elevators' states if they're alive.
			var peers []ButtonState
			for e := 0; e < NUM_ELEVATORS; e++ {
				if e != myID && alive[e] {
					peers = append(peers, orders[e][f][b])
				}
			}

			switch myState {

			// IDLE -> BUTTON_PRESSED if:
			//    1) All peers are either IDLE or BUTTON_PRESSED
			//    2) At least one peer is BUTTON_PRESSED
			case Idle:
				if allIn(peers, Idle, ButtonPressed) && anyIs(peers, ButtonPressed) {
					myState = ButtonPressed
				}

			// BUTTON_PRESSED -> ORDER_ASSIGNED if:
			//    1) All peers are either BUTTON_PRESSED or ORDER_ASSIGNED
			case ButtonPressed:
				if allIn(peers, ButtonPressed, OrderAssigned) {
					myState = OrderAssigned
				}

			// ORDER_ASSIGNED -> ORDER_COMPLETE if:
			//    1) All peers are either ORDER_ASSIGNED or ORDER_COMPLETE
			//    2) At least one peer is ORDER_COMPLETE
			case OrderAssigned:
				if allIn(peers, OrderAssigned, OrderComplete) && anyIs(peers, OrderComplete) {
					myState = OrderComplete
				}

			// ORDER_COMPLETE -> remain ORDER_COMPLETE if all peers are ORDER_COMPLETE or IDLE
			case OrderComplete:
				if allIn(peers, OrderComplete, Idle) {
					myState = Idle
				}
			}
			orders[myID][f][b] = myState
		}
	}
	return orders
}

// allIn checks whether every peer is either optA or optB.
func allIn(peers []ButtonState, optA, optB ButtonState) bool {
	for _, p := range peers {
		if p != optA && p != optB {
			return false
		}
	}
	return true
}

// anyIs checks if at least one peer matches a specific state.
func anyIs(peers []ButtonState, target ButtonState) bool {
	for _, p := range peers {
		if p == target {
			return true
		}
	}
	return false
}
