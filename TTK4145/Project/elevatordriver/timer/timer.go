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
	var motorTimeout time.Time
	timeCounter := time.NewTimer(time.Hour)
	timeCounter.Stop()

	for {
		select {
		case startDoor = <-doorOpenChan:
			timeCounter = time.NewTimer(3 * time.Second)

		case startMotor = <-motorActiveChan:
			motorTimeout = time.Now().Add(3 * time.Second)
			
		case <-timeCounter.C:
			if startDoor{
				fmt.Println("Door timeout")
				startDoor = false
				doorClosedChan <- true
			}
					
		default:
			if startMotor && time.Now().After(motorTimeout) {
				motorInactiveChan <- true
			}
			time.Sleep(3 * time.Millisecond)
		}
	}
}