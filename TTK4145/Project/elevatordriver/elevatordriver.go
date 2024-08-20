package elevatordriver

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"time"
)

func ElevatorDriver(
	fromOrderAssignerChannel <-chan bool,
	toOrderAssignerChannel chan<- bool,
	lifelineChannel chan<- bool,
	nodeID int,
) {
	print("Elevator module initiated with name: ", nodeID)

	var (
		elevator = initelevator()
	)
	print("hei")
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//drv_motorActivity := make(chan bool)

	go hwelevio.PollButtons(drv_buttons)
	go hwelevio.PollFloorSensor(drv_floors)
	go hwelevio.PollObstructionSwitch(drv_obstr)
	go hwelevio.PollStopButton(drv_stop)
	//go hwelevio.MontitorMotorActivity(drv_motorActivity, 3.0)
	print("hei")
	for {
		select {
		case <-drv_obstr:
			print("obst")
		case <-drv_buttons:
			ElevatorPrint(elevator)
			toOrderAssignerChannel <- true
			print("buttonevent")
		case <-drv_floors:
			print("floor")

		default:
			time.Sleep(10 * time.Millisecond) // Prevent busy loop
		}
	}
}
