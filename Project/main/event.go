package main

import (
	def "definitions"
	hw "hardware"
	"fsm"
	"network"
	"time"
)


func EventHandler(){
	//Check for all events in loop
	//Make convinient variables
	//Fix lights
	
	//Convenient channels
	var BtnChan = make(chan def.BtnPress, 10)
	var IncomingMessageChan =make(chan def.Message, 10)
	var FloorChan = make(chan int)
	var DeadElevatorChan = make(chan int)
	//var incomingUdpMsgChan =(make chan network.udpMessage)

	//Convenient variables/structs

	//Threads
	go eventBtnPressed(BtnChan)
	go eventIncommingMessage(IncommingMessageChan)
	go eventCabAtFloor(FloorChan)
	//go network.forwardIncoming(IncomingMessageChan, incomingUdpMsgChan)

	for{
		select{
			case btnPress:= <- BtnChan:
				//Do something :P
				//Check if there is an order here already
				//
			case currFloor := <- FloorChan:
				//Handle floor
			case deadElevator := <- DeadElevatorChan:
				//Handle dead elevator
				//Check the whole queue for the dead lifts requests
				//Send them out as new requests
				//how to stop multiple elevators doing this?
			case incomingMsg := <- IncomingMessageChan:
				//Handle message
				//AlivePing
				if incomingMsg.Category == 1{
					//Check IP
					//update alive ping timer

				}
				//NewRequest
				if incomingMsg.Category == 2{

				}
				//CompleteRequest
				if incomingMsg.Category == 3{

				}
				//Cost
				if incomingMsg.Category == 4 {
					
				}
		}
		time.Millisecond(10)
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
		time.Millisecond(1)
	}
}

func eventCabAtFloor(ch chan int){
	//initialize with invalid values
	var floorReached = -2
	var prevFloor = -3
	for{
		if hw.GetFloor != -1{
			if prevFloor != floorReached{
				floorReached = hw.GetFloor
				ch <-floorReached
			}
		}
		time.Millisecond(1)
	}
}
//This is handled by network.forwardIncoming() really
func eventIncommingMessage(ch chan def.Message){
	network.forwardIncoming()
}

func eventExternRequestTimeout(ch chan ...){

}

func eventDeadElevator(ch chan int){
	//Check elevator array for dead elevators
	//every 5 seconds
}