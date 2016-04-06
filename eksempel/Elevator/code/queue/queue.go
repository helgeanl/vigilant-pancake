// Package queue defines and stores queues for the lift. Two queues are used:
// One local queue containing all orders assigned to this particular lift,
// and one remote containing all external orders assigned to any lift on the
// network. The remote order also stores information about the IP of the lift
// assigned to each order, and it has a timer attached to each order. The
// timer makes sure that if an order is assigned to a lift, and we never
// receive an 'order complete' message, the order will still be handled.
package queue

import (
	def "config"
	"fmt"
	"log"
	"time"
)

// queue defines a queue, a 2D array of orderStatuses representing the
// buttons on the lift panel.
type queue struct {
	matrix [def.NumFloors][def.NumButtons]orderStatus
}

// orderStatus defines the status of an order: Whether it is active, which
// lift is assigned to take it, and how long it has been active. (The latter
// two are only used in the remote queue.)
type orderStatus struct {
	active bool
	addr   string      `json:"-"`
	timer  *time.Timer `json:"-"`
}

var inactive = orderStatus{active: false, addr: "", timer: nil}

var local queue
var remote queue

var updateLocal = make(chan bool)
var takeBackup = make(chan bool, 10)
var OrderTimeoutChan = make(chan def.Keypress)
var newOrder chan bool

func Init(newOrderTemp chan bool, outgoingMsg chan def.Message) {
	newOrder = newOrderTemp
	go updateLocalQueue()
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialised.", def.ColN)
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, "", nil})
	newOrder <- true
}

// AddRemoteOrder adds an order to the remote queue, and spawns a timer
// for the order. (If the order times out, it will be taken care of.)
func AddRemoteOrder(floor, button int, addr string) {
	alreadyExist := IsRemoteOrder(floor, button)
	remote.setOrder(floor, button, orderStatus{true, addr, nil})
	if !alreadyExist {
		go remote.startTimer(floor, button)
	}
	updateLocal <- true
}

// RemoveRemoteOrdersAt removes all orders at the given floor from the remote
// queue.
func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		remote.stopTimer(floor, b)
		remote.setOrder(floor, b, inactive)
	}
	updateLocal <- true
}

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int, outgoingMsg chan<- def.Message) {
	for b := 0; b < def.NumButtons; b++ {
		remote.stopTimer(floor, b)
		local.setOrder(floor, b, inactive)
		remote.setOrder(floor, b, inactive)
	}
	outgoingMsg <- def.Message{Category: def.CompleteOrder, Floor: floor}
}

// ShouldStop returns whether the lift should stop when it reaches the given
// floor, going in the given direction.
func ShouldStop(floor, dir int) bool {
	return local.shouldStop(floor, dir)
}

// ChooseDirection returns the direction the lift should continue after the
// current floor, going in the given direction.
func ChooseDirection(floor, dir int) int {
	return local.chooseDirection(floor, dir)
}

// IsLocalOrder returns whether there in an order with the given floor and
// button in the local queue.
func IsLocalOrder(floor, button int) bool {
	return local.isOrder(floor, button)
}

// IsRemoteOrder returns true if there is a order with the given floor and
// button in the remote queue.
func IsRemoteOrder(floor, button int) bool {
	return remote.isOrder(floor, button)
}

// ReassignOrders finds all orders assigned to a dead lift, removes them from
// the remote queue, and sends them on the network as new, unassigned orders.
func ReassignOrders(deadAddr string, outgoingMsg chan<- def.Message) {
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if remote.matrix[f][b].addr == deadAddr {
				remote.setOrder(f, b, inactive)
				outgoingMsg <- def.Message{Category: def.NewOrder, Floor: f, Button: b}
			}
		}
	}
}

// printQueues prints local and remote queue to screen in a somewhat legible
// manner.
func printQueues() {
	fmt.Printf(def.ColC)
	fmt.Println("Local   Remote")
	for f := def.NumFloors - 1; f >= 0; f-- {

		s1 := ""
		if local.isOrder(f, def.BtnUp) {
			s1 += "↑"
		} else {
			s1 += " "
		}
		if local.isOrder(f, def.BtnInside) {
			s1 += "×"
		} else {
			s1 += " "
		}
		fmt.Printf(s1)
		if local.isOrder(f, def.BtnDown) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}

		s2 := "   "
		if remote.isOrder(f, def.BtnUp) {
			fmt.Printf("↑")
			s2 += "(↑ " + remote.matrix[f][def.BtnUp].addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isOrder(f, def.BtnDown) {
			fmt.Printf("↓")
			s2 += "(↓ " + remote.matrix[f][def.BtnDown].addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", s2)
		fmt.Println()
	}
	fmt.Printf(def.ColN)
}

// updateLocalQueue checks remote queue for new orders assigned to this lift
// and copies them to the local queue.
func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if remote.isOrder(f, b) {
					if b != def.BtnInside && remote.matrix[f][b].addr == def.Laddr {
						if !local.isOrder(f, b) {
							local.setOrder(f, b, orderStatus{true, "", nil})
							newOrder <- true
						}
					}
				}
			}
		}
	}
}
