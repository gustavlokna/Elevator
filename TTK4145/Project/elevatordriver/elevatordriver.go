package elevatordriver

import (
	. "Project/dataenums"
	"Project/elevatordriver/timer"
	"Project/hwelevio"
	"fmt"
	"time"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	payloadToLights chan<- PayloadFromDriver,
) {
	var (
		elevator = initelevator()
		prevelevator = elevator
		floorChannel       = make(chan int)
		obstructionChannel = make(chan bool)
		doorOpenChan       = make(chan bool, 1)
		doorClosedChan     = make(chan bool, 1)
		motorActiveChan    = make(chan bool, 10)
		motorInactiveChan  = make(chan bool, 10)
		clearedRequests    = [NFloors][NButtons]bool{}
		obstruction        bool
	)

	go hwelevio.PollFloorSensor(floorChannel)
	go hwelevio.PollObstructionSwitch(obstructionChannel)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)

	hwelevio.SetMotorDirection(elevator.Dirn)
	elevator.CurrentFloor = <-floorChannel
	elevator.Dirn = MDStop
	hwelevio.SetMotorDirection(elevator.Dirn)

	payload := PayloadFromElevator{ Elevator: elevator, CompletedOrders: clearedRequests}
	payloadFromElevator <- payload
	payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

	for {
		prevelevator = elevator
		var clearedRequests [NFloors][NButtons]bool
		select {
		case elevator.CurrentFloor = <-floorChannel:
			elevator.ActiveSatus = true
			motorActiveChan <- true
			fmt.Println("FLOOR SENSOR TRIGGERED")
			fmt.Println("FLOOR SENSOR TRIGGERED")
			fmt.Println("FLOOR SENSOR TRIGGERED")
			fmt.Println("FLOOR SENSOR TRIGGERED")
			switch {
			case elevator.Requests[elevator.CurrentFloor][BCab]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				doorOpenChan <- true

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				doorOpenChan <- true

			case elevator.Dirn == MDUp && requestsAbove(elevator):
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				doorOpenChan <- true

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				doorOpenChan <- true

			case elevator.Dirn == MDDown && requestsBelow(elevator):
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				motorActiveChan <- false
				doorOpenChan <- true

			default:
				fmt.Println("PENIS ")
				fmt.Println("PENIS ")
				fmt.Println("PENIS ")
				fmt.Println("PENIS ")
				fmt.Println("PENIS ")
				fmt.Println("PENIS ")
				fmt.Println("PENIS ")
				ElevatorPrint(elevator)
				elevator = chooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)
				ElevatorPrint(elevator)
				payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
			}
			fmt.Println("COOOOOCK")
			fmt.Println("COOOOOCK")
			fmt.Println("COOOOOCK")
			// payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		case elevator.Requests = <-newOrderChannel:
			ElevatorPrint(elevator)
			// if elevator.CurrentBehaviour == EBIdle {
			// 		elevator = chooseDirection(elevator)
			// 		hwelevio.SetMotorDirection(elevator.Dirn)
			// }
			// // fmt.Println("NEW ORDER RECEIVED")
			// fmt.Println("NEW ORDER RECEIVED")
			// fmt.Println("NEW ORDER RECEIVED")
			// ElevatorPrint(elevator)
			// if elevator.CurrentBehaviour == EBIdle {
			// 	switch {
			// 	case elevator.Requests[elevator.CurrentFloor][BHallUp]:
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		elevator.Dirn = MDUp
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			// 	case elevator.Requests[elevator.CurrentFloor][BHallDown]:
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		elevator.Dirn = MDDown
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			// 	case elevator.Requests[elevator.CurrentFloor][BCab]:
			// 		elevator.CurrentBehaviour = EBDoorOpen
			// 		doorOpenChan <- true
			// 		payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			// 	default:
			// 		motorActiveChan <- true
			// 		fmt.Println("WE ARE HERE ")
			// 		fmt.Println("WE ARE HERE ")
			// 		fmt.Println("WE ARE HERE ")
			// 		fmt.Println("WE ARE HERE ")
			// 		ElevatorPrint(elevator)
			// 		elevator = chooseDirection(elevator)
			// 		hwelevio.SetMotorDirection(elevator.Dirn)
			// 		ElevatorPrint(elevator)
			// 	}
			// }
			// fmt.Println("HELLO OUSIDE SWITCH")
			// fmt.Println("HELLO OUSIDE SWITCH")
			// fmt.Println("HELLO OUSIDE SWITCH")
			// fmt.Println("HELLO OUSIDE SWITCH")
			// switch elevator.CurrentBehaviour {
			// 	case EBIdle:
			// 		elevator = chooseDirection(elevator)

			// 	case EBMoving:
			// 		hwelevio.SetMotorDirection(elevator.Dirn)
			// 	}
			// payloadFromElevator <- PayloadFromElevator{ Elevator: elevator, CompletedOrders: clearedRequests}

		case <-doorClosedChan:
			if obstruction {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- true
				// payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
				continue
			}

			switch {
			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("Case 1 ")
				ElevatorPrint(elevator)
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BCab] && requestsAbove(elevator):
				fmt.Println("Case 2 ")
				continue

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				fmt.Println("Case 3 ")
				clearedRequests[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				fmt.Println("Case 4 ")
				clearedRequests[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BCab] && requestsBelow(elevator):
				fmt.Println("Case 5 ")
				continue

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("Case 6 ")
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false
			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				clearedRequests[elevator.CurrentFloor][BCab] = true
				elevator.Requests[elevator.CurrentFloor][BCab] = false
			}
			elevator.CurrentBehaviour = EBIdle
			elevator = chooseDirection(elevator)
			//elevator.Dirn = MDStop
			payloadToLights <- PayloadFromDriver{
				CurrentFloor: elevator.CurrentFloor,
				DoorLight:    false,
			}
			// payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		case <-motorInactiveChan:
			if elevator.CurrentBehaviour == EBMoving {
				elevator.ActiveSatus = false
				// payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case obstruction = <-obstructionChannel:
			if elevator.CurrentBehaviour == EBDoorOpen {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- !obstruction
			}
			// payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		default:
			time.Sleep(10 * time.Millisecond)
		}

		switch elevator.CurrentBehaviour {
		case EBIdle:
			elevator = chooseDirection(elevator)

		case EBMoving:
			hwelevio.SetMotorDirection(elevator.Dirn)
		}
		if prevelevator != elevator {
			payload = PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: clearedRequests,
			}
			payloadFromElevator <- payload
		}
	}
}
