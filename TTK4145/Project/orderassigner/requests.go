package orderassigner

import (
	. "Project/dataenums"
)

func buttonPressed(btnEvent ButtonEvent) [NFloors][NButtons]bool {
	var buttons [NFloors][NButtons]bool
	buttons[btnEvent.Floor][btnEvent.Button] = true
	return buttons
}
func buttonAlreadyActive(HRAInput HRAInput,elevatorName string,
	btnEvent ButtonEvent) bool {
	switch btnEvent.Button {
	case BHallUp:
		return HRAInput.HallRequests[btnEvent.Floor][BHallUp]
	case BHallDown:
		return HRAInput.HallRequests[btnEvent.Floor][BHallDown]
	case BCab:
		return HRAInput.States[elevatorName].CabRequests[btnEvent.Floor]
	}
	return false
}
