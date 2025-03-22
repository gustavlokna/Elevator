package driver

import (
	. "Project/dataenums"
)

func initelevator() Elevator {
	return Elevator{
		CurrentFloor:     -1,
		Dirn:             MDDown,
		CurrentBehaviour: Moving,
		ActiveSatus:      true, 
	}
}
