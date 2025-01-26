package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"time"
)

func ElevatorDriver(
	newOrderChannel <-chan [NFloors][NButtons]bool,
	payloadFromElevator chan<- PayloadFromElevator,
	payloadToLights chan <-PayloadFromDriver,
	nodeID string,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator       = initelevator()
		prevelevator   = elevator
		completedOrders = [NFloors][NButtons]bool{}
		obstruction    = false
		// somewhat compied fromØ 
		doorTimeout  time.Time 
		toggledoorLight = false 
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
	elevator.Dirn = MDStop
	hwelevio.SetMotorDirection(elevator.Dirn)
	
	// Initialize and send initial PayloadFromElevator
	// Initialize and send initial PayloadFromElevator
	payload := PayloadFromElevator{
		Elevator:        elevator,
		CompletedOrders: completedOrders,
	}
	payloadFromElevator <- payload


	for {
		prevelevator = elevator
		select {
		case obstruction = <-drv_obstr:
			print("obst: ", obstruction)
		case elevator.CurrentFloor = <-drv_floors:
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
			hwelevio.SetMotorDirection(elevator.Dirn)

		case EBMoving:
			if ShouldStop(elevator) {
				hwelevio.SetMotorDirection(MDStop)
				elevator.Dirn = MDStop
				elevator.CurrentBehaviour = EBDoorOpen
				// send to lights boolean that doorlight shows 
				toggledoorLight = true 
				print("SEND MSG TO LIGHS")
				doorTimeout = time.Now().Add(3*time.Second)
				payloadToLights <- PayloadFromDriver{
					CurrentFloor : elevator.CurrentFloor,
					DoorLight : toggledoorLight, 
				}
				continue
			}

		case EBDoorOpen: // recive back from lights 
			//ElevatorPrint(elevator)
			if obstruction {
				doorTimeout = time.Now().Add(3*time.Second)
				//add state called obst ? 
				print("hello we have a obst")
			} else {
				if time.Now().After(doorTimeout){
					
					completedOrders = ClearAtCurrentFloor(elevator)
					// time.Sleep(3 * time.Second) // Simulate door open time
					elevator.CurrentBehaviour = EBIdle
					toggledoorLight = false 
					payloadToLights <- PayloadFromDriver{
						CurrentFloor : elevator.CurrentFloor,
						DoorLight : toggledoorLight, 
					}
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
			
		}
	}
}
