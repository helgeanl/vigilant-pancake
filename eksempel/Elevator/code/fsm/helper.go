package fsm

// Direction returns the current direction as stored in the state machine.
func Direction() int {
	return dir
}

// Floor returns the current floor as stored in the state machine.
// If the lift is between floors, it returns the most recent floor.
func Floor() int {
	return floor
}

func stateString(state int) string {
	switch state {
	case idle:
		return "idle"
	case moving:
		return "moving"
	case doorOpen:
		return "door open"
	default:
		return "error: bad state"
	}
}
