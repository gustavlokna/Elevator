package elevatordriver


import (
	. "Project/dataenums"
)


func initelevator() Elevator {
	return Elevator{
		CurrentFloor:     -1,
		CurrentBehaviour: EBIdle,
		Config: ElevatorConfig{
			ClearRequestVariant: ClearRequestVariantConfig,
			DoorOpenDurationS:   DoorOpenDurationSConfig,
		},
	}
}