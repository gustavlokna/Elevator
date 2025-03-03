package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"Project/elevatordriver/timer"
	"fmt"
	"time"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	payloadToLights chan<- PayloadFromDriver,
	nodeID string,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator        = initelevator()
		prevelevator    = elevator
		completedOrders = [NFloors][NButtons]bool{}
		obstruction     = false
		timerActive     = false
		motorTimeout    time.Time
		toggledoorLight = false
	)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	doorClosedChan := make(chan bool, 1) // buffered channel
	doorOpenChan := make(chan bool, 1)
	motorActiveChan    := make(chan bool,10)
	motorInactiveChan  := make(chan bool,10)

	go hwelevio.PollFloorSensor(drv_floors)
	go hwelevio.PollObstructionSwitch(drv_obstr)
	go hwelevio.PollStopButton(drv_stop)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)

	fmt.Println("HALLO")

	// Initialization of elevator
	hwelevio.SetMotorDirection(elevator.Dirn)
	elevator.CurrentFloor = <-drv_floors
	elevator.Dirn = MDStop
	hwelevio.SetMotorDirection(elevator.Dirn)

	// Initialize and send initial PayloadFromElevator
	// Initialize and send initial PayloadFromElevator
	// TODO NO NEED TO DeFINE THIS
	payload := PayloadFromElevator{
		Elevator:        elevator,
		CompletedOrders: completedOrders,
	}
	payloadFromElevator <- payload

	payloadToLights <- PayloadFromDriver{
		CurrentFloor: elevator.CurrentFloor,
		DoorLight:    toggledoorLight,
	}

	for {
		prevelevator = elevator
		//dobel init : not good
		//we need smart way to avouid
		// what i want is to set it false on the start of the loop :)
		// else prev is stored
		var completedOrders [NFloors][NButtons]bool

		select {
		case obstruction = <-drv_obstr:
			//ElevatorPrint(elevator)
			print("obst: ", obstruction)

		case elevator.CurrentFloor = <-drv_floors:
			motorTimeout = time.Now().Add(3 * time.Second)
			print("etasje: ", elevator.CurrentFloor)

			//send to lightshandler
			//hwelevio.SetFloorIndicator(elevator.CurrentFloor)
			payloadToLights <- PayloadFromDriver{
				CurrentFloor: elevator.CurrentFloor,
				DoorLight:    toggledoorLight,
			}
			elevator.ActiveSatus = true
			motorActiveChan <- true
			fmt.Println("FLOOR SENSOR TRIGGERED")
			fmt.Println("FLOOR SENSOR TRIGGERED")
			fmt.Println("FLOOR SENSOR TRIGGERED")
			fmt.Println("FLOOR SENSOR TRIGGERED")
			switch {
			case elevator.Requests[elevator.CurrentFloor][BCab]:
				fmt.Println("CAB REQUEST TRIGGERED")
				fmt.Println("CAB REQUEST TRIGGERED")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("UP REQUEST TRIGGERED")
				fmt.Println("UP REQUEST TRIGGERED")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true

			case elevator.Dirn == MDUp && requestsAbove(elevator):
				fmt.Println("UP REQUEST ABOVE TRIGGERED")
				fmt.Println("UP REQUEST ABOVE TRIGGERED")

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				fmt.Println("MDUP DOWN REQUEST TRIGGERED")
				fmt.Println("MDUP DOWN REQUEST TRIGGERED")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				fmt.Println("DOWN REQUEST TRIGGERED")
				fmt.Println("DOWN REQUEST TRIGGERED")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true

			case elevator.Dirn == MDDown && requestsBelow(elevator):
				fmt.Println("DOWN REQUEST BELOW TRIGGERED")
				fmt.Println("DOWN REQUEST BELOW TRIGGERED")

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("DOWN REQUEST UP TRIGGERED")
				fmt.Println("DOWN REQUEST UP TRIGGERED")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
			}
		case elevator.Requests = <-newOrderChannel:
			ElevatorPrint(elevator)
			

		case <-doorClosedChan:
			fmt.Println("DOR CLOSE")

			elevator.ActiveSatus = !obstruction
			if obstruction {
				doorOpenChan <- true
				continue
			}

			switch {
			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("Case 1 ")
				completedOrders[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BCab] && requestsAbove(elevator):
				fmt.Println("Case 2 ")
				//do nothing

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				fmt.Println("Case 3 ")
				completedOrders[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				fmt.Println("Case 4 ")
				completedOrders[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BCab] && requestsBelow(elevator):
				fmt.Println("Case 5 ")
				//do nothing

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				fmt.Println("Case 6 ")
				completedOrders[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				completedOrders[elevator.CurrentFloor][BCab] = true
				elevator.Requests[elevator.CurrentFloor][BCab] = false

			}
			fmt.Print("PENIS")
			elevator.CurrentBehaviour = EBIdle
			elevator = ChooseDirection(elevator)
			//elevator.Dirn = MDStop
			toggledoorLight = false
			payloadToLights <- PayloadFromDriver{
				CurrentFloor: elevator.CurrentFloor,
				DoorLight:    toggledoorLight,
			}
			//doorOpen = false
			fmt.Print("SEND MSG TO LIGHS")

		default:
			time.Sleep(10 * time.Millisecond)
		}

		switch elevator.CurrentBehaviour {
		case EBIdle:
			elevator = ChooseDirection(elevator)

		case EBMoving:
			hwelevio.SetMotorDirection(elevator.Dirn)
			if timerActive && time.Now().After(motorTimeout) {
				//elevator.ActiveSatus = false
				//print("motor timeout")
			}
			if elevator.CurrentBehaviour == EBMoving && !timerActive {
				motorTimeout = time.Now().Add(3 * time.Second)
				timerActive = true
			}

			// bolea is copied from Ø and if sentence can just be put in case elevator.CurrentFloor = <-drv_floors:
			
			if ShouldStop(elevator) && elevator.CurrentFloor != prevelevator.CurrentFloor {
				/*
				elevator.ActiveSatus = true
				timerActive = false
				hwelevio.SetMotorDirection(MDStop)
				// TODO TEST THIS
				//elevator.Dirn = MDStop
				elevator.CurrentBehaviour = EBDoorOpen
				// send to lights boolean that doorlight shows
				toggledoorLight = true
				//print("SEND MSG TO LIGHS")
				doorOpenChan <- true
				payloadToLights <- PayloadFromDriver{
					CurrentFloor: elevator.CurrentFloor,
					DoorLight:    toggledoorLight,
				}
				*/
				continue
			}
			/*
				case EBDoorOpen: // recive back from lights
					if obstruction {

						elevator.ActiveSatus = false
						doorTimeout = time.Now().Add(3*time.Second)
						//add state called obst ?

					}
					//This logic is copied from Ø
					if !doorOpen{
						fmt.Println("START TIMER")
						doorOpen = true
						doorTimeout = time.Now().Add(3*time.Second)
						toggledoorLight = true
						payloadToLights <- PayloadFromDriver{
							CurrentFloor : elevator.CurrentFloor,
							DoorLight : toggledoorLight,
						}
					}  else {
						if time.Now().After(doorTimeout){
							fmt.Println("DOOR CLOSED")
							doorClosedChan <- true
							fmt.Println("DOOR CLOSED")
							elevator.ActiveSatus = true

						}
					}
			*/

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
