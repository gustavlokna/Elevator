package elevatordriver

import (
	. "Project/dataenums"
	"Project/elevatordriver/timer"
	"Project/hwelevio"
	"fmt"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	payloadToLights chan<- PayloadFromDriver,
) {
	var (
		floorChannel       = make(chan int)
		obstructionChannel = make(chan bool)
		doorOpenChan       = make(chan bool)
		doorClosedChan     = make(chan bool)
		motorActiveChan    = make(chan bool)
		motorInactiveChan  = make(chan bool)
		clearedRequests    = [NFloors][NButtons]bool{}
		obstruction        bool
	)

	go hwelevio.PollFloorSensor(floorChannel)
	go hwelevio.PollObstructionSwitch(obstructionChannel)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)

	elevator := initelevator()
	hwelevio.SetMotorDirection(elevator.Dirn)

	payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
	payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

	for {
		var clearedRequests [NFloors][NButtons]bool //TODO: Remove
		select {
		case elevator.CurrentFloor = <-floorChannel:
			elevator.ActiveSatus = true
			motorActiveChan <- true
			switch {
			case elevator.Requests[elevator.CurrentFloor][BCab]:
				fmt.Println("")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				doorOpenChan <- true

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case elevator.Dirn == MDUp && requestsAbove(elevator):
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case elevator.Dirn == MDDown && requestsBelow(elevator):
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			default:
				elevator = chooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
			}
			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		case <-doorClosedChan:
			if obstruction {
				elevator.ActiveSatus = !obstruction
				fmt.Println(!obstruction)
				doorOpenChan <- true
				payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
				continue
			}

			switch {
			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BCab] && requestsAbove(elevator):

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				clearedRequests[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				clearedRequests[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BCab] && requestsBelow(elevator):

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			//case elevator.Requests[elevator.CurrentFloor][BCab]: 
			// This case was not necessary after changing chooseDirection
			// but this can have induced other errors. I have not tried yet. 

			default: 
				elevator = chooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)		
			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				clearedRequests[elevator.CurrentFloor][BCab] = true
				elevator.Requests[elevator.CurrentFloor][BCab] = false
			}

			elevator.CurrentBehaviour = EBIdle

			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
			payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

		case <-motorInactiveChan:

			if elevator.CurrentBehaviour == EBMoving {
				elevator.ActiveSatus = false
				payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case obstruction = <-obstructionChannel:
			if elevator.CurrentBehaviour == EBDoorOpen {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- !obstruction
			}
			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		case elevator.Requests = <-newOrderChannel:
			ElevatorPrint(elevator)
			switch elevator.CurrentBehaviour {
			case EBIdle:
				switch {
				case elevator.Requests[elevator.CurrentFloor][BHallUp]: 
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDUp
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case elevator.Requests[elevator.CurrentFloor][BHallDown]:
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDDown
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case elevator.Requests[elevator.CurrentFloor][BCab]:
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				
				// TODO THE CASES BELOW ARE WRONG They should be in combination and use choose and set motor dir
				case requestsAbove(elevator): 
					motorActiveChan <- true
					elevator.CurrentBehaviour = EBMoving
					elevator.Dirn = MDUp
					hwelevio.SetMotorDirection(elevator.Dirn)

				case requestsBelow(elevator):
					motorActiveChan <- true
					elevator.CurrentBehaviour = EBMoving
					elevator.Dirn = MDDown
					hwelevio.SetMotorDirection(elevator.Dirn)
				default:
					elevator.Dirn = MDStop
				}

			case EBMoving:
			case EBDoorOpen:
			}
			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
		}
	}
}
