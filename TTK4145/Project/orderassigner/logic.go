package orderassigner

import (
	. "Project/dataenums"
	//"os"
	"encoding/json"
	"os/exec"
)

func worldViewsAlign(hraInput HRAInput) bool{
	WorldView := true
	//logic to check if all counters are equal
	return WorldView
}

func assignOrders(hraInput HRAInput, elevatorName string) [NFloors][NButtons]bool {
	var orderList [NFloors][NButtons]bool

	if len(hraInput.States) == 0 {
		print("HRAInput.States is empty, skipping order assignment")
		return orderList
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

	if orders, ok := output[elevatorName]; ok {
		for floor := 0; floor < NFloors && floor < len(orders); floor++ {
			for btn := BHallUp; btn < BCab; btn++ {
				orderList[floor][btn] = orders[floor][btn]
			}
			orderList[floor][BCab] = hraInput.States[elevatorName].CabRequests[floor]
		}
	}
	
	return orderList
}

func InitialiseHRAInput() HRAInput {
	hraInput := HRAInput{
		HallRequests: make([][2]bool, NFloors),
		States:       make(map[string]HRAElevState),
		CounterHallRequests : make([][2]int, NFloors),
	}
	return hraInput
}
func handlePayloadFromElevator(hraInput HRAInput, e Elevator,
	elevatorName string) HRAInput{
	behavior, direction, cabRequests := convertElevatorState(e)
	hraInput.States[elevatorName] = HRAElevState{
		Behavior:    behavior,
		Floor:       e.CurrentFloor,
		Direction:   direction,
		CabRequests: cabRequests,
	}
	return hraInput
}

func handlePayloadFromNetwork(localHRA HRAInput, Incoming HRAInput)HRAInput{
	return Incoming
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

