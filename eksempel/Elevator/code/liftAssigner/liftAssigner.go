// Package liftAssigner gathers the cost values of the lifts on the network,
// and assigns the best candidate to each order.
package liftAssigner

import (
	def "config"
	"log"
	"queue"
	"time"
)

type reply struct {
	cost int
	lift string
}
type order struct {
	floor  int
	button int
	timer  *time.Timer
}

// Run collects cost values from all lifts for each new order, and attempts
// to find the best lift for each order, when either all online lifts have
// replied or after a timeout.
func Run(costReply <-chan def.Message, numOnline *int) {
	// Gathered cost data for each order is stored here until a lift is
	// assigned to the order.
	unassigned := make(map[order][]reply)

	var timeout = make(chan *order)
	const timeoutDuration = 10 * time.Second

	for {
		select {
		case message := <-costReply:
			newOrder := order{floor: message.Floor, button: message.Button}
			newReply := reply{cost: message.Cost, lift: message.Addr}

			for oldOrder := range unassigned {
				if equal(oldOrder, newOrder) {
					newOrder = oldOrder
				}
			}

			// Check if order in queue.
			if replyList, exist := unassigned[newOrder]; exist {
				// Check if newReply already is registered.
				found := false
				for _, reply := range replyList {
					if reply == newReply {
						found = true
					}
				}
				// Add it if it wasn't.
				if !found {
					unassigned[newOrder] = append(unassigned[newOrder], newReply)
					newOrder.timer.Reset(timeoutDuration)
				}
			} else {
				// If order not in queue at all, init order list with it
				newOrder.timer = time.NewTimer(timeoutDuration)
				unassigned[newOrder] = []reply{newReply}
				go costTimer(&newOrder, timeout)
			}
			chooseBestLift(unassigned, numOnline, false)

		case <-timeout:
			log.Println(def.ColR, "Not all costs received in time!", def.ColN)
			chooseBestLift(unassigned, numOnline, true)
		}
	}
}

// chooseBestLift checks if any of the orders waiting for a lift assignment
// have collected enough information to have a lift assigned. For all orders
// that have, it selects a lift, and adds it to the queue.
// It assumes that all lifts always make the same decision, but if they do not,
// a timer for each order assured that this never gives unhandled orders.
func chooseBestLift(unassigned map[order][]reply, numOnline *int, orderTimedOut bool) {
	const maxInt = int(^uint(0) >> 1)
	// Loop through all lists.
	for order, replyList := range unassigned {
		// Check if the list is complete or the timer has timed out.
		if len(replyList) == *numOnline || orderTimedOut {
			lowestCost := maxInt
			var bestLift string

			// Loop through costs in each complete list.
			for _, reply := range replyList {
				if reply.cost < lowestCost {
					lowestCost = reply.cost
					bestLift = reply.lift
				} else if reply.cost == lowestCost {
					// Prioritise on lowest IP value if cost is the same.
					if reply.lift < bestLift {
						lowestCost = reply.cost
						bestLift = reply.lift
					}
				}
			}
			queue.AddRemoteOrder(order.floor, order.button, bestLift)
			order.timer.Stop()
			delete(unassigned, order)
		}
	}
}

func costTimer(newOrder *order, timeout chan<- *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}

func equal(o1, o2 order) bool {
	return o1.floor == o2.floor && o1.button == o2.button
}
