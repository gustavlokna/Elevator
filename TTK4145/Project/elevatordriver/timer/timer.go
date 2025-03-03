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
	MotorTimer := time.NewTimer(time.Hour)
	MotorTimer.Stop()
	DoorTimer := time.NewTimer(time.Hour)
	DoorTimer.Stop()

	for {
		select {
		case startDoor = <-doorOpenChan:
			DoorTimer = time.NewTimer(3 * time.Second)

		case startMotor = <-motorActiveChan:
			MotorTimer = time.NewTimer(3 * time.Second)
			
		case <-DoorTimer.C:
			if startDoor{
				fmt.Println("Door timeout")
				startDoor = false
				doorClosedChan <- true
			}
		case <-MotorTimer.C:
			if startMotor{
				motorInactiveChan <- true
			}
		}
	}
}