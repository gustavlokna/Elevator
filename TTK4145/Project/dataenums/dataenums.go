package dataenums

const (
	NFloors    int = 4
	NButtons   int = 3
	PollRateMS     = 20
 	NUM_ELEVATORS int = 3
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

type ButtonState int

const (
	Initial ButtonState = iota
	Idle 
	ButtonPressed
	OrderAssigned
	OrderComplete
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

// TODO find out what ClearRequestVarient is!
// this is copied code from last project
// TODO THIS IS ASS CODE 
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
	//todo do not need this ? 
	Requests         [NFloors][NButtons]bool
	CurrentBehaviour ElevatorBehaviour
	//what is this condfig++
	ActiveSatus 	bool 
	Config           ElevatorConfig
}

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	//TODO:  Make variabels HallRequests        [][2]ButtonState `json:"hallRequests"` 
	//HallRequests        [NFloors][NButtons]bool `json:"hallRequests"` 
	HallRequests        [NFloors][2]bool `json:"hallRequests"` 
	//NOTE THIS WAS BOOL AND TO USE ASSIGNER MUST BE CONVERTED TO INTEGER 
	//Coul BE SIMPLE IN ASSIGNER IF HallRequests == OrderAssigned then 1, else 0 
	States              map[string]HRAElevState `json:"states"`
}

type Message struct {
	//TODO: Make int
	SenderId     string // IPv4
	ElevatorList  [NUM_ELEVATORS]HRAElevState
	HallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState
	OnlineStatus bool
}

type PayloadFromElevator struct {
	Elevator        Elevator
	CompletedOrders [NFloors][NButtons]bool
}

type PayloadFromassignerToNetwork struct{
	//TODO Is not just hallRequests. Name does not fit is also cab 
	HallRequests        [NFloors][NButtons]ButtonState `json:"hallRequests"` 
	States              map[string]HRAElevState `json:"states"`
	ActiveSatus 		bool 
}

type PayloadFromNetworkToAssigner struct {
	AliveList     [NUM_ELEVATORS]bool
	ElevatorList  [NUM_ELEVATORS]HRAElevState
	//TODO IS NOT just HallORders ? 
	HallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState
}

type PayloadFromDriver struct {
	CurrentFloor int
	DoorLight bool 
}

type NetworkNodeRegistry struct {
	Nodes []string
	New   string
	Lost  []string
}