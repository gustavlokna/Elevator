package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"time"
)

func ElevatorDriver(
	fromOrderAssignerChannel <-chan [NFloors][NButtons]bool,
	toOrderAssignerChannel chan<- Elevator,
	lifelineChannel chan<- bool,
	nodeID int,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator = initelevator()
		prevelevator = elevator
		obstruction = false
	)

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//drv_motorActivity := make(chan bool)

	go hwelevio.PollFloorSensor(drv_floors)
	go hwelevio.PollObstructionSwitch(drv_obstr)
	go hwelevio.PollStopButton(drv_stop)
	//go hwelevio.MontitorMotorActivity(drv_motorActivity, 3.0)
	for {
		prevelevator = elevator
		select {
		case obstruction = <-drv_obstr:
			// if true set obstr variable true and
			// else set false
			print("obst: ", obstruction)
		case elevator.CurrentFloor = <-drv_floors:
			hwelevio.SetFloorIndicator(elevator.CurrentFloor)
		case elevator.Requests = <-fromOrderAssignerChannel:
			ElevatorPrint(elevator)
			
		default:
			// Prevent busy loop
			time.Sleep(10 * time.Millisecond)
		}

		switch elevator.CurrentBehaviour {
		case EBIdle:
			//move on assigned orders
		case EBDoorOpen:
			if obstruction {
				print("hello we have a obst")
				//stop motor
				//restart timer 
			}else{
				//stop motor
				// start timer 
				print("wihuu")
			}
		case EBMoving:
			ElevatorPrint(elevator)
			if ShouldStop(elevator){
				print("HALLO DU MÅ STOPPE")
				hwelevio.SetMotorDirection(MDStop)
			}
				// set elev.behavior = stop 

		}
		/*
		default:
			if timer.TimedOut() 
				//stop timer 
				// set state as idle
			time.Sleep(10 * time.Millisecond)
		*/
		if prevelevator != elevator {
			toOrderAssignerChannel <- elevator
			}
	}
}
