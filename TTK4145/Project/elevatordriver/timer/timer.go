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
			fmt.Println("Motor timer started")
			MotorTimer = time.NewTimer(3 * time.Second)
			
		case <-DoorTimer.C:
			if startDoor{
				fmt.Println("Door timeout")
				startDoor = false
				doorClosedChan <- true
			}
		case <-MotorTimer.C:
			fmt.Println("Motor timeout")
			if startMotor{
				fmt.Println("Motor timeout")
				fmt.Println("Motor timeout")
				fmt.Println("Motor timeout")
				fmt.Println("Motor timeout")

				motorInactiveChan <- true
			}
		}
	}
}