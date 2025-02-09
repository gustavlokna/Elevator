package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"Project/network/local"
	"Project/network/nodes"
	"fmt"
	//"os"
	"time"
	"strconv"
	"reflect"
)

const lifelinePort int = 1337
const messagePort int = lifelinePort + 1



func Network(messagefromOrderAssigner <-chan PayloadFromassignerToNetwork,
	messagetoOrderAssignerChannel chan<- PayloadFromNetworkToAssigner,
	nodeID string) {

	// **WAIT FOR INITIALIZATION** BEFORE STARTING MAIN LOOP
	fmt.Println("Waiting for network initialization...")
	nodeIP, err := local.GetIP()
	nodeIDInt,_ := strconv.Atoi(nodeID)
	if err != nil {
		print("Unable to get the IP address")
	}
	//TODO: MAKE THIS BETTER 
	nodeIPint, err := strconv.Atoi(nodeIP)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Converted number:", nodeIPint)
	}

	//nodeUid := fmt.Sprintf("peer-%s-%d", nodeIP, os.Getpid())

	// setup lifeline for network node registry
	nodeRegistryChannel := make(chan nodes.NetworkNodeRegistry)
	TransmissionEnableChannel := make(chan bool)
	go nodes.Sender(lifelinePort, nodeID, TransmissionEnableChannel)

	go nodes.Receiver(lifelinePort, nodeRegistryChannel)


	// setup broadcast for message transmission
	broadcastTransmissionChannel := make(chan Message)
	broadcastReceiverChannel := make(chan Message)
	go broadcast.Sender(messagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(nodeIP, messagePort, broadcastReceiverChannel)
	
	var (
		//onlineStatus      = false
		//messageInstance   Message
		//TODO: DO WE NEED THIS? 
		//lastMessage       Message
		
		aliveList         [NUM_ELEVATORS]bool
		ackMap 			  [NUM_ELEVATORS]bool
		//TODO THIS CAN BE A [NUM_ELEVATORS]HRAInput
		elevatorList      [NUM_ELEVATORS]HRAElevState
		hallOrderList     [NUM_ELEVATORS][NFloors][NButtons]ButtonState
	)
	
	// Periodic broadcast of the last updated message

	// TODO: This is copied ?
	go func() {
		for {
			broadcastTransmissionChannel <- Message{
				SenderId: nodeID,
				ElevatorList: elevatorList, 
				HallOrderList: hallOrderList, 
				OnlineStatus: aliveList[nodeIDInt], 
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	
	for {
		select {
		case reg := <-nodeRegistryChannel:
			for _, lostNode := range reg.Lost {
				// TODO REMOVE 
				fmt.Printf("Node lost connection: %s\n", lostNode)
				lostNodeInt,_ := strconv.Atoi(lostNode)
				aliveList[lostNodeInt] = false
				
				//TODO if only one node is alive
				// assigning will not work, but this is outside specs ? 
				hallOrderList[lostNodeInt] = resetHallCalls()
			}
			for _, activeNode := range reg.Nodes {
				fmt.Printf("Node active: %s\n", activeNode)
				activeNodeInt,_ := strconv.Atoi(activeNode)
				hallOrderList[activeNodeInt] = resetHallCalls()
				aliveList[activeNodeInt] = true	
			}
			//if offline send to orderassigner! 
			// send btn to ass? 


		case msg := <-broadcastReceiverChannel:
			/*
			fmt.Println("My elevatorList: ", elevatorList)
			fmt.Println("Incoming elevatorList: ", msg.ElevatorList)

			if !reflect.DeepEqual(elevatorList, msg.ElevatorList) {
				fmt.Println("Mismatch found!")
			}
			*/
			senderId, _ := strconv.Atoi(msg.SenderId)
			aliveList[senderId] = msg.OnlineStatus 
		
			// Directly check if the incoming elevator list matches mine
			ackMap[senderId] = reflect.DeepEqual(elevatorList, msg.ElevatorList)
		
			// Only update state if alive
			elevatorList[senderId] = msg.ElevatorList[senderId]
			hallOrderList[senderId] = msg.HallOrderList[senderId]
		
			// Run cyclic logic to update local hallOrderList
			hallOrderList = cyclicCounter(hallOrderList, nodeIDInt)
		
			// Check if all active elevators acknowledge the same states
			
			allAcknowledged := true
			for i := 0; i < NUM_ELEVATORS; i++ {
				if nodeIDInt==i{
					continue 
				}
				print(ackMap[i] )
				if aliveList[i] && !ackMap[i] {
					allAcknowledged = false
					fmt.Println("NOOO")
					break
				}
			}
			
			// Only send message to assigner if all acknowledgments are true
			
			if allAcknowledged {
				messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
					AliveList:     aliveList,
					ElevatorList:  elevatorList,
					HallOrderList: hallOrderList,
				}
			}
			/*
			messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
				AliveList:     aliveList,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
			}
			*/ 
		case payload := <-messagefromOrderAssigner:
			/*
			messageInstance.SenderId = nodeID
			messageInstance.HallOrderList[nodeIDInt] = payload.HallRequests
			//TODO BURDE VÆRE SAMME 
			messageInstance.ElevatorList[nodeIDInt] = payload.States[nodeID]
			messageInstance.OnlineStatus = payload.ActiveSatus
			if !messageInstance.OnlineStatus {
				// TODO set btn_pressed = assign and send to assigner 
				print("sending msg back")
				//messagetoOrderAssignerChannel <- messageInstance
			}
			*/

			hallOrderList[nodeIDInt]= payload.HallRequests
			elevatorList[nodeIDInt] =  payload.States[nodeID]
			aliveList[nodeIDInt] = payload.ActiveSatus


			broadcastTransmissionChannel <- Message{
				SenderId: nodeID,
				ElevatorList: elevatorList, 
				HallOrderList: hallOrderList, 
				OnlineStatus: aliveList[nodeIDInt], 
			}
		}
	}
}
