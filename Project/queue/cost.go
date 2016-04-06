package queue

import (
	"config"
	"fmt"
)

//CalcCost calculates the total amount of work it takes to reach the 
//given floor. 
func CalcCost(currDir, currFloor, prevFloor, targetFloor, targetBtn int) int{
	totCost := 0
	dir 	:= currDir
	//Is the target above or below
	targetDir := targetFloor - prevFloor


	if currFloor == -1 {
		//Elevator is between floors, +1 cost
		totCost++
	else if dir != config.DirStop {
		//Elevator is at floor but not IDLE +2 cost
		totCost += 2
	}
	if dir != config.DirStop {
		if Targetdir != dir {
			//if the elevator must switch direction, +10 cost
			totCost += 10
		}
	}
	
	//Add +1 cost for every stop on the way to target
	//I assume here that currFloor==prevFloor if the elevator is at
	// a floor
	if targetDir < 0 && dir == config.DirUp{
		for floor := prevFloor; floor < targetFloor || floor == config.Numfloors; floor++{
			if isOrder(floor, config.Btn) || isOrder(floor, config.BtnInside){
				totCost ++
			}
			totCost++
		}
	}
	if targetDir > 0 && dir == config.DirDown{
		for floor := prevFloor; floor > targetFloor || floor == 0 ; floor--{
			if isOrder(floor, targetBtn) || isOrder(floor, config.BtnInside){
				totCost ++
			}
			totCost++
		}
	}

	return totCost
}