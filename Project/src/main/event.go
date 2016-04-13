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
			log.Println("Event: Button pressed",def.ColW)
			if !q.HasRequest(btnPress.Floor,btnPress.Button){
				msgCh.Outgoing<-def.Message{Category:def.NewRequest,Floor:btnPress.Floor,Button:btnPress.Button}
			}
		case incomingMsg := <-msgCh.Incoming:
			log.Println("Event: New message incomming",def.ColW)
			handleMessage(incomingMsg, msgCh.Outgoing)
		
		case btnLightUpdate := <-hwCh.BtnLightChan:
			log.Println("Event: Update light",def.ColW)
			hw.SetBtnLamp(btnLightUpdate)

		case requestTimeout := <-q.RequestTimeoutChan:
			log.Println("Event: Request is timeout",def.ColW)
			q.ReassignRequest(requestTimeout.Floor,requestTimeout.Button, msgCh.Outgoing)
		
		case motorDir := <-hwCh.MotorDir:
			log.Println("Event: Set motor direction",def.ColW)
			hw.SetMotorDir(motorDir)
		
		case floorLamp := <-hwCh.FloorLamp:
			log.Println("Event: Set floor lamp",def.ColW)
			hw.SetFloorLamp(floorLamp)
		
		case doorLamp := <-hwCh.DoorLamp:
			log.Println("Event: Set door lamp",def.ColW)
			hw.SetDoorLamp(doorLamp)
		case <-q.NewRequest:
			log.Println("Event: New Request",def.ColW)
			fsm.OnNewRequest(msgCh.Outgoing, hwCh)
		
		case currFloor := <-eventCh.FloorReached:
			log.Println("Event: New floor",def.ColW)
			fsm.OnFloorArrival(hwCh,msgCh.Outgoing, currFloor)
		
		case <-eventCh.DoorTimeout:
			log.Println("Event: Door timeout",def.ColW)
			fsm.OnDoorTimeout(hwCh)
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

func eventBtnPressed(ch chan<- def.BtnPress) {
	lastBtnPressed := def.BtnPress{
		Button: -1,
		Floor:  -1,
	}
	btnPressed := def.BtnPress{
		Button: -2,
		Floor:  -2,
	}
	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {
				if floor == 0 && btn == def.BtnHallDown{
					//invalid
				}else if floor == def.NumFloors-1 && btn == def.BtnHallUp{
					//invalid
				}else if hw.ReadBtn(floor, btn){	
					btnPressed.Floor =floor
					btnPressed.Button =btn
					if lastBtnPressed != btnPressed {
						ch <- btnPressed
					}
					lastBtnPressed = btnPressed
				}
			}
		}
		time.Sleep(100*time.Millisecond)
	}
}

func eventCabAtFloor(ch chan int) {
	//initialize with invalid values
	var FloorReached = -2
	var prevFloor = -3
	for {
		if hw.GetFloor() != -1 {
			if prevFloor != FloorReached {
				FloorReached = hw.GetFloor()
				ch <- FloorReached
			}
		}
		time.Sleep(100*time.Millisecond)
	}
}

func handleMessage(incomingMsg def.Message, outgoingMsg chan def.Message){
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
			}
		case def.NewRequest:
			cost := q.CalcCost(fsm.Elevator.Dir, hw.GetFloor(),fsm.Elevator.Floor,incomingMsg.Floor, incomingMsg.Button)
			outgoingMsg<-def.Message{Category: def.Cost, Cost: cost}
		case def.CompleteRequest:
			q.RemoveRequest(incomingMsg.Floor, incomingMsg.Button)
		case def.Cost:
			q.CostReply<-incomingMsg 
		default:
			//Do nothing, invalid msg
	}
}