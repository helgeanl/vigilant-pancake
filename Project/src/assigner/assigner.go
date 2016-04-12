// Package liftAssigner gathers the cost values of the lifts on the network,
// and assigns the best candidate to each order.

// Recive costs from elevators and compare

package assigner

import (
	def "definitions"
	"log"
	"queue"
	"time"
)

type reply struct {
	cost int
	elevator string
}
type request struct {
	floor  int
	button int
	timer  *time.Timer
}

func CollectCosts(message, *numOnline int){
	requestList := make( map[request][]reply)
	var timeout = make(chan *request)

	for{
		select{
		case message := <-costReply:
			newRequest := request{floor: message.Floor,button: message.Button}
			newReply := reply{cost: message.Cost, elevator: message.Addr}

			// Check if request is in queue
			if replyList, exist := requestList[newRequest]; exist {
				// Check if newReply already is registered.
				found := false
				for _, reply := range replyList {
					if reply == newReply {
						found = true
					}
				}
				// Add to list if not found
				if !found {
					requestList[newRequest] = append(requestList[newRequest], newReply)
					newRequest.timer.Reset(def.CostReplyTimeoutDuration)
				}
			} else {
				// If order not in queue at all, init order list with it
				newRequest.timer = time.NewTimer(def.CostReplyTimeoutDuration)
				requestList[newRequest] = []reply{newReply}
				go costTimer(&newRequest, timeout)
			}
			chooseBestLift(unassigned, numOnline, false)
		case <- timeout:
			choose best elevator
		}
	}
}

func chooseBestElevator(requestList map[request][]reply, *numOnline int, timeout bool){
	var bestElevatorAddr string
	// Go through list of requests and find the best elevator in each replyList
	for request,replyList := range requestList{
		if len(replyList) == *numOnline || timeout{
			lowestCost := 10000
			for _,reply := range replyList{
				if reply.cost < lowestCost{
					lowestCost = reply.cost
					bestElevatorAddr = reply.Addr
				}else if reply.cost == lowestCost{
					if reply.Addr < bestElevator{
						bestElevatorAddr = reply.Addr
					}
				}
			}
			queue.AddRequest(request.floor, request.button, bestElevatorAddr)
			request.timer.Stop()
			delete(unassigned, request)
		}
	}
}

func costTimer(newRequest *request, timeout chan<- *request) {
	<-newRequest.timer.C
	timeout <- newRequest
}
