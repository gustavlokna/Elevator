package timer

import (
	. "Project/config"
	"time"
)

type TimerType int

const (
	DoorTimer TimerType = iota
	MotorWatchdogTimer
)

func Timer(
	doorOpenChan <- chan bool,
	motorActiveChan <- chan bool,
	doorClosedChan chan <- bool,
	motorInactiveChan chan <- bool,
) {
	var startDoor, startMotor bool
	MotorTimer := time.NewTimer(time.Hour)
	MotorTimer.Stop()
	DoorTimer := time.NewTimer(time.Hour)
	DoorTimer.Stop()

	for {
		select {
		case startDoor = <-doorOpenChan:
			DoorTimer = time.NewTimer(DoorOpenDuration)

		case startMotor = <-motorActiveChan:
			MotorTimer = time.NewTimer(MotorTimeout)

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
