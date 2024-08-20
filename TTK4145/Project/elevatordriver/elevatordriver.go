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

	hwelevio.SetMotorDirection(elevator.Dirn)
	elevator.CurrentFloor = <-drv_floors
	elevator.Dirn = MDStop
	hwelevio.SetMotorDirection(elevator.Dirn)

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
			ElevatorPrint(elevator)
		case elevator.Requests = <-fromOrderAssignerChannel:
			ElevatorPrint(elevator)
		default:
			// Prevent busy loop
			time.Sleep(10 * time.Millisecond)
		}
		print(elevator.CurrentBehaviour)
		switch elevator.CurrentBehaviour {
		case EBIdle:
			print("Switching to EBIdle")
			elevator = ChooseDirection(elevator)
			hwelevio.SetMotorDirection(elevator.Dirn)
			ElevatorPrint(elevator)
		
		case EBMoving:
			print("Switching to EBMoving")
			//ElevatorPrint(elevator)
			if ShouldStop(elevator) {
				print("HALLO DU MÅ STOPPE")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				print(elevator.CurrentBehaviour)
				//print("Set elevator.CurrentBehaviour to EBDoorOpen")
				ElevatorPrint(elevator)
				continue
			}
			
			print("trengte ikke stoppe")
			print("elevator.CurrentBehaviour", EBToString(elevator.CurrentBehaviour))
		
		case EBDoorOpen:
			print("døren er åpen (EBDoorOpen case)")
			//outputDevice.DoorLight(true)
			// Todo set doorlight
			//startTimer(elevator.Config.DoorOpenDurationS)
			elevator = ClearAtCurrentFloor(elevator)
			if obstruction {
				print("hello we have a obst")
				// Handle obstruction
			} else {
				print("wihuu")
				time.Sleep(3 * time.Second) // Simulerer dørens åpningstid
				elevator.CurrentBehaviour = EBIdle
				//hwelevio.SetDoorLight(false)
				print("Switching back to EBIdle from EBDoorOpen")
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
			toOrderAssignerChannel <- elevator
		}
	}
}
