package assigner

import (
	def "definitions"
	"log"
	"queue"
	"time"
)

type reply struct {
	cost     int
	elevator string
}
type request struct {
	floor  int
	button int
	timer  *time.Timer
}

// CollectCosts recive all cost from online elevators
func CollectCosts(costReply <-chan def.Message, numOnlineCh <-chan int) {
	requestMap := make(map[request][]reply)
	var timeout = make(chan *request)
	var numOnline = 1
	for {
		select {
		case message := <-costReply:
			handleCostReply(requestMap, message, numOnline, timeout)
		case numOnlineUpdate := <-numOnlineCh:
			numOnline = numOnlineUpdate
		case <-timeout:
			log.Println(def.ColR, "Not all costs received in time!", def.ColN)
			chooseBestElevator(requestMap, numOnline, true)
		}
	}
}
//handleCostReply stores in a map mapping request to costReply
func handleCostReply(requestMap map[request][]reply, message def.Message, numOnline int, timeout chan *request) {
	newRequest := request{floor: message.Floor, button: message.Button}
	newReply := reply{cost: message.Cost, elevator: message.Addr}
	log.Println(def.ColR, "New Cost incomming from: ", message.Addr, " for cost: ", message.Cost, def.ColN)

	// Compare newRequests with existingRequest without the timer
	for existingRequest := range requestMap {
		if equal(existingRequest, newRequest) {
			newRequest = existingRequest
		}
	}
	if replyList, exist := requestMap[newRequest]; exist {
		// Check if newReply already is registered.
		found := false
		for _, reply := range replyList {
			if reply == newReply {
				found = true
			}
		}
		// Add to list if not found
		if !found {
			requestMap[newRequest] = append(requestMap[newRequest], newReply)
			newRequest.timer.Reset(def.CostReplyTimeoutDuration)
		}
	} else {
		// If newRequest not in requestMap, make new replyList
		newRequest.timer = time.NewTimer(def.CostReplyTimeoutDuration)
		requestMap[newRequest] = []reply{newReply}
		go costTimer(&newRequest, timeout)
	}
	chooseBestElevator(requestMap, numOnline, false)
}

// chooseBestElevator goes through a map of requests and finds the best elevator in each replyList
func chooseBestElevator(requestMap map[request][]reply, numOnline int, isTimeout bool) {
	var bestElevator string

	for request, replyList := range requestMap {
		if len(replyList) == numOnline || isTimeout {
			lowestCost := 10000
			for _, reply := range replyList {
				if reply.cost < lowestCost {
					lowestCost = reply.cost
					bestElevator = reply.elevator
				} else if reply.cost == lowestCost {
					// On equal cost, the elevator with lowest IP get the request
					if reply.elevator < bestElevator {
						bestElevator = reply.elevator
					}
				}
			}
			queue.AddRequest(request.floor, request.button, bestElevator)
			request.timer.Stop()
			delete(requestMap, request)
		}
	}
}

func equal(r1, r2 request) bool {
	return r1.floor == r2.floor && r1.button == r2.button
}

func costTimer(newRequest *request, timeout chan<- *request) {
	<-newRequest.timer.C
	timeout <- newRequest
}
