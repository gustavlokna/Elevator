package elevatordriver

import (
	. "Project/dataenums"
)



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

func ClearAtCurrentFloor(e Elevator) (Elevator, [NFloors][NButtons]bool) {
	clearedRequests := [NFloors][NButtons]bool{}

	switch e.Config.ClearRequestVariant {
	case CRVAll:
		for btn := BHallUp; btn <= BCab; btn++ {
			if e.Requests[e.CurrentFloor][btn] {
				clearedRequests[e.CurrentFloor][btn] = true
				e.Requests[e.CurrentFloor][btn] = false
			}
		}

	case CRVInDirn:
		if e.Requests[e.CurrentFloor][BCab] {
			clearedRequests[e.CurrentFloor][BCab] = true
			e.Requests[e.CurrentFloor][BCab] = false
		}
		switch e.Dirn {
		case MDUp:
			if !requestsAbove(e) && !e.Requests[e.CurrentFloor][BHallUp] {
				if e.Requests[e.CurrentFloor][BHallDown] {
					clearedRequests[e.CurrentFloor][BHallDown] = true
					e.Requests[e.CurrentFloor][BHallDown] = false
				}
			}
			if e.Requests[e.CurrentFloor][BHallUp] {
				clearedRequests[e.CurrentFloor][BHallUp] = true
				e.Requests[e.CurrentFloor][BHallUp] = false
			}

		case MDDown:
			if !requestsBelow(e) && !e.Requests[e.CurrentFloor][BHallDown] {
				if e.Requests[e.CurrentFloor][BHallUp] {
					clearedRequests[e.CurrentFloor][BHallUp] = true
					e.Requests[e.CurrentFloor][BHallUp] = false
				}
			}
			if e.Requests[e.CurrentFloor][BHallDown] {
				clearedRequests[e.CurrentFloor][BHallDown] = true
				e.Requests[e.CurrentFloor][BHallDown] = false
			}

		default:
			for btn := BHallUp; btn <= BCab; btn++ {
				if e.Requests[e.CurrentFloor][btn] {
					clearedRequests[e.CurrentFloor][btn] = true
					e.Requests[e.CurrentFloor][btn] = false
				}
			}
		}
	}

	return e, clearedRequests
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