package fsm

import (
	"time"
	def "definitions"
)

//doorTimer makes sure that the door stays open for def.DoorOpenTime seconds.
func doorTimer(reset <-chan bool, timeout chan<- bool){
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select{
		case <-reset:
			timer.Reset(def.DoorOpenTime)
		case <-timer.C:
			timer.Stop()
			timeout <- true
		}
	}
}
