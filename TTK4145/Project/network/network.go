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
)

const lifelinePort int = 1337
const messagePort int = lifelinePort + 1



func Network(messagefromOrderAssigner <-chan PayloadFromassignerToNetwork,
	messagetoOrderAssignerChannel chan<- PayloadFromNetworkToAssigner,
	nodeID string) {
	nodeIP, err := local.GetIP()
	nodeIDInt,_ := strconv.Atoi(nodeID)
	if err != nil {
		print("Unable to get the IP address")
	}
	//TODO: MAKE THIS BETTER 
	/*
	nodeIPint, err := strconv.Atoi(nodeIP)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Converted number:", nodeIPint)
	}
	*/
	fmt.Printf("Node initialized with ID: %s\n", nodeID)

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
		onlineStatus      = false
		messageInstance   Message
		//TODO: DO WE NEED THIS? 
		lastMessage       Message
		
		aliveList         [NUM_ELEVATORS]bool
		//TODO THIS CAN BE A [NUM_ELEVATORS]HRAInput
		elevatorList      [NUM_ELEVATORS]HRAElevState
		hallOrderList     [NUM_ELEVATORS][NFloors][NButtons]ButtonState
	)
	
	// Periodic broadcast of the last updated message

	// TODO: This is copied ?
	
	go func() {
		for {
			broadcastTransmissionChannel <- lastMessage
			time.Sleep(10 * time.Millisecond)
		}
	}()
	
	for {
		select {
		case reg := <-nodeRegistryChannel:
			for _, lostNode := range reg.Lost {

				fmt.Printf("Node lost connection: %s\n", lostNode)
				lostNodeInt,_ := strconv.Atoi(lostNode)
				//TODO: let this be overwritte by incommin msg but since broadcast 
				aliveList[lostNodeInt] = false 	
				//TODO if only one node is alive
				// assigning will not work, but this is outside specs
				
				// Handle lost nodes (e.g., update aliveList or notify assigner)
			}
		
			for _, activeNode := range reg.Nodes {
				fmt.Printf("Node active: %s\n", activeNode)
				activeNodeInt,_ := strconv.Atoi(activeNode)
				//TODO: let this be overwritte by incommin msg but since broadcast 
				//freqency is so high i do not belive this will be a problem 
				// this however can be problem if we set our elevator to online. 
				// but we are obstructed. This Will need some better logic. 
				aliveList[activeNodeInt] = true	
				// set all states of node to garbage 
				// Handle active nodes as needed
			}
		
			//if offline send to orderassigner! 
			// send btn to ass? 


		case msg := <-broadcastReceiverChannel:
			//we cant just set equal
			// fmt.Println("hallo vi er på nettet")
			// fmt.Println("msg id: ", msg.SenderId)
			// Convert SenderId (string) to an integer
			senderId, _ := strconv.Atoi(msg.SenderId)
			
			aliveList[senderId] = msg.OnlineStatus
			elevatorList[senderId]= msg.ElevatorList[senderId]
			hallOrderList[senderId]= msg.HallOrderList[senderId]
			//printHallOrderList(hallOrderList)
			//Cyclic counter logic updates local world view
			hallOrderList = cyclicCounter(hallOrderList,aliveList,nodeIDInt)
			//printHallOrderList(hallOrderList)

			lastMessage.HallOrderList = hallOrderList
			

			messagetoOrderAssignerChannel <- PayloadFromNetworkToAssigner{
				AliveList:     aliveList,
				ElevatorList:  elevatorList,
				HallOrderList: hallOrderList,
			}
			
			
			
			
			//send msg to assigner with function 

		case payload := <-messagefromOrderAssigner:
			messageInstance.SenderId = nodeID
			messageInstance.HallOrderList[nodeIDInt] = payload.HallRequests
			//TODO BURDE VÆRE SAMME 
			messageInstance.ElevatorList[nodeIDInt] = payload.States[nodeID]
			
			// TODO should contain info also abot motorstop and obst
			messageInstance.OnlineStatus = onlineStatus
			lastMessage = messageInstance
			hallOrderList[nodeIDInt]= payload.HallRequests
			//printHallOrderList(hallOrderList)
			//printElevatorList(messageInstance.ElevatorList)

			// TODO this should be simpler (just add everything to elevatorList
			//  hallOrderList etc and brodcast thoe variables )
			elevatorList[nodeIDInt] =  payload.States[nodeID]
			//fmt.Println("Broadcast transmitted to network")
			if !messageInstance.OnlineStatus {
				// TODO set btn_pressed = assign and send to assigner 
				print("sending msg back")
				//messagetoOrderAssignerChannel <- messageInstance
			}
			broadcastTransmissionChannel <- messageInstance
		}
	}
}
