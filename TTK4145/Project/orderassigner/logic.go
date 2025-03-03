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
			orderList[floor][BCab] = PayloadFromNetworkToAssigner.ElevatorList[nodeID].CabRequests[floor]
		}
		return orderList
	}
	hraInput := convertPayloadToHRAInput(PayloadFromNetworkToAssigner, nodeID)

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
			orderList[floor][BCab] = hraInput.States[elevatorID].CabRequests[floor]
		}
	}
	return orderList
}

func convertPayloadToHRAInput(payload PayloadFromNetworkToAssigner, nodeID int) HRAInput {

	hraInput := initHRAInput()
	for i, alive := range payload.AliveList {

		if alive {
			elevatorID := fmt.Sprintf("elevator_%d", i) 
			hraInput.States[elevatorID] = payload.ElevatorList[i]
		}
	}
	for floor := 0; floor < NFloors; floor++ {
		for btn := BHallUp; btn <= BHallDown; btn++ {
			allAssigned := true
			for i, alive := range payload.AliveList {
				if alive && payload.HallOrderList[i][floor][btn] != OrderAssigned {
					allAssigned = false
					break
				}
			}
			hraInput.HallRequests[floor][btn] = allAssigned
		}
	}

	return hraInput
}

func initHRAInput() HRAInput {
	hraInput := HRAInput{
		HallRequests: [NFloors][NButtons - 1]bool{},
		States:       make(map[string]HRAElevState),
	}
	return hraInput
}

func initPayloadToNetwork() PayloadFromassignerToNetwork {
	payloadFromassignerToNetwork := PayloadFromassignerToNetwork{
		HallRequests: [NFloors][NButtons]ButtonState{},
		States:       make(map[string]HRAElevState),
	}
	return payloadFromassignerToNetwork
}

func handlePayloadFromElevator(fromElevator PayloadFromElevator,
	toNetwork PayloadFromassignerToNetwork, nodeID string) PayloadFromassignerToNetwork {

	behavior, direction, cabRequests := convertElevatorState(fromElevator.Elevator)
	toNetwork.States[nodeID] = HRAElevState{
		Behavior:    behavior,
		Floor:       fromElevator.Elevator.CurrentFloor,
		Direction:   direction,
		CabRequests: cabRequests,
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
			incomingState := netPayload.HallOrderList[nodeID][f][b]
			localState := payload.HallRequests[f][b]
			
			if localState == OrderComplete && incomingState == OrderAssigned {
				// do nothing; stay OrderComplete.
			} else {
				payload.HallRequests[f][b] = incomingState
			}
		}
	}
	return payload
}

func convertElevatorState(e Elevator) (string, string, []bool) {
	behavior := EBToString(e.CurrentBehaviour)
	direction := ElevDirToString(e.Dirn)
	cabRequests := make([]bool, NFloors)

	for f := 0; f < NFloors; f++ {
		cabRequests[f] = e.Requests[f][BCab]
	}
	return behavior, direction, cabRequests
}

func updateLightStates(payload PayloadFromNetworkToAssigner, myID int) [NFloors][NButtons]ButtonState {
	var updatedLights [NFloors][NButtons]ButtonState
	updatedLights = payload.HallOrderList[myID]

	for floor := 0; floor < NFloors; floor++ {
		if payload.ElevatorList[myID].CabRequests[floor] {
			updatedLights[floor][BCab] = OrderAssigned
		}
	}
	return updatedLights
}
