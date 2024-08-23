package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"time"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	nodeID string,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator       = initelevator()
		prevelevator   = elevator
		completedOrders = [NFloors][NButtons]bool{}
		obstruction    = false
	)

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go hwelevio.PollFloorSensor(drv_floors)
	go hwelevio.PollObstructionSwitch(drv_obstr)
	go hwelevio.PollStopButton(drv_stop)

	// Initialization of elevator
	hwelevio.SetMotorDirection(elevator.Dirn)
	elevator.CurrentFloor = <-drv_floors
	elevator.Dirn = MDStop
	hwelevio.SetMotorDirection(elevator.Dirn)
	
	// Initialize and send initial PayloadFromElevator
// Initialize and send initial PayloadFromElevator
	payload := PayloadFromElevator{
		Elevator:        elevator,
		CompletedOrders: completedOrders,
	}
	payloadFromElevator <- payload


	for {
		prevelevator = elevator
		select {
		case obstruction = <-drv_obstr:
			print("obst: ", obstruction)
		case elevator.CurrentFloor = <-drv_floors:
			print("etasje: ", elevator.CurrentFloor)
			hwelevio.SetFloorIndicator(elevator.CurrentFloor)
		case elevator.Requests = <-newOrderChannel:
		default:
			time.Sleep(10 * time.Millisecond)
		}

		switch elevator.CurrentBehaviour {
		case EBIdle:
			elevator = ChooseDirection(elevator)
			hwelevio.SetMotorDirection(elevator.Dirn)

		case EBMoving:
			if ShouldStop(elevator) {
				hwelevio.SetMotorDirection(MDStop)
				elevator.Dirn = MDStop
				elevator.CurrentBehaviour = EBDoorOpen
				continue
			}

		case EBDoorOpen:
			ElevatorPrint(elevator)
			if obstruction {
				print("hello we have a obst")
			} else {
				completedOrders = ClearAtCurrentFloor(elevator)
				time.Sleep(3 * time.Second) // Simulate door open time
				elevator.CurrentBehaviour = EBIdle
			}
		}

		// Update and send PayloadFromElevator if elevator state changes
		if prevelevator != elevator {
			payload = PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: completedOrders,
			}
			payloadFromElevator <- payload
			
		}
	}
}
