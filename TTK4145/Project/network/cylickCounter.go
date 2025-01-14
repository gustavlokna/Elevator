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
    // I want to only progress if all are equal 
    //or i am behind unless transition from complete to idle 
	switch myOrder {
	case Idle:
		switch nodeOrder {
		case Idle:
			myOrder = Idle
		case ButtonPressed:
			myOrder = ButtonPressed
		case OrderAssigned: 
			// Error should not happen
			myOrder = Idle
		case OrderComplete: 
			myOrder = Idle

	}  
	case ButtonPressed: 
		switch nodeOrder {
		case Idle:
			myOrder = ButtonPressed 
		case ButtonPressed:
			myOrder = OrderAssigned
		case OrderAssigned: 
			myOrder = OrderAssigned
		case OrderComplete: 
			// Error should not happen 
			// I have set it to ButtonPressed
			// such that we waith for that elevator to catch up 
			myOrder = ButtonPressed

	}  

	case OrderAssigned:
		switch nodeOrder {
		case Idle:
			// Error Should not happen
			myOrder = ButtonPressed
		case ButtonPressed:
			myOrder = OrderAssigned
		case OrderAssigned: 
			myOrder = OrderAssigned
		case OrderComplete: 
			myOrder = OrderComplete

	}  

	case OrderComplete:
		switch nodeOrder {
		case Idle:
			myOrder = Idle
		case ButtonPressed:
			myOrder = ButtonPressed
		case OrderAssigned: 
			myOrder = OrderComplete
		case OrderComplete: 
			myOrder = Idle
	}  
	}
	return myOrder

/*
func cyclicLogic(myOrder ButtonState, 
    nodeOrder ButtonState,
    ) ButtonState {
    // I want to only progress if all are equal 
    //or i am behind unless transition from complete to idle 
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

    // I CANNOT progress from order assigned before i get msg form assigner! 

	return (myOrder + 1) % (OrderComplete + 1)
}
*/


