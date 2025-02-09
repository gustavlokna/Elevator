package elevatordriver

import (
	. "Project/dataenums"
	"fmt"
)






// TODO THIS IS COPIED FROM LAST YEARS PROJECT
// IT IS ALSO ASS
/*
func ClearAtCurrentFloor(e Elevator) [NFloors][NButtons]bool{
	clearedRequests := [NFloors][NButtons]bool{}

	switch e.Config.ClearRequestVariant {
	case CRVAll:
		for btn := BHallUp; btn <= BCab; btn++ {
			if e.Requests[e.CurrentFloor][btn] {
				clearedRequests[e.CurrentFloor][btn] = true
			}
		}

	case CRVInDirn:
		if e.Requests[e.CurrentFloor][BCab] {
			clearedRequests[e.CurrentFloor][BCab] = true
		}
		switch e.Dirn {
		case MDUp:
			if !requestsAbove(e) && !e.Requests[e.CurrentFloor][BHallUp] {
				if e.Requests[e.CurrentFloor][BHallDown] {
					clearedRequests[e.CurrentFloor][BHallDown] = true
				}
			}
			if e.Requests[e.CurrentFloor][BHallUp] {
				clearedRequests[e.CurrentFloor][BHallUp] = true
			}

		case MDDown:
			if !requestsBelow(e) && !e.Requests[e.CurrentFloor][BHallDown] {
				if e.Requests[e.CurrentFloor][BHallUp] {
					clearedRequests[e.CurrentFloor][BHallUp] = true
				}
			}
			if e.Requests[e.CurrentFloor][BHallDown] {
				clearedRequests[e.CurrentFloor][BHallDown] = true
			}

		default:
			for btn := BHallUp; btn <= BCab; btn++ {
				if e.Requests[e.CurrentFloor][btn] {
					clearedRequests[e.CurrentFloor][btn] = true
				}
			}
		}
	}
	return clearedRequests
}
*/
func clearAtCurrentFloor(e Elevator) ([NFloors][NButtons]bool, Elevator) {
	clearedRequests := [NFloors][NButtons]bool{}
	fmt.Println("WE ARE WORKING")
	// Always clear cab order
	if e.Requests[e.CurrentFloor][BCab] {
		clearedRequests[e.CurrentFloor][BCab] = true
		e.Requests[e.CurrentFloor][BCab]= false 
	}

	switch e.Dirn {
	case MDUp:
		fmt.Println("WE MOVE UP in life ")
		if e.Requests[e.CurrentFloor][BHallUp] {
			fmt.Println("WE CLEAR MDUP")
			clearedRequests[e.CurrentFloor][BHallUp] = true
			e.Requests[e.CurrentFloor][BHallUp] = false 
		}
	case MDDown:
		fmt.Println("WE MOVE down in life ")
		if e.Requests[e.CurrentFloor][BHallDown] {
			fmt.Println("WE CLEAR MDOWN")
			clearedRequests[e.CurrentFloor][BHallDown] = true
			e.Requests[e.CurrentFloor][BHallDown] = false 
		}
	case MDStop: 
		if e.Requests[e.CurrentFloor][BHallUp] {
			clearedRequests[e.CurrentFloor][BHallUp] = true
			e.Requests[e.CurrentFloor][BHallUp] = false  
		}
		if e.Requests[e.CurrentFloor][BHallDown] && !clearedRequests[e.CurrentFloor][BHallUp] {
			clearedRequests[e.CurrentFloor][BHallDown] = true
			e.Requests[e.CurrentFloor][BHallDown] = false 
		}
	}

	return clearedRequests, e
}


// TODO BELOW HERE IS COPIED FROM LAST YEARS PROJECT 
func ShouldStop(e Elevator) bool {
	switch e.Dirn {
	case MDDown:
		return e.Requests[e.CurrentFloor][BHallDown] ||
			e.Requests[e.CurrentFloor][BCab] ||
			!requestsBelow(e)
	case MDUp:
		return e.Requests[e.CurrentFloor][BHallUp] ||
			e.Requests[e.CurrentFloor][BCab] ||
			!requestsAbove(e)
	default:
		return true
	}
}

func ChooseDirection(el Elevator) Elevator {
	dirnBehaviour := decideDirection(el)

	el.Dirn = dirnBehaviour.Dirn
	el.CurrentBehaviour = dirnBehaviour.Behaviour
	return el
}

func decideDirection(el Elevator) DirnBehaviourPair {
	switch el.Dirn {
	case MDUp:
		return decideDirectionUp(el)
	case MDDown:
		return decideDirectionDown(el)
	case MDStop:
		return decideDirectionStop(el)
	default:
		return DirnBehaviourPair{
			MDStop,
			EBIdle,
		}
	}
}


func decideDirectionUp(el Elevator) DirnBehaviourPair {
	if requestsAbove(el) {
		return DirnBehaviourPair{MDUp, EBMoving}
	} else if requestsHere(el) {
		return DirnBehaviourPair{MDDown, EBDoorOpen}
	} else if requestsBelow(el) {
		return DirnBehaviourPair{MDDown, EBMoving}
	}
	return DirnBehaviourPair{MDStop, EBIdle}
}

func decideDirectionDown(e Elevator) DirnBehaviourPair {
	if requestsBelow(e) {
		return DirnBehaviourPair{MDDown, EBMoving}
	} else if requestsHere(e) {
		return DirnBehaviourPair{MDUp, EBDoorOpen}
	} else if requestsAbove(e) {
		return DirnBehaviourPair{MDUp, EBMoving}
	}
	return DirnBehaviourPair{MDStop, EBIdle}
}

func decideDirectionStop(e Elevator) DirnBehaviourPair {
	if requestsHere(e) {
		return DirnBehaviourPair{MDStop, EBDoorOpen}
	} else if requestsAbove(e) {
		return DirnBehaviourPair{MDUp, EBMoving}
	} else if requestsBelow(e) {
		return DirnBehaviourPair{MDDown, EBMoving}
	}
	return DirnBehaviourPair{MDStop, EBIdle}
}



func requestsAbove(e Elevator) bool {
	for f := e.CurrentFloor + 1; f < NFloors; f++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.CurrentFloor; f++ {
		for btn := BHallUp; btn <= BCab; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for btn := BHallUp; btn <= BCab; btn++ {
		if e.Requests[e.CurrentFloor][btn] {
			return true

		}
	}
	return false
}


