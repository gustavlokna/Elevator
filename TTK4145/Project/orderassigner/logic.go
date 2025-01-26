package orderassigner

import (
	. "Project/dataenums"
	//"os"
	"encoding/json"
	"os/exec"
	"fmt"
)


func assignOrders(PayloadFromNetworkToAssigner PayloadFromNetworkToAssigner, 
	nodeID int) [NFloors][NButtons]bool {
	var orderList [NFloors][NButtons]bool

	if !(PayloadFromNetworkToAssigner.AliveList[nodeID]) {
		print("local elevator not alive")
		return orderList
	}
	hraInput := convertPayloadToHRAInput(PayloadFromNetworkToAssigner, nodeID) 

	jsonBytes, err := json.Marshal(hraInput)
	if err != nil {
		print("Failed to marshal HRAInput: %v\n", err)
		return orderList
	}
	fmt.Printf("Serialized HRAInput: %s\n", string(jsonBytes))

	
	// TODO SOME LOGIC TO MAKE THE ButtonState to bool 
	fmt.Println("all god here")
	ret, err := exec.Command("hall_request_assigner", "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		print("exec.Command error: %v\nOutput: %s\n", err, string(ret))
		return orderList
	}
	
	fmt.Printf("Raw output: %s\n", string(ret))

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
			orderList[floor][BCab] = hraInput.States[elevatorID].CabRequests[floor]
		}
	}
	
	return orderList
}


func convertPayloadToHRAInput(payload PayloadFromNetworkToAssigner, nodeID int) HRAInput {
	
	hraInput := InitialiseHRAInput()
	for i, alive := range payload.AliveList {

		if alive {

			elevatorID := fmt.Sprintf("elevator_%d", i) // Convert index to string key
			hraInput.States[elevatorID] = payload.ElevatorList[i]
		}
	}
	//print("printer hra")
	//PrintHraInput(hraInput)
	// Iterate over all floors and buttons
	for floor := 0; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BHallDown; btn++ { 
			allAssigned := true
	
			// Check all alive elevators for the specific button and floor
			for i, alive := range payload.AliveList {
				if alive {
					if payload.HallOrderList[i][floor][btn] != OrderAssigned {
						allAssigned = false
						break
					}
				}
			}
			if allAssigned {
				hraInput.HallRequests[floor][btn] = true
			} else {
				hraInput.HallRequests[floor][btn] = false 
			}
		}
	}
	
	return hraInput
}

// TODO DO NOT LIKE HAVING BOTH THESE TWO 
func InitialiseHRAInput() HRAInput {
	hraInput := HRAInput{
		HallRequests: [NFloors][NButtons-1]bool{},
		States:       make(map[string]HRAElevState),
	}
	return hraInput
}

func InitialisePayloadFromassignerToNetwork() PayloadFromassignerToNetwork {
	payloadFromassignerToNetwork := PayloadFromassignerToNetwork{
		HallRequests: [NFloors][NButtons]ButtonState{},
		States:       make(map[string]HRAElevState),
	}
	return payloadFromassignerToNetwork
}


func handlePayloadFromElevator(fromElevator  PayloadFromElevator, 
	toNetwork PayloadFromassignerToNetwork, nodeID string ) PayloadFromassignerToNetwork{

	behavior, direction, cabRequests := convertElevatorState(fromElevator.Elevator)
	toNetwork.States[nodeID] = HRAElevState{
		Behavior:    behavior,
		Floor:       fromElevator.Elevator.CurrentFloor,
		Direction:   direction,
		CabRequests: cabRequests,
	}
	// Todo we set the cabrequests twice. this is because we have to make them :=0
	// do this smarter
	toNetwork = orderComplete(toNetwork, nodeID, fromElevator.CompletedOrders)

	return toNetwork
}



func handlePayloadFromNetwork(payload PayloadFromassignerToNetwork, 
	PayloadFromNetwork PayloadFromNetworkToAssigner,
	nodeID int)PayloadFromassignerToNetwork{
	payload.HallRequests = PayloadFromNetwork.HallOrderList[nodeID]
	//payload.States[nodeID] = PayloadFromNetwork.ElevatorList[nodeID]
	//payload.States = PayloadFromNetwork.ElevatorList
	return payload
}

func convertElevatorState(e Elevator) (string, string, []bool) {
	behavior := EBToString(e.CurrentBehaviour)
	direction := ElevDirToString(e.Dirn)

	// Convert cab requests
	cabRequests := make([]bool, NFloors)
	for f := 0; f < NFloors; f++ {
		cabRequests[f] = e.Requests[f][BCab]
	}
	return behavior, direction, cabRequests
}



