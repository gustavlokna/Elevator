package elevatordriver

import (
	. "Project/dataenums"
	"Project/elevatordriver/timer"
	"Project/hwelevio"
	//"fmt"
	//"fmt"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	payloadToLights chan<- PayloadFromDriver,
) {
	floorChannel := make(chan int)
	obstructionChannel := make(chan bool)

	doorOpenChan := make(chan bool)
	doorClosedChan := make(chan bool)

	motorActiveChan := make(chan bool)
	motorInactiveChan := make(chan bool)

	go hwelevio.PollFloorSensor(floorChannel)
	go hwelevio.PollObstructionSwitch(obstructionChannel)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)

	elevator := initelevator()
	hwelevio.SetMotorDirection(elevator.Dirn)

	clearedRequests := [NFloors][NButtons]bool{}
	var obstruction bool

	payloadFromElevator <- PayloadFromElevator{
		Elevator:        elevator,
		CompletedOrders: clearedRequests,
	}

	payloadToLights <- PayloadFromDriver{
		CurrentFloor: elevator.CurrentFloor,
		DoorLight:    false,
	}

	for {
		select {
		case elevator.CurrentFloor = <-floorChannel:
			ElevatorPrint(elevator)
			elevator.ActiveSatus = true
			motorActiveChan <- true
			switch elevator.CurrentBehaviour {
			case EBMoving:

				switch elevator.Dirn {
				case MDUp:

					if elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BCab] || !requestsAbove(elevator) {
						if requestsHere(elevator) {
							hwelevio.SetMotorDirection(MDStop)
							motorActiveChan <- false

							payloadToLights <- PayloadFromDriver{
								CurrentFloor: elevator.CurrentFloor,
								DoorLight:    true,
							}

							doorOpenChan <- true
							elevator.CurrentBehaviour = EBDoorOpen
						}
					}

				case MDDown:
					if elevator.Requests[elevator.CurrentFloor][BHallDown] || elevator.Requests[elevator.CurrentFloor][BCab] || !requestsBelow(elevator) {
						if requestsHere(elevator) {
							hwelevio.SetMotorDirection(MDStop)
							motorActiveChan <- false

							payloadToLights <- PayloadFromDriver{
								CurrentFloor: elevator.CurrentFloor,
								DoorLight:    true,
							}

							doorOpenChan <- true
							elevator.CurrentBehaviour = EBDoorOpen

						}
					}
				}
				// -------------------------------- DUE TO INITIALIZING ERRORS --------------------------------------
			default:
				hwelevio.SetMotorDirection(MDStop)
				//TODO: Maybe write this into a struct from the beginning to make it more clean

				payloadFromElevator <- PayloadFromElevator{
					Elevator:        elevator,
					CompletedOrders: clearedRequests,
				}
				payloadToLights <- PayloadFromDriver{
					CurrentFloor: elevator.CurrentFloor,
					DoorLight:    false,
				}
				//---------------------------------------------------------------------------------------------------
			}

		case <-doorClosedChan:
			if obstruction {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- true
				payloadFromElevator <- PayloadFromElevator{
					Elevator:        elevator,
					CompletedOrders: clearedRequests,
				}
				continue
			}
			var clearedRequests [NFloors][NButtons]bool
			if elevator.CurrentBehaviour == EBDoorOpen {
				if elevator.Requests[elevator.CurrentFloor][BCab] {
					clearedRequests[elevator.CurrentFloor][BCab] = true
					elevator.Requests[elevator.CurrentFloor][BCab] = false
				}
				switch elevator.Dirn {
				case MDUp:
					if !requestsAbove(elevator) && !elevator.Requests[elevator.CurrentFloor][BHallUp] {
						if elevator.Requests[elevator.CurrentFloor][BHallDown] {
							clearedRequests[elevator.CurrentFloor][BHallDown] = true
							elevator.Requests[elevator.CurrentFloor][BHallDown] = false
						}
					}
					if elevator.Requests[elevator.CurrentFloor][BHallUp] {
						clearedRequests[elevator.CurrentFloor][BHallUp] = true
						elevator.Requests[elevator.CurrentFloor][BHallUp] = false
					}
				case MDDown:
					if !requestsBelow(elevator) && !elevator.Requests[elevator.CurrentFloor][BHallDown] {
						if elevator.Requests[elevator.CurrentFloor][BHallUp] {
							clearedRequests[elevator.CurrentFloor][BHallUp] = true
							elevator.Requests[elevator.CurrentFloor][BHallUp] = false
						}
					}
					if elevator.Requests[elevator.CurrentFloor][BHallDown] {
						clearedRequests[elevator.CurrentFloor][BHallDown] = true
						elevator.Requests[elevator.CurrentFloor][BHallDown] = false
					}
				}
				elevator = ChooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)
				payloadFromElevator <- PayloadFromElevator{
					Elevator:        elevator,
					CompletedOrders: clearedRequests,
				}
				payloadToLights <- PayloadFromDriver{
					CurrentFloor: elevator.CurrentFloor,
					DoorLight:    false,
				}
			}

		case <-motorInactiveChan:
			if elevator.CurrentBehaviour == EBMoving {
				elevator.ActiveSatus = false
				payloadFromElevator <- PayloadFromElevator{
					Elevator:        elevator,
					CompletedOrders: clearedRequests,
				}
			}

		case obstruction = <-obstructionChannel:
			if elevator.CurrentBehaviour == EBDoorOpen {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- !obstruction
			}

			payloadFromElevator <- PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
			}

		case elevator.Requests = <-newOrderChannel:
			//TODO: Wrte into switch-case
			if elevator.CurrentBehaviour == EBIdle {
				switch  {
				case elevator.Requests[elevator.CurrentFloor][BCab]:
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				case elevator.Requests[elevator.CurrentFloor][BHallUp] && elevator.Dirn == MDUp:
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				case elevator.Requests[elevator.CurrentFloor][BHallDown] && elevator.Dirn == MDDown:
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				default:
					motorActiveChan <- true
					elevator = ChooseDirection(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)
				}
				// if requestsHere(elevator) {
				// 	elevator.CurrentBehaviour = EBDoorOpen
				// 	doorOpenChan <- true
				// 	payloadToLights <- PayloadFromDriver{
				// 		CurrentFloor: elevator.CurrentFloor,
				// 		DoorLight:    true,
				// 	}
				// }
				// if !requestsHere(elevator) {
				// 	motorActiveChan <- true
				// 	elevator = ChooseDirection(elevator)
				// 	hwelevio.SetMotorDirection(elevator.Dirn)

				// }
			}
			payloadFromElevator <- PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
			}
		}
	}
}
