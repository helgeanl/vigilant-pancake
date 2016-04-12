// FMS for Elevator
// *** some comment
// events: timeout, floorArrived, newRequest
// state: idle, moving, doorOpen
package fsm

import (
	def "definitions"
	"log"
	"queue"
)

// Enumeration of Elevator behaviour
const (
	idle int = iota
	moving
	doorOpen
)

var Elevator struct {
	floor     int
	dir       int
	behaviour int
}

type Channels struct {
	// Events
	NewRequest   chan bool
	FloorReached chan int
	doorTimeout  chan bool
	// Hardware interaction
	MotorDir  chan int
	FloorLamp chan int
	DoorLamp  chan bool
	// Door timer
	doorTimerReset chan bool
	// Network interaction
	OutgoingMsg chan def.Message
}

//TAKE IN NESSECARY CHANNELS
//
func Init(eventCh def.EventChan, hwCh def.HardwareChan, msgCh def.MessageChan, startFloor int) {
	Elevator.behaviour = idle
	Elevator.dir = def.DirStop
	Elevator.floor = startFloor

	go doorTimer(hwCh.DoorTimerReset, eventCh.DoorTimeout)
}

func OnNewRequest(OutgoingMsg chan def.Message, hwCh def.HardwareChan) {
	switch Elevator.behaviour {
	case doorOpen:
		if queue.ShouldStop(Elevator.floor, Elevator.dir) {
			hwCh.DoorTimerReset <- true
			queue.RemoveLocalRequestsAt(Elevator.floor, OutgoingMsg)
		}
	case moving:
		//Do nothing
	case idle:
		Elevator.dir = queue.ChooseDirection(Elevator.floor, Elevator.dir)
		if Elevator.dir == def.DirStop {
			hwCh.DoorLamp <- true
			hwCh.DoorTimerReset <- true
			queue.RemoveLocalRequestsAt(Elevator.floor, OutgoingMsg)
			Elevator.behaviour = doorOpen
		} else {
			hwCh.MotorDir <- Elevator.dir
			Elevator.behaviour = moving
		}
	default: // Error handling
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf(def.ColR, "This state doesn't exist", def.ColN)
	}
	// set all lights
}

func OnFloorArrival(hwCh def.HardwareChan, OutgoingMsg chan def.Message, newFloor int) {
	Elevator.floor = newFloor
	hwCh.FloorLamp <- Elevator.floor

	switch Elevator.behaviour {
	case moving:
		if queue.ShouldStop(newFloor, Elevator.dir) {
			hwCh.MotorDir <- def.DirStop
			hwCh.DoorLamp <- true
			hwCh.DoorTimerReset <- true
			queue.RemoveLocalRequestsAt(Elevator.floor, OutgoingMsg)
			//setAllLights(Elevator);
			Elevator.behaviour = doorOpen
		}
	case doorOpen:
		// do nothing
	case idle:
		// Don´t care
	default: // Error handling
	}
}

func OnDoorTimeout(hwCh def.HardwareChan) {
	switch Elevator.behaviour {
	case doorOpen:
		Elevator.dir = queue.ChooseDirection(Elevator.floor, Elevator.dir)
		hwCh.DoorLamp <- false
		hwCh.MotorDir <- Elevator.dir
		if Elevator.dir == def.DirStop {
			Elevator.behaviour = idle
		} else {
			Elevator.behaviour = moving
		}
	case moving:
		// Don´t care
	case idle:
		// Don´t care
	default: // Error handling
	}
}
