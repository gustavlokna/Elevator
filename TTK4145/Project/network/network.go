package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"Project/network/local"
	"Project/network/nodes"
	"fmt"
	"os"
	"time"
	"strconv"
)

const lifelinePort int = 1337
const messagePort int = lifelinePort + 1
const NUM_ELEVATORS int = 3


func Network(messagefromOrderAssigner <-chan HRAInput,
	messagetoOrderAssignerChannel chan<- Message,
	nodeID string) {
	nodeIP, err := local.GetIP()
	if err != nil {
		print("Unable to get the IP address")
	}
	//TODO CONVERT THIS SMARTER DO NOT USE THIS 
	nodeIDINT, err := strconv.Atoi(nodeID)

	nodeUid := fmt.Sprintf("peer-%s-%d", nodeIP, os.Getpid())

	// setup lifeline for network node registry
	nodeRegistryChannel := make(chan nodes.NetworkNodeRegistry)
	TransmissionEnableChannel := make(chan bool)
	go nodes.Sender(lifelinePort, nodeUid, TransmissionEnableChannel)
	go nodes.Receiver(lifelinePort, nodeRegistryChannel)

	// setup broadcast for message transmission
	broadcastTransmissionChannel := make(chan Message)
	broadcastReceiverChannel := make(chan Message)
	go broadcast.Sender(messagePort, broadcastTransmissionChannel)
	go broadcast.Receiver(nodeIP, messagePort, broadcastReceiverChannel)
	
	var (
		onlineStatus      = false
		messageInstance   Message
		lastMessage       Message
		
		aliveList         [NUM_ELEVATORS]bool
		/*
		elevatorList      [NUM_ELEVATORS]Elevator
		cabOrderList      [NUM_ELEVATORS][NFloors]Elevator
		hallOrderList     [NUM_ELEVATORS][NFloors][NButtons]ButtonState
		*/
	)
	
	// Periodic broadcast of the last updated message
	// Periodic broadcast of the last updated message

	// TODO: This is copied ?
	
	go func() {
		for {
			if !isEmptyHRAInput(lastMessage.Payload) { // Check if lastMessage.Payload is not empty
				broadcastTransmissionChannel <- lastMessage
				//print("Broadcasting last message to network")
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	
	for {
		select {
		case reg := <-nodeRegistryChannel:

			// on state change, pass to main process
			if contains(reg.Lost, nodeUid) {
				fmt.Println("Node lost connection:", nodeUid)
				onlineStatus = false
				
				aliveList[nodeIDINT] = false
				//if i lose connection update aliveList

			} else if reg.New == nodeUid {
				fmt.Println("Node connected:", nodeUid)
				onlineStatus = true
				aliveList[nodeIDINT] = true
			}
			//if offline send to orderassigner! 
			// send btn to ass? 


		case msg := <-broadcastReceiverChannel:
			//we cant just set equal
			fmt.Println("hallo vi er på nettet")
			fmt.Println("msg id: ", msg.SenderId)

			messagetoOrderAssignerChannel <- msg

			//handle incoming msg
			//Cyclic counter logic updates local world view
			//send msg to assigner with function 

		case payload := <-messagefromOrderAssigner:
			fmt.Println("msg from assigmer")
			messageInstance.SenderId = nodeID
			messageInstance.Payload = payload
			messageInstance.OnlineStatus = onlineStatus
			lastMessage = messageInstance
			//fmt.Println("Broadcast transmitted to network")
			if !messageInstance.OnlineStatus {
				print("sending msg back")
				messagetoOrderAssignerChannel <- messageInstance
			}
			broadcastTransmissionChannel <- messageInstance
		}
	}
}
