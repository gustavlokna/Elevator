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
		}
	}

	return orderList
}
