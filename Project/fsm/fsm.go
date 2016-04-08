// FMS for elevator
// *** some comment
// events: timeout, floorArrived, newRequest
// state: idle, moving, doorOpen
package fsm

import (
	def "definitions"
	"log"
	"queue"
)

// Enumeration of elevator behaviour
const (
	idle int = iota
	moving
	doorOpen
)

var elevator Elevator

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


func Init(ch Channels, startFloor int) {
	elevator.behaviour = idle
	elevator.dir = def.DirStop
	elevator.floor = startFloor

	ch.doorTimeout = make(chan bool)
	ch.doorTimerReset = make(chan bool)

	go doorTimer(ch.doorTimeout, ch.doorTimerReset)
	go monitorEvents(ch)
}

func monitorEvents(ch Channels) {
	for {
		select {
		case <-ch.NewRequest:
			onNewRequest(ch)
		case floor := <-ch.FloorReached:
			onFloorReached(ch, floor)
		case <-ch.doorTimeout:
			onDoorTimeout(ch)
		}
	}
}


func onNewRequest(ch Channels) {
	// print queue
	switch elevator.behaviour {
	case doorOpen:
		//if at ordered floor, start timer again
		// else add order to queue
		if queue.ShouldStop(floor,dir){
			ch.doorTimerReset <- true
			queue.RemoveOrder(floor, ch.OutgoingMsg)
		}
		// else: add order if not done before
	case moving:
		// add request to queue if not done elsewhere
	case idle:
		// add request to queue, if not done before
		// if request at current floor ,
		//		open door,start timer, state = doorOpen
		// else start moving towards requested floor
		// 		state = moving
		elevator.dir = queue.ChooseDirction(floor,dir)
		if elevator.dir = def.DirStop {
			ch.DoorLamp <- true
			ch.doorTimerStart
			queue.RemoveOrder(....)
			elevator.behaviour = doorOpen
		}else{
			ch.MotorDir <- elevator.dir
			elevator.behaviour = moving
		}
	default: // Error handling
		def.CloseConnectionChan <- true
		def.Restart(... some error ...)
		//log.Fatalf(def.ColR, "This state doesn't exist", def.ColN)
	}
	// set all lights
}

func onFloorArrival(ch Channels, newFloor int) {
	elevator.floor = newFloor
	ch.FloorLamp <- elevator.floor

	switch elevator.behaviour {
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
            elevator = requests_clearAtCurrentFloor(elevator);
            timer_start(elevator.config.doorOpenDuration_s);
            setAllLights(elevator);
            elevator.behaviour = doorOpen;
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

		elevator.dir = queue.ChooseDirection(floor,dir);
        ch.DoorLamp <- false
        ch.MotorDir <- elevator.dir
        if elevator.dir == def.DirStop {
            elevator.behaviour = idle;
        } else {
            elevator.behaviour = moving;
        }
	case moving:
		// Don´t care
	case idle:
		// Don´t care
	default: // Error handling
	}
}
