// This is a complete rewrite of the queue package
package queue

import (
	"../defs"
	"fmt"
	"log"
)

var _ = fmt.Printf
var _ = log.Printf

type orderStatus struct {
	active bool
	addr   string
}

var inactive = orderStatus{false, ""}

type queue struct {
	q [nF][nB]orderStatus
}

var local queue
var shared queue

// AddInternalOrder adds an order to the local queue.
func AddInternalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, ""})
}

// RemoveInternalOrder removes an order from the local queue.
func RemoveInternalOrder(floor int, button int) {
	local.setOrder(floor, button, blankOrder)
}

// AddSharedOrder adds the given order to the shared queue.
func AddSharedOrder(floor, button int, addr string) {
	shared.q[floor][button] = orderStatus{true, addr}
	updateLocalQueue() // bad abstraction
}

func RemoveSharedOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		shared.setOrder(floor, b, blankOrder)
	}
	updateLocalQueue() // bad abstraction
}

// ChooseDirection returns the direction the lift should continue after the
// current floor.
func ChooseDirection(currFloor, currDir int) int {
	return local.chooseDirection(currFloor, currDir)
}

// ShouldStop returns whether the lift should stop at the given floor, if
// going in the given direction.
func ShouldStop(floor, dir int) {
	return local.shouldStop(floor, dir)
}

// RemoveOrdersAt removes all orders at the given floor in local and shared queue.
func RemoveOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		local.setOrder(floor, b, blankOrder)
		shared.setOrder(floor, b, blankOrder)
	}
	SendOrderCompleteMessage(floor) // bad abstraction
}

// IsOrder returns whether there in an order with the given floor and button
// in the local queue.
func IsOrder(floor, button int) bool { // Rename to IsLocalOrder
	return local.isActiveOrder(floor, button)
}

// IsSharedOrder returns true if there is a order with the given floor and
// button in the shared queue.
func IsSharedOrder(floor, button int) bool {
	return shared.isActiveOrder(floor, button)
}

// ReassignOrders finds all orders assigned to the given dead lift, removes
// them from the shared queue, and sends them on the network as new, un-
// assigned orders.
func ReassignOrders(deadAddr string) {
	// loop thru shared queue
	// remove all orders assigned to the dead lift
	// send neworder-message for each removed order
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if shared.q[f][b].addr == deadAddr {
				shared.setOrder(f, b, blankOrder)
				reassignMessage := &defs.Message{
					Kind:   defs.NewOrder,
					Floor:  f,
					Button: b}
				defs.MessageChan <- *reassignMessage
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given floor.
func SendOrderCompleteMessage(floor int) {
	message := &defs.Message{Kind: defs.CompleteOrder, Floor: floor}
	defs.MessageChan <- *message
}

func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	return local.deepCopy().calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir)
}

/*
 * Methods on queue struct:
 */

func (q *queue) isEmpty() bool {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.q[f][b].active {
				return false
			}
		}
	}
	return true
}

func (q *queue) setOrder(floor, button int, status orderStatus) {
	q.q[floor][button] = status
}

func (q *queue) isActiveOrder(floor, button int) {
	return q.q[floor][button].active
}

func (q *queue) chooseDirection(floor, dir int) int {
	if q.isEmpty() {
		log.Println("ChooseDirection(): empty queue, returning stop")
		return defs.DirnStop
	}
	switch dir {
	case defs.DirnDown:
		if q.isOrdersBelow(floor) && floor > 0 {
			return defs.DirnDown
		} else {
			return defs.DirnUp
		}
	case defs.DirnUp:
		if q.isOrdersAbove(floor) && floor < defs.NumFloors-1 {
			return defs.DirnUp
		} else {
			return defs.DirnDown
		}
	case defs.DirnStop:
		if q.isOrdersAbove(floor) {
			return defs.DirnUp
		} else if q.isOrdersBelow(floor) {
			return defs.DirnDown
		} else {
			return defs.DirnStop
		}
	default:
		log.Printf("ChooseDirection(): called with invalid direction %d, returning stop\n", currDir)
		return defs.DirnStop
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case defs.DirnDown:
		return q.isActiveOrder(floor, defs.ButtonCallDown) ||
			q.isActiveOrder(floor, defs.ButtonCommand) ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case defs.DirnUp:
		return q.isActiveOrder(floor, defs.ButtonCallUp) ||
			q.isActiveOrder(floor, defs.ButtonCommand) ||
			floor == defs.NumFloors-1 ||
			!isOrdersAbove(floor)
	case defs.DirnStop:
		return q.isActiveOrder(floor, defs.ButtonCallDown) ||
			q.isActiveOrder(floor, defs.ButtonCallUp) ||
			q.isActiveOrder(floor, defs.ButtonCommand)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", direction)
		return false
	}
}

func (q *queue) print() {
	var status string
	var lifts string

	for f := nF - 1; f >= 0; f-- {
		if q.q[f][buttUp].active {
			status += "↑"
			lifts += "(↑ " + q.q[f][buttUp].addr + ")"
		} else {
			status += " "
		}
		if q.q[f][buttInt].active {
			status += "×"
			lifts += "(× " + q.q[f][buttInt].addr + ")"
		} else {
			status += " "
		}
		if q.q[f][buttDown].active {
			status += "↓   "
			lifts += "(↓ " + q.q[f][buttDown].addr + ")"
		} else {
			status += " "
		}
		status += lifts + "\n"
		lifts = ""
	}
	fmt.Printf(status)
}

func (q *queue) deepCopy() *queue {
	var copy queue
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			copy[f][b] = q[f][b]
		}
	}
	return copy
}

// this should run on a copy of local queue
func (q *queue) calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q.setOrder(targetFloor, targetButton, orderStatus{true, ""})

	cost := 0
	floor := prevFloor
	dir := currDir

	// Go to valid state (a floor/dir that mirrors a button)
	if currFloor == -1 {
		// Between floors, add 1 cost
		cost++
	} else if dir != defs.DirnStop {
		// At floor, but moving, add 2 cost
		cost += 2
		// Pass på at currFloor == prevFloor!
	}
	floor, dir = incrementFloor(floor, dir)

	for !(floor == targetFloor && q.shouldStop(floor, dir)) {
		if q.shouldStop(floor, dir) {
			cost += 2
		}
		dir = q.chooseDirection(floor, dir)
		floor, dir = incrementFloor(floor, dir)
		cost += 2
	}
	return cost
}

func incrementFloor(floor, dir int) (floor, dir int) {
	switch dir {
	case defs.DirnDown:
		floor--
	case defs.DirnUp:
		floor++
	case defs.DirnStop:
		fmt.Println("incrementFloor(): direction stop, not incremented (this is okay)")
	default:
		fmt.Println("incrementFloor(): invalid direction, not incremented")
	}
	if floor == 0 && dir == defs.DirnDown {
		dir = defs.DirnUp
	}
	if floor == defs.NumFloors-1 && dir == defs.DirnUp {
		dir = defs.DirnDown
	}
}

func updateLocalQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if shared.isActiveOrder(f, b) {
				if b != defs.ButtonCommand && shared.q[f][b].addr == defs.Laddr.String() {
					// set local order f b
					local.setOrder(f, b, orderStatus{true, ""})
				}
			}
		}
	}
}
