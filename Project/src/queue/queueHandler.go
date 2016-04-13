package queue

import (
	def "definitions"
	"fmt"
	"strconv"
	"time"
	"log"
)

type requestStatus struct {
	status bool
	addr   string      `json:"-"`
	timer  *time.Timer `json:"-"`
}

type queueType struct {
	matrix [def.NumFloors][def.NumButtons]requestStatus
}

var queue queueType
var RequestTimeoutChan = make(chan def.BtnPress,10)
var NewRequest = make(chan bool,10)
var CostReply = make(chan def.Message,10)
var takeBackup = make(chan bool,10)


func Init(outgoingMsg chan def.Message) {
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialized.", def.ColN)

}

func AddRequest(floor int, btn int, addr string) {

	//if !queue.hasRequest(floor, btn) {
		
		queue.setRequest(floor, btn, requestStatus{true, addr, nil})

		go queue.startTimer(floor, btn)
		if addr == def.LocalIP {
			log.Println(def.ColW,"Request is local",def.ColN)
			NewRequest <- true
		}
		log.Println(def.ColW,"Request is not local",def.ColN)
	//}
}

func RemoveRequest(floor, btn int) {
	queue.setRequest(floor, btn, requestStatus{status: false, addr: "", timer: nil})
}

func RemoveLocalRequestsAt(floor int, outgoingMsgCh chan<- def.Message) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if queue.matrix[floor][btn].addr == def.LocalIP {
			queue.setRequest(floor, btn, requestStatus{status: false, addr: "", timer: nil})
			if btn != def.BtnCab {
				outgoingMsgCh <- def.Message{Category: def.CompleteRequest, Floor: floor, Button: btn}
			}
		}
	}
}

// Go through queue, and resend requests belonging to dead elevator
func ReassignAllRequestsFrom(addr string, outgoingMsgCh chan<- def.Message) {
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if queue.matrix[floor][btn].addr == addr {
				ReassignRequest(floor, btn, outgoingMsgCh)
			}
		}
	}
}

func ReassignRequest(floor, btn int, outgoingMsg chan<- def.Message) {
	outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
}

// Set status of request, sync request lights, take backup
func (q *queueType) setRequest(floor, btn int, request requestStatus) {
	q.matrix[floor][btn] = request
	takeBackup <- true
	printQueue()
}

// Start timer for request in queue
func (q *queueType) startTimer(floor, btn int) {
	log.Println(def.ColW,"Start request timer",def.ColN)
	q.matrix[floor][btn].timer = time.NewTimer(def.RequestTimeoutDuration)
	<-q.matrix[floor][btn].timer.C
	// Wait until timeout
	RequestTimeoutChan <- def.BtnPress{floor, btn}
	log.Println(def.ColW,"Request timer is done!",def.ColN)
}

func (q *queueType) stopTimer(floor, btn int) {
	if q.matrix[floor][btn].timer != nil {
		q.matrix[floor][btn].timer.Stop()
	}
}

func printQueue() {
	fmt.Printf(def.ColB)
	fmt.Println("****** Queue ****** ")
	for f := def.NumFloors - 1; f >= 0; f-- {
		s := strconv.Itoa(f+1)

		if queue.hasRequest(f, def.BtnHallUp) {
			s += " (↑ " + queue.matrix[f][def.BtnHallUp].addr[12:15] + " ) "
		} else {
			fmt.Printf("(     )")
		}
		if queue.hasRequest(f, def.BtnHallDown) {
			s += " (↓ " + queue.matrix[f][def.BtnHallDown].addr[12:15] + " ) "
		} else {
			fmt.Printf("(     )")
		}
		if queue.hasRequest(f, def.BtnCab) {
			s += " (  x  )"
		} else {
			s += " (     )"
		}
		fmt.Printf("%s", s)
		fmt.Println()
	}
	fmt.Printf(def.ColN)
}
