package timer

import (
	"fmt"
	"time"
)

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
	var startDoor, startMotor bool
	var doorTimeout, motorTimeout time.Time	
	for {
		select {
		case startDoor = <-startDoorTimer:
			fmt.Println("RESTART TIMERS")
			fmt.Println(startDoor)
			doorTimeout = time.Now().Add(3 * time.Second)
		case startMotor = <-startMotorTimer:
			motorTimeout = time.Now().Add(3 * time.Second)
			
		default:
			
			if startDoor && time.Now().After(doorTimeout) {
				startDoor = false 
				doorTimeoutChan <- true // Notify that motor watchdog has expired
			}
			if startMotor && time.Now().After(motorTimeout) {
				motorTimeoutChan <- true // Notify that motor watchdog has expired
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
