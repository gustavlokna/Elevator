package orderassigner
import (
	. "Project/dataenums"
)


//these are double up and also in utils for elevatordriver. 
//should be deleted from there when project is complete 
func EBToString(behaviour ElevatorBehaviour) string {
	switch behaviour {
	case EBIdle:
		return "idle"
	case EBDoorOpen:
		return "doorOpen"
	case EBMoving:
		return "moving"
	default:
		return "Unknown"
	}
}
func ElevDirToString(d HWMotorDirection) string {
	switch d {
	case MDDown:
		return "down"
	case MDStop:
		return "stop"
	case MDUp:
		return "up"
	default:
		return "DirUnknown"
	}
}

