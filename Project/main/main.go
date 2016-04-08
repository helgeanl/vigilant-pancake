package main

import (
	hw "hardware"
	def "definitions"
	"fsm"
)


func main() {
	
	//Variables
	var startFloor int
	var err error

	//Structs
	ch := fsm.Channels{
		NewRequest:		make(chan bool),
		FloorReached:	make(chan int),
		MotorDir: 		make(chan int),
		FloorLamp:		make(chan int),
		DoorLamp:		make(chan bool),
		OutgoingMag:	make(chan def.Message, 10),
	}

	startFloor, err := hw.Init()
	if err!= nil{
		def.Restart(err)
	}

	fsm.Init(ch, startFloor)

	for{

	}
}
