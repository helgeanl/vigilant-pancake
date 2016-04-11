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

			// if not, append request to requestList and start timer
			// choose best elevator

		case <- timeout:
			choose best elevator
		}
	}
}

func chooseBestElevator(requestList map[request][]reply, *numOnline int, timeout bool){
	var lowestCost int
	var bestElevator string
	// Go through list of requests and find the best elevator in each replyList
	for request,replyList := range requestList{
		if len(replyList) == *numOnline || timeout{
			lowestCost = 10000
			for _,reply := range replyList{
				if reply.cost < lowestCost{
					lowestCost = reply.cost
					bestElevator = reply.Addr
				}else if reply.cost == lowestCost{
					if reply.Addr < bestElevator{
						bestElevator = reply.Addr
					}
				}
			}
			// add request to queue
			// stop timer
			// delete request from requestList
		}
	}
}











// Run collects cost values from all lifts for each new order, and attempts
// to find the best lift for each order, when either all online lifts have
// replied or after a timeout.
func CollectCosts(costReply <-chan def.Message, numOnline *int) {
	// Gathered cost data for each order is stored here until a lift is
	// assigned to the order.
	unassigned := make(map[request][]reply)

	var timeout = make(chan *request)
	const timeoutDuration = 10 * time.Second

	for {
		select {
		case message := <-costReply:
			newRequest := request{floor: message.Floor, button: message.Button}
			newReply := reply{cost: message.Cost, lift: message.Addr}

			for oldRequest := range unassigned {
				if equal(oldRequest, newRequest) {
					newRequest = oldRequest
				}
			}

			// Check if order in queue.
			if replyList, exist := unassigned[newRequest]; exist {
				// Check if newReply already is registered.
				found := false
				for _, reply := range replyList {
					if reply == newReply {
						found = true
					}
				}
				// Add it if it wasn't.
				if !found {
					unassigned[newRequest] = append(unassigned[newRequest], newReply)
					newRequest.timer.Reset(timeoutDuration)
				}
			} else {
				// If order not in queue at all, init order list with it
				newRequest.timer = time.NewTimer(timeoutDuration)
				unassigned[newRequest] = []reply{newReply}
				go costTimer(&newRequest, timeout)
			}
			chooseBestLift(unassigned, numOnline, false)

		case <-timeout:
			log.Println(def.ColR, "Not all costs received in time!", def.ColN)
			chooseBestLift(unassigned, numOnline, true)
		}
	}
}

// chooseBestElevator checks if any of the requests waiting for a elevator assignment
// have collected enough information to have a elevator assigned. For all requests
// that have, it selects a elevator, and adds it to the queue.
// It assumes that all elevators always make the same decision, but if they do not,
// a timer for each request assured that this never gives unhandled requests.
func chooseBestElevator(unassigned map[request][]reply, numOnline *int, requestTimedOut bool) {
	const maxInt = int(^uint(0) >> 1)
	// Loop through all lists.
	for request, replyList := range unassigned {
		// Check if the list is complete or the timer has timed out.
		if len(replyList) == *numOnline || requestTimedOut {
			lowestCost := maxInt
			var bestElevator string

			// Loop through costs in each complete list.
			for _, reply := range replyList {
				if reply.cost < lowestCost {
					lowestCost = reply.cost
					bestElevator = reply.elevator
				} else if reply.cost == lowestCost {
					// Prioritise on lowest IP value if cost is the same.
					if reply.elevator < bestElevator {
						lowestCost = reply.cost
						bestElevator = reply.elevator
					}
				}
			}
			queue.AddRemoteRequest(request.floor, request.button, bestElevator)
			request.timer.Stop()
			delete(unassigned, request)
		}
	}
}

func costTimer(newRequest *request, timeout chan<- *request) {
	<-newRequest.timer.C
	timeout <- newRequest
}

func equal(r1, r2 request) bool {
	return r1.floor == r2.floor && r1.button == r2.button
}
