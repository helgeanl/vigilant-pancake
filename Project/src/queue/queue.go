package queue

import (
	def "definitions"
	"fmt"
	"time"
)

var queue queue
var RequestTimeoutChan chan def.BtnPress
//make a request inactive
const inactive = requestStatus{status: false, addr: "", timer: nil}

type requestStatus struct {
	status bool
	addr   string       `json:"-"`
	timer  *timer.Timer `json:"-"`
}

type queue struct {
	matrix [def.Numfloors][def.NumButtons]requestStatus
}





/// -------------------




func Init(newRequestTemp chan bool, outgoingMsg chan def.Message) {
	newRequest = newRequestTemp /// ??????
	go updateLocalQueue()
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialised.", def.ColN)
}

func AddRequestAt(floor int, btn int, addr string){
	if !queue.hasRequest(floor,btn){
		queue.setRequest(floor,btn,requestStatus{true,addr,nil})
		queue.startTimer(floor, btn)
	}
}

// Go through queue, and resend requests belonging to dead elevator
func ReassignRequests(addr string, outgoingMsg chan<-def.Message){
	for floor := 0; floor < def.NumFloors; floor++{
		for btn := 0; btn < def.NumButtons; btn++{
			if queue.matrix[floor][btn].addr = addr{
				outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
			}
		}
	}
}

func RemoveOrderAt(floor,btn int){
	queue.setRequest(floor,btn,inactive)
}





// Set status of request, sync request lights, take backup
func (q *queue) setRequest(floor, btn int, request requestStatus){
	q.matrix[floor][btn] = request
	btnLightCh <- def.BtnPress{floor,btn}
	// take backup
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



// -------------------
