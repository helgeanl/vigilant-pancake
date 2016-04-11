package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	"network"
)

func main() {

	//Variables
	var startFloor int
	var err error

	//Structs
	ch := fsm.Channels{
		NewRequest:   make(chan bool),
		FloorReached: make(chan int),
		MotorDir:     make(chan int),
		FloorLamp:    make(chan int),
		DoorLamp:     make(chan bool),
		OutgoingMsg:  make(chan def.Message, 10),
	}

	startFloor, err := hw.Init()
	if err != nil {
		def.Restart(err)
	}
	//"fsm.Channels is now devided into def.HardwareChannels, def.EventChannels"
	//and def.MessageChannels
	fsm.Init(ch, startFloor)
	network.Init()

	for {

	}
}
