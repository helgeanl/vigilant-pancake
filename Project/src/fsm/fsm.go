package fsm

import (
	def "definitions"
	"log"
	"queue"
)

const (
	idle int = iota
	moving
	doorOpen
)

var Elevator struct {
	Floor     int
	Dir       int
	Behaviour int
}

func Init(eventCh def.EventChan, hwCh def.HardwareChan, startFloor int) {
	Elevator.Behaviour = idle
	Elevator.Dir = def.DirStop
	Elevator.Floor = startFloor

	go doorTimer(hwCh.DoorTimerReset, eventCh.DoorTimeout)
	log.Println(def.ColG, "FSM initialized.", def.ColN)
}

func OnNewRequest(OutgoingMsg chan<- def.Message, hwCh def.HardwareChan) {
	switch Elevator.Behaviour {
	case doorOpen:
		if queue.ShouldStop(Elevator.Floor, Elevator.Dir) {
			hwCh.DoorTimerReset <- true
			queue.RemoveLocalRequestsAt(Elevator.Floor, OutgoingMsg)
		}
	case moving:
		//Do nothing
	case idle:
		Elevator.Dir = queue.ChooseDirection(Elevator.Floor, Elevator.Dir)
		if Elevator.Dir == def.DirStop {
			hwCh.DoorLamp <- true
			hwCh.DoorTimerReset <- true
			queue.RemoveLocalRequestsAt(Elevator.Floor, OutgoingMsg)
			Elevator.Behaviour = doorOpen
		} else {
			hwCh.MotorDir <- Elevator.Dir
			Elevator.Behaviour = moving
		}
	}
}

func OnFloorArrival(hwCh def.HardwareChan, OutgoingMsg chan<- def.Message, newFloor int) {
	Elevator.Floor = newFloor
	hwCh.FloorLamp <- Elevator.Floor
	switch Elevator.Behaviour {
	case moving:
		if queue.ShouldStop(newFloor, Elevator.Dir) {
			hwCh.MotorDir <- def.DirStop
			hwCh.DoorLamp <- true
			hwCh.DoorTimerReset <- true
			queue.RemoveLocalRequestsAt(Elevator.Floor, OutgoingMsg)
			Elevator.Behaviour = doorOpen
		}
	}
}

func OnDoorTimeout(hwCh def.HardwareChan) {
	switch Elevator.Behaviour {
	case doorOpen:
		Elevator.Dir = queue.ChooseDirection(Elevator.Floor, Elevator.Dir)
		hwCh.DoorLamp <- false
		hwCh.MotorDir <- Elevator.Dir
		if Elevator.Dir == def.DirStop {
			Elevator.Behaviour = idle
		} else {
			Elevator.Behaviour = moving
		}
	}
}
