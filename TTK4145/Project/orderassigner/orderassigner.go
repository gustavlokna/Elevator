package orderassigner

import (
	. "Project/dataenums"
	"Project/elevatordriver"
	"time"
)
func OrderAssigner(
	fromOrderAssignerChannel chan<- [NFloors][NButtons]bool,
	toOrderAssignerChannel <-chan Elevator,
	lifelineChannel chan<- bool,
	nodeID int,
) {
	for {
		select {
		case el:=<-toOrderAssignerChannel: 
			print("gott the new order")
			elevatordriver.ElevatorPrint(el)
			fromOrderAssignerChannel <- el.Requests
		default:
			time.Sleep(10 * time.Millisecond) // Prevent busy loop
		}
	}
}