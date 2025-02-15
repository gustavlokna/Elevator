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
			origState := orders[myID][f][b]
			myState := origState

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
				if e != myID && alive[e]  && orders[e][f][b] != Initial {
					peers = append(peers, orders[e][f][b])
				}
			}

			// Attempt a valid transition.
			switch origState {
			case Idle:
				if allIn(peers, Idle, ButtonPressed) && anyIs(peers, ButtonPressed) {
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
				if allIn(peers, OrderComplete, Idle) {
					myState = Idle
				}
			}

			// If no valid transition occurred, check for an illegal combination.
			// For each state, we define the allowed peer states:
			// Idle: allowed peers are Idle or ButtonPressed.
			// ButtonPressed: allowed peers are ButtonPressed or OrderAssigned.
			// OrderAssigned: allowed peers are OrderAssigned or OrderComplete.
			// OrderComplete: allowed peers are OrderComplete or Idle.
			if myState == origState {
				switch origState {
				case Idle:
					if !allIn(peers, Idle, ButtonPressed) || !allIn(peers, Idle, OrderComplete) {
						myState = Initial // Illegal combination detected.
					}
				case ButtonPressed:
					if !allIn(peers, ButtonPressed, OrderAssigned) || !allIn(peers, Idle, ButtonPressed) {
						myState = Initial
					}
				case OrderAssigned:
					if !allIn(peers, OrderAssigned, OrderComplete) || !allIn(peers, OrderAssigned, ButtonPressed){
						myState = Initial
					}
				case OrderComplete:
					if !allIn(peers, OrderComplete, Idle)  || !allIn(peers, OrderAssigned, OrderComplete) {
						myState = Initial
					}
				}
			}

			orders[myID][f][b] = myState
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
