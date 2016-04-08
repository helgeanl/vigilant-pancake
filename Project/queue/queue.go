package queue

import (
	def "definitions"
	"fmt"
	"time"
)

type orderStatus struct {
	active bool
	addr   string       `json:"-"`
	timer  *timer.Timer `json:"-"`
}

type queue struct {
	qMatrix [def.Numfloors][def.NumButtons]ordreStatus
}

//make an order inactive
var inactive = ordreStatus{active: false, addr: "", timer: nil}

func isOrder(floor, btn int) int {
	if "Den har en bestilling" {
		return true
	}
	return false
}

// requests_above
func (q *queue) hasOrdersAbove(floor int) bool {
	for f := floor + 1; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.hasOrder(f, b) { //q.matrix[f][b]
				return true
			}
		}
	}
	return false
}

// requests_below
func (q *queue) hasOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(floor, dir int) int {
	switch dir {
	case config.DirUp:
		if q.hasOrdersAbove(floor){
			return def.DirUp
		} else if q.hasOrdersBelow(floor){
			return def.DirDown
		} else{
			return def.DirStop
		}
	case def.DirDown, def.DirStop:
		if q.hasOrdersBelow(floor) {
			return def.DirDown
		} else if q.hasOrdersAbove(floor) {
			return def.DirUp
		} else {
			return def.DirStop
		}
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Printf("%sChooseDirection(): called with invalid direction %d, returning stop%s\n", def.ColR, dir, def.ColN)
		return 0
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return
			q.isOrder(floor, def.BtnDown) ||
			q.isOrder(floor, def.BtnInside) ||
			!q.isOrdersBelow(floor)
	case def.DirUp:
		return
			q.isOrder(floor, def.BtnUp) ||
			q.isOrder(floor, def.BtnInside) ||
			!q.isOrdersAbove(floor)
	case def.DirStop:
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalln(def.ColR, "This direction doesn't exist", def.ColN)
	}
	return false
}
