package main

import (
	"Project/hwelevio"
	"Project/elev"
)

func main() {

	numFloors := 4

	hwelevio.Init("localhost:15657", numFloors)

	var d hwelevio.MotorDirection = hwelevio.MD_Up
	//hwelevio.SetMotorDirection(d)

	elevator.Init_elevator_logic(numFloors, d)

	select {}
}