package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"Project/network/local"
	"Project/network/nodes"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const lifelinePort int = 1337
const messagePort int = lifelinePort + 1

func Network(messagefromOrderAssigner <-chan PayloadFromassignerToNetwork,
	messagetoOrderAssignerChannel chan<- PayloadFromNetworkToAssigner,
	nodeID string) {

	// **WAIT FOR INITIALIZATION** BEFORE STARTING MAIN LOOP
	fmt.Println("Waiting for network initialization...")
	// TODO NOT NECESSARY TO GET IP HERE, MOVE TO FUNC AND SEE IF NECESSARY  
	nodeIP, err := local.GetIP()
	if err != nil {
		print("Unable to get the IP address")
	}
	// TODO MAKE CODE COMPATIBLE WITHOUT THIS "STR TO INT CONV"
	nodeIDInt, _ := strconv.Atoi(nodeID)

	//TODO: MAKE THIS BETTER
	// TODO REMOVE ? 
	nodeIPint, err := strconv.Atoi(nodeIP)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Converted number:", nodeIPint)
	}

	nodeRegistryChannel := make(chan nodes.NetworkNodeRegistry)
	TransmissionEnableChannel := make(chan bool)
	// TODO COMBINE THE "NODES" AND "BRADCASTS" "SENDER" and "RECIVER", into one "SENDER" and one "RECIVER"
	go nodes.Sender(lifelinePort, nodeID, TransmissionEnableChannel)
	go nodes.Receiver(lifelinePort, nodeRegistryChannel)

	broadcastTransmissionChannel := make(chan Message)
	broadcastReceiverChannel := make(chan Message)
	go broadcast.Sender(messagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(nodeIP, messagePort, broadcastReceiverChannel)

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
			//if offline send to orderassigner!
			// send btn to ass?

		case msg := <-broadcastReceiverChannel:
			senderId, _ := strconv.Atoi(msg.SenderId)
			aliveList[senderId] = msg.OnlineStatus

			ackMap[senderId] = reflect.DeepEqual(elevatorList, msg.ElevatorList)

			// Only update state if alive
			// TODO THIS CAN BE A FUNC
			elevatorList[senderId] = msg.ElevatorList[senderId]
			if !init {
				fmt.Println(senderId)
				elevatorList[nodeIDInt] = msg.ElevatorList[nodeIDInt]
				printElevatorList(elevatorList)
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
			printElevatorList(elevatorList)
			hallOrderList[nodeIDInt] = payload.HallRequests
			elevatorList[nodeIDInt] = payload.States[nodeID]
			aliveList[nodeIDInt] = payload.ActiveSatus
			//s TODO MOVE THIS TO THE DEFULT CASE
			if !online {
				printElevatorList(elevatorList)
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
