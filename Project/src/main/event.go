package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	q "queue"
	"time"
	"assigner"
	"log"
)

var onlineElevatorMap = make(map[string]*time.Timer)

func EventHandler(eventCh def.EventChan, msgCh def.MessageChan, hwCh def.HardwareChan) {
	//Check for all events in loop
	//Make convinient variables



	//Threads
	go eventBtnPressed(hwCh.BtnPressed)
	go eventCabAtFloor(eventCh.FloorReached)

	for {
		select {
		case btnPress := <-hwCh.BtnPressed:
			log.Println(def.ColW,"Event: Button pressed",def.ColN)
			if !q.HasRequest(btnPress.Floor,btnPress.Button){
				if btnPress.Button == def.BtnCab{
					q.AddRequest(btnPress.Floor,btnPress.Button,def.LocalIP)
				}else{
					msgCh.Outgoing<-def.Message{Category:def.NewRequest,Floor:btnPress.Floor,Button:btnPress.Button, Cost:0}
					log.Println(def.ColM,"Sent new request on the network",def.ColN)
				}
			}
		case incomingMsg := <-msgCh.Incoming:
			go handleMessage(incomingMsg, msgCh.Outgoing)
		
		case btnLightUpdate := <- q.LightUpdate://<-hwCh.BtnLightChan:
			log.Println(def.ColB,"Event: Update light",def.ColN)
			hw.SetBtnLamp(btnLightUpdate)

		case requestTimeout := <-q.RequestTimeoutChan:
			log.Println(def.ColW,"Event: Request is timeout",def.ColN)
			q.ReassignRequest(requestTimeout.Floor,requestTimeout.Button, msgCh.Outgoing)
		
		case motorDir := <-hwCh.MotorDir:
			log.Println(def.ColW,"Event: Set motor direction",def.ColN)
			hw.SetMotorDir(motorDir)
		
		case floorLamp := <-hwCh.FloorLamp:
			log.Println(def.ColW,"Event: Set floor lamp",def.ColN)
			hw.SetFloorLamp(floorLamp)
		
		case doorLamp := <-hwCh.DoorLamp:
			log.Println(def.ColW,"Event: Set door lamp to: ",doorLamp,def.ColN)
			hw.SetDoorLamp(doorLamp)

		case <-q.NewRequest:
			log.Println(def.ColW,"Event: New Request",def.ColN)
			fsm.OnNewRequest(msgCh.Outgoing, hwCh)
		
		case currFloor := <-eventCh.FloorReached:
			log.Println(def.ColW,"Event: New floor",def.ColN)
			fsm.OnFloorArrival(hwCh,msgCh.Outgoing, currFloor)
		case <-eventCh.DoorTimeout:
			log.Println(def.ColW,"Event: Door timeout",def.ColN)
			fsm.OnDoorTimeout(hwCh)
		}
		time.Sleep(10*time.Millisecond)
	}
}

func eventBtnPressed(ch chan<- def.BtnPress) {
	lastBtnPressed := def.BtnPress{
		Floor:  -1,
		Button: -1,
	}
	btnPressed := def.BtnPress{
		Floor:  -2,
		Button: -2,
	}
	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {
				if floor == 0 && btn == def.BtnHallDown{
					//invalid
				}else if floor == def.NumFloors-1 && btn == def.BtnHallUp{
					//invalid
				}else if hw.ReadBtn(floor, btn){	
					btnPressed.Floor = floor
					btnPressed.Button = btn
					if lastBtnPressed != btnPressed {
						ch <- btnPressed
						log.Println(def.ColR,"Button pressed: ",btnPressed,def.ColN)
					}
					lastBtnPressed = btnPressed
				}
			}
		}
		time.Sleep(100*time.Millisecond)
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
				log.Println(def.ColB,"New floor",def.ColN)
				prevFloor=FloorReached
			}
		}
		time.Sleep(100*time.Millisecond)
	}
}

func handleMessage(incomingMsg def.Message, outgoingMsg chan<- def.Message){
	switch incomingMsg.Category {
		case def.Alive:
			IP := incomingMsg.Addr
			if t, exists := onlineElevatorMap[IP]; exists {
				t.Reset(def.ElevTimeoutDuration)
			} else {
				f := func(){
					q.ReassignAllRequestsFrom(IP, outgoingMsg)
					delete(onlineElevatorMap,IP)
					assigner.NumOnlineCh <- len(onlineElevatorMap)
				}
				onlineElevatorMap[IP] = time.AfterFunc(def.ElevTimeoutDuration, f)
				assigner.NumOnlineCh <- len(onlineElevatorMap)
				log.Println(def.ColG,"New elevator: ",IP," | Number online: ",len(onlineElevatorMap),def.ColN)
			}
		case def.NewRequest:
			log.Println(def.ColC,"New request incomming",def.ColN)
			cost := q.CalcCost(fsm.Elevator.Dir, hw.GetFloor(),fsm.Elevator.Floor,incomingMsg.Floor, incomingMsg.Button)
			outgoingMsg<-def.Message{Category: def.Cost,Floor:incomingMsg.Floor,Button:incomingMsg.Button, Cost: cost}
		case def.CompleteRequest:
			log.Println(def.ColG,"Request is completed",def.ColN)
			q.RemoveRequest(incomingMsg.Floor, incomingMsg.Button)
		case def.Cost:
			log.Println(def.ColC,"Cost reply",def.ColN)
			q.CostReply<-incomingMsg 
		default:
			//Do nothing, invalid msg
	}
}