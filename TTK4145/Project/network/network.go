package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Network(messagefromOrderAssigner <-chan FromAssignerToNetwork,
	messagetoOrderAssignerChannel chan<- FromNetworkToAssigner,
	nodeID string) {

	// TODO MAKE CODE COMPATIBLE WITHOUT THIS "STR TO INT CONV"
	nodeIDInt, _ := strconv.Atoi(nodeID)

	nodeRegistryChannel := make(chan NetworkNodeRegistry)
	broadcastTransmissionChannel := make(chan Message)
	broadcastReceiverChannel := make(chan Message)
	go broadcast.Sender(MessagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(MessagePort, nodeID, broadcastReceiverChannel, nodeRegistryChannel)

	var (
		elevatorList   = initializeElevatorList()
		hallOrderList  [NElevators][NFloors][NButtons]ButtonState
		aliveList      [NElevators]bool
		ackMap         [NElevators]bool
		online         bool
		init           bool
	)
	//INIT
	payload := <-messagefromOrderAssigner
	hallOrderList[nodeIDInt] = payload.HallRequests
	aliveList[nodeIDInt] = payload.ActiveSatus
	elevatorList[nodeIDInt] = payload.States[nodeID]

	for {
		select {
		case reg := <-nodeRegistryChannel:
			// TODO THIS reg/ or it can be a double variable CAN ALSO CONTAIN THE ONLIE STATUS :)
			for _, lostNode := range reg.Lost {
				fmt.Printf("Node lost connection: %s\n", lostNode)
				lostNodeInt, _ := strconv.Atoi(lostNode)
				if lostNodeInt == nodeIDInt {
					online = false
				} else {
					fmt.Println("WE SET AN ELEVATOR INACTIVE")
					// check if newOrder = true must be set (but i do not think so)
					aliveList[lostNodeInt] = false
					hallOrderList[lostNodeInt] = resetHallCalls()
				}

			}
			for _, connectedNode := range reg.New {
				fmt.Printf("Node active: %s\n", connectedNode)
				activeNodeInt, _ := strconv.Atoi(connectedNode)
				if activeNodeInt == nodeIDInt {
					hallOrderList[activeNodeInt] = resetHallCalls()
					online = true
				}
				aliveList[activeNodeInt] = true
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
			for i := 0; i < NElevators; i++ {
				if nodeIDInt == i {
					continue
				}
				if aliveList[i] && !ackMap[i] {
					allAcknowledged = false
					break
				}
			}
			if allAcknowledged{
				for i := 0; i < NElevators; i++ {
					if i != nodeIDInt {
						ackMap[i] = false
					}
				}
				printHallOrderList(hallOrderList)
				messagetoOrderAssignerChannel <- FromNetworkToAssigner{
					AliveList:     aliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}
			}

		case payload := <-messagefromOrderAssigner:
			hallOrderList[nodeIDInt] = payload.HallRequests
			aliveList[nodeIDInt] = payload.ActiveSatus
			elevatorList[nodeIDInt] = payload.States[nodeID]

		case <-time.After(10 * time.Millisecond):
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
				messagetoOrderAssignerChannel <- FromNetworkToAssigner{
					AliveList:     newAliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}

			}

		}
	}
}
