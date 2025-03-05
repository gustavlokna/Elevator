package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	//"github.com/google/go-cmp/cmp"
	"github.com/cespare/xxhash/v2"
)

func hashStruct(v interface{}) uint64 {
	b, _ := json.Marshal(v) // Convert struct to JSON
	return xxhash.Sum64(b)  // Compute fast hash
}

var mu sync.Mutex

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
		elevatorList = initializeElevatorList()
		//TODO THIS CAN BE A [NUM_ELEVATORS]HRAInput
		hallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState
		aliveList     [NUM_ELEVATORS]bool
		ackMap        [NUM_ELEVATORS]bool
		online        bool
		init          bool
		proession     bool
		
	)
	oldCabRequests := make([]bool, len(elevatorList[nodeIDInt].CabRequests))
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

			if !reflect.DeepEqual(hallOrderList, msg.HallOrderList) || !reflect.DeepEqual(aliveList, msg.AliveList) {
				proession = true
			}
			
			if !reflect.DeepEqual(elevatorList[senderId].CabRequests, msg.ElevatorList[senderId].CabRequests) {
				proession = true
			}

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
			for i := 0; i < NUM_ELEVATORS; i++ {
				if nodeIDInt == i {
					continue
				}
				if aliveList[i] && !ackMap[i] {
					allAcknowledged = false
					break
				}
			}
			if allAcknowledged && proession {
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
				proession = false
			}

		case payload := <-messagefromOrderAssigner:
			hallOrderList[nodeIDInt] = payload.HallRequests
			aliveList[nodeIDInt] = payload.ActiveSatus
			elevatorList[nodeIDInt] = payload.States[nodeID]

		case <-time.After(10 * time.Millisecond):		
			
			if !reflect.DeepEqual(oldCabRequests, elevatorList[nodeIDInt].CabRequests){
				copy(oldCabRequests, elevatorList[nodeIDInt].CabRequests)
				proession = true
				fmt.Println("HELLO")
			}

			broadcastTransmissionChannel <- Message{
				SenderId:      nodeID,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
				OnlineStatus:  aliveList[nodeIDInt],
				AliveList:     aliveList,
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
