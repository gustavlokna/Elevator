package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// TODO MOVE DATA ENUMS ?
const messagePort int = 1338

func Network(messagefromOrderAssigner <-chan PayloadFromassignerToNetwork,
	messagetoOrderAssignerChannel chan<- PayloadFromNetworkToAssigner,
	nodeID string) {

	// TODO MAKE CODE COMPATIBLE WITHOUT THIS "STR TO INT CONV"
	nodeIDInt, _ := strconv.Atoi(nodeID)

	nodeRegistryChannel := make(chan NetworkNodeRegistry)
	broadcastTransmissionChannel := make(chan Message)
	broadcastReceiverChannel := make(chan Message)
	go broadcast.Sender(messagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(messagePort, nodeID, broadcastReceiverChannel, nodeRegistryChannel)

	var (
		//TODO THIS CAN BE A [NUM_ELEVATORS]HRAInput
		elevatorList =   initializeElevatorList()
		//elevatorList  [NUM_ELEVATORS]HRAElevState
		hallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState
		aliveList     [NUM_ELEVATORS]bool
		ackMap        [NUM_ELEVATORS]bool
		online        bool
		init          bool
	)

	for {
		select {
		case reg := <-nodeRegistryChannel:
			// TODO THIS reg/ or it can be a double variable CAN ALSO CONTAIN THE ONLIE STATUS :)
			fmt.Println("HALLO")
			for _, lostNode := range reg.Lost {
				fmt.Printf("Node lost connection: %s\n", lostNode)
				lostNodeInt, _ := strconv.Atoi(lostNode)
				if lostNodeInt == nodeIDInt {
					online = false
				} else {
					fmt.Println("WE SET AN ELEVATOR INACTIVE")
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
			// TODO THIS CAN BE A FUNC

			if !init {
				elevatorList[nodeIDInt] = msg.ElevatorList[nodeIDInt]
				init = true
			}

			
			aliveList[senderId] = msg.OnlineStatus
			elevatorList[senderId] = msg.ElevatorList[senderId]
			hallOrderList[senderId] = msg.HallOrderList[senderId]
			hallOrderList = cyclicCounter(hallOrderList, nodeIDInt)
			//printHallOrderList(hallOrderList)
			//printElevatorList(elevatorList)
			//TODO THIS CAN BE FUNC
			allAcknowledged := true
			for i := 0; i < NUM_ELEVATORS; i++ {
				if nodeIDInt == i {
					continue
				}
				if aliveList[i] && !ackMap[i] {
					allAcknowledged = false
					break
				}
			}
			if allAcknowledged {
				for i := 0; i < NUM_ELEVATORS; i++ {
					if i != nodeIDInt {
						ackMap[i] = false
					}
				}
				messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
					AliveList:     aliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}
			}

		case payload := <-messagefromOrderAssigner:
			hallOrderList[nodeIDInt] = payload.HallRequests
			elevatorList[nodeIDInt] = payload.States[nodeID]
			aliveList[nodeIDInt] = payload.ActiveSatus

		case <-time.After(10 * time.Millisecond):
			broadcastTransmissionChannel <- Message{
				SenderId:      nodeID,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
				OnlineStatus:  aliveList[nodeIDInt],
			}

			if !online {
				newAliveList := [NUM_ELEVATORS]bool{}
				newAliveList[nodeIDInt] = aliveList[nodeIDInt]
				messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
					AliveList:     newAliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}

			}

		}
	}
}
