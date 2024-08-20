package orderassigner

import (
	. "Project/dataenums"
)

func buttonPressed(btnEvent ButtonEvent) [NFloors][NButtons]bool {
	var buttons [NFloors][NButtons]bool
	buttons[btnEvent.Floor][btnEvent.Button] = true
	return buttons
}