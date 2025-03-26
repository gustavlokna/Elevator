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
			currentOrder := orders[myID][floor][btn]
			updatedOrder := currentOrder

			if updatedOrder == Initial {
				for elevator := 0; elevator < NElevators; elevator++ {
					if elevator != myID && orders[elevator][floor][btn] != Initial {
						updatedOrder = orders[elevator][floor][btn]
						break
					}
				}
				orders[myID][floor][btn] = updatedOrder
				continue
			}

			var peers []ButtonState
			for elevator := 0; elevator < NElevators; elevator++ {
				if elevator != myID && orders[elevator][floor][btn] != Initial {
					peers = append(peers, orders[elevator][floor][btn])
				}
			}

			// Attempt a valid transition.
			switch currentOrder {
			case Standby:
				if allIn(peers, Standby, ButtonPressed) && anyIs(peers, ButtonPressed) {
					updatedOrder = ButtonPressed
				}
			case ButtonPressed:
				if allIn(peers, ButtonPressed, OrderAssigned) {
					updatedOrder = OrderAssigned
				}
			case OrderAssigned:
				if allIn(peers, OrderAssigned, OrderComplete) && anyIs(peers, OrderComplete) {
					updatedOrder = OrderComplete
				}
			case OrderComplete:
				if allIn(peers, OrderComplete, Standby) {
					updatedOrder = Standby
				}
			}

			// If no valid transition occurred, check for an illegal combination.
			if updatedOrder == currentOrder {
				switch currentOrder {
				case Standby:
					if !(allIn(peers, Standby, ButtonPressed) || allIn(peers, Standby, OrderComplete)) {
						updatedOrder = Initial
					}
				case ButtonPressed:
					if !(allIn(peers, ButtonPressed, OrderAssigned) || allIn(peers, Standby, ButtonPressed)) {
						updatedOrder = Initial
					}
				case OrderAssigned:
					if !(allIn(peers, OrderAssigned, OrderComplete) || allIn(peers, OrderAssigned, ButtonPressed)) {
						updatedOrder = Initial
					}
				case OrderComplete:
					if !(allIn(peers, OrderComplete, Standby) || allIn(peers, OrderAssigned, OrderComplete)) {
						updatedOrder = Initial
					}
				}
			}

			orders[myID][floor][btn] = updatedOrder
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
