package queue

import (
	"../defs"
	"log"
	"fmt"
)

var _ = fmt.Printf

type sharedOrder struct {
	isOrderActive    bool
	assignedLiftAddr string
}
var blankOrder = sharedOrder{isOrderActive: false, assignedLiftAddr: ""}

var (
	localQueue [defs.NumFloors][defs.NumButtons]bool
	sharedQueue [defs.NumFloors][defs.NumButtons]sharedOrder
)

// --------------- PUBLIC: ---------------

func Init() {
	resetLocalQueue()
	resetSharedQueue()
}

// AddInternalOrder adds internal orders to the local queue.
func AddInternalOrder(floor int, button int) {
	/*if button != defs.ButtonCommand {
		// error!
		return
	}*/
	localQueue[floor][button] = true
}

// RemoveInternalOrder adds internal orders to the local queue.
func RemoveInternalOrder(floor int, button int) {
	/*if button != defs.ButtonCommand {
		// error!
		return
	}*/
	localQueue[floor][button] = false
}

// ChooseDirection returns the direction the lift should continue after the
// current floor.
func ChooseDirection(currFloor int, currDir int) int {
	if !IsAnyOrders() {
		log.Println("ChooseDirection(): no orders!")
		return defs.DirnStop
	}
	switch currDir {
	case defs.DirnDown:
		if isOrdersBelow(currFloor) && currFloor > 0 {
			return defs.DirnDown
		} else {
			return defs.DirnUp
		}
	case defs.DirnUp:
		if isOrdersAbove(currFloor) && currFloor < defs.NumFloors-1 {
			return defs.DirnUp
		} else {
			return defs.DirnDown
		}
	case defs.DirnStop:
		if isOrdersAbove(currFloor) {
			return defs.DirnUp
		} else if isOrdersBelow(currFloor) {
			return defs.DirnDown
		} else {
			return defs.DirnStop
		}
	default:
		log.Printf("localQueue: ChooseDirection called with invalid direction %d!\n", currDir)
		return defs.DirnStop
	}
}

// ShouldStop returns whether the lift should stop at the given floor, if
// going in the given direction.
func ShouldStop(floor int, direction int) bool {
	switch direction {
	case defs.DirnDown:
		return localQueue[floor][defs.ButtonCallDown] ||
			localQueue[floor][defs.ButtonCommand] ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case defs.DirnUp:
		return localQueue[floor][defs.ButtonCallUp] ||
			localQueue[floor][defs.ButtonCommand] ||
			floor == defs.NumFloors-1 ||
			!isOrdersAbove(floor)
	case defs.DirnStop:
		return localQueue[floor][defs.ButtonCallDown] ||
			localQueue[floor][defs.ButtonCallUp] ||
			localQueue[floor][defs.ButtonCommand]
	default:
		log.Printf("localQueue: ShouldStop called with invalid direction %d!\n", direction)
		return false
	}
}

// RemoveOrdersAt removes all orders in the local queue at the given floor.
func RemoveOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		localQueue[floor][b] = false
		sharedQueue[floor][b] = blankOrder
	}
	SendOrderCompleteMessage(floor)
}

// IsOrder returns whether there in an order with the given floor and button
// in the local queue.
func IsOrder(floor int, button int) bool {
	return localQueue[floor][button]
}

// ReassignOrders finds all orders assigned to the given dead lift, removes
// them from the shared queue, and sends them on the network as new, un-
// assigned orders.
func ReassignOrders(deadAddr string) { // better name plz
	// loop thru shared queue
	// remove all orders assigned to the dead lift
	// send neworder-message for each removed order
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if sharedQueue[f][b].assignedLiftAddr == deadAddr {
				sharedQueue[f][b] = blankOrder
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

// AddSharedOrder adds the given order to the shared queue.
func AddSharedOrder(floor, button int, addr string) {
	sharedQueue[floor][button] = sharedOrder{true, addr}
	updateLocalQueue()
}

func CopyLocalQueue() [defs.NumFloors][defs.NumButtons]bool {
	var copy [defs.NumFloors][defs.NumButtons]bool
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			copy[f][b] = localQueue[f][b]
		}
	}
	return copy
}

// --------------- PRIVATE: ---------------

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func IsAnyOrders() bool {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func updateLocalQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if sharedQueue[f][b].isOrderActive {
				//fmt.Printf("updateLocalQueue(): laddr = %s, addr in que = %s\n",
				//	defs.Laddr.String(), sharedQueue[f][b].assignedLiftAddr)
				if b != defs.ButtonCommand &&
				sharedQueue[f][b].isOrderActive &&
				sharedQueue[f][b].assignedLiftAddr == defs.Laddr.String() {
						//Laddr gets changed in UdpInit, i think this is fine
						localQueue[f][b] = true
				}
			}
		}
	}
}

// RemoveSharedOrder removes the giver order from the shared queue. This is
// done when an order is completed.
func RemoveSharedOrdersAt(floor int) { // Rename to RemoveOrder()!
	for b := 0; b < defs.NumButtons; b++ {
		sharedQueue[floor][b] = blankOrder
	}
	updateLocalQueue()
}

func resetLocalQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			localQueue[f][b] = false
		}
	}
}

func resetSharedQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			sharedQueue[f][b] = blankOrder
		}
	}
}

func PrintQueues() {
	fmt.Println("Local   Shared")
	for f := defs.NumFloors-1; f >= 0; f-- {
		lifts := "   "

		if localQueue[f][defs.ButtonCallUp] {
			fmt.Printf("↑")
		} else {
			fmt.Printf(" ")
		}
		if localQueue[f][defs.ButtonCommand] {
			fmt.Printf("×")
		} else {
			fmt.Printf(" ")
		}
		if localQueue[f][defs.ButtonCallDown] {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}
		if sharedQueue[f][defs.ButtonCallUp].isOrderActive {
			fmt.Printf("↑")
			lifts += "(↑ " + defs.LastPartOfIp(sharedQueue[f][defs.ButtonCallUp].assignedLiftAddr) + ")"
		} else {
			fmt.Printf(" ")
		}
		if sharedQueue[f][defs.ButtonCallDown].isOrderActive {
			fmt.Printf("↓")
			lifts += "(↓ " + defs.LastPartOfIp(sharedQueue[f][defs.ButtonCallDown].assignedLiftAddr) + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", lifts)
		fmt.Println()
	}
}

/*func updateSharedQueue(floor int, button int) {
	// If order completed was assigned to this elevator: Remove from shared queue
	if button == defs.ButtonCommand {
		// error
		return
	}

	if sharedQueue[floor][button].isOrderActive
	&& sharedQueue[floor][button].assignedLiftAddr == laddr {
		sharedQueue[floor][button].isOrderActive = false
		sharedQueue[floor][button].assignedLiftAddr = invalidAddr
	}
}*/
