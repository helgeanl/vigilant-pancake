package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	"network"
	q "queue"
	"time"
)

func EventHandler(eventCh def.EventChan, msgCh def.MessageChan, hwCh def.HardwareChan) {
	//Check for all events in loop
	//Make convinient variables
	//Fix lights

	onlineElevatorMap := make(map[string]time.Timer)

	//Threads
	go eventBtnPressed(hwCh.BtnPressed)
	go eventCabAtFloor(eventCh.FloorReached)

	for {
		select {
		case btnPress := <-hwCh.BtnPressed:
			//Do something :P
			//Check if there is an order here already?
			//
		case currFloor := <-eventCh.FloorReached:
			//Handle floor
			//Kjør fsm.onFloorArrival
		case incomingMsg := <-msgCh.Incoming:
			//AlivePing

			//MAKE THIS A FUNCTION

			switch incomingMsg.Category {
			case def.Alive:
				IP := incomingMsg.Addr
				if t, ok := onlineElevatorMap[IP]; ok {
					t.Reset()
				} else {
					onlineElevatorMap[IP] = time.AfterFunc(def.ElevTimeoutDuration, q.ReassignRequest(IP))
				}
			case def.NewRequest:
				//Kjør fsm.onNewRequest
			case def.CompleteRequest:
				q.RemoveOrderAt(incomingMsg.Floor, incomingMsg.Button)
			case def.Cost:
				//Send message to Assigner
			default:
				//Burde ikke skje...
			}
			//MAKE LIGHT CASE
		}
		time.Millisecond(10)
	}
}

func eventBtnPressed(ch chan def.BtnPress) {
	//Check for a button beeing pressed
	lastBtnPressed := def.BtnPress{
		Button: -1,
		Floor:  -1,
	}
	btnPressed := def.BtnPress{
		Button: -2,
		Floor:  -2,
	}
	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {
				if hw.ReadBtn(floor, btn) {
					btnPressed{btn, floor}
					if lastBtnPressed != btnPressed {
						ch <- btnPressed
					}
					lastBtnPressed = btnPressed
				}
			}
		}
		time.Millisecond(1)
	}
}

func eventCabAtFloor(ch chan int) {
	//initialize with invalid values
	var floorReached = -2
	var prevFloor = -3
	for {
		if hw.GetFloor() != -1 {
			if prevFloor != floorReached {
				floorReached = hw.GetFloor()
				ch <- floorReached
			}
		}
		time.Millisecond(1)
	}
}

func eventRequestTimeout(ch chan BtnPress) {

}
