package elevatordriver


import (
	. "Project/dataenums"
)


func initelevator() Elevator {
	return Elevator{
		CurrentFloor:     -1, //TODO need to go down one floor
		Dirn: MDDown,
		CurrentBehaviour: EBIdle,
		Config: ElevatorConfig{
			ClearRequestVariant: ClearRequestVariantConfig,
			DoorOpenDurationS:   DoorOpenDurationSConfig,
		},
	}
}

/*
func init() {
	elevator = elev.ElevatorInit()
	nodeIP, _ = local.GetIP()
	SetAllLights()
	elevio.OutputDevice.DoorLight(false)
	elevio.OutputDevice.StopButtonLight(false)
}
*/