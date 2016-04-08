package main

import (
	def "config"
	"fsm"
	"hw"
	"liftAssigner"
	"log"
	"network"
	"os"
	"os/signal"
	"queue"
	"time"
)

var onlineLifts = make(map[string]network.UdpConnection)
var numOnline int

var deadChan = make(chan network.UdpConnection)
var costChan = make(chan def.Message)
var outgoingMsg = make(chan def.Message, 10)
var incomingMsg = make(chan def.Message, 10)

func main() {
	var floor int
	var err error
	floor, err = hw.Init()
	if err != nil {
		def.Restart.Run()
		log.Fatal(err)
	}

	ch := fsm.Channels{
		NewOrder:     make(chan bool),
		FloorReached: make(chan int),
		MotorDir:     make(chan int, 10),
		FloorLamp:    make(chan int, 10),
		DoorLamp:     make(chan bool, 10),
		OutgoingMsg:  outgoingMsg,
	}
	fsm.Init(ch, floor)

	network.Init(outgoingMsg, incomingMsg)

	go liftAssigner.Run(costChan, &numOnline)
	go eventHandler(ch)
	go syncLights()

	queue.Init(ch.NewOrder, outgoingMsg)

	go safeKill()

	balboa := make(chan bool)
	<-balboa
}

func eventHandler(ch fsm.Channels) {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case key := <-buttonChan:
			switch key.Button {
			case def.BtnInside:
				queue.AddLocalOrder(key.Floor, key.Button)
			case def.BtnUp, def.BtnDown:
				outgoingMsg <- def.Message{Category: def.NewOrder, Floor: key.Floor, Button: key.Button}
			}
		case floor := <-floorChan:
			ch.FloorReached <- floor
		case message := <-incomingMsg:
			handleMessage(message)

		case connection := <-deadChan:
			handleDeadLift(connection.Addr)
		case order := <-queue.OrderTimeoutChan:
			log.Println(def.ColR, "Order timeout, I will do it myself!", def.ColN)
			queue.RemoveRemoteOrdersAt(order.Floor)
			queue.AddRemoteOrder(order.Floor, order.Button, def.Laddr)

		case dir := <-ch.MotorDir:
			hw.SetMotorDir(dir)
		case floor := <-ch.FloorLamp:
			hw.SetFloorLamp(floor)
		case value := <-ch.DoorLamp:
			hw.SetDoorLamp(value)
		}
	}
}

func pollButtons() <-chan def.Keypress {
	c := make(chan def.Keypress)
	go func() {
		var buttonState [def.NumFloors][def.NumButtons]bool

		for {
			for f := 0; f < def.NumFloors; f++ {
				for b := 0; b < def.NumButtons; b++ {
					if (f == 0 && b == def.BtnDown) ||
						(f == def.NumFloors-1 && b == def.BtnUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- def.Keypress{Button: b, Floor: f}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()
	return c
}

func pollFloors() <-chan int {
	c := make(chan int)
	go func() {
		oldFloor := hw.Floor()

		for {
			newFloor := hw.Floor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()
	return c
}

// handleMessage handles incoming messages from the network.
func handleMessage(msg def.Message) {
	const aliveTimeout = 2 * time.Second

	switch msg.Category {
	case def.Alive:
		if connection, exist := onlineLifts[msg.Addr]; exist {
			connection.Timer.Reset(aliveTimeout)
		} else {
			newConnection := network.UdpConnection{msg.Addr, time.NewTimer(aliveTimeout)}
			onlineLifts[msg.Addr] = newConnection
			numOnline = len(onlineLifts)
			go connectionTimer(&newConnection)
			log.Printf("%sConnection to IP %s established!%s", def.ColG, msg.Addr[0:15], def.ColN)
		}
	case def.NewOrder:
		cost := queue.CalculateCost(msg.Floor, msg.Button, fsm.Floor(), hw.Floor(), fsm.Direction())
		outgoingMsg <- def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: cost}
	case def.CompleteOrder:
		queue.RemoveRemoteOrdersAt(msg.Floor)
	case def.Cost:
		costChan <- msg
	}
}

// handleDeadLift removes a dead lift from the list of online lifts, and
// reassigns any orderes assigned to it.
func handleDeadLift(deadAddr string) {
	log.Printf("%sConnection to IP %s is dead!%s", def.ColR, deadAddr[0:15], def.ColN)
	delete(onlineLifts, deadAddr)
	numOnline = len(onlineLifts)
	queue.ReassignOrders(deadAddr, outgoingMsg)
}

// connectionTimer finds out when any lifts are lost from the network.
func connectionTimer(connection *network.UdpConnection) {
	<-connection.Timer.C
	deadChan <- *connection
}

// safeKill turns the motor off if the program is killed with CTRL+C.
func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hw.SetMotorDir(def.DirStop)
	log.Fatal(def.ColR, "User terminated program.", def.ColN)
}

// syncLights checks the queues and updates all order lamps accordingly.
func syncLights() {
	for {
		<-def.SyncLightsChan
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if (b == def.BtnUp && f == def.NumFloors-1) || (b == def.BtnDown && f == 0) {
					continue
				} else {
					switch b {
					case def.BtnInside:
						hw.SetButtonLamp(f, b, queue.IsLocalOrder(f, b))
					case def.BtnUp, def.BtnDown:
						hw.SetButtonLamp(f, b, queue.IsRemoteOrder(f, b))
					}
				}
			}
		}
	}
}
