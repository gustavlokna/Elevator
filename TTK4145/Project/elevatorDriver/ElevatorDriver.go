package elevatorDriver

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/elevatorDriver/timer"
	"Project/hwelevio"
)

func ElevatorDriver(
	newOrder <-chan [NFloors][NButtons]bool,
	driverEvents chan<- FromDriverToAssigner,
	localLights chan<- FromDriverToLight,
) {
	var (
		floorChan         = make(chan int, ChannelBufferSize)
		obstructionChan   = make(chan bool, ChannelBufferSize)
		doorOpenChan      = make(chan bool, ChannelBufferSize)
		doorClosedChan    = make(chan bool, ChannelBufferSize)
		motorActiveChan   = make(chan bool, ChannelBufferSize)
		motorInactiveChan = make(chan bool, ChannelBufferSize)
		clearedRequests   = [NFloors][NButtons]bool{}
		obstruction       bool
	)

	go hwelevio.PollFloorSensor(floorChan)
	go hwelevio.PollObstructionSwitch(obstructionChan)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)

	elevator := initelevator()
	hwelevio.SetMotorDirection(elevator.Direction)

	for {

		select {
		case elevator.CurrentFloor = <-floorChan:
			elevator.ActiveStatus = true
			motorActiveChan <- true
			switch {
			case orderAtCurrentFloorInDirection(elevator) || elevator.Requests[elevator.CurrentFloor][BCab]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderCurrentDirection(elevator):
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

			case orderAtCurrentFloorOppositeDirection(elevator):
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderOppositeDirection(elevator):
				elevator.Direction = setMotorOppositeDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Direction)
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

			default:
				elevator.Direction = MDStop
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = Idle
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case <-doorClosedChan:
			if obstruction {
				elevator.ActiveStatus = !obstruction
				doorOpenChan <- true
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
				continue
			}

			switch {
			case orderAtCurrentFloorInDirection(elevator):
				clearedRequests[elevator.CurrentFloor][directionToButton(elevator.Direction)] = true

			case orderAtCurrentFloorOppositeDirection(elevator) && !orderCurrentDirection(elevator):
				elevator.Direction = setMotorOppositeDirection(elevator)
				clearedRequests[elevator.CurrentFloor][directionToButton(elevator.Direction)] = true

			case !elevator.Requests[elevator.CurrentFloor][BCab]:
				elevator.Direction = MDStop
				hwelevio.SetMotorDirection(MDStop)
			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				clearedRequests[elevator.CurrentFloor][BCab] = true
			}
			elevator.CurrentBehaviour = Idle

			driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
			localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
			clearedRequests = [NFloors][NButtons]bool{}

		case <-motorInactiveChan:
			if elevator.CurrentBehaviour == Moving {
				elevator.ActiveStatus = false
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
			}

		case obstruction = <-obstructionChan:
			if elevator.CurrentBehaviour == DoorOpen {
				elevator.ActiveStatus = !obstruction
				doorOpenChan <- !obstruction
			}
			driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

		case elevator.Requests = <-newOrder:
			switch elevator.CurrentBehaviour {
			case Idle:
				switch {
				case orderAtCurrentFloorInDirection(elevator) || elevator.Requests[elevator.CurrentFloor][BCab]:
					elevator.Direction = buttonToDirection(elevator)
					elevator.CurrentBehaviour = DoorOpen
					doorOpenChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case orderCurrentDirection(elevator):
					elevator.CurrentBehaviour = Moving
					hwelevio.SetMotorDirection(elevator.Direction)
					motorActiveChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
					driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

				case orderAtCurrentFloorOppositeDirection(elevator):
					elevator.Direction = setMotorOppositeDirection(elevator)
					elevator.CurrentBehaviour = DoorOpen
					doorOpenChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case orderOppositeDirection(elevator):
					elevator.CurrentBehaviour = Moving
					elevator.Direction = setMotorOppositeDirection(elevator)
					hwelevio.SetMotorDirection(elevator.Direction)
					motorActiveChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
					driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

				default:
					elevator.Direction = MDStop
					hwelevio.SetMotorDirection(MDStop)
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
					driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}
				}
			case Moving:
			case DoorOpen:
			}

		}
	}
}
