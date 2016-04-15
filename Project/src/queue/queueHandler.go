package queue

import (
	def "definitions"
	"fmt"
	"log"
	"strconv"
	"time"
)

type RequestStatus struct {
	Status bool
	Addr   string      `json:"-"`
	Timer  *time.Timer `json:"-"`
}

type QueueType struct {
	Matrix [def.NumFloors][def.NumButtons]RequestStatus
}

var queue QueueType
var RequestTimeoutChan = make(chan def.BtnPress, 10)
var NewRequest = make(chan bool, 10)
var CostReply = make(chan def.Message, 10)
var LightUpdate = make(chan def.LightUpdate, 10)
var takeBackup = make(chan bool, 10)

func Init(outgoingMsg chan def.Message) {
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialized.", def.ColN)

}

func AddRequest(floor int, btn int, addr string) {

	//if !queue.hasRequest(floor, btn) {

	queue.setRequest(floor, btn, RequestStatus{true, addr, nil})

	go queue.startTimer(floor, btn)
	LightUpdate <- def.LightUpdate{floor, btn, true}
	if addr == def.LocalIP {
		log.Println(def.ColW, "Request is local", def.ColN)
		NewRequest <- true
	} else {
		log.Println(def.ColW, "Request is not local", def.ColN)
	}
	//}
}

func RemoveRequest(floor, btn int) {
	LightUpdate <- def.LightUpdate{floor, btn, false}
	queue.stopTimer(floor, btn)
	queue.setRequest(floor, btn, RequestStatus{Status: false, Addr: "", Timer: nil})
}

func RemoveLocalRequestsAt(floor int, outgoingMsgCh chan<- def.Message) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if queue.Matrix[floor][btn].Addr == def.LocalIP {
			RemoveRequest(floor, btn)
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
			if queue.Matrix[floor][btn].Addr == addr { /////////////////////////////////Maybe we need to stop the timer??
				ReassignRequest(floor, btn, outgoingMsgCh)
			}
		}
	}
}

func ReassignRequest(floor, btn int, outgoingMsg chan<- def.Message) {
	RemoveRequest(floor, btn) ////////////////////////////////////////
	log.Println(def.ColB, "Reassigning request", def.ColN)
	outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
}

// Set status of request, sync request lights, take backup
func (q *QueueType) setRequest(floor, btn int, request RequestStatus) { ////////////////////////////////////////////////////////////////////////////////////////////////////////
	q.Matrix[floor][btn] = request
	takeBackup <- true
	printQueue()
}

// Start timer for request in queue
func (q *QueueType) startTimer(floor, btn int) {
	log.Println(def.ColW, "Start request timer", def.ColN)
	q.Matrix[floor][btn].Timer = time.NewTimer(def.RequestTimeoutDuration)
	<-q.Matrix[floor][btn].Timer.C
	// Wait until timeout
	if q.Matrix[floor][btn].Status {
		RequestTimeoutChan <- def.BtnPress{floor, btn}
		log.Println(def.ColW, "Request timer is done!", def.ColN)
	}
}

func (q *QueueType) stopTimer(floor, btn int) {
	if q.Matrix[floor][btn].Timer != nil {
		q.Matrix[floor][btn].Timer.Stop()
		log.Println(def.ColR, "Timer on Floor: ", floor, " Button: ", btn, "stopped.")
	}
}

func printQueue() {
	fmt.Println(def.ColB)
	fmt.Println("********************************")
	fmt.Println("*       Up      Down     Cab   *")
	for f := def.NumFloors - 1; f >= 0; f-- {
		s := "* " + strconv.Itoa(f+1) + "  "

		if queue.hasRequest(f, def.BtnHallUp) {
			s += "( ." + queue.Matrix[f][def.BtnHallUp].Addr[12:15] + " ) "
		} else {
			s += "(      ) "
		}
		if queue.hasRequest(f, def.BtnHallDown) {
			s += "( ." + queue.Matrix[f][def.BtnHallDown].Addr[12:15] + " ) "
		} else {
			s += "(      ) "
		}
		if queue.hasRequest(f, def.BtnCab) {
			s += "(  x  ) *"
		} else {
			s += "(     ) *"
		}
		fmt.Println(s)
	}
	fmt.Println("********************************")
	fmt.Println(def.ColN)
}
