package network

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/network/broadcast"
	"reflect"
	"strconv"
	"time"
)

func Network(worldview <-chan FromAssignerToNetwork,
	stateBroadcast chan<- FromNetworkToAssigner,
	nodeID string) {

	nodeIDInt, _ := strconv.Atoi(nodeID)

	nodeRegistryChannel := make(chan NetworkNodeRegistry)
	broadcastTransmissionChannel := make(chan Message)
	broadcastReceiverChannel := make(chan Message)
	go broadcast.Sender(MessagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(MessagePort, nodeID, broadcastReceiverChannel, nodeRegistryChannel)

	var (
		elevatorList  = initializeElevatorList()
		hallOrderList [NElevators][NFloors][NButtons]ButtonState
		aliveList     [NElevators]bool
		ackMap        [NElevators]bool
		online        bool
		init          bool
	)

	for {
		select {
		case reg := <-nodeRegistryChannel:
			for _, lostNode := range reg.Lost {
				lostNodeInt, _ := strconv.Atoi(lostNode)
				switch {
				case lostNodeInt == nodeIDInt:
					online = false
				default:
					aliveList[lostNodeInt] = false
					hallOrderList[lostNodeInt] = resetHallCalls()
				}
			}
			for _, connectedNode := range reg.New {
				activeNodeInt, _ := strconv.Atoi(connectedNode)
				switch {
				case activeNodeInt == nodeIDInt:
					online = true
					aliveList[activeNodeInt] = true
					hallOrderList[activeNodeInt] = resetHallCalls()
				default:
					aliveList[activeNodeInt] = true
				}
			}

		case msg := <-broadcastReceiverChannel:
			senderId, _ := strconv.Atoi(msg.SenderId)
			ackMap[senderId] = reflect.DeepEqual(elevatorList, msg.ElevatorList) && reflect.DeepEqual(hallOrderList, msg.HallOrderList)

			if !init {
				elevatorList[nodeIDInt] = msg.ElevatorList[nodeIDInt]
				init = true
			}

			aliveList[senderId] = msg.OnlineStatus
			elevatorList[senderId] = msg.ElevatorList[senderId]
			hallOrderList[senderId] = msg.HallOrderList[senderId]
			hallOrderList = cyclicCounter(hallOrderList, nodeIDInt)

			//TODO THIS CAN BE FUNC
			allAcknowledged := true
			for elevator := 0; elevator < NElevators; elevator++ {
				if nodeIDInt == elevator {
					continue
				}
				if aliveList[elevator] && !ackMap[elevator] {
					allAcknowledged = false
					break
				}
			}
			if allAcknowledged {
				for elevator := 0; elevator < NElevators; elevator++ {
					if elevator != nodeIDInt {
						ackMap[elevator] = false
					}
				}
				//printHallOrderList(hallOrderList)
				stateBroadcast <- FromNetworkToAssigner{
					AliveList:     aliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}
			}

		case payload := <-worldview:
			hallOrderList[nodeIDInt] = payload.HallRequests
			aliveList[nodeIDInt] = payload.ActiveStatus
			elevatorList[nodeIDInt] = payload.States[nodeID]

		case <-time.After(BroadcastRate):
			broadcastTransmissionChannel <- Message{
				SenderId:      nodeID,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
				OnlineStatus:  aliveList[nodeIDInt],
				AliveList:     aliveList,
			}

			if !online {
				newAliveList := [NElevators]bool{}
				newAliveList[nodeIDInt] = aliveList[nodeIDInt]
				stateBroadcast <- FromNetworkToAssigner{
					AliveList:     newAliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}

			}

		}
	}
}
