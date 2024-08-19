package elevatordriver
/*
import (
	. "Project/dataenums"
)
*/
func ElevatorDriver(
	fromOrderAssignerChannel <-chan bool,
	toOrderAssignerChannel chan<- bool,
	lifelineChannel chan<- bool,
) {
	var (
		_ = initelevator() 
	)
	print("jeg er inne i driver")
}
