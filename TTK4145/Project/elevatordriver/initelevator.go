package elevatordriver

import (
	. "Project/dataenums"
)
// TODO DO W ENED THE CONFIG ? 
func initelevator() Elevator {
	return Elevator{
		CurrentFloor:     -1,
		Dirn:             MDDown,
		CurrentBehaviour: EBIdle,
		ActiveSatus:      true,
		Config: ElevatorConfig{
			ClearRequestVariant: ClearRequestVariantConfig,
			DoorOpenDurationS:   DoorOpenDurationSConfig,
		},
	}
}