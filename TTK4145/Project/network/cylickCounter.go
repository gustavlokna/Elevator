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
				nodeOrder := hallOrderList[node][floor][btn]
				if !aliveList[node] || node == myID || nodeOrder == Initial{
					continue
				}
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
	if myOrder == Initial {
        return nodeOrder
    }

	switch myOrder {
	case Idle:
		switch nodeOrder {

		case Idle:
			//print("penis")
			myOrder = Idle
		case ButtonPressed:
			//print("hallo")
			myOrder = ButtonPressed
		case OrderAssigned: 
			//print("kuk")
			// Error should not happen
			// was myOrder = IDLE
			myOrder = ButtonPressed
		case OrderComplete: 
			//print("pikk")
			myOrder = Idle
	}  
	case ButtonPressed: 
		switch nodeOrder {
		case Idle:
			//print("fack")
			myOrder = ButtonPressed 
		case ButtonPressed:
			//print("tissefant")
			myOrder = OrderAssigned
		case OrderAssigned: 
			//print("tullbal")
			myOrder = OrderAssigned
		case OrderComplete: 
			//print("hva faen ButtonPressed, OrderComplete  ")
			// Error should not happen 
			// I have set it to ButtonPressed
			// such that we waith for that elevator to catch up 
			myOrder = ButtonPressed

	}  

	case OrderAssigned:
		switch nodeOrder {
		case Idle:
			//print("hva faen, OrderAssigned Idle")
			// Error Should not happen
			myOrder = ButtonPressed
		case ButtonPressed:
			//print("hei, OrderAssigned ButtonPressed ")
			myOrder = OrderAssigned
		case OrderAssigned: 
			//print("hei, OrderAssigned OrderAssigned")
			myOrder = OrderAssigned
		case OrderComplete: 
		//print("hva faen, OrderAssigned OrderComplete")
			myOrder = OrderComplete

	}  

	case OrderComplete:
		switch nodeOrder {
		case Idle:
			print("hva faen, OrderComplete Idle")
			myOrder = Idle
		case ButtonPressed:
			print("hva faen, OrderComplete ButtonPressed")
			myOrder = ButtonPressed
		case OrderAssigned: 
			print("hva faen, OrderComplete OrderAssigned")
			myOrder = OrderComplete
		case OrderComplete: 
			print("hva faen, OrderComplete OrderComplete")
			myOrder = Idle
	}  
	}
	//print("myOrder: ", myOrder)
	return myOrder
}
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


