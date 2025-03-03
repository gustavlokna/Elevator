package timer

import (
	"fmt"
	"time"
)

type TimerType int

const (
	DoorTimer TimerType = iota
	MotorWatchdogTimer
)

func Timer(
	doorOpenChan <-chan bool,
	motorActiveChan <-chan bool,
	doorClosedChan chan<- bool,
	motorInactiveChan chan<- bool,
) {
	var startDoor, startMotor bool
	var doorTimeout, motorTimeout time.Time
	for {
		select {
		case startDoor = <-doorOpenChan:
			doorTimeout = time.Now().Add(3 * time.Second)

		case startMotor = <-motorActiveChan:
			motorTimeout = time.Now().Add(3 * time.Second)

		default:
			if startDoor && time.Now().After(doorTimeout) {
				fmt.Println("Door timeout")
				startDoor = false
				doorClosedChan <- true
			}
			if startMotor && time.Now().After(motorTimeout) {
				motorInactiveChan <- true
			}
			time.Sleep(3 * time.Millisecond)
		}
	}
}