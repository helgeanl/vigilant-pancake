package network

import (
	def "definitions"
	"encoding/json"
	"log"
	"time"
)

func Init(outgoingMsg chan def.Message, incomingMsg chan def.Message ) {
	// Ports randomly chosen to reduce likelihood of port collision.
	const localListenPort = 37115
	const broadcastListenPort = 37116

	const messageSize = 1024

	var udpSend = make(chan udpMessage)
	var udpReceive = make(chan udpMessage, 10)
	err := udpInit(localListenPort, broadcastListenPort, messageSize, udpSend, udpReceive)
	if err != nil {
		log.Println(def.ColR, "UdpInit() error: ", err, def.ColN)
	}

	go aliveSpammer(outgoingMsg)
	go forwardOutgoing(outgoingMsg, udpSend)
	go forwardIncoming(incomingMsg, udpReceive)

	log.Println(def.ColG, "Network initialized - IP: ", def.LocalIP, def.ColN)
}

// aliveSpammer periodically sends messages on the network to notify all
// lifts that this lift is still online ("alive").
func aliveSpammer(outgoingMsg chan def.Message) {
	alive := def.Message{Category: def.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		outgoingMsg <- alive
		time.Sleep(def.SpamInterval)
	}
}

// forwardOutgoing continuosly checks for messages to be sent on the network
// by reading the OutgoingMsg channel. Each message read is sent to the udp channel
// as JSON.
func forwardOutgoing(outgoingMsg <-chan def.Message, udpSend chan<- udpMessage) {
	for {
		msg := <-outgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println(def.ColR, "json.Marshal error: ", def.ColG, err, def.ColN)
		}
		udpSend <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}

func forwardIncoming(incomingMsg chan<- def.Message, udpReceive <-chan udpMessage) {
	for {
		udpMessage := <-udpReceive
		var message def.Message

		if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
			log.Println(def.ColR, "json.Unmarshal error: ", def.ColG, err, def.ColN)
		}
		message.Addr = udpMessage.raddr
		incomingMsg <- message
	}
}
