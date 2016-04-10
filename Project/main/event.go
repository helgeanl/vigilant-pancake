package main

import (
	def "definitions"
	hw "hardware"
	"fsm"
	"network"
)


func EventHandler(){
	//Check for all events in loop
	//Make convinient variables
	//Fix lights
	
	//Convenient channels
	var BtnChan = make(chan def.BtnPress, 10)
	var IncomingMessageChan =make(chan def.Message, 10)
	var FloorChan = make(chan int)
	var DeadElevatorChan = make(chan Byte)

	//Convenient variables/structs


	//Threads
	go eventBtnPressed(BtnChan)
	go eventIncommingMessage(IncommingMessageChan)
	go eventCabAtFloor(FloorChan)

	for{
		select{
			case BtnPress:= <- BtnChan:
				//Do something :P
			case Message := <- IncommingMessageChan:
				//Do something
			case currFloor := <- FloorChan
		}
	}
}

func eventBtnPressed(ch chan def.BtnPress){
	//Check for a button beeing pressed
	lastBtnPressed := def.BtnPress{
		Button: -1,
		Floor: -1,
	} 
	btnPressed := def.BtnPress{
		Button: -2,
		Floor: -2,
	}
	for{
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {
				if hw.ReadBtn(floor, btn){
					btnPressed{btn,floor}
					if lastBtnPressed != btnPressed{
						ch <- btnPressed
					}
					lastBtnPressed = btnPressed
				}
			}
		}
	}
}

func eventCabAtFloor(ch chan int){
	
}

func eventIncommingMessage(ch chan def.Message){

}

func eventExternRequestTimeout(ch chan ...){

}

func eventDeadElevator(ch chan ){

}