package queue

import def "definitions"

//CalcCost calculates the total amount of work it takes to reach the given floor.
func CalcCost(currDir, currFloor, prevFloor, targetFloor, targetBtn int) int {
	totCost := 0
	dir := currDir
	//Is the target above or below
	targetDir := targetFloor - prevFloor
	if currFloor == -1 {
		// Elevator is between floors, +1 cost
		totCost++
	} else if dir != def.DirStop {
		// Elevator is at floor but not IDLE +2 cost
		totCost += 2
	}
	if dir != def.DirStop {
		if targetDir != dir {
			// If the elevator must switch direction, +10 cost
			totCost += 10
		}
	}

	//Add +1 cost for every stop on the way to target
	if targetDir > 0 && dir == def.DirUp || dir == def.DirStop {
		for floor := prevFloor; floor < targetFloor || floor == def.NumFloors; floor++ {
			if queue.hasLocalRequest(floor, targetBtn) || queue.hasLocalRequest(floor, def.BtnCab) {
				totCost++
			}
			totCost++
		}
	}
	if targetDir < 0 && dir == def.DirDown || dir == def.DirStop {
		for floor := prevFloor; floor > targetFloor || floor == 0; floor-- {
			if queue.hasLocalRequest(floor, targetBtn) || queue.hasLocalRequest(floor, def.BtnCab) {
				totCost++
			}
			totCost++
		}
	}
	return totCost
}
