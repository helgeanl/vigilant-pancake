package fsm

import (
	"time"
	"config"
)

//doorTimer makes sure that the door stays open for config.DoorOpenTime
//seconds.
func doorTimer(reset <-chan bool, timeout chan<- bool){
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select{
		case <-reset:
			timer.Reset(config.DoorOpenTime)
		case <-timer.C:
			timer.Stop()
			timeout <- true
		}
	}
}