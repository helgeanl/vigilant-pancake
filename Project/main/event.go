package main

import (
	def "definitions"
	hw "hardware"
)


func EventHandler(){
	//Check for all events in loop
	//Make convinient variables
	//Fix lights
	
	var BtnChan = make(chan def.BtnPress, 10)

	for{
	}
}

func eventBtnPressed(){
	//Check for a button beeing pressed
	lastBtnPresse := def.BtnPress{
				Button: btn,
				Floor: floor,
	} 
	var BtnPressed:=-2
	
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if hw.ReadBtn(floor, btn){

			}
		}
	}
}