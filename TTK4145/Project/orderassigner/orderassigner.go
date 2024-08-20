package orderassigner

import (
	"time"
)
func OrderAssigner(
	fromOrderAssignerChannel chan<- bool,
	toOrderAssignerChannel <-chan bool,
	lifelineChannel chan<- bool,
	nodeID int,
) {
	for {
		select {
		case <-toOrderAssignerChannel: 
			print("gott the new order")
		default:
			time.Sleep(10 * time.Millisecond) // Prevent busy loop
		}
	}
}