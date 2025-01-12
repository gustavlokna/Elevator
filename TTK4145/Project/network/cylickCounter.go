package network

/*
case msg := <-broadcastReceiverChannel:
    for buttonIdx, hallRequest := range msg.Payload.HallRequests {
        if hallRequest[0] { // Example condition for a button press
            // Update the local state for this button
            globalStates := extractGlobalStates(buttonIdx, msg)
            localState := localStates[nodeID][buttonIdx]
            localStates[nodeID][buttonIdx] = nextState(localState, globalStates, buttonIdx)
        }
    }

    // Check for synchronization for `OrderAssigned` or other states
    for buttonIdx := range msg.Payload.HallRequests {
        if allElevatorsInState(localStates, buttonIdx, OrderAssigned) {
            assignOrder(buttonIdx)
            broadcastTransmissionChannel <- Message{
                SenderId:     nodeID,
                Payload:      createPayload(localStates),
                OnlineStatus: onlineStatus,
            }
        }
    }

    // Forward the message to the order assigner
    messagetoOrderAssignerChannel <- msg
*/






/*

func cyclicCounterLogic(localState int, globalStates map[string]int, button int) int {
	// do nothing if in idle and all others in idle and order complete 
	if localState == Idle  { 
		for _, state := range globalStates {
            if state != Idle && state != OrderComplete {
				// If another elevator is ahead, allow progression
				return (localState + 1) % (OrderComplete + 1)
            }
        }
        return localState
    }

	for _, state := range globalStates {
		if state < localState && state != Idle {
			// If any elevator is behind, stay in the current state
			return localState
		}
	}

	// Progress to the next state if all are synchronized or ahead
	return (localState + 1) % (OrderComplete+1) 
	}


*/