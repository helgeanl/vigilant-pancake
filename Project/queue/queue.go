package queue

import (
	"config"
	"fmt"
	"time"
)

type orderStatus struct {
	active bool
	addr   string       `json:"-"`
	timer  *timer.Timer `json:"-"`
}

type queue struct {
	qMatrix [config.Numfloors][config.NumButtons]ordreStatus
}

//make an order inactive
var inactive = ordreStatus{active: false, addr: "", timer: nil}

func isOrder(floor, btn int) int {
	if "Den har en bestilling" {
		return true
	}
	return false
}
