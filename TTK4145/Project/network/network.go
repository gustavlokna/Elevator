package network

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/network/broadcast"
	"reflect"
	"time"
)

func Network(worldview <-chan FromAssignerToNetwork,
	stateBroadcast chan<- FromNetworkToAssigner,
	nodeID int) {

	var (
		nodeRegistryChannel          = make(chan NetworkNodeRegistry)
		broadcastTransmissionChannel = make(chan Message)
		broadcastReceiverChannel     = make(chan Message)
		elevatorList                 = initializeElevatorList()
		hallOrderList                [NElevators][NFloors][NButtons]ButtonState
		aliveList                    [NElevators]bool
		ackMap                       [NElevators]bool
		online                       bool
		init                         bool
	)

	go broadcast.Sender(MessagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(MessagePort, nodeID, broadcastReceiverChannel, nodeRegistryChannel)

	for {
		select {
		case reg := <-nodeRegistryChannel:
			for _, lostNode := range reg.Lost {
				switch {
				case lostNode == nodeID:
					online = false
				default:
					aliveList[lostNode] = false
					hallOrderList[lostNode] = resetHallCalls()
				}
			}
			for _, connectedNode := range reg.New {
				switch {
				case connectedNode == nodeID:
					online = true
					aliveList[connectedNode] = true
					hallOrderList[connectedNode] = resetHallCalls()
				default:
					aliveList[connectedNode] = true
				}
			}

		case msg := <-broadcastReceiverChannel:
			ackMap[msg.SenderId] = reflect.DeepEqual(elevatorList, msg.ElevatorList) && reflect.DeepEqual(hallOrderList, msg.HallOrderList)

			if !init {
				elevatorList[nodeID] = msg.ElevatorList[nodeID]
				init = true
			}

			aliveList[msg.SenderId] = msg.OnlineStatus
			elevatorList[msg.SenderId] = msg.ElevatorList[msg.SenderId]
			hallOrderList[msg.SenderId] = msg.HallOrderList[msg.SenderId]
			hallOrderList = cyclicCounter(hallOrderList, nodeID)

			if allAcknowledged(ackMap, aliveList, nodeID) {
				for elevator := 0; elevator < NElevators; elevator++ {
					if elevator != nodeID {
						ackMap[elevator] = false
					}
				}
				stateBroadcast <- FromNetworkToAssigner{
					AliveList:     aliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}
			}

		case payload := <-worldview:
			hallOrderList[nodeID] = payload.HallRequests
			aliveList[nodeID] = payload.ActiveStatus
			elevatorList[nodeID] = payload.States[nodeID]

		case <-time.After(BroadcastRate):
			broadcastTransmissionChannel <- Message{
				SenderId:      nodeID,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
				OnlineStatus:  aliveList[nodeID],
				AliveList:     aliveList,
			}

			if !online {
				newAliveList := [NElevators]bool{}
				newAliveList[nodeID] = aliveList[nodeID]
				stateBroadcast <- FromNetworkToAssigner{
					AliveList:     newAliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}

			}

		}
	}
}
