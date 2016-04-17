package main

import (
	"assigner"
	def "definitions"
	"fsm"
	hw "hardware"
	"log"
	q "queue"
	"time"
)

var onlineElevatorMap = make(map[string]def.UdpConnection)

func EventHandler(eventCh def.EventChan, msgCh def.MessageChan, hwCh def.HardwareChan) {

	go eventBtnPressed(hwCh.BtnPressed)
	go eventCabAtFloor(eventCh.FloorReached)

	for {
		select {
		case btnPress := <-hwCh.BtnPressed:
			handleBtnPress(btnPress, msgCh.Outgoing)
		case incomingMsg := <-msgCh.Incoming:
			go handleMessage(incomingMsg, msgCh)
		case btnLightUpdate := <-q.LightUpdate:
			log.Println(def.ColW, "Light update", def.ColN)
			hw.SetBtnLamp(btnLightUpdate)
		case requestTimeout := <-q.RequestTimeoutChan:
			q.ReassignRequest(requestTimeout.Floor, requestTimeout.Button, msgCh.Outgoing)
		case motorDir := <-hwCh.MotorDir:
			hw.SetMotorDir(motorDir)
		case floorLamp := <-hwCh.FloorLamp:
			hw.SetFloorLamp(floorLamp)
		case doorLamp := <-hwCh.DoorLamp:
			hw.SetDoorLamp(doorLamp)
		case <-q.NewRequest:
			log.Println(def.ColW, "Event: New Request", def.ColN)
			fsm.OnNewRequest(msgCh.Outgoing, hwCh)
		case currFloor := <-eventCh.FloorReached:
			fsm.OnFloorArrival(hwCh, msgCh.Outgoing, currFloor)
		case <-eventCh.DoorTimeout:
			fsm.OnDoorTimeout(hwCh)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func eventBtnPressed(ch chan<- def.BtnPress) {
	lastBtnPressed := def.BtnPress{Floor: -1, Button: -1}
	btnPressed := def.BtnPress{Floor: -2, Button: -2}
	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {
				if floor == 0 && btn == def.BtnHallDown || floor == def.NumFloors-1 && btn == def.BtnHallUp {
					//invalid
				} else if hw.ReadBtn(floor, btn) {
					btnPressed.Floor = floor
					btnPressed.Button = btn
					if lastBtnPressed != btnPressed {
						ch <- btnPressed
					}
					lastBtnPressed = btnPressed
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func eventCabAtFloor(ch chan<- int) {
	//initialize with invalid values
	var FloorReached = -2
	var prevFloor = -3
	for {
		if hw.GetFloor() != -1 {
			FloorReached = hw.GetFloor()
			if prevFloor != FloorReached {
				ch <- FloorReached
				prevFloor = FloorReached
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func handleBtnPress(btnPress def.BtnPress, outgoingMsg chan<- def.Message) {
	if btnPress.Button == def.BtnCab {
		q.AddRequest(btnPress.Floor, btnPress.Button, def.LocalIP)
	} else {
		outgoingMsg <- def.Message{Category: def.NewRequest, Floor: btnPress.Floor, Button: btnPress.Button, Cost: 0}
	}
}

func handleMessage(incomingMsg def.Message, msgCh def.MessageChan) {
	switch incomingMsg.Category {
	case def.Alive:
		addr := incomingMsg.Addr
		if connection, exist := onlineElevatorMap[addr]; exist {
			connection.Timer.Reset(def.ElevTimeoutDuration)
		} else {
			newConnection := def.UdpConnection{Addr: addr, Timer: time.NewTimer(def.ElevTimeoutDuration)}
			onlineElevatorMap[addr] = newConnection
			msgCh.NumOnline <- len(onlineElevatorMap)
			go connectionTimer(&newConnection, msgCh.Outgoing, msgCh.NumOnline)
			log.Println(def.ColG, "New elevator: ", addr, " | Number online: ", len(onlineElevatorMap), def.ColN)
		}
	case def.NewRequest:
		log.Println(def.ColC, "New request incomming", def.ColN)
		cost := q.CalcCost(fsm.Elevator.Dir, hw.GetFloor(), fsm.Elevator.Floor, incomingMsg.Floor, incomingMsg.Button)
		msgCh.Outgoing <- def.Message{Category: def.Cost, Floor: incomingMsg.Floor, Button: incomingMsg.Button, Cost: cost}
	case def.CompleteRequest:
		log.Println(def.ColG, "Request is completed", def.ColN)
		q.RemoveRequest(incomingMsg.Floor, incomingMsg.Button)
	case def.Cost:
		log.Println(def.ColC, "Cost reply", def.ColN)
		msgCh.CostReply <- incomingMsg
	}
}

func handleDeadElevator(con def.UdpConnection, outgoingMsg chan<- def.Message,numOnlineCh chan<- int) {
	log.Println(def.ColR, "Connection to ", def.ColG, con.Addr, def.ColR, " is lost| Number online: ", len(onlineElevatorMap), def.ColN)
	delete(onlineElevatorMap, con.Addr)
	numOnlineCh <- len(onlineElevatorMap)
	q.ReassignAllRequestsFrom(con.Addr, outgoingMsg)
}

func connectionTimer(connection *def.UdpConnection, outgoingMsg chan<- def.Message,numOnlineCh chan<- int) {
	<-connection.Timer.C
	if (*connection).Addr != def.LocalIP {
		handleDeadElevator(*connection, outgoingMsg,numOnlineCh)
	}
}
