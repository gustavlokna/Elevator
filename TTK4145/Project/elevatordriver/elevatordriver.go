package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"time"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	newStateChanel chan<- Elevator,
	orderDoneChannel chan <- [NFloors][NButtons]bool,
	nodeID string,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator     = initelevator()
		prevelevator = elevator
		obstruction  = false
	)

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//drv_motorActivity := make(chan bool)

	go hwelevio.PollFloorSensor(drv_floors)
	go hwelevio.PollObstructionSwitch(drv_obstr)
	go hwelevio.PollStopButton(drv_stop)

	//initilazation of elevator
	hwelevio.SetMotorDirection(elevator.Dirn)
	elevator.CurrentFloor = <-drv_floors
	elevator.Dirn = MDStop
	hwelevio.SetMotorDirection(elevator.Dirn)
	newStateChanel <- elevator

	//go hwelevio.MontitorMotorActivity(drv_motorActivity, 3.0)
	for {
		prevelevator = elevator
		select {
		case obstruction = <-drv_obstr:
			// if true set obstr variable true and
			// else set false
			print("obst: ", obstruction)
		case elevator.CurrentFloor = <-drv_floors:
			print("etasje: ", elevator.CurrentFloor)
			hwelevio.SetFloorIndicator(elevator.CurrentFloor)
			//ElevatorPrint(elevator)
		case elevator.Requests = <-newOrderChannel:
			//ElevatorPrint(elevator)
		default:
			// Prevent busy loop
			time.Sleep(10 * time.Millisecond)
		}
		//print(elevator.CurrentBehaviour)
		switch elevator.CurrentBehaviour {
		case EBIdle:
			elevator = ChooseDirection(elevator)
			hwelevio.SetMotorDirection(elevator.Dirn)
			//ElevatorPrint(elevator)

		case EBMoving:
			if ShouldStop(elevator) {
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				continue
			}

		case EBDoorOpen:
			//outputDevice.DoorLight(true)
			// Todo set doorlight
			//startTimer(elevator.Config.DoorOpenDurationS)
			elevator, clearedRequests := ClearAtCurrentFloor(elevator)
			if obstruction {
				print("hello we have a obst")
				// Handle obstruction
			} else {
				time.Sleep(3 * time.Second) // Simulerer dørens åpningstid
				elevator.CurrentBehaviour = EBIdle
				orderDoneChannel <- clearedRequests
				//hwelevio.SetDoorLight(false)
			}
		}

		/*		/*
				default:
					if timer.TimedOut()
						//stop timer
						// set state as idle
					time.Sleep(10 * time.Millisecond)
		*/
		if prevelevator != elevator {
			newStateChanel <- elevator
		}
	}
}
