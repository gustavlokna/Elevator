package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//TODO MOVE DATA ENUMS ? 
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
		aliveList [NUM_ELEVATORS]bool
		ackMap    [NUM_ELEVATORS]bool
		//TODO THIS CAN BE A [NUM_ELEVATORS]HRAInput
		// TODO INITILIZE THE LIST, else a crash may occur
		elevatorList  [NUM_ELEVATORS]HRAElevState
		hallOrderList [NUM_ELEVATORS][NFloors][NButtons]ButtonState

		online bool
		init   bool
	)

	// Periodic broadcast of the last updated message

	// TODO: This is copied ?
	//TODO MOVE THIS TO SENDER. DO NOT BRADCAST UNTIL FIRST MESSAGE :) 
	//I:E SAVE THIS IN SENDER AND LOOP 
	// WAIT TILL MSG :) 
	go func() {
		for {
			// NOTE REKEFØLGEN PÅ TIMEREN ER VIKTIG I HENHOLD TIL REINIT
			// TODO MAY BE A PROBLEM WHEN PACETLOSS, BUT UNSURE
			// IF SO, ADD LOGIC THAT FIRST SLEEP IS LONGER THAN ALL OTHETS
			time.Sleep(10 * time.Millisecond)
			broadcastTransmissionChannel <- Message{
				SenderId:      nodeID,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
				OnlineStatus:  aliveList[nodeIDInt],
			}

		}
	}()

	for {
		select {
		case reg := <-nodeRegistryChannel:
			// TODO THIS reg/ or it can be a double variable CAN ALSO CONTAIN THE ONLIE STATUS :)
			fmt.Println("HALLO")
			for _, lostNode := range reg.Lost {
				// TODO REMOVE
				fmt.Printf("Node lost connection: %s\n", lostNode)
				lostNodeInt, _ := strconv.Atoi(lostNode)

				// TODO MOVE THIS TO THE DEFULT CASE
				if lostNodeInt == nodeIDInt {
					online = false
					newAliveList := [NUM_ELEVATORS]bool{}
					newAliveList[nodeIDInt] = aliveList[nodeIDInt]
					messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
						AliveList:     newAliveList,
						ElevatorList:  elevatorList,
						HallOrderList: hallOrderList,
					}
				} else {
					aliveList[lostNodeInt] = false
					hallOrderList[lostNodeInt] = resetHallCalls()
				}

			}
			// TODO HERE WE REALY JUST CARE ABOUT US
			// AS WE WANT TO RESET OUR CALLS ON OUR REINT
			for _, activeNode := range reg.Nodes {
				fmt.Printf("Node active: %s\n", activeNode)
				activeNodeInt, _ := strconv.Atoi(activeNode)
				if activeNodeInt == nodeIDInt {
					hallOrderList[activeNodeInt] = resetHallCalls()
					online = true
				}
				//hallOrderList[activeNodeInt] = resetHallCalls()
				aliveList[activeNodeInt] = true
			}

		case msg := <-broadcastReceiverChannel:
			
			senderId, _ := strconv.Atoi(msg.SenderId)
			aliveList[senderId] = msg.OnlineStatus

			ackMap[senderId] = reflect.DeepEqual(elevatorList, msg.ElevatorList)

			// Only update state if alive
			// TODO THIS CAN BE A FUNC
			elevatorList[senderId] = msg.ElevatorList[senderId]
			if !init {
				elevatorList[nodeIDInt] = msg.ElevatorList[nodeIDInt]
				init = true
			}
			hallOrderList[senderId] = msg.HallOrderList[senderId]
			hallOrderList = cyclicCounter(hallOrderList, aliveList, nodeIDInt)

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
			//s TODO MOVE THIS TO THE DEFULT CASE
			if !online {
				newAliveList := [NUM_ELEVATORS]bool{}
				newAliveList[nodeIDInt] = aliveList[nodeIDInt]
				messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
					AliveList:     newAliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}
			}

			broadcastTransmissionChannel <- Message{
				SenderId:      nodeID,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
				OnlineStatus:  aliveList[nodeIDInt],
			}
		}
	}

	// TODO ADD DEFULT CASE THAT IF WE ARE OFLINE SEND TO ASS AND SLEEP :)
}
