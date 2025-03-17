package elevator

import (
	"Project/hwelevio"
	"fmt"
	"time"
)

type Order struct {
	UP     [4]int
	DOWN   [4]int
	INSIDE [4]int
}

const (
	STATE_MOVING_UP            = iota
	STATE_MOVING_DOWN          = iota
	STATE_MOVING_UP_TO_GO_DOWN = iota
	STATE_MOVING_DOWN_TO_GO_UP = iota
	STATE_UNDEFINED            = iota

	//STATE_STANDING_STILL_MOVING_UP   = iota
	//STATE_STANDING_STILL_MOVING_DOWN = iota
)

func addOrder(bEvent hwelevio.ButtonEvent, order *Order) {
	switch bEvent.Button {
	case hwelevio.BT_HallUp:
		order.UP[bEvent.Floor] = 1
	case hwelevio.BT_HallDown:
		order.DOWN[bEvent.Floor] = 1
	case hwelevio.BT_Cab:
		order.INSIDE[bEvent.Floor] = 1
	}
	hwelevio.SetButtonLamp(bEvent.Button, bEvent.Floor, true)
}

func removeOrder(bEvent hwelevio.ButtonEvent, order *Order) {
	switch bEvent.Button {
	case hwelevio.BT_HallUp:
		order.UP[bEvent.Floor] = 0
	case hwelevio.BT_HallDown:
		order.DOWN[bEvent.Floor] = 0
	case hwelevio.BT_Cab:
		order.INSIDE[bEvent.Floor] = 0
	}
	hwelevio.SetButtonLamp(bEvent.Button, bEvent.Floor, false)

}

func getMotorDirection(currentFloor int, destinationFloor int) hwelevio.MotorDirection {
	if currentFloor < destinationFloor {
		return hwelevio.MD_Up
	} else if currentFloor > destinationFloor {
		return hwelevio.MD_Down
	} else {
		return hwelevio.MD_Stop
	}
}

// Vi må endre slik at *state settes som en del av getMotorDirection slik at den blir satt riktig
// se andre for loop - state er ikke nødvendigvis "STATE_MOVING_DOWN", gleder også for andre retning.

func chooseDirection(state *int, order *Order, currentFloor int) hwelevio.MotorDirection {
	// Check if the elevator is moving up or needs to go up to serve a down request
	if *state == STATE_MOVING_UP || *state == STATE_MOVING_UP_TO_GO_DOWN {
		// Check for any orders above the current floor
		for f := currentFloor + 1; f < 4; f++ {
			if order.UP[f] == 1 || order.INSIDE[f] == 1 {
				*state = STATE_MOVING_UP
				return hwelevio.MD_Up
			}
		}
		// Check for any orders below the current floor
		for f := 3; f >= 0; f-- {
			if order.DOWN[f] == 1 || order.INSIDE[f] == 1 {
				*state = STATE_MOVING_DOWN
				return getMotorDirection(currentFloor, f)
			}
		}
		// Check for any up orders below the current floor
		for f := 0; f < currentFloor; f++ {
			if order.UP[f] == 1 {
				*state = STATE_MOVING_DOWN_TO_GO_UP
				return hwelevio.MD_Down
			}
		}
	} else if *state == STATE_MOVING_DOWN || *state == STATE_MOVING_DOWN_TO_GO_UP {
		// Check for any orders below the current floor
		for f := currentFloor - 1; f >= 0; f-- {
			if order.DOWN[f] == 1 || order.INSIDE[f] == 1 {
				*state = STATE_MOVING_DOWN
				return hwelevio.MD_Down
			}
		}
		// Check for any orders above the current floor
		for f := 0; f < 4; f++ {
			if order.UP[f] == 1 || order.INSIDE[f] == 1 {
				*state = STATE_MOVING_UP
				return getMotorDirection(currentFloor, f)
			}
		}
		// Check for any down orders above the current floor
		for f := 3; f > currentFloor; f-- {
			if order.DOWN[f] == 1 {
				*state = STATE_MOVING_UP_TO_GO_DOWN
				return hwelevio.MD_Up
			}
		}
	}

	// If no orders are found, stop the elevator
	return hwelevio.MD_Stop

}
func pickup_from_curr_floor(current_floor int, order *Order, state int) bool {
	// Checks if the current floor has an order that fits with the elevators direction. Should we pick up?
	if order.INSIDE[current_floor] == 1 {
		return true
	}

	if (state == STATE_MOVING_UP || state == STATE_MOVING_DOWN_TO_GO_UP) && order.UP[current_floor] == 1 {
		return true
	}

	if (state == STATE_MOVING_DOWN || state == STATE_MOVING_UP_TO_GO_DOWN) && order.DOWN[current_floor] == 1 {
		return true
	}

	return false
}

