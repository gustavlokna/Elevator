package elevatordriver
import (
	"time"
)
//this folder is totaly copied
// TODO NOTING HERE IS BEING USED. 
// REMOVE ?!
func GetCurrentTimeAsFloat() float64 {
	now := time.Now()
	return float64(now.Unix()) + float64(now.Nanosecond())*1e-9
}

var endTime float64
var isActive bool
var IsInfinite bool

func startTimer(duration float64) {
	endTime = GetCurrentTimeAsFloat() + duration
	isActive = true
}