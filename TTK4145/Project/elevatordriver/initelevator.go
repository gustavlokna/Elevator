package elevatordriver

import (
	. "Project/dataenums"
)

func initelevator() Elevator {
	// Default initialization if file doesn't exist or decoding fails
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
