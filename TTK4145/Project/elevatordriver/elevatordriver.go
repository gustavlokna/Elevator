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
			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		case <-doorClosedChan:

			if obstruction {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- true
				payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
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
				//do nothing

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
				//do nothing

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("Case 6 ")
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false
			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				clearedRequests[elevator.CurrentFloor][BCab] = true
				elevator.Requests[elevator.CurrentFloor][BCab] = false
			}
			ElevatorPrint(elevator)
			elevator = chooseDirection(elevator)
			ElevatorPrint(elevator)
			//hwelevio.SetMotorDirection(elevator.Dirn)

			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
			payloadToLights <- PayloadFromDriver{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

		case <-motorInactiveChan:
			if elevator.CurrentBehaviour == EBMoving {
				elevator.ActiveSatus = false
				payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case obstruction = <-obstructionChannel:
			fmt.Println("OBSTRUCTION SWITCH TRIGGERED")
			fmt.Println(obstruction)
			if elevator.CurrentBehaviour == EBDoorOpen {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- !obstruction
			}
			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}

		case elevator.Requests = <-newOrderChannel:
			fmt.Println("NEW ORDER RECEIVED")
			fmt.Println("NEW ORDER RECEIVED")
			fmt.Println("NEW ORDER RECEIVED")
			ElevatorPrint(elevator)
			if elevator.CurrentBehaviour == EBIdle{
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

				default:
					motorActiveChan <- true
					fmt.Println("WE ARE HERE ")
					fmt.Println("WE ARE HERE ")
					fmt.Println("WE ARE HERE ")
					fmt.Println("WE ARE HERE ")
					ElevatorPrint(elevator)
					elevator = chooseDirection(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)
					ElevatorPrint(elevator)
				}
			}
			fmt.Println("HELLO OUSIDE SWITCH")
			fmt.Println("HELLO OUSIDE SWITCH")
			fmt.Println("HELLO OUSIDE SWITCH")
			fmt.Println("HELLO OUSIDE SWITCH")
			payloadFromElevator <- PayloadFromElevator{Elevator: elevator, CompletedOrders: clearedRequests}
		}

	}
}
