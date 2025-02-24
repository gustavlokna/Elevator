package timer

import "time"

// TimerType defines the type of timer
type TimerType int

const (
	DoorTimer TimerType = iota
	MotorWatchdogTimer
)

// Timer handles timing for door opening and motor inactivity (watchdog).
func Timer(
	startDoorTimer <-chan bool,
	startMotorTimer <-chan bool,
	doorTimeoutChan chan<- bool,
	motorTimeoutChan chan<- bool,
) {
	for {
		select {
		case start := <-startDoorTimer:
			if start {
				go func() {
					time.Sleep(3 * time.Second)
					doorTimeoutChan <- true // Notify that door timer has expired
				}()
			}

		case start := <-startMotorTimer:
			if start {
				go func() {
					time.Sleep(5 * time.Second)
					motorTimeoutChan <- true // Notify that motor watchdog has expired
				}()
			}
		}
	}
}

// package elevatordriver
// import (
// 	"time"
// )
// //this folder is totaly copied
// // TODO NOTING HERE IS BEING USED.
// // REMOVE ?!
// func GetCurrentTimeAsFloat() float64 {
// 	now := time.Now()
// 	return float64(now.Unix()) + float64(now.Nanosecond())*1e-9
// }

// var endTime float64
// var isActive bool
// var IsInfinite bool

// func startTimer(duration float64) {
// 	endTime = GetCurrentTimeAsFloat() + duration
// 	isActive = true
// }
