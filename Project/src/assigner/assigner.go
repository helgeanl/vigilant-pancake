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

var NumOnlineCh = make(chan int)

func CollectCosts(costReply chan def.Message, numOnlineCh chan int){
	requestMap := make( map[request][]reply)
	var timeout = make(chan *request)
	var numOnline int = 1
	for{
		select{
		case message := <-costReply:
			newRequest := request{floor: message.Floor,button: message.Button}
			newReply := reply{cost: message.Cost, elevator: message.Addr}
			log.Println(def.ColR, "NewRequest på", newRequest,def.ColN)
			// Check if request is in queue
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
				// If order not in queue at all, init order list with it
				newRequest.timer = time.NewTimer(def.CostReplyTimeoutDuration)
				requestMap[newRequest] = []reply{newReply}
				go costTimer(&newRequest, timeout)
			}
			chooseBestElevator(requestMap,numOnline,false)
		case numOnlineUpdate := <- numOnlineCh:
			numOnline = numOnlineUpdate
		case <- timeout:
			log.Println(def.ColR,"Not all costs received in time!",def.ColN)
			chooseBestElevator(requestMap,numOnline,true)
		}
	}
}

func chooseBestElevator(requestMap map[request][]reply, numOnline int, timeout bool){
	var bestElevator string
	// Go through list of requests and find the best elevator in each replyList
	for request,replyList := range requestMap{
		log.Println(def.ColB,"Number online: ", numOnline)
		if len(replyList) == numOnline || timeout{
			log.Println(def.ColB,"All costs collected, timeout: ", timeout)
			lowestCost := 10000
			for _,reply := range replyList{
				if reply.cost < lowestCost{
					lowestCost = reply.cost
					bestElevator = reply.elevator
				}else if reply.cost == lowestCost{
					if reply.elevator < bestElevator{
						bestElevator = reply.elevator
					}
				}
			}
			log.Println(def.ColB,"Will now add order to Floor:",request.floor," Button",request.button," To Elevator:",bestElevator, def.ColN )
			queue.AddRequest(request.floor, request.button, bestElevator)
			request.timer.Stop()
			delete(requestMap, request)
		}
	}
}

func costTimer(newRequest *request, timeout chan<- *request) {
	<-newRequest.timer.C
	timeout <- newRequest
}
