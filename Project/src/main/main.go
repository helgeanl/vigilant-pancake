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
	eventCh := def.EventChan{
		FloorReached: 	make(chan int),
		DoorTimeout:  	make(chan bool),
		DeadElevator: 	make(chan int, 10),
	}
	hwCh := def.HardwareChan{
		MotorDir:       make(chan int, 10),
		FloorLamp:      make(chan int, 10),
		DoorLamp:       make(chan bool, 10),
		BtnPressed:     make(chan def.BtnPress, 10),
		DoorTimerReset: make(chan bool, 10),
	}
	msgCh := def.MessageChan{
		Outgoing: 		make(chan def.Message, 10),
		Incoming: 		make(chan def.Message, 10),
		CostReply: 		make(chan def.Message, 10),
		NumOnline: 		make(chan int),
	}

	//Initialization
	startFloor := hw.Init()
	fsm.Init(eventCh, hwCh, startFloor)
	network.Init(msgCh.Outgoing, msgCh.Incoming)
	q.RunBackup(msgCh.Outgoing)

	go EventHandler(eventCh, msgCh, hwCh)
	go assigner.CollectCosts(msgCh.CostReply, msgCh.NumOnline)
	go safeKill()

	hold := make(chan bool)
	<-hold
}

// Turn off motor if program is terminated by user
func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hw.SetMotorDir(def.DirStop)
	log.Fatal(def.Col0, "User terminated program.", def.ColN)
}
