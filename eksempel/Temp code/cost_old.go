package cost

import (
	"../defs"
	"../queue"
	"errors"
	"fmt"
	"log"
)

var _ = errors.New

// CalculateCost calculates how much effort it takes this lift to carry out
// the given order. Each sheduled stop on the way there and each travel
// between adjacent floors will add cost 2. Cost 1 is added if the lift
// starts between floors.
func CalculateCost(targetFloor, targetButton, fsmFloor, fsmDir, currFloor int) (int, error) {
	// Bug! Lift between floor 2 and 1 going down: Gives cost 1 to floor 2, should be 5.
	if (targetButton != defs.ButtonCallUp) && (targetButton != defs.ButtonCallDown) {
		return 0, fmt.Errorf("CalculateCost() called with invalid order: %d\n", targetButton)
	}

	fmt.Printf("CalculateCost called with parameters %d, %d, %d, %d, %d\n",
		targetFloor, targetButton, fsmFloor, fsmDir, currFloor)

	fmt.Printf("Cost floor sequence: ")

	cost := 0
	var err error

	// Between floors
	if currFloor == -1 {
		fmt.Printf("%d >>> ", currFloor)
		cost += 1
		fsmFloor, err = incrementFloor(fsmFloor, fsmDir)
		if err != nil {
			defer log.Println(err)
		}
	} else if (fsmDir != defs.DirnStop) && (fsmFloor != targetFloor) {
		// Not between floors but moving (i.e. departing or arriving)
		cost += 2
		fsmFloor, err = incrementFloor(fsmFloor, fsmDir)
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Printf("%d", fsmFloor)
	for !(fsmFloor == targetFloor && queue.ShouldStop(fsmFloor, fsmDir)) {
		if queue.ShouldStop(fsmFloor, fsmDir) {
			cost += 2
			fmt.Printf("(S)")
		}
		fsmDir = queue.ChooseDirection(fsmFloor, fsmDir)
		fsmFloor, err = incrementFloor(fsmFloor, fsmDir)
		if err != nil {
			defer log.Println(err)
		}
		cost += 2
		fmt.Printf(" >>> %d", fsmFloor)
	}
	fmt.Printf(" = cost %d\n", cost)

	return cost, nil
}

func incrementFloor(floor int, direction int) (int, error) {
	switch direction {
	case defs.DirnDown:
		floor--
	case defs.DirnUp:
		floor++
	case defs.DirnStop:
		return floor, errors.New("Error(ish): Direction stop, floor not incremented.")
	default:
		return floor, errors.New("Error: Invalid direction, floor not incremented.")
	}

	if (floor < 0) or (floor >= defs.NumFloors) {
		return floor, fmt.Printf("Error: Floor incremented to invalid floor %d.", floor)
	}

	return floor, nil
}
