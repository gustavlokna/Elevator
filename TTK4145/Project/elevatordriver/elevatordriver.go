package elevatordriver

import (
	. "Project/dataenums"
	"Project/elevatordriver/timer"
	"Project/hwelevio"
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

// ------------------------------------ TODO: Make these into one-liners --------------------------------------
	payloadFromElevator <- PayloadFromElevator{
		Elevator:        elevator,
		CompletedOrders: clearedRequests,
	}

	payloadToLights <- PayloadFromDriver{
		CurrentFloor: elevator.CurrentFloor,
		DoorLight:    false,
	}
// -------------------------------------------------------------------------------------------------------------
	for {
		select {
		case elevator.CurrentFloor = <-floorChannel:
			elevator.ActiveSatus = true
			motorActiveChan <- true
			switch elevator.CurrentBehaviour {
			case EBMoving:
				switch {
				case elevator.Requests[elevator.CurrentFloor][BCab]:
					hwelevio.SetMotorDirection(MDStop)
					motorActiveChan <- false
					elevator.CurrentBehaviour = EBDoorOpen

					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}

					doorOpenChan <- true

				case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
					hwelevio.SetMotorDirection(MDStop)
					motorActiveChan <- false
					elevator.CurrentBehaviour = EBDoorOpen

					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}

					doorOpenChan <- true

				case elevator.Dirn == MDUp && requestsAbove(elevator):
					// do nothing (no requests at floor)
				case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
					hwelevio.SetMotorDirection(MDStop)
					motorActiveChan <- false
					elevator.CurrentBehaviour = EBDoorOpen

					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}

					doorOpenChan <- true

				case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
					hwelevio.SetMotorDirection(MDStop)
					motorActiveChan <- false
					elevator.CurrentBehaviour = EBDoorOpen

					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}

					doorOpenChan <- true

				case elevator.Dirn == MDDown && requestsBelow(elevator):
					// DO NOTHING
				case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
					hwelevio.SetMotorDirection(MDStop)
					motorActiveChan <- false
					elevator.CurrentBehaviour = EBDoorOpen
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
					doorOpenChan <- true

				default:
					elevator = ChooseDirection(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    false,
					}
				}

			default:
				elevator = ChooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)
			}
			// payloadToLights <- PayloadFromDriver{
			// 	CurrentFloor: elevator.CurrentFloor,
			// 	DoorLight:    false,
			// }
			payloadFromElevator <- PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
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

			var clearedRequests [NFloors][NButtons]bool //TODO: Remove
			
			if elevator.CurrentBehaviour == EBDoorOpen {
				if elevator.Requests[elevator.CurrentFloor][BCab] {
					clearedRequests[elevator.CurrentFloor][BCab] = true
					elevator.Requests[elevator.CurrentFloor][BCab] = false
				}
				switch {
				case elevator.Dirn == MDUp && !requestsAbove(elevator) && !elevator.Requests[elevator.CurrentFloor][BHallUp] && elevator.Requests[elevator.CurrentFloor][BHallDown]:
					clearedRequests[elevator.CurrentFloor][BHallDown] = true
					elevator.Requests[elevator.CurrentFloor][BHallDown] = false

				case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
					clearedRequests[elevator.CurrentFloor][BHallUp] = true
					elevator.Requests[elevator.CurrentFloor][BHallUp] = false

				case elevator.Dirn == MDDown && !requestsBelow(elevator) && !elevator.Requests[elevator.CurrentFloor][BHallDown] && elevator.Requests[elevator.CurrentFloor][BHallUp]:
						clearedRequests[elevator.CurrentFloor][BHallUp] = true
						elevator.Requests[elevator.CurrentFloor][BHallUp] = false

				case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
					clearedRequests[elevator.CurrentFloor][BHallDown] = true
					elevator.Requests[elevator.CurrentFloor][BHallDown] = false
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
			// if elevator.CurrentBehaviour == EBDoorOpen {
			// 	if elevator.Requests[elevator.CurrentFloor][BCab] {
			// 		clearedRequests[elevator.CurrentFloor][BCab] = true
			// 		elevator.Requests[elevator.CurrentFloor][BCab] = false
			// 	}
			// 	switch elevator.Dirn {
			// 	case MDUp:
			// 		if !requestsAbove(elevator) && !elevator.Requests[elevator.CurrentFloor][BHallUp] {
			// 			if elevator.Requests[elevator.CurrentFloor][BHallDown] {
			// 				clearedRequests[elevator.CurrentFloor][BHallDown] = true
			// 				elevator.Requests[elevator.CurrentFloor][BHallDown] = false
			// 			}
			// 		}
			// 		if elevator.Requests[elevator.CurrentFloor][BHallUp] {
			// 			clearedRequests[elevator.CurrentFloor][BHallUp] = true
			// 			elevator.Requests[elevator.CurrentFloor][BHallUp] = false
			// 		}
			// 	case MDDown:
			// 		if !requestsBelow(elevator) && !elevator.Requests[elevator.CurrentFloor][BHallDown] {
			// 			if elevator.Requests[elevator.CurrentFloor][BHallUp] {
			// 				clearedRequests[elevator.CurrentFloor][BHallUp] = true
			// 				elevator.Requests[elevator.CurrentFloor][BHallUp] = false
			// 			}
			// 		}
			// 		if elevator.Requests[elevator.CurrentFloor][BHallDown] {
			// 			clearedRequests[elevator.CurrentFloor][BHallDown] = true
			// 			elevator.Requests[elevator.CurrentFloor][BHallDown] = false
			// 		}
			// 	}
			// 	elevator = ChooseDirection(elevator)
			// 	hwelevio.SetMotorDirection(elevator.Dirn)
			// 	payloadFromElevator <- PayloadFromElevator{
			// 		Elevator:        elevator,
			// 		CompletedOrders: clearedRequests,
			// 	}
			// 	payloadToLights <- PayloadFromDriver{
			// 		CurrentFloor: elevator.CurrentFloor,
			// 		DoorLight:    false,
			// 	}
			// }

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
			if elevator.CurrentBehaviour == EBIdle {
				switch {
				case elevator.Requests[elevator.CurrentFloor][BHallUp]:
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDUp
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				case elevator.Requests[elevator.CurrentFloor][BHallDown]:
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDDown
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				case elevator.Requests[elevator.CurrentFloor][BCab]:
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
			}
			ElevatorPrint(elevator)
			payloadFromElevator <- PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
			}
		}
	}
}
