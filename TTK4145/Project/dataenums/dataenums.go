package dataenums

const (
	NFloors                 int = 4
	NButtons                int = 3
	PollRateMS                  = 20
	NUM_ELEVATORS           int = 3
	DoorOpenDurationSConfig int = 3
	MotorTimeoutS           int = 4 // TODO MAKE 3 s (worked on slow elevs)
)

// THE addr used to int helwvio
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

type ElevatorBehaviour int

const (
	EBIdle ElevatorBehaviour = iota
	EBDoorOpen
	EBMoving
)

type DirnBehaviourPair struct {
	Dirn      HWMotorDirection
	Behaviour ElevatorBehaviour
}

type Elevator struct {
	CurrentFloor     int
	Dirn             HWMotorDirection
	Requests         [NFloors][NButtons]bool
	CurrentBehaviour ElevatorBehaviour
	ActiveSatus      bool
}

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [NFloors][2]bool `json:"hallRequests"`
	States map[string]HRAElevState `json:"states"`
}

// TODO MUST GIVE NEW NAMES 
type Message struct {
	//TODO: Make int
	SenderId      string // IPv4
	ElevatorList  [NUM_ELEVATORS]HRAElevState
	HallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState
	OnlineStatus  bool
	AliveList     [NUM_ELEVATORS]bool
}

type PayloadFromElevator struct {
	Elevator        Elevator
	CompletedOrders [NFloors][NButtons]bool
}

type PayloadFromassignerToNetwork struct {
	//TODO Is not just hallRequests. Name does not fit is also cab
	HallRequests [NFloors][NButtons]ButtonState `json:"hallRequests"`
	States       map[string]HRAElevState        `json:"states"`
	ActiveSatus  bool
}

type PayloadFromNetworkToAssigner struct {
	AliveList    [NUM_ELEVATORS]bool
	ElevatorList [NUM_ELEVATORS]HRAElevState
	//TODO IS NOT just HallORders ?
	HallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState
}

type PayloadFromDriver struct {
	CurrentFloor int
	DoorLight    bool
}

type NetworkNodeRegistry struct {
	Nodes []string
	New   []string
	Lost  []string
}
