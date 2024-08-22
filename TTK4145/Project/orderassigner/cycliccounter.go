package orderassigner
/*
import (
	. "Project/dataenums"
)

// Define a cyclic counter type for elevator orders
type CCounter int

const (
	Completed CCounter = iota
	NoOrder
	UnConfirmed
	Confirmed
)

type CyclicElevator struct {
	ID       int
	Requests [NFloors][NButtons]CCounter
}


func AddElevator(existingElevs []CyclicElevator, 
	newElev CyclicElevator) []CyclicElevator {
	updatedElevs := append(existingElevs, newElev)
	return updatedElevs
}

// RemoveElevator removes a CyclicElevator from the slice based on its ID
func RemoveElevator(existingElevs []CyclicElevator, 
	idToRemove int) []CyclicElevator {
	var updatedElevs []CyclicElevator
	for _, elev := range existingElevs {
		if elev.ID != idToRemove {
			updatedElevs = append(updatedElevs, elev)
		}
	}
	return updatedElevs
}

// CyclicCounter iterates over the floors, buttons, and elevator IDs to update states
func CyclicCounter(existingElevs []CyclicElevator, myId int) []CyclicElevator {
	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			allEqual := true
			otherElevatorsFurtherProgressed := false
			firstValue := existingElevs[0].Requests[floor][btn]
			for i := 0; i < len(existingElevs); i++ {
				if existingElevs[i].Requests[floor][btn] != firstValue {
					allEqual = false
					break
				}
				if existingElevs[i].ID == myId {
					// Skip the elevator if it's the one corresponding to myId
					continue
				}

				// Implement your logic here, e.g., updating requests or handling orders
				if existingElevs[i].Requests[floor][btn] > existingElevs[myId].Requests[floor][btn] {
					otherElevatorsFurtherProgressed = true 
				}
			if allEqual {
				if existingElevs[myId].Requests[floor][btn] != Confirmed{
					existingElevs[myId].Requests[floor][btn]++
				}
				else {
					existingElevs[myId].Requests[floor][btn] = Completed
				}
			if otherElevatorsFurtherProgressed{
				if existingElevs[myId].Requests[floor][btn] != Completed{
					existingElevs[myId].Requests[floor][btn]++
				}	
			}
			}

				// Additional logic can be added here
			}
		}
	}

	// Return the updated list of elevators
	return existingElevs
}
*/