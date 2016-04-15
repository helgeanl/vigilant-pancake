package queue

import (
	def "definitions"
	"log"
)

func ChooseDirection(floor, dir int) int {
	switch dir {
	case def.DirUp:
		if queue.hasRequestsAbove(floor) {
			return def.DirUp
		} else if queue.hasRequestsBelow(floor) {
			return def.DirDown
		} else {
			return def.DirStop
		}
	case def.DirDown, def.DirStop:
		if queue.hasRequestsBelow(floor) {
			return def.DirDown
		} else if queue.hasRequestsAbove(floor) {
			return def.DirUp
		} else {
			return def.DirStop
		}
	default:
		return 0
	}
}

func ShouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return queue.hasLocalRequest(floor, def.BtnHallDown) ||
			queue.hasLocalRequest(floor, def.BtnCab) ||
			!queue.hasRequestsBelow(floor)
	case def.DirUp:
		return queue.hasLocalRequest(floor, def.BtnHallUp) ||
			queue.hasLocalRequest(floor, def.BtnCab) ||
			!queue.hasRequestsAbove(floor)
	case def.DirStop:
	default:
		log.Fatalln(def.ColR, "This direction doesn't exist", def.ColN)
	}
	return false
}

func HasRequest(floor, btn int) bool {
	return queue.hasRequest(floor, btn)
}

func (q *QueueType) hasRequest(floor, btn int) bool {
	return q.Matrix[floor][btn].Status
}

func (q *QueueType) hasLocalRequest(floor, btn int) bool {
	return q.Matrix[floor][btn].Status && q.Matrix[floor][btn].Addr == def.LocalIP
}

func (q *QueueType) hasRequestsAbove(floor int) bool {
	for f := floor + 1; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.hasLocalRequest(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *QueueType) hasRequestsBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.hasLocalRequest(f, b) {
				return true
			}
		}
	}
	return false
}
