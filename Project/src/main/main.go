package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	"network"
	q "queue"
	"assigner"
	"log"
	"time"
)

func main() {
	//Structs
	eventCh := def.EventChan{
		NewRequest:     make(chan bool,10),
		FloorReached:   make(chan int),
		DoorTimeout:    make(chan bool),
		DeadElevator:   make(chan int,10),// Really needed??
		RequestTimeout: make(chan def.BtnPress,10),

	}
	hwCh := def.HardwareChan{
		MotorDir:       make(chan int,10),
		FloorLamp:      make(chan int, 10),
		DoorLamp:       make(chan bool,10),
		BtnPressed:     make(chan def.BtnPress,10),
		BtnLightChan:   make(chan def.LightUpdate,10),
		DoorTimerReset: make(chan bool,10),
	}
	msgCh := def.MessageChan{
		Outgoing: make(chan def.Message,10),
		Incoming: make(chan def.Message,10),
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
		time.Sleep(time.Hour)// FIX THIS
	}
}
