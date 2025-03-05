package orderassigner

import (
	. "Project/dataenums"
	"encoding/json"
	"fmt"
	"os/exec"
)

func assignOrders(PayloadFromNetworkToAssigner PayloadFromNetworkToAssigner,
	nodeID int) [NFloors][NButtons]bool {
	var orderList [NFloors][NButtons]bool

	if !(PayloadFromNetworkToAssigner.AliveList[nodeID]) {
		print("local elevator not alive")
		for floor := 0; floor < NFloors; floor++ {
			orderList[floor][BCab] = (PayloadFromNetworkToAssigner.Orders[nodeID][floor][BCab]==OrderAssigned)
		}
		return orderList
	}
	hraInput := HRAInput{
		HallRequests: [NFloors][NButtons - 1]bool{},
		States:       make(map[string]HRAElevState),
	}

	// Process each elevator
	for i, alive := range PayloadFromNetworkToAssigner.AliveList {
		if !alive {
			continue
		}
		elevatorID := fmt.Sprintf("elevator_%d", i)

		elevState := PayloadFromNetworkToAssigner.ElevatorList[i]
		elevState.CabRequests = make([]bool, NFloors)
		for floor := 0; floor < NFloors; floor++ {
			elevState.CabRequests[floor] = (PayloadFromNetworkToAssigner.Orders[i][floor][BCab] == OrderAssigned)
		}
		hraInput.States[elevatorID] = elevState
	}
	for floor := 0; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BHallDown; btn++ {
			allAssigned := true
			for i, alive := range PayloadFromNetworkToAssigner.AliveList {
				if alive && PayloadFromNetworkToAssigner.Orders[i][floor][btn] != OrderAssigned {
					allAssigned = false
					break
				}
			}
			hraInput.HallRequests[floor][btn] = allAssigned
		}
	}

	jsonBytes, err := json.Marshal(hraInput)
	if err != nil {
		print("Failed to marshal HRAInput: %v\n", err)
		return orderList
	}
	ret, err := exec.Command("hall_request_assigner", "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		print("exec.Command error: %v\nOutput: %s\n", err, string(ret))
		return orderList
	}

	output := make(map[string][][2]bool)
	if err := json.Unmarshal(ret, &output); err != nil {
		print("json.Unmarshal error: %v\n", err)
		return orderList
	}
	elevatorID := fmt.Sprintf("elevator_%d", nodeID)
	if orders, ok := output[elevatorID]; ok {
		for floor := 0; floor < NFloors && floor < len(orders); floor++ {
			for btn := BHallUp; btn < BCab; btn++ {
				orderList[floor][btn] = orders[floor][btn]
			}
			orderList[floor][BCab] = (PayloadFromNetworkToAssigner.Orders[nodeID][floor][BCab]==OrderAssigned)
		}
	}
	return orderList
}




func initPayloadToNetwork() PayloadFromassignerToNetwork {
	payloadFromassignerToNetwork := PayloadFromassignerToNetwork{
		Orders: [NFloors][NButtons]ButtonState{},
		States:       make(map[string]HRAElevState),
	}
	return payloadFromassignerToNetwork
}

func handlePayloadFromElevator(fromElevator PayloadFromElevator,
	toNetwork PayloadFromassignerToNetwork, nodeID string) PayloadFromassignerToNetwork {

	toNetwork.States[nodeID] = HRAElevState{
		Behavior:    EBToString(fromElevator.Elevator.CurrentBehaviour),
		Floor:       fromElevator.Elevator.CurrentFloor,
		Direction:   ElevDirToString(fromElevator.Elevator.Dirn),
	}
	toNetwork.ActiveSatus = fromElevator.Elevator.ActiveSatus
	toNetwork = orderComplete(toNetwork, nodeID, fromElevator.CompletedOrders)

	return toNetwork
}


func handlePayloadFromNetwork(
	payload PayloadFromassignerToNetwork,
	netPayload PayloadFromNetworkToAssigner,
	nodeID int,
) PayloadFromassignerToNetwork {
	for f := 0; f < NFloors; f++ {
		for b := 0; b < NButtons; b++ {
			incomingState := netPayload.Orders[nodeID][f][b]
			localState := payload.Orders[f][b]
			
			if localState == OrderComplete && incomingState == OrderAssigned {
				// do nothing; stay OrderComplete.
			} else {
				payload.Orders[f][b] = incomingState
			}
			
		}
	}
	return payload
}



func updateLightStates(payload PayloadFromNetworkToAssigner, myID int) [NFloors][NButtons]ButtonState {
	var updatedLights [NFloors][NButtons]ButtonState = payload.Orders[myID]
	// updatedLights = payload.HallOrderList[myID]
	// TOD IS IN HALL LIST ABOVE SO REMOVE 
	/*
	for floor := 0; floor < NFloors; floor++ {
		
		if payload.ElevatorList[myID].CabRequests[floor] {
			updatedLights[floor][BCab] = OrderAssigned
		}
	}
	*/
	return updatedLights
}
