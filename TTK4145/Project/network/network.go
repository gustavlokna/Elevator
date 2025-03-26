package network

import (
	. "Project/config"
	. "Project/dataenums"
	"Project/network/broadcast"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Network(worldview <-chan FromAssignerToNetwork,
	stateBroadcast chan<- FromNetworkToAssigner,
	nodeID int) {

	var (
		nodeRegistryChannel          = make(chan NetworkNodeRegistry)
		broadcastTransmissionChannel = make(chan Message)
		broadcastReceiverChannel     = make(chan Message)
		elevatorList                 = initElevatorList()
		hallOrderList                [NElevators][NFloors][NButtons]ButtonState
		aliveList                    [NElevators]bool
		ackMap                       [NElevators]bool
		online                       bool
		init                         bool
	)

	go broadcast.Sender(MessagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(MessagePort, strconv.Itoa(nodeID), broadcastReceiverChannel, nodeRegistryChannel)

	for {
		select {
		case reg := <-nodeRegistryChannel:
			for _, lostNode := range reg.Lost {
				lostNodeInt, _ := strconv.Atoi(lostNode)
				switch {
				case lostNodeInt == nodeID:
					fmt.Printf("⚠️  This node (ID %d) marked offline\n", nodeID)
					online = false
				default:
					fmt.Printf("❌ Node %d lost connection\n", lostNodeInt)
					aliveList[lostNodeInt] = false
					hallOrderList[lostNodeInt] = resetHallCalls()
				}
			}

			for _, connectedNode := range reg.New {
				activeNodeInt, _ := strconv.Atoi(connectedNode)
				switch {
				case activeNodeInt == nodeID:
					fmt.Printf("✅ This node (ID %d) is now online\n", nodeID)
					online = true
					aliveList[activeNodeInt] = true
					hallOrderList[activeNodeInt] = resetHallCalls()
				default:
					fmt.Printf("➕ Node %d joined the network\n", activeNodeInt)
					aliveList[activeNodeInt] = true
				}
			}

		case msg := <-broadcastReceiverChannel:
			senderId, _ := strconv.Atoi(msg.SenderId)
			ackMap[senderId] = reflect.DeepEqual(elevatorList, msg.ElevatorList) && reflect.DeepEqual(hallOrderList, msg.HallOrderList)

			if !init {
				elevatorList[nodeID] = msg.ElevatorList[nodeID]
				init = true
			}

			aliveList[senderId] = msg.OnlineStatus
			elevatorList[senderId] = msg.ElevatorList[senderId]
			hallOrderList[senderId] = msg.HallOrderList[senderId]
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
				SenderId:      strconv.Itoa(nodeID),
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
