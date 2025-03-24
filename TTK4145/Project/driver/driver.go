package driver

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/driver/timer"
	"Project/hwelevio"
)

func Driver(
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
	hwelevio.SetMotorDirection(elevator.Dirn)

	for {

		select {
		case elevator.CurrentFloor = <-floorChan:
			elevator.ActiveStatus = true
			motorActiveChan <- true

			switch {
			case elevator.Requests[elevator.CurrentFloor][BCab]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderAtCurrentFloorInDirn(elevator):
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderCurrentDirn(elevator):
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

			case orderAtCurrentFloorOppositeDirn(elevator):
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = DoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case orderOppositeDirn(elevator):
				localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
				driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

			default:
				elevator.Dirn = MDStop
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
			case orderAtCurrentFloorInDirn(elevator):
				clearedRequests[elevator.CurrentFloor][dirnToBtn(elevator.Dirn)] = true

			case orderCurrentDirn(elevator):

			case orderAtCurrentFloorOppositeDirn(elevator):
				elevator.Dirn = setMotorOppositeDirn(elevator)
				clearedRequests[elevator.CurrentFloor][dirnToBtn(elevator.Dirn)] = true

			case orderOppositeDirn(elevator):

			default:
				elevator.Dirn = MDStop
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
				case orderAtCurrentFloorInDirn(elevator) || elevator.Requests[elevator.CurrentFloor][BCab]:
					elevator.Dirn = btnToDirn(elevator)
					elevator.CurrentBehaviour = DoorOpen
					doorOpenChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case orderCurrentDirn(elevator):
					elevator.CurrentBehaviour = Moving
					hwelevio.SetMotorDirection(elevator.Dirn)
					motorActiveChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

				case orderAtCurrentFloorOppositeDirn(elevator):
					elevator.Dirn = setMotorOppositeDirn(elevator)
					elevator.CurrentBehaviour = DoorOpen
					doorOpenChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case orderOppositeDirn(elevator):
					elevator.CurrentBehaviour = Moving
					elevator.Dirn = setMotorOppositeDirn(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)
					motorActiveChan <- true
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

				default:
					elevator.Dirn = MDStop
					hwelevio.SetMotorDirection(MDStop)
					localLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

				}

			case Moving:
			case DoorOpen:

			}
			driverEvents <- FromDriverToAssigner{Elevator: elevator, CompletedOrders: clearedRequests}

		}
	}
}
