package config

import "time"

const (
	NFloors           = 4
	NButtons          = 3
	NElevators        = 3
	PollRateMS        = 20 * time.Millisecond
	DoorOpenDurationS = 3 * time.Second
	MotorTimeoutS     = 4 * time.Second
	BufferSize        = 4 * 1024

	HeartbeatInterval = 150 * time.Millisecond
	HeartbeatTimeout  = 3000 * time.Millisecond

	Addr        = "localhost:15657"
	MessagePort = 1338
)
