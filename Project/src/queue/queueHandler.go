package queue

import (
	def "definitions"
	"log"
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
var LightUpdate = make(chan def.LightUpdate, 10)
var takeBackup = make(chan bool, 10)

func AddRequest(floor int, btn int, addr string, NewRequest) {
	if !queue.hasRequest(floor,btn){
		queue.setRequest(floor, btn, RequestStatus{Status: true, Addr: addr, Timer: nil})
		if addr == def.LocalIP {
			NewRequest <- true
		} else {
			go queue.startTimer(floor, btn)
		}
	}
}

func RemoveRequest(floor, btn int) {
	queue.stopTimer(floor, btn)
	queue.setRequest(floor, btn, RequestStatus{Status: false, Addr: "", Timer: nil})
}

func RemoveLocalRequestsAt(floor int, outgoingMsgCh chan<- def.Message) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if queue.Matrix[floor][btn].Addr == def.LocalIP {
			if btn == def.BtnCab {
				RemoveRequest(floor, btn)
			} else {
				outgoingMsgCh <- def.Message{Category: def.CompleteRequest, Floor: floor, Button: btn}
			}
		}
	}
}

func ReassignRequest(floor, btn int, outgoingMsg chan<- def.Message) {
	RemoveRequest(floor, btn)
	log.Println(def.ColB, "Reassigning request", def.ColN)
	outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
}

// ReassignAllRequestsFrom goes through queue, and resend requests belonging to dead elevator
func ReassignAllRequestsFrom(addr string, outgoingMsgCh chan<- def.Message) {
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if queue.Matrix[floor][btn].Addr == addr {
				ReassignRequest(floor, btn, outgoingMsgCh)
			}
		}
	}
}

func (q *QueueType) setRequest(floor, btn int, request RequestStatus) {
	q.Matrix[floor][btn] = request
	LightUpdate <- def.LightUpdate{Floor: floor, Button: btn, UpdateTo: request.Status}
	takeBackup <- true
	printQueue()
}

func (q *QueueType) startTimer(floor, btn int) {
	q.Matrix[floor][btn].Timer = time.NewTimer(def.RequestTimeoutDuration)
	<-q.Matrix[floor][btn].Timer.C
	// Wait until timeout
	if q.Matrix[floor][btn].Status {
		RequestTimeoutChan <- def.BtnPress{floor, btn}
	}
}

func (q *QueueType) stopTimer(floor, btn int) {
	if q.Matrix[floor][btn].Timer != nil {
		q.Matrix[floor][btn].Timer.Stop()
	}
}
