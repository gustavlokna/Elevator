package elevatordriver

import (
	. "Project/dataenums"
)
func buttonPressed(elevator Elevator, btnEvent ButtonEvent) Elevator{
	if shouldClearImmediately(elevator, btnEvent) && (elevator.CurrentBehaviour == EBDoorOpen) {
		startTimer(elevator.Config.DoorOpenDurationS)
	} else {
		elevator.Requests[btnEvent.Floor][btnEvent.Button] = true
	}
	return elevator
}