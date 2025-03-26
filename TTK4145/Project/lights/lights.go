package lights

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/hwelevio"
)

func Lights(
	orderLights <-chan [NFloors][NButtons]ButtonState,
	elevatorLights <-chan FromDriverToLight,
) {
	for {
		select {
		case elevator := <-elevatorLights:
			hwelevio.SetFloorIndicator(elevator.CurrentFloor)
			hwelevio.SetDoorOpenLamp(elevator.DoorLight)

		case orders := <-orderLights:
			for floor := 0; floor < NFloors; floor++ {
				for btn := 0; btn < NButtons; btn++ {
					hwelevio.SetButtonLamp(Button(btn), floor, orders[floor][btn] == OrderAssigned)
				}
			}
		}
	}
}
