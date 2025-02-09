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

			// Collect peers who are alive
			var peers []ButtonState
			for e := 0; e < NUM_ELEVATORS; e++ {
				if e != myID && alive[e] {
					peers = append(peers, orders[e][f][b])
				}
			}

			switch myState {
			case Idle:
				// If all peers are Idle or ButtonPressed, move to ButtonPressed
				if allIn(peers, Idle, ButtonPressed) {
					myState = ButtonPressed
				}
			case ButtonPressed:
				// If all peers are ButtonPressed or OrderAssigned, move to OrderAssigned
				if allIn(peers, ButtonPressed, OrderAssigned) {
					myState = OrderAssigned
				}
			case OrderAssigned:
				// If all peers are OrderAssigned or OrderComplete, move to OrderComplete
				if allIn(peers, OrderAssigned, OrderComplete) {
					myState = OrderComplete
				}
			case OrderComplete:
				// If all peers are OrderComplete or Idle, remain OrderComplete
				if allIn(peers, OrderComplete, Idle) {
					myState = OrderComplete
				}
			}
			orders[myID][f][b] = myState
		}
	}
	return orders
}

// allIn checks whether every state in peers is either optA or optB.
func allIn(peers []ButtonState, optA, optB ButtonState) bool {
	for _, p := range peers {
		if p != optA && p != optB {
			return false
		}
	}
	return true
}
