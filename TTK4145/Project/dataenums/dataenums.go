package dataenums

const (
	NFloors    int = 4
	NButtons   int = 3
	PollRateMS     = 20
)
/*
"not necsessery"
*/
const Addr string = "localhost:15657"


type Button int

const (
	BHallUp Button = iota
	BHallDown
	BCab
)

type HWMotorDirection int 

const (
	MDDown HWMotorDirection = iota - 1
	MDStop
	MDUp
)

type ButtonEvent struct {
	Floor  int
	Button Button
}

type inputDevice struct {
	FloorSensor   func() int
	RequestButton func(f int, btn Button) bool
	StopButton    func() bool
	Obstruction   func() bool
}

type outputDevice struct {
	FloorIndicator     func(f int)
	RequestButton      func(f int, btn Button) bool
	RequestButtonLight func(f int, btn Button, v bool)
	DoorLight          func(v bool)
	StopButtonLight    func(v bool)
	MotorDirection     func(d HWMotorDirection)
}


const (
	ClearRequestVariantConfig = CRVInDirn
	DoorOpenDurationSConfig   = 3
	InputPollRateMsConfig     = 25
	MotorTimeoutS             = 3
)

type ElevatorBehaviour int

const (
	EBIdle ElevatorBehaviour = iota
	EBDoorOpen
	EBMoving
)

//TODO find out what ClearRequestVarient is!
//this is copied code from last project 
type ClearRequestVarient int

const (
	CRVAll ClearRequestVarient = iota
	CRVInDirn
)

type ElevatorConfig struct {
	ClearRequestVariant ClearRequestVarient
	DoorOpenDurationS   float64
}

type DirnBehaviourPair struct {
	Dirn      HWMotorDirection
	Behaviour ElevatorBehaviour
}

type Elevator struct {
	CurrentFloor     int
	Dirn             HWMotorDirection
	Requests         [NFloors][NButtons]bool
	CurrentBehaviour ElevatorBehaviour
	Config           ElevatorConfig
}

type ElevInputDevice struct {
	FloorSensor   func() int
	RequestButton func(f int, btn Button) bool
	StopButton    func() bool
	Obstruction   func() bool
}

type ElevOutputDevice struct {
	FloorIndicator     func(f int)
	RequestButton      func(f int, btn Button) bool
	RequestButtonLight func(f int, btn Button, v bool)
	DoorLight          func(v bool)
	StopButtonLight    func(v bool)
	MotorDirection     func(d HWMotorDirection)
}

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	CounterHallRequests [][2]int 
	States       map[string]HRAElevState `json:"states"`

}