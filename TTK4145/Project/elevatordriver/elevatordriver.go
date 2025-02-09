package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"time"
	"fmt"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	payloadToLights chan <-PayloadFromDriver,
	nodeID string,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator, wasReset = initelevator()
		prevelevator       = elevator
		completedOrders    = [NFloors][NButtons]bool{}
		obstruction        = false
		timerActive        = false
		motorTimeout       time.Time
		doorTimeout        time.Time 
		toggledoorLight    = false 
		doorOpen 			= false
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
	if wasReset{
		elevator.Dirn = MDStop
	}
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
		CurrentFloor : elevator.CurrentFloor,
		DoorLight : toggledoorLight, 
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
			
			print("obst: ", obstruction)
			
		case elevator.CurrentFloor = <-drv_floors:
			motorTimeout = time.Now().Add(3 * time.Second)
			print("etasje: ", elevator.CurrentFloor)

			//send to lightshandler 
			//hwelevio.SetFloorIndicator(elevator.CurrentFloor)
			payloadToLights <- PayloadFromDriver{
				CurrentFloor : elevator.CurrentFloor,
				DoorLight : toggledoorLight, 
			}
		case elevator.Requests = <-newOrderChannel:
			
		default:
			time.Sleep(10 * time.Millisecond)
		}

		switch elevator.CurrentBehaviour {
		case EBIdle:
			elevator = ChooseDirection(elevator)

		case EBMoving:
			hwelevio.SetMotorDirection(elevator.Dirn)
			if timerActive && time.Now().After(motorTimeout){
				elevator.ActiveSatus = false 
				print("motor timeout")
			}
			if elevator.CurrentBehaviour == EBMoving && !timerActive{
				motorTimeout = time.Now().Add(3 * time.Second)
				timerActive = true
			}

			// bolea is copied from Ø and if sentence can just be put in case elevator.CurrentFloor = <-drv_floors:
			if ShouldStop(elevator) && elevator.CurrentFloor != prevelevator.CurrentFloor {
				elevator.ActiveSatus = true 
				timerActive = false
				hwelevio.SetMotorDirection(MDStop)
				// TODO TEST THIS
				//elevator.Dirn = MDStop
				elevator.CurrentBehaviour = EBDoorOpen
				// send to lights boolean that doorlight shows 
				toggledoorLight = true 
				//print("SEND MSG TO LIGHS")
				doorTimeout = time.Now().Add(3*time.Second)
				payloadToLights <- PayloadFromDriver{
					CurrentFloor : elevator.CurrentFloor,
					DoorLight : toggledoorLight, 
				}
				continue
			}

		case EBDoorOpen: // recive back from lights
			if obstruction {
					
				elevator.ActiveSatus = false 
				doorTimeout = time.Now().Add(3*time.Second)
				//add state called obst ? 
				//print("hello we have a obst")
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
					elevator.ActiveSatus = true 
					completedOrders, elevator  = clearAtCurrentFloor(elevator)

					//elevator.Dirn  =  decideDirection(elevator).Dirn
					//elevator.CurrentBehaviour =  decideDirection(elevator).CurrentBehaviour
					// time.Sleep(3 * time.Second) // Simulate door open time
					// SHOULD NOT BE EBIDLE ? 
					//elevator.CurrentBehaviour = EBIdle
					elevator = ChooseDirection(elevator)
					//elevator.Dirn = MDStop
					toggledoorLight = false 
					payloadToLights <- PayloadFromDriver{
						CurrentFloor : elevator.CurrentFloor,
						DoorLight : toggledoorLight, 
					}
					doorOpen = false 
				}
			}


			
		}

		// Update and send PayloadFromElevator if elevator state changes
		if prevelevator != elevator {
			payload = PayloadFromElevator{
				Elevator:        elevator,
				CompletedOrders: completedOrders,
			}
			payloadFromElevator <- payload
			saveElevator(elevator)
		}
		
	}
}
