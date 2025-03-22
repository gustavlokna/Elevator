package assigner

import (
	. "Project/dataenums"
	"Project/hwelevio"
	"fmt"
	"reflect"
	"strconv"
)

func Assigner(
	newOrderChannel chan<- [NFloors][NButtons]bool,
	payloadFromElevator <-chan FromDriverToAssigner,
	toNetworkChannel chan<- FromAssignerToNetwork,
	fromNetworkChannel <-chan FromNetworkToAssigner,
	fromAsstoLight chan<- [NFloors][NButtons]ButtonState,
	nodeID string,
) {
	var (
		PayloadFromassignerToNetwork = initPayloadToNetwork()
		prevAssignedOrders           [NFloors][NButtons]bool
		drv_buttons                  = make(chan ButtonEvent)
		prevMsg                      FromNetworkToAssigner
		newOrder                     bool
	)
	// Convert nodeID to int
	myID, err := strconv.Atoi(nodeID)
	if err != nil {
		fmt.Printf("Invalid nodeID: %v\n", err)
		return
	}
	payload := <-payloadFromElevator
	PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
		PayloadFromassignerToNetwork, nodeID)

	toNetworkChannel <- PayloadFromassignerToNetwork

	go hwelevio.PollButtons(drv_buttons)
	for {
		select {
		case btnEvent := <-drv_buttons:
			PayloadFromassignerToNetwork = buttonPressed(PayloadFromassignerToNetwork,
				nodeID, btnEvent)
			toNetworkChannel <- PayloadFromassignerToNetwork

		case payload := <-payloadFromElevator:
			PayloadFromassignerToNetwork = handlePayloadFromElevator(payload,
				PayloadFromassignerToNetwork, nodeID)

			toNetworkChannel <- PayloadFromassignerToNetwork

		case msg := <-fromNetworkChannel:

			if !reflect.DeepEqual(prevMsg.HallOrderList, msg.HallOrderList) || !reflect.DeepEqual(prevMsg, msg.AliveList) {
				newOrder = true
			}

			for i := 0; i < NElevators; i++ {
				if !reflect.DeepEqual(prevMsg.ElevatorList[i].CabRequests, msg.ElevatorList[i].CabRequests) {
					newOrder = true
					break
				}
			}			
			
			PayloadFromassignerToNetwork = handlePayloadFromNetwork(PayloadFromassignerToNetwork,
				msg, myID)
			if newOrder{
				newOrders := assignOrders(msg, myID)
				if newOrders != prevAssignedOrders {
					newOrderChannel <- newOrders
					prevAssignedOrders = newOrders
				}
			}
			fromAsstoLight <- updateLightStates(msg, myID)

		}
	}

}
