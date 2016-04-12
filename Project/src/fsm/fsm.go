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

var Elevator def.Elevator

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

	go doorTimer(hwCh.doorTimerReset, eventCh.DoorTimeout)
}

func onNewRequest(OutgoingMsg chan def.Message, hwCh def.HardwareChan) {
	floor = Elevator.floor
	dir = Elevator.dir
	switch Elevator.behaviour {
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			hwCh.doorTimerReset <- true
			queue.RemoveLocalRequestsAt(floor, OutgoingMsg)
		}
	case moving:
		//Do nothing
	case idle:
		Elevator.dir = queue.ChooseDirction(floor, dir)
		if Elevator.dir == def.DirStop {
			hwCh.DoorLamp <- true
			hwCh.doorTimerReset <- true
			queue.RemoveLocalRequestsAt(floor, OutgoingMsg)
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

func onFloorArrival(hwCh def.HardwareChan, newFloor int) {
	Elevator.floor = newFloor
	ch.FloorLamp <- Elevator.floor

	switch Elevator.behaviour {
	case moving:
		if queue.ShouldStop(Newfloor, Elevator.dir) {
			ch.MotorDir <- def.DirStop
			ch.DoorLamp <- true
			hwCh.doorTimerReset <- true
			queue.RemoveLocalRequestsAt(floor, OutgoingMsg)
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

func onDoorTimeout(hwCh def.HardwareChan) {
	switch state {
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
