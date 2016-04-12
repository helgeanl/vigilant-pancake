package queue

import (
	def "definitions"
	"fmt"
	"time"
)




func (q *queue) ChooseDirection(floor, dir int) int {
	switch dir {
	case def.DirUp:
		if q.hasRequestsAbove(floor){
			return def.DirUp
		} else if q.hasRequestsBelow(floor){
			return def.DirDown
		} else{
			return def.DirStop
		}
	case def.DirDown, def.DirStop:
		if q.hasRequestsBelow(floor) {
			return def.DirDown
		} else if q.hasRequestsAbove(floor) {
			return def.DirUp
		} else {
			return def.DirStop
		}
	default:
		def.CloseConnectionChan <- true
		//// !!!!!!!!!!!!!!!
		def.Restart.Run()
		log.Printf("%sChooseDirection(): called with invalid direction %d, returning stop%s\n", def.ColR, dir, def.ColN)
		return 0
	}
}

func (q *queue) ShouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return
			q.hasLocalRequest(floor, def.BtnHallDown) ||
			q.hasLocalRequest(floor, def.BtnCab) ||
			!q.hasRequestsBelow(floor)
	case def.DirUp:
		return
			q.hasLocalRequest(floor, def.BtnHallUp) ||
			q.hasLocalRequest(floor, def.BtnCab) ||
			!q.hasRequestsAbove(floor)
	case def.DirStop:
	default:
		def.CloseConnectionChan <- true
		//// !!!!!!!!!!!!!!!
		def.Restart.Run()
		log.Fatalln(def.ColR, "This direction doesn't exist", def.ColN)
	}
	return false
}

func HasRequest(floor,btn int)bool{
	return queue.hasRequest(floor,btn)
}
func (q *queue) hasRequest(floor, btn int) bool {
	return q.matrix[floor][btn].status
}
func (q *queue) hasLocalRequest(floor, btn int) bool {
	return q.matrix[floor][btn].status && q.matrix[floor][btn].addr == localIP
}

func (q *queue) hasRequestAbove(floor int) bool {
	for f := floor + 1; f < def.NumFloors; f++ {
		for b:= 0; b < def.NumButtons; b++ {
			if q.hasLocalRequest(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) hasRequestsBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.hasLocalRequest(f, b) {
				return true
			}
		}
	}
	return false
}
