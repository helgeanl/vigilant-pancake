package queue

import (
	def "definitions"
	"fmt"
	"time"
)

var localIP string
var queue queue
var RequestTimeoutChan = make(chan def.BtnPress)
var costReply = make(chan def.Message )
var takeBackup chan bool
var newRequest chan bool


//make a request inactive !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
const inactive = requestStatus{status: false, addr: "", timer: nil}

type requestStatus struct {
	status bool
	addr   string       `json:"-"`
	timer  *timer.Timer `json:"-"`
}

type queue struct {
	matrix [def.Numfloors][def.NumButtons]requestStatus
}

func Init(addr string, outgoingMsg chan def.Message) {
	localIP = addr
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialised.", def.ColN)

}


func AddRequest(floor int, btn int, addr string){
	if !HasRequest(floor,btn){
		if addr = localIP{
			newRequest <- true
		}
		queue.setRequest(floor,btn,requestStatus{true,addr,nil})
		queue.startTimer(floor, btn)
	}
}

func RemoveRequest(floor,btn int){
	queue.setRequest(floor,btn,inactive)
}

// Go through queue, and resend requests belonging to dead elevator
func ReassignAllRequestsFrom(addr string, outgoingMsg chan<-def.Message){
	for floor := 0; floor < def.NumFloors; floor++{
		for btn := 0; btn < def.NumButtons; btn++{
			if queue.matrix[floor][btn].addr = addr{
				outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
			}
		}
	}
}

func ReassignRequest(floor,btn int, outgoingMsg chan<-def.Message){
	outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
}

// Set status of request, sync request lights, take backup
func (q *queue) setRequest(floor, btn int, request requestStatus){
	q.matrix[floor][btn] = request
	//btnLightCh <- def.BtnPress{floor,btn}
	takeBackup <- true
	// printt
}

// Start timer for request in queue
func (q * queue) startTimer(floor, btn int){
	q.matrix[floor][btn].timer = time.NewTimer(def.RequestTimeoutDuration)
	<-q.matrix[floor][btn].timer.C
	// Wait until timeout
	RequestTimeoutChan <- def.Btnpress{floor, btn}
}

func (q * queue) stopTimer(floor, btn int){
	if q.matrix[floor][btn].timer != nil{
		q.matrix[floor][btn].timer.Stop()
	}
}
