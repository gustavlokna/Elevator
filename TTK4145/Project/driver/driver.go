package driver

import (
	. "Project/dataenums"
	"Project/driver/timer"
	"Project/hwelevio"
	"fmt"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- FromDriverToAssigner,
	payloadToLights chan<- FromDriverToLight,
) {
	var (
		floorChannel       = make(chan int, BufferSize)
		obstructionChannel = make(chan bool, BufferSize)
		doorOpenChan       = make(chan bool, BufferSize)
		doorClosedChan     = make(chan bool, BufferSize)
		motorActiveChan    = make(chan bool, BufferSize)
		motorInactiveChan  = make(chan bool, BufferSize)
		clearedRequests    = [NFloors][NButtons]bool{}
		obstruction        bool
	)

	go hwelevio.PollFloorSensor(floorChannel)
	go hwelevio.PollObstructionSwitch(obstructionChannel)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)

	elevator := initelevator()
	hwelevio.SetMotorDirection(elevator.Dirn)

	for {

		select {
		case elevator.CurrentFloor = <-floorChannel:
			elevator.ActiveSatus = true
			motorActiveChan <- true

			switch {
			case elevator.Requests[elevator.CurrentFloor][BCab]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderAtCurrentFloorInDir(elevator):
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderInCurrentDir(elevator):
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

			case orderAtCurrentFloorOppositeDir(elevator):
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderOppositeDir(elevator):
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

			default:
				elevator.Dirn = MDStop
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = Idle
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case <-doorClosedChan:
			if obstruction {
				elevator.ActiveSatus = !obstruction
				fmt.Println(!obstruction)
				doorOpenChan <- true
				payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
				continue
			}

			switch {
			case orderAtCurrentFloorInDir(elevator):
				clearedRequests[elevator.CurrentFloor][dirToBtn(elevator.Dirn)] = true
				elevator.Requests[elevator.CurrentFloor][dirToBtn(elevator.Dirn)] = false

			case orderInCurrentDir(elevator):

			case orderAtCurrentFloorOppositeDir(elevator):
				elevator.Dirn = setMotorOppositeDir(elevator)
				clearedRequests[elevator.CurrentFloor][dirToBtn(elevator.Dirn)] = true
				elevator.Requests[elevator.CurrentFloor][dirToBtn(elevator.Dirn)] = false

			case orderOppositeDir(elevator):

			default:
				elevator.Dirn = MDStop
				hwelevio.SetMotorDirection(MDStop)
			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				clearedRequests[elevator.CurrentFloor][BCab] = true
			}
			// if move itpo reevant cases
			elevator.CurrentBehaviour = Idle

			payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
			payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
			// Reset clearedRequests to all false
			clearedRequests = [NFloors][NButtons]bool{}

		case <-motorInactiveChan:

			if elevator.CurrentBehaviour == Moving {
				elevator.ActiveSatus = false
				payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case obstruction = <-obstructionChannel:
			if elevator.CurrentBehaviour == DoorOpen {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- !obstruction
			}
			payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

		case elevator.Requests = <-newOrderChannel:
			ElevatorPrint(elevator)
			switch elevator.CurrentBehaviour {
			case Idle:
				switch {
				case orderAtCurrentFloorInDir(elevator) || elevator.Requests[elevator.CurrentFloor][BCab]:
					elevator.Dirn = btnToDirn(elevator)
					elevator.CurrentBehaviour = DoorOpen
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case orderInCurrentDir(elevator):
					elevator.CurrentBehaviour = Moving
					hwelevio.SetMotorDirection(elevator.Dirn)
					motorActiveChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

				case orderAtCurrentFloorOppositeDir(elevator):
					elevator.Dirn = setMotorOppositeDir(elevator)
					elevator.CurrentBehaviour = DoorOpen
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case orderOppositeDir(elevator):
					elevator.CurrentBehaviour = Moving
					elevator.Dirn = setMotorOppositeDir(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)
					motorActiveChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

				default:
					elevator.Dirn = MDStop
					hwelevio.SetMotorDirection(MDStop)
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

				}

			case Moving:
			case DoorOpen:
			}
			payloadFromElevator <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

		}
	}
}