func serve_order(current_floor int, order *Order, state int, obstruction chan bool, ready_to_move chan bool) {
	// Remove cab order
	removeOrder(hwelevio.ButtonEvent{Floor: current_floor, Button: hwelevio.BT_Cab}, order)

	// Remove hall Down order
	if state == STATE_MOVING_DOWN {
		removeOrder(hwelevio.ButtonEvent{Floor: current_floor, Button: hwelevio.BT_HallDown}, order)
	}
	// Remove hall Up order
	if state == STATE_MOVING_UP {
		removeOrder(hwelevio.ButtonEvent{Floor: current_floor, Button: hwelevio.BT_HallUp}, order)
	}

	// Door open
	hwelevio.SetDoorOpenLamp(true)
	three_sec_delay := time.NewTimer(time.Second * 3)
	is_obstructed := false

	// Obstruction
	for {
		select {
		case is_obstructed = <-obstruction:
			if !is_obstructed {
				three_sec_delay.Reset(time.Second * 3)
			}

		case <-three_sec_delay.C:
			if !is_obstructed {
				hwelevio.SetDoorOpenLamp(false)
				ready_to_move <- true
				return
			}
		}
	}
}

func elevator_logic_loop(state int, order Order, current_floor int,
	drv_buttons chan hwelevio.ButtonEvent, drv_floors chan int, drv_obstr chan bool, drv_stop chan bool, ready_to_move chan bool) {

	standing_still := true
	can_move := true

	for {
		// Wait on multiple communcation operations.
		select {
		// Receive operation on the channel drv_buttons
		case buttonEvent := <-drv_buttons:
			addOrder(buttonEvent, &order)

			// Set motor direction to Stop and serve the order
			if standing_still && buttonEvent.Floor == current_floor {
				hwelevio.SetMotorDirection(hwelevio.MD_Stop)
				can_move = false
				go serve_order(current_floor, &order, state, drv_obstr, ready_to_move)
			}

			fmt.Printf("%+v\n", order)
			//
			if can_move {
				direction := chooseDirection(&state, &order, current_floor)
				hwelevio.SetMotorDirection(direction)
				//standing_still = (direction == hwelevio.MD_Stop) // PROBLEM: Når vil dette egentlig inntreffe
			}

		case floorSensorReading := <-drv_floors:
			// Sensor reading
			fmt.Printf("%+v\n", floorSensorReading)
			current_floor = floorSensorReading
			hwelevio.SetFloorIndicator(current_floor)
			if pickup_from_curr_floor(current_floor, &order, state) {
				hwelevio.SetMotorDirection(hwelevio.MD_Stop)
				can_move = false
				go serve_order(current_floor, &order, state, drv_obstr, ready_to_move)

			}
		// Elevator is ready to move. No obstruction is active
		case <-ready_to_move:
			can_move = true
			direction := chooseDirection(&state, &order, current_floor)
			hwelevio.SetMotorDirection(direction)
			standing_still = direction == hwelevio.MD_Stop

			if pickup_from_curr_floor(current_floor, &order, state) {
				hwelevio.SetMotorDirection(hwelevio.MD_Stop)
				can_move = false
				go serve_order(current_floor, &order, state, drv_obstr, ready_to_move)
			}

		case stopButtonEvent := <-drv_stop:
			fmt.Printf("%+v\n", stopButtonEvent)

		}
	}

}

func Init_elevator_logic(numFloors int, d hwelevio.MotorDirection) {

	var state int = STATE_UNDEFINED
	var order Order = Order{[4]int{0, 0, 0, 0}, [4]int{0, 0, 0, 0}, [4]int{0, 0, 0, 0}}
	var currentFloor int

	drv_buttons := make(chan hwelevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ready_to_move := make(chan bool)

	go hwelevio.PollButtons(drv_buttons)
	go hwelevio.PollFloorSensor(drv_floors)
	go hwelevio.PollObstructionSwitch(drv_obstr)
	go hwelevio.PollStopButton(drv_stop)

	for button := 0; button < 3; button++ {
		for floor := 0; floor < 4; floor++ {
			hwelevio.SetButtonLamp(hwelevio.ButtonType(button), floor, false)
		}
	}

	hwelevio.SetDoorOpenLamp(false)
	hwelevio.SetMotorDirection(hwelevio.MD_Down)

	currentFloor = <-drv_floors // Wait for the elevator to reach a floor before starting the elevator logic

	state = STATE_MOVING_UP
	hwelevio.SetMotorDirection(hwelevio.MD_Stop)
	hwelevio.SetFloorIndicator(currentFloor)

	go elevator_logic_loop(state, order, currentFloor, drv_buttons, drv_floors, drv_obstr, drv_stop, ready_to_move)
}

//mangler stoppknapp og ideell håndtering av ordre
