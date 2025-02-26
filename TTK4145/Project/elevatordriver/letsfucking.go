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
			elevator.ActiveSatus = true
			motorActiveChan <- true
			switch elevator.CurrentBehaviour {
			case EBMoving:

				switch elevator.Dirn {
				case MDUp:
					//If there is a request above this if statement will not be true, even if there is a request here that needs to be handeled
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

				// switch elevator.Dirn {
				// case MDUp:
				// 	//If there is a request above this if statement will not be true, even if there is a request here that needs to be handeled
					
				// 	//If both elevators are in floor 0 and BHallUP is presed in both floor 1 and 2 the one elevator takes both, while one of them FAILES


				// 	switch {
				// 	case elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BCab]:
				// 			hwelevio.SetMotorDirection(MDStop)
				// 			motorActiveChan <- false

				// 			payloadToLights <- PayloadFromDriver{
				// 				CurrentFloor: elevator.CurrentFloor,
				// 				DoorLight:    true,
				// 			}

				// 			doorOpenChan <- true
				// 			elevator.CurrentBehaviour = EBDoorOpen

				// 	case elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BCab] || !requestsAbove(elevator):
				// 		if requestsHere(elevator) {
				// 			hwelevio.SetMotorDirection(MDStop)
				// 			motorActiveChan <- false

				// 			payloadToLights <- PayloadFromDriver{
				// 				CurrentFloor: elevator.CurrentFloor,
				// 				DoorLight:    true,
				// 			}

				// 			doorOpenChan <- true
				// 			elevator.CurrentBehaviour = EBDoorOpen
				// 		}
				// 	}
				// case MDDown:
				// 	switch {
				// 	case elevator.Requests[elevator.CurrentFloor][BHallDown] || elevator.Requests[elevator.CurrentFloor][BCab]:
				// 		hwelevio.SetMotorDirection(MDStop)
				// 			motorActiveChan <- false

				// 			payloadToLights <- PayloadFromDriver{
				// 				CurrentFloor: elevator.CurrentFloor,
				// 				DoorLight:    true,
				// 			}

				// 			doorOpenChan <- true
				// 			elevator.CurrentBehaviour = EBDoorOpen
 
				// 	case elevator.Requests[elevator.CurrentFloor][BHallDown] || elevator.Requests[elevator.CurrentFloor][BCab] || !requestsBelow(elevator):
				// 			if requestsHere(elevator) {
				// 				hwelevio.SetMotorDirection(MDStop)
				// 				motorActiveChan <- false

				// 				payloadToLights <- PayloadFromDriver{
				// 					CurrentFloor: elevator.CurrentFloor,
				// 					DoorLight:    true,
				// 				}

				// 				doorOpenChan <- true
				// 				elevator.CurrentBehaviour = EBDoorOpen
				// 			}
				// 	}
				// }

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
			if elevator.CurrentBehaviour == EBIdle {
				
			//TODO: Does not like it when two orders are placed in the same floor that they both are in

			// -------------------------------------- UNSURE ABOUT THE NEED TO HANDLE THIS -------------------------------------------
			// 	switch  {
			// 	case elevator.Requests[elevator.CurrentFloor][BCab]:
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{
			// 			CurrentFloor: elevator.CurrentFloor,
			// 			DoorLight:    true,
			// 		}
			// 	case elevator.Requests[elevator.CurrentFloor][BHallUp] && elevator.Dirn == MDUp: // Does one elevator know that the other is handeling this order when pressed at the same time? 
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{
			// 			CurrentFloor: elevator.CurrentFloor,
			// 			DoorLight:    true,
			// 		}
			// 	case elevator.Requests[elevator.CurrentFloor][BHallDown] && elevator.Dirn == MDDown:
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{
			// 			CurrentFloor: elevator.CurrentFloor,
			// 			DoorLight:    true,
			// 		}
			// 	case elevator.Requests[elevator.CurrentFloor][BHallDown] && elevator.Dirn == MDUp && !requestsAbove(elevator):
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{
			// 			CurrentFloor: elevator.CurrentFloor,
			// 			DoorLight:    true,
			// 		}
			// 	case elevator.Requests[elevator.CurrentFloor][BHallUp] && elevator.Dirn == MDDown && !requestsBelow(elevator):
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{
			// 			CurrentFloor: elevator.CurrentFloor,
			// 			DoorLight:    true,
			// 		}
			// 	case elevator.Requests[elevator.CurrentFloor][BHallDown] && elevator.Dirn == MDUp && requestsAbove(elevator):
			// 		motorActiveChan <- true
			// 		elevator = ChooseDirection(elevator)
			// 		hwelevio.SetMotorDirection(elevator.Dirn)

			// 	case elevator.Requests[elevator.CurrentFloor][BHallUp] && elevator.Dirn == MDDown && requestsBelow(elevator):
			// 		motorActiveChan <- true
			// 		elevator = ChooseDirection(elevator)
			// 		hwelevio.SetMotorDirection(elevator.Dirn)

			// 	default:
			// 		motorActiveChan <- true
			// 		elevator = ChooseDirection(elevator)
			// 		hwelevio.SetMotorDirection(elevator.Dirn)
			// 	}

				// ----------------------------------------------------------------------------------------------------------------

				if requestsHere(elevator) {
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				}
				if !requestsHere(elevator) {
					motorActiveChan <- true
					elevator = ChooseDirection(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)

				}
			}
			payloadFromElevator <- PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
			}
			ElevatorPrint(elevator)
		}
	}
}
