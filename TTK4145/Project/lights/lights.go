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
				for button := 0; button < NButtons; button++ {
					hwelevio.SetButtonLamp(Button(button), floor, orders[floor][button] == OrderAssigned)
				}
			}
		}
	}
}
