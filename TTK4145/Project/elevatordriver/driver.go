package elevatordriver

import (
	. "Project/dataenums"
)

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


