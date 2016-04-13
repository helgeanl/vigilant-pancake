package main

import (
	def "definitions"
	"fsm"
	hw "hardware"
	q "queue"
	"time"
)

func EventHandler(eventCh def.EventChan, msgCh def.MessageChan, hwCh def.HardwareChan) {
	//Check for all events in loop
	//Make convinient variables

	onlineElevatorMap := make(map[string]time.Timer)

	//Threads
	go eventBtnPressed(hwCh.BtnPressed)
	go eventCabAtFloor(eventCh.FloorReached)

	for {
		select {
		case btnPress := <-hwCh.BtnPressed:
			if !q.HasRequest(btnPress.Floor,btnPress.Button){
				msgCh.Outgoing<-def.Message{Category:def.NewRequest,Floor:btnPress.Floor,Button:btnPress.Button}
			}
		case incomingMsg := <-msgCh.Incoming:
			go handleMessage(incomingMsg, msgCh.Outgoing)
		
		case btnLightUpdate := <-hwCh.BtnLightChan:
			hw.SetBtnLamp(btnLightUpdate)

		case requestTimeout := <-q.RequestTimeoutChan:
			q.ReassignRequest(requestTimeout.Floor,requestTimeout.Button, msgCh.Outgoing)
		
		case motorDir := <-hwCh.MotorDir:
			hw.SetMotorDir(motorDir)
		
		case floorLamp := <-hwCh.FloorLamp:
			hw.SetFloorLamp(floorLamp)
		
		case doorLamp := <-hwCh.DoorLamp:
			hw.SetDoorLamp(doorLamp)
		case <-q.NewRequest:
			fsm.OnNewRequest(msgCh.Outgoing, hwCh)
		
		case currFloor := <-eventCh.FloorReached:
			fsm.OnFloorArrival(hwCh,msgCh.Outgoing, currFloor)
		
		case <-eventCh.DoorTimeout:
			fsm.OnDoorTimeout(hwCh)
		}
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
				if hw.ReadBtn(floor, btn) {
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
	var floorReached = -2
	var prevFloor = -3
	for {
		if hw.GetFloor() != -1 {
			if prevFloor != floorReached {
				floorReached = hw.GetFloor()
				ch <- floorReached
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
				t.Reset()
			} else {
				onlineElevatorMap[IP] = time.AfterFunc(def.ElevTimeoutDuration, q.ReassignAllRequestsFrom(IP, outgoingMsg))
			}
		case def.NewRequest:
			cost := q.CalcCost(fsm.Elevator.dirn, hw.GetFloor(),fsm.Elevator.floor,incomingMsg.Floor, incomingMsg.Button)
			outgoingMsg<-def.Message{Category: def.Cost, Cost: cost}
		case def.CompleteRequest:
			q.RemoveRequest(incomingMsg.Floor, incomingMsg.Button)
		case def.Cost:
			q.costReply<-incomingMsg 
		default:
			//Do nothing, invalid msg
	}
}