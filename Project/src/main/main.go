package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	"network"
	q "queue"
)

func main() {

	//Variables
	var startFloor int
	var err error

	//Structs
	eventCh := def.EventChan{
		NewRequest:     make(chan bool),
		FloorReached:   make(chan int),
		DoorTimeout:    make(chan bool), //Really needed??
		DeadElevator:   make(chan int),
		RequestTimeout: make(chan BtnPress),
	}
	hwCh := def.HardwareChan{
		MotorDir:       make(chan int),
		FloorLamp:      make(chan int),
		DoorLamp:       make(chan bool),
		BtnPressed:     make(chan BtnPress),
		BtnLightChan: 	make(chan LightUpdate),
		doorTimerReset: make(chan bool),
	}
	msgCh := def.MessageChan{
		Outgoing: make(chan Message),
		Incoming: make(chan Message),
	}

	//initialization
	startFloor, err := hw.Init()
	if err != nil {
		def.Restart(err)
	}
	//"fsm.Channels is now devided into def.HardwareChannels, def.EventChannels"
	//and def.MessageChannels
	fsm.Init(ch, startFloor)
	localIP := network.Init(msgCh.Outgoing, msgCh.Incoming)
	q.Init(localIP,msgCh.Outgoing)
	//Threads
	go EventHandler(eventCh, msgCh, hwCh)

	for { //Or a channel that holds until it gets kill signal

	}
}
