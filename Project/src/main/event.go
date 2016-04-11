package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	"network"
	"time"
)

func EventHandler(eventCh def.EventChan, msgCh def.MessageChan, hwCh def.HardwareChan) {
	//Check for all events in loop
	//Make convinient variables
	//Fix lights

	onlineElevator := make(map[string]time.Timer)

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
		case deadElevator := <-eventCh.DeadElevatorChan:
			//Handle dead elevator
			//Check the whole queue for the dead lifts requests
			//Send them out as new requests
			//how to stop multiple elevators doing this?
		case incomingMsg := <-msgCh.Incoming:
			//AlivePing
			switch incomingMsg.Category {
			case def.Alive:
				IP := incomingMsg.Addr
				if t, ok:=onlineElevator[IP];ok{
					t.Reset()
				}
				else{
					onlineElevator[IP] = time.AfterFunc(d, f)
				}

			case def.NewRequest:

			case def.CompleteRequest:

			case def.Cost:
				//Send message to Assigner
			default:
				//Burde ikke skje...
			}
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

func eventDeadElevator(ch chan int, m map[string]time.Timer) {
	//Check elevator array for dead elevators
	//every 5 seconds
	for {

		time.Second(5)
	}
}
