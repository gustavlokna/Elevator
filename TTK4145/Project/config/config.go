package config

import "time"

const (
	NFloors             = 4
	NButtons            = 3
	NElevators          = 3
	PollRateMS          = 20 * time.Millisecond
	BroadcastRate       = 50 * time.Millisecond
	DoorOpenDuration    = 3 * time.Second
	MotorTimeout        = 4 * time.Second
	BroadcastBufferSize = 4 * 1024
	ChannelBufferSize   = 100
	HeartbeatInterval   = 1000 * time.Millisecond
	HeartbeatTimeout    = 3000 * time.Millisecond
	Addr                = "localhost:15657"
	MessagePort         = 1338
)
