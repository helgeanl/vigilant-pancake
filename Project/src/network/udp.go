package network

import (
	def "definitions"
	"log"
	"net"
	"strconv"
)

var udpConnection bool
var baddr *net.UDPAddr //Broadcast address

type udpMessage struct {
	raddr  string //if receiving raddr=senders address, if sending raddr should be set to "broadcast" or an ip:port
	data   []byte
	length int //length of received data, in #bytes // N/A for sending
}

func udpInit(localListenPort, broadcastListenPort, message_size int, send_ch, receive_ch chan udpMessage) (err error) {
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
	go udpReceiveServer(localListenConn, broadcastListenConn, message_size, receive_ch)
	go udpTransmitServer(localListenConn, broadcastListenConn, send_ch)

	return err
}

func udpTransmitServer(lconn, bconn *net.UDPConn, send_ch <-chan udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(def.ColR, "ERROR in udp_transmit_server: ", r, " Closing connection.", def.ColN)
		}
	}()

	var err error
	var n int

	for {
		msg := <-send_ch
		if msg.raddr == "broadcast" {
			n, err = lconn.WriteToUDP(msg.data, baddr)
		} else {
			raddr, err := net.ResolveUDPAddr("udp", msg.raddr)
			if err != nil {
			}
			n, err = lconn.WriteToUDP(msg.data, raddr)
		}
		if err != nil || n < 0 {
			log.Println(def.ColR, "Error: udp_transmit_server: writing", def.ColN)
		}
	}
}

func udpReceiveServer(lconn, bconn *net.UDPConn, messageSize int, receiveCh chan<- udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(def.ColR, "ERROR in udp_receive_server: ", r, "Closing connection.", def.ColN)
			log.Println(def.ColR, "Trying to reconennect", def.ColN)
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			log.Println(def.ColR, "ERROR in udp_connection_reader: ", r, "Closing connection.", def.ColN)
			log.Println(def.ColR, "Trying to reconennect", def.ColN)
		}
	}()

	for {
		buf := make([]byte, messageSize)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil || n < 0 {
			log.Println(def.ColR, "Error: udp_connection_readerDialUDP: reading", def.ColN)
			log.Println(def.ColR, "Trying to reconennect", def.ColN)
		}
		rcvCh <- udpMessage{raddr: raddr.String(), data: buf, length: n}
	}
}
