package dataenums

import . "Project/config"

type Button int

const (
	BHallUp Button = iota
	BHallDown
	BCab
)

type ButtonState int

const (
	Initial ButtonState = iota
	Standby
	ButtonPressed
	OrderAssigned
	OrderComplete
)

type MotorDirection int

const (
	MDDown MotorDirection = iota - 1
	MDStop
	MDUp
)

type ButtonEvent struct {
	Floor  int
	Button Button
}

type ElevatorBehaviour int

const (
	Idle ElevatorBehaviour = iota
	DoorOpen
	Moving
)

type Elevator struct {
	CurrentFloor     int
	Direction        MotorDirection
	Requests         [NFloors][NButtons]bool
	CurrentBehaviour ElevatorBehaviour
	ActiveStatus     bool
}

type HRAElevState struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [NFloors][2]bool        `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

type Message struct {
	SenderId      string
	ElevatorList  [NElevators]HRAElevState
	HallOrderList [NElevators][NFloors][NButtons]ButtonState
	OnlineStatus  bool
	AliveList     [NElevators]bool
}

type FromAssignerToNetwork struct {
	HallRequests [NFloors][NButtons]ButtonState
	States       map[int]HRAElevState
	ActiveStatus bool
}

type FromNetworkToAssigner struct {
	AliveList     [NElevators]bool
	ElevatorList  [NElevators]HRAElevState
	HallOrderList [NElevators][NFloors][NButtons]ButtonState
}

type FromDriverToLight struct {
	CurrentFloor int
	DoorLight    bool
}

type FromDriverToAssigner struct {
	Elevator        Elevator
	CompletedOrders [NFloors][NButtons]bool
}

type NetworkNodeRegistry struct {
	Nodes []string
	New   []string
	Lost  []string
}
