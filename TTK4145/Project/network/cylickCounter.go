package network

import (
	. "Project/dataenums"
)

func cyclicCounter(
	hallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState,
	aliveList [NUM_ELEVATORS]bool,
	myID int,
) [NUM_ELEVATORS][NFloors][NButtons]ButtonState {
	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			myOrder := hallOrderList[myID][floor][btn]
			for node := 0; node < NUM_ELEVATORS; node++ {
				if !aliveList[node] {
					continue
				}
				nodeOrder := hallOrderList[node][floor][btn]
				myOrder = cyclicLogic(myOrder, nodeOrder)
			}
			hallOrderList[myID][floor][btn] = myOrder
		}
	}
	return hallOrderList
}

func cyclicLogic(myOrder ButtonState, 
    nodeOrder ButtonState,
    ) ButtonState {
	if myOrder == Idle {
		if nodeOrder != Idle && nodeOrder != OrderComplete {
			// Progress to next state if another elevator is ahead
			return (nodeOrder + 1) % (OrderComplete + 1)
		}
	}
	if nodeOrder < myOrder && nodeOrder != Idle {
		// Stay in the current state if another elevator is behind
		return myOrder
	}
	// Progress to the next state
	return (myOrder + 1) % (OrderComplete + 1)
}

	


