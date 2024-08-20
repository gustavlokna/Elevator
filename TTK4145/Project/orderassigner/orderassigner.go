package orderassigner

import (
	. "Project/dataenums"
	//"Project/elevatordriver"
	"Project/hwelevio"
	"time"
)

// buttonlights must be set inside this 
func OrderAssigner(
	fromOrderAssignerChannel chan<- [NFloors][NButtons]bool,
	toOrderAssignerChannel <-chan Elevator,
	lifelineChannel chan<- bool,
	nodeID int,
) {
	drv_buttons := make(chan ButtonEvent)
	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btn:= <-drv_buttons:
			fromOrderAssignerChannel <- buttonPressed(btn)
		case <-toOrderAssignerChannel:
			print("elevator was changed")
		default:
			time.Sleep(10 * time.Millisecond) // Prevent busy loop
		}
		
	}
}