package queue

import (
	def "definitions"
	"fmt"
	"time"
)


var queue queue
var RequestTimeoutChan = make(chan def.BtnPress)
var NewRequest = make(chan bool)
var CostReply = make(chan def.Message )
var takeBackup = make(chan bool)


type requestStatus struct {
	status bool
	addr   string       `json:"-"`
	timer  *timer.Timer `json:"-"`
}

type queue struct {
	matrix [def.Numfloors][def.NumButtons]requestStatus
}

func Init(outgoingMsg chan def.Message) {
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialised.", def.ColN)

}

func AddRequest(floor int, btn int, addr string){
	if !queue.hasRequest(floor,btn){
		if addr = def.LocalIP{
			NewRequest <- true
		}
		queue.setRequest(floor,btn,requestStatus{true,addr,nil})
		queue.startTimer(floor, btn)
	}
}

func RemoveRequest(floor, btn int){
	queue.setRequest(floor,btn,requestStatus{status: false, addr: "", timer: nil})
}

func RemoveLocalRequestsAt(floor int, outgoingMsgCh chan def.Message){
	for btn :=0; btn < def.NumButtons; btn++{
		if queue[floor][btn].addr == def.LocalIP{
			queue.setRequest(floor,btn,requestStatus{status: false, addr: "", timer: nil})
			if btn != def.BtnCab{
				outgoingMsg <- def.Message{Category: def.CompleteRequest, Floor: floor, Button: btn}
			}
		}
	}
}

// Go through queue, and resend requests belonging to dead elevator
func ReassignAllRequestsFrom(addr string, outgoingMsgCh chan def.Message){
	for floor := 0; floor < def.NumFloors; floor++{
		for btn := 0; btn < def.NumButtons; btn++{
			if queue.matrix[floor][btn].addr = addr{
				ReassignRequest(floor,btn, outgoingMsgCh)
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
	takeBackup <- true
	printQueue()
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

func printQueue() {
	fmt.Printf(def.ColC)
	fmt.Println("****** Queue ****** ")
	for f := def.NumFloors - 1; f >= 0; f-- {
		s := strconv.Itoa(f-1)

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
			s1 += " (     )"
		}
		fmt.Printf("%s", s)
		fmt.Println()
	}
	fmt.Printf(def.ColN)
}
