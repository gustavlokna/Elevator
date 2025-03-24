package hwelevio

import (
	. "Project/dataenums"
	. "Project/config"
	"fmt"
	"net"
	"sync"
	"time"
)

var _initialize bool = false
var _mtx sync.Mutex
var _conn net.Conn

func Init(addr string) {
	print("_initialize", _initialize)
	if _initialize {
		fmt.Println("Driver already _initialize!")
		return
	}
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	_initialize = true
}

func SetMotorDirection(dirn MotorDirection) {
	write([4]byte{1, byte(dirn), 0, 0})
}

func SetButtonLamp(button Button, floor int, value bool) {
	write([4]byte{2, byte(button), byte(floor), toByte(value)})
}

func SetFloorIndicator(floor int) {
	write([4]byte{3, byte(floor), 0, 0})
}

func SetDoorOpenLamp(value bool) {
	write([4]byte{4, toByte(value), 0, 0})
}

func PollButtons(receiver chan<- ButtonEvent) {
	prev := make([][3]bool, NFloors)
	for {
		time.Sleep(PollRateMS)
		for floor := 0; floor < NFloors; floor++ {
			for btn := BHallUp; btn <= BCab; btn++ {
				v := getButton(btn, floor)
				if v != prev[floor][btn] && v {
					receiver <- ButtonEvent{floor, Button(btn)}
				}
				prev[floor][btn] = v
			}
		}
	}
}

func PollFloorSensor(receiver chan<- int) {
	prev := -1
	for {
		time.Sleep(PollRateMS)
		v := getFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}

func PollObstructionSwitch(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(PollRateMS)
		v := getObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func getButton(button Button, floor int) bool {
	a := read([4]byte{6, byte(button), byte(floor), 0})
	return toBool(a[1])
}

func getFloor() int {
	a := read([4]byte{7, 0, 0, 0})
	if a[1] != 0 {
		return int(a[2])
	} else {
		return -1
	}
}

func getObstruction() bool {
	a := read([4]byte{9, 0, 0, 0})
	return toBool(a[1])
}

func read(in [4]byte) [4]byte {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	var out [4]byte
	_, err = _conn.Read(out[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	return out
}

func write(in [4]byte) {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}
}

func toByte(a bool) byte {
	var btn byte = 0
	if a {
		btn = 1
	}
	return btn
}

func toBool(a byte) bool {
	var btn bool = false
	if a != 0 {
		btn = true
	}
	return btn
}
