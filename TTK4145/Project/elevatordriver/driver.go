package elevatordriver

import (
	. "Project/dataenums"
)

var outputDevice ElevOutputDevice

func buttonPressed(elevator Elevator, btnEvent ButtonEvent) Elevator{
	if shouldClearImmediately(elevator, btnEvent) && (elevator.CurrentBehaviour == EBDoorOpen) {
		startTimer(elevator.Config.DoorOpenDurationS)
	} else {
		elevator.Requests[btnEvent.Floor][btnEvent.Button] = true
	}
	return elevator
}

func MoveOnActiveOrders(el Elevator) {
	switch el.CurrentBehaviour {
	case EBIdle:
		pair := ChooseDirection(el)
		el.Dirn = pair.Dirn
		el.CurrentBehaviour = pair.Behaviour
		switch pair.Behaviour {
		case EBDoorOpen:
			outputDevice.DoorLight(true)
			startTimer(el.Config.DoorOpenDurationS)
			el = ClearAtCurrentFloor(el)
		case EBMoving:
			outputDevice.MotorDirection(el.Dirn)
		}
	}
	//SetAllLights()
}

func ClearAtCurrentFloor(e Elevator) Elevator {

	beforeClear := make(map[Button]bool)
	for btn := BHallUp; btn <= BCab; btn++ {
		beforeClear[btn] = e.Requests[e.CurrentFloor][btn]
	}

	switch e.Config.ClearRequestVariant {
	case CRVAll:
		for btn := BHallUp; btn <= BCab; btn++ {
			e.Requests[e.CurrentFloor][btn] = false
		}

	case CRVInDirn:
		e.Requests[e.CurrentFloor][BCab] = false
		switch e.Dirn {
		case DirUp:
			if !requestsAbove(e) && !e.Requests[e.CurrentFloor][BHallUp] {
				e.Requests[e.CurrentFloor][BHallDown] = false
			}
			e.Requests[e.CurrentFloor][BHallUp] = false

		case DirDown:
			if !requestsBelow(e) && !e.Requests[e.CurrentFloor][BHallDown] {
				e.Requests[e.CurrentFloor][BHallUp] = false
			}
			e.Requests[e.CurrentFloor][BHallDown] = false
		default:
			e.Requests[e.CurrentFloor][BHallUp]   = false
			e.Requests[e.CurrentFloor][BHallDown] = false

		}
	}
	//Todo find out what this is for 
	for btn, wasPressed := range beforeClear {
		if wasPressed && !e.Requests[e.CurrentFloor][btn] {
			statehandler.UpdateStateOnCompletedHallOrder(e, elevatorName, e.CurrentFloor, btn)
		}
	}
	return e
}


func ChooseDirection(el Elevator) DirnBehaviourPair {
	switch el.Dirn {
	case DirUp:
		return decideDirectionUp(el)
	case DirDown:
		return decideDirectionDown(el)
	case DirStop:
		return decideDirectionStop(el)
	default:
		return DirnBehaviourPair{
			DirStop,
			EBIdle,
		}
	}
}


func decideDirectionUp(el Elevator) DirnBehaviourPair {
	if requestsAbove(el) {
		return DirnBehaviourPair{DirUp, EBMoving}
	} else if requestsHere(el) {
		return DirnBehaviourPair{DirDown, EBDoorOpen}
	} else if requestsBelow(el) {
		return DirnBehaviourPair{DirDown, EBMoving}
	}
	return DirnBehaviourPair{DirStop, EBIdle}
}

func decideDirectionDown(e Elevator) DirnBehaviourPair {
	if requestsBelow(e) {
		return DirnBehaviourPair{DirDown, EBMoving}
	} else if requestsHere(e) {
		return DirnBehaviourPair{DirUp, EBDoorOpen}
	} else if requestsAbove(e) {
		return DirnBehaviourPair{DirUp, EBMoving}
	}
	return DirnBehaviourPair{DirStop, EBIdle}
}

func decideDirectionStop(e Elevator) DirnBehaviourPair {
	if requestsHere(e) {
		return DirnBehaviourPair{DirStop, EBDoorOpen}
	} else if requestsAbove(e) {
		return DirnBehaviourPair{DirUp, EBMoving}
	} else if requestsBelow(e) {
		return DirnBehaviourPair{DirDown, EBMoving}
	}
	return DirnBehaviourPair{DirStop, EBIdle}
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