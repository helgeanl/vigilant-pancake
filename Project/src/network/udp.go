package network

import (
	def "definitions"
	"log"
	"net"
	"strconv"
)

var udpConnection bool
var baddr *net.UDPAddr

type udpMessage struct {
	raddr  string
	data   []byte
	length int
}

func udpInit(localListenPort, broadcastListenPort, messageSize int, sendCh, receiveCh chan udpMessage) (err error) {
	//Generating broadcast address
	baddr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(broadcastListenPort))
	if err != nil {
		return err
	}

	//Generating localaddress
	tempConn, err := net.DialUDP("udp4", nil, baddr)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	laddr, err := net.ResolveUDPAddr("udp4", tempAddr.String())
	laddr.Port = localListenPort
	def.LocalIP = laddr.String()

	//Creating local listening connections
	localListenConn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return err
	}

	//Creating listener on broadcast connection
	broadcastListenConn, err := net.ListenUDP("udp", baddr)
	if err != nil {
		localListenConn.Close()
		return err
	}
	udpConnection = true
	go udpReceiveServer(localListenConn, broadcastListenConn, messageSize, receiveCh)
	go udpTransmitServer(localListenConn, broadcastListenConn, sendCh)

	return err
}

func udpTransmitServer(lconn, bconn *net.UDPConn, sendCh <-chan udpMessage) {
	for {
		msg := <-sendCh
		if msg.raddr == "broadcast" {
			lconn.WriteToUDP(msg.data, baddr)
		} else {
			raddr, _ := net.ResolveUDPAddr("udp", msg.raddr)
			lconn.WriteToUDP(msg.data, raddr)
		}
	}
}

func udpReceiveServer(lconn, bconn *net.UDPConn, messageSize int, receiveCh chan<- udpMessage) {
	bconnRcvCh := make(chan udpMessage)
	lconnRcvCh := make(chan udpMessage)

	go udpConnectionReader(lconn, messageSize, lconnRcvCh)
	go udpConnectionReader(bconn, messageSize, bconnRcvCh)

	for {
		select {
		case buf := <-bconnRcvCh:
			receiveCh <- buf
		case buf := <-lconnRcvCh:
			receiveCh <- buf
		}
	}
}

func udpConnectionReader(conn *net.UDPConn, messageSize int, rcvCh chan<- udpMessage) {
	for {
		buf := make([]byte, messageSize)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil || n < 0 {
			log.Println(def.ColR, "Trying to reconennect", def.ColN)
		}
		rcvCh <- udpMessage{raddr: raddr.String(), data: buf, length: n}
	}
}
