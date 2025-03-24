package lights

import (
	. "Project/dataenums"
	. "Project/config"
	"Project/hwelevio"
)

func LightsHandler(
	orderList <-chan [NFloors][NButtons]ButtonState,
	localLights <-chan FromDriverToLight,
) {
	for {
		select {
		case payload := <- localLights:
			hwelevio.SetFloorIndicator(payload.CurrentFloor)
			hwelevio.SetDoorOpenLamp(payload.DoorLight)

		case orders := <-orderList:
			for floor := 0; floor < NFloors; floor++ {
				for btn := 0; btn < NButtons; btn++ {
					hwelevio.SetButtonLamp(Button(btn), floor, orders[floor][btn] == OrderAssigned)
				}
			}
		}
	}
}
