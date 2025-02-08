package lights

import (
	. "Project/dataenums"
	"Project/hwelevio"
)



func LightsHandler(
	orderList <-chan [NFloors][NButtons]ButtonState, 
	payloadFromDriver <-chan PayloadFromDriver,
) {
	for {
		select {
		case payload := <-payloadFromDriver:
			hwelevio.SetFloorIndicator(payload.CurrentFloor)
			hwelevio.SetDoorOpenLamp(payload.DoorLight)

		case orders := <-orderList:
			for floor := 0; floor < NFloors; floor++ {
				for button := 0; button < NButtons; button++ {
					hwelevio.SetButtonLamp(Button(button), floor, orders[floor][button]==OrderAssigned)
				}
			}
		}
	}
}

/*
func LightsHandler(
	orderList <-chan [NFloors][NButtons]ButtonState, 
	payloadFromDriver <-chan PayloadFromDriver,
) {
	for {
		select {
		case payload := <-payloadFromDriver:
			hwelevio.SetFloorIndicator(payload.CurrentFloor)
			hwelevio.SetDoorOpenLamp(payload.DoorLight)

		case orders := <-orderList:
			for floor := 0; floor < NFloors; floor++ {
				for button := 0; button < NButtons; button++ {
					// Light is ON if the order is either OrderAssigned or OrderComplete
					if orders[floor][button] == OrderAssigned || orders[floor][button] == OrderComplete {
						hwelevio.SetButtonLamp(Button(button), floor, true)
					} else {
						hwelevio.SetButtonLamp(Button(button), floor, false)
					}
				}
			}
		}
	}
}
*/