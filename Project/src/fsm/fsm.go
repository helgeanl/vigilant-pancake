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
	NewRequest     chan bool
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
func Init(eventCh def.EventChan,hwCh def.HardwareChan,msgCh def.MessageChan, startFloor int) {
	Elevator.behaviour = idle
	Elevator.dir = def.DirStop
	Elevator.floor = startFloor

	go doorTimer(eventCh.DoorTimeout, hwCh.doorTimerReset)
}



func onNewRequest() {
	// print queue
	floor = Elevator.floor
	dir = Elevator.dir


	queue.printQueue()
	switch Elevator.behaviour {
	case doorOpen:
		//if at ordered floor, start timer again
		// else add order to queue
		if queue.ShouldStop(floor,dir){
			ch.doorTimerReset <- true
			queue.RemoveOrder(floor, ch.OutgoingMsg)
			// send complete
		}
	case moving:
		// add request to queue if not done elsewhere
	case idle:
		// add request to queue, if not done before
		// if request at current floor ,
		//		open door,start timer, state = doorOpen
		// else start moving towards requested floor
		// 		state = moving
		Elevator.dir = queue.ChooseDirction(floor,dir)
		if Elevator.dir = def.DirStop {
			ch.DoorLamp <- true
			ch.doorTimerStart
			queue.RemoveOrder(floor, ch.OutgoingMsg)
			Elevator.behaviour = doorOpen
		}else{
			ch.MotorDir <- Elevator.dir
			Elevator.behaviour = moving
		}
	default: // Error handling
		def.CloseConnectionChan <- true
		def.Restart(... some error ...)
		//log.Fatalf(def.ColR, "This state doesn't exist", def.ColN)
	}
	// set all lights
}

func onFloorArrival(ch Channels, newFloor int) {
	Elevator.floor = newFloor
	ch.FloorLamp <- Elevator.floor

	switch Elevator.behaviour {
	case moving:
		// if floor is in queue
		// then stop MOTOR,
		// turn on doorlight and start timer
		// clear request
		// Turn off button lights
		// state = doorOpen

		// semi-pseudokode
		if queue.ShouldStop(floor, dir){
			ch.MotorDir <- def.DirStop
            ch.DoorLamp <- true
            Elevator = requests_clearAtCurrentFloor(Elevator);
            timer_start(Elevator.config.doorOpenDuration_s);
            setAllLights(Elevator);
            Elevator.behaviour = doorOpen;
		}
	case doorOpen:
		// do nothing
	case idle:
		// Don´t care
	default: // Error handling
	}
}

func onDoorTimeout(ch Channels) {
	switch state {
	case doorOpen:
		// Check for new direction
		// if new direction:
		// 		Move towards new request
		// 		state = moving
		// else: state = idle
		// turn off doorLamp

		Elevator.dir = queue.ChooseDirection(floor,dir);
        ch.DoorLamp <- false
        ch.MotorDir <- Elevator.dir
        if Elevator.dir == def.DirStop {
            Elevator.behaviour = idle;
        } else {
            Elevator.behaviour = moving;
        }
	case moving:
		// Don´t care
	case idle:
		// Don´t care
	default: // Error handling
	}
}
