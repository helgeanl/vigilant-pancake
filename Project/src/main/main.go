package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	"network"
	q "queue"
	"assigner"
	"log"
)

func main() {
	//Structs
	eventCh := def.EventChan{
		NewRequest:     make(chan bool),
		FloorReached:   make(chan int),
		DoorTimeout:    make(chan bool), //Really needed??
		DeadElevator:   make(chan int),
		RequestTimeout: make(chan def.BtnPress),
	}
	hwCh := def.HardwareChan{
		MotorDir:       make(chan int),
		FloorLamp:      make(chan int),
		DoorLamp:       make(chan bool),
		BtnPressed:     make(chan def.BtnPress),
		BtnLightChan:   make(chan def.LightUpdate),
		DoorTimerReset: make(chan bool),
	}
	msgCh := def.MessageChan{
		Outgoing: make(chan def.Message),
		Incoming: make(chan def.Message),
	}

	//initialization
	startFloor, err := hw.Init()
	if err != nil {
		def.Restart.Run()
		log.Fatalf(def.ColR, "Error in HW", def.ColN)
	}

	//"fsm.Channels is now devided into def.HardwareChannels, def.EventChannels"
	//and def.MessageChannels
	fsm.Init(eventCh, hwCh, msgCh, startFloor)
	network.Init(msgCh.Outgoing, msgCh.Incoming)
	q.Init(msgCh.Outgoing)
	//Threads
	go EventHandler(eventCh, msgCh, hwCh)
	go assigner.CollectCosts(q.CostReply, assigner.NumOnlineCh)
	for { //Or a channel that holds until it gets kill signal

	}
}
