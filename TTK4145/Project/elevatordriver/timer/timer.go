package timer

import (
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
	MotorTimer := time.NewTimer(time.Hour)
	MotorTimer.Stop()
	DoorTimer := time.NewTimer(time.Hour)
	DoorTimer.Stop()

	for {
		select {
		case startDoor = <-doorOpenChan:
			// TODO MAKE VARIABLE IN CONFIG FILE
			DoorTimer = time.NewTimer(3 * time.Second)

		case startMotor = <-motorActiveChan:
			// WHEN WORKING ON SLOW ELEV USE 4 Sec, but should be 3
			// TODO MAKE VARIABLE IN CONFIG FILE
			MotorTimer = time.NewTimer(4 * time.Second)

		case <-DoorTimer.C:
			if startDoor {
				startDoor = false
				doorClosedChan <- true
			}
		case <-MotorTimer.C:
			if startMotor {
				motorInactiveChan <- true
			}
		}
	}
}
