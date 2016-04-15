package main

import (
	"assigner"
	def "definitions"
	"fsm"
	hw "hardware"
	"log"
	"network"
	"os"
	"os/signal"
	q "queue"
)

func main() {
	//Structs
	eventCh := def.EventChan{
		FloorReached:   make(chan int),
		DoorTimeout:    make(chan bool),
		DeadElevator:   make(chan int, 10),
		RequestTimeout: make(chan def.BtnPress, 10),
	}
	hwCh := def.HardwareChan{
		MotorDir:       make(chan int, 10),
		FloorLamp:      make(chan int, 10),
		DoorLamp:       make(chan bool, 10),
		BtnPressed:     make(chan def.BtnPress, 10),
		BtnLightChan:   make(chan def.LightUpdate, 10),
		DoorTimerReset: make(chan bool, 10),
	}
	msgCh := def.MessageChan{
		Outgoing: make(chan def.Message, 10),
		Incoming: make(chan def.Message, 10),
	}

	//initialization
	startFloor := hw.Init()

	fsm.Init(eventCh, hwCh, msgCh, startFloor)
	network.Init(msgCh.Outgoing, msgCh.Incoming)
	q.RunBackup(msgCh.Outgoing)

	//Threads
	go EventHandler(eventCh, msgCh, hwCh)
	go assigner.CollectCosts(q.CostReply, assigner.NumOnlineCh)

	go safeKill()
	hold := make(chan bool)
	<-hold
}

func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hw.SetMotorDir(def.DirStop)
	log.Fatal(def.Col0, "User terminated program.", def.ColN)
}
