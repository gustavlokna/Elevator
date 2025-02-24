package elevatordriver

import (
	. "Project/dataenums"
	"Project/elevatordriver/timer"
	"Project/hwelevio"
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
	// var obstruction = <-obstructionChannel
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
			switch elevator.CurrentBehaviour {
			case EBMoving:
				motorActiveChan <- true
				switch elevator.Dirn {
				case MDUp:
					
					//TODO: Understand why we are chechking for request above and below in this case?

					//if elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BCab] || requestsAbove(elevator) {
					if elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BCab] {
						hwelevio.SetMotorDirection(MDStop)
						motorActiveChan <- false

						payloadFromElevator <- PayloadFromElevator{
							Elevator:        elevator,
							CompletedOrders: clearedRequests,
						}
						payloadToLights <- PayloadFromDriver{
							CurrentFloor: elevator.CurrentFloor,
							DoorLight:    true,
						}

						doorOpenChan <- true	
						elevator.CurrentBehaviour = EBDoorOpen
					}
				case MDDown:
					//if elevator.Requests[elevator.CurrentFloor][BHallDown] || elevator.Requests[elevator.CurrentFloor][BCab] || requestsBelow(elevator) {
					if elevator.Requests[elevator.CurrentFloor][BHallDown] || elevator.Requests[elevator.CurrentFloor][BCab] {
						hwelevio.SetMotorDirection(MDStop)
						motorActiveChan <- false

						payloadFromElevator <- PayloadFromElevator{
							Elevator:        elevator,
							CompletedOrders: clearedRequests,
						}
						payloadToLights <- PayloadFromDriver{
							CurrentFloor: elevator.CurrentFloor,
							DoorLight:    true,
						} 

						doorOpenChan <- true
						elevator.CurrentBehaviour = EBDoorOpen

					}
				}

// -------------------------------- DUE TO INITIALIZING ERRORS --------------------------------------
			default:
				hwelevio.SetMotorDirection(MDStop)
				//TODO: Maybe write this into a struct from the beginning to make it more clean
				
				elevator = ChooseDirection(elevator)

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
			//TODO: Might be kicking in to late if obstruction is turned on and of while the initial dooropen-timer is counting
			if obstruction {
				// elevator.ActiveSatus = !obstruction
				// payloadToLights <- PayloadFromDriver{
				// 	CurrentFloor: elevator.CurrentFloor,
				// 	DoorLight:    true,
				// }
				// payloadFromElevator <- PayloadFromElevator{
				// 	Elevator:        elevator,
				// 	CompletedOrders: clearedRequests,
				// }
				continue
			}
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
			//TODO: This does not work. 
			elevator.ActiveSatus = !obstruction
			if obstruction || elevator.CurrentBehaviour == EBDoorOpen {
				payloadToLights <- PayloadFromDriver{
					CurrentFloor: elevator.CurrentFloor,
					DoorLight:    true,
				}
			} else {
				payloadToLights <- PayloadFromDriver{
					CurrentFloor: elevator.CurrentFloor,
					DoorLight:    false,
				}
			}
			payloadFromElevator <- PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
			}


		case elevator.Requests = <-newOrderChannel:
			if elevator.CurrentBehaviour == EBIdle {
				if elevator.Requests[elevator.CurrentFloor][BCab] {
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- PayloadFromDriver{
						CurrentFloor: elevator.CurrentFloor,
						DoorLight:    true,
					}
				} else {
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
			} 
		}
	}
}