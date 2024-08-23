package network

import (
	. "Project/dataenums"
	"Project/network/broadcast"
	"Project/network/local"
	"Project/network/nodes"
	"fmt"
	"os"
	//"time"
)

const lifelinePort int = 1337
const messagePort int = lifelinePort + 1

func Network(messagefromOrderAssigner <-chan HRAInput,
	messagetoOrderAssignerChannel chan<- HRAInput,
	ipChannel chan<- string,
	nodeID string) {
	nodeIP, err := local.GetIP()
	if err != nil {
		print("Unable to get the IP address")
	}
	ipChannel <- nodeIP // pass the IP address to main process

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
		onlineStatus = false
		onlineStatusTest = false
		Message  Message
	)
	// Periodic broadcast of the last updated message
	// Periodic broadcast of the last updated message
	/*
		go func() {
			for {
				if !isEmptyHRAInput(lastMessage.Payload) { // Check if lastMessage.Payload is not empty
					broadcastTransmissionChannel <- lastMessage
					print("Broadcasting last message to network")
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()
	*/
	for {
		select {
		case reg := <-nodeRegistryChannel:

			// on state change, pass to main process
			if contains(reg.Lost, nodeUid) {
				fmt.Println("Node lost connection:", nodeUid)
				onlineStatus = false
			} else if reg.New == nodeUid {
				fmt.Println("Node connected:", nodeUid)
				onlineStatus = true
			}

			//if offline remove yourself from hra
			//send hra to assigner

		case msg := <-broadcastReceiverChannel:
			//we cant just set equal
			fmt.Println("hallo vi er på nettet")
			messagetoOrderAssignerChannel <- msg.Payload
			//handle incoming msg
			//send msg to assigner

		case payload := <-messagefromOrderAssigner:
			fmt.Println("msg from assigmer")
			Message.SenderId = nodeIP
			Message.Payload = payload
			Message.OnlineStatus = onlineStatus
			fmt.Println("Broadcast transmitted to network")
			if !onlineStatusTest {
				print("sending msg back")
				messagetoOrderAssignerChannel <- Message.Payload
			}
			broadcastTransmissionChannel <- Message
		}
	}
}
