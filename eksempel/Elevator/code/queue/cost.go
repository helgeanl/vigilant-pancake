package queue

import (
	def "config"
	"log"
)

// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// floors on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between floors.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q := local.deepCopy()
	q.setOrder(targetFloor, def.BtnInside, orderStatus{true, "", nil})

	cost := 0
	floor := prevFloor
	dir := currDir

	if currFloor == -1 {
		// Between floors, add 1 cost.
		cost++
	} else if dir != def.DirStop {
		// At floor, but moving, add 2 cost.
		cost += 2
	}
	floor, dir = incrementFloor(floor, dir)

	// Simulate how the lift will move, and accumulate cost until it 'reaches' target.
	// Break after 10 iterations to assure against a stuck loop.
	for n := 0; !(floor == targetFloor && q.shouldStop(floor, dir)) && n < 10; n++ {
		if q.shouldStop(floor, dir) {
			cost += 2
			q.setOrder(floor, def.BtnUp, inactive)
			q.setOrder(floor, def.BtnDown, inactive)
			q.setOrder(floor, def.BtnInside, inactive)
		}
		dir = q.chooseDirection(floor, dir)
		floor, dir = incrementFloor(floor, dir)
		cost += 2
	}
	return cost
}

func incrementFloor(floor, dir int) (int, int) {
	switch dir {
	case def.DirDown:
		floor--
	case def.DirUp:
		floor++
	case def.DirStop:
		// Don't increment.
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalln(def.ColR, "incrementFloor(): invalid direction, not incremented", def.ColN)
	}

	if floor <= 0 && dir == def.DirDown {
		dir = def.DirUp
		floor = 0
	}
	if floor >= def.NumFloors-1 && dir == def.DirUp {
		dir = def.DirDown
		floor = def.NumFloors - 1
	}
	return floor, dir
}
