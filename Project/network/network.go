package network

import (
	"config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func Init(outgoingMsg, incomingMsg chan config.Message) {
	// Ports randomly chosen to reduce likelihood of port collision.
	const localListenPort = 37103
	const broadcastListenPort = 37104

	const messageSize = 1024

	var udpSend = make(chan udpMessage)
	var udpReceive = make(chan udpMessage, 10)
	err := udpInit(localListenPort, broadcastListenPort, messageSize, udpSend, udpReceive)
	if err != nil {
		fmt.Print("UdpInit() error: %v \n", err)
	}

	go aliveSpammer(outgoingMsg)
	go forwardOutgoing(outgoingMsg, udpSend)
	go forwardIncoming(incomingMsg, udpReceive)

	log.Println(config.ColG, "Network initialised.", config.ColN)
}

// aliveSpammer periodically sends messages on the network to notify all
// lifts that this lift is still online ("alive").
func aliveSpammer(outgoingMsg chan<- config.Message) {
	const spamInterval = 400 * time.Millisecond
	alive := config.Message{Category: config.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		outgoingMsg <- alive
		time.Sleep(spamInterval)
	}
}

// forwardOutgoing continuosly checks for messages to be sent on the network
// by reading the OutgoingMsg channel. Each message read is sent to the udp file
// as JSON.
func forwardOutgoing(outgoingMsg <-chan config.Message, udpSend chan<- udpMessage) {
	for {
		msg := <-outgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("%sjson.Marshal error: %v\n%s", config.ColR, err, config.ColN)
		}

		udpSend <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}

func forwardIncoming(incomingMsg chan<- config.Message, udpReceive <-chan udpMessage) {
	for {
		udpMessage := <-udpReceive
		var message config.Message

		if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
			fmt.Printf("json.Unmarshal error: %s\n", err)
		}

		message.Addr = udpMessage.raddr
		incomingMsg <- message
	}
}
