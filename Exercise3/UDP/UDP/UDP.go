// Exercise 3 - UDP
package main

import (
  "fmt"
  "net"
  "strconv"
  //"time"
)
/*
func receiveTest(){
  ln, err := net.Listen("tcp", ":30000")
  if err != nil {
	   // handle error
  }
  for {
	   conn, err := ln.Accept()
	    if err != nil {
		      // handle error
	    }
	    go handleConnection(conn)
  }
}
*/
var laddr *net.UDPAddr //Local address
var baddr *net.UDPAddr //Broadcast address



func udp_init( )(err error){
  //Generating broadcast address
	baddr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(30000))
	if err != nil {
		return err
	}

	//Generating localaddress
	tempConn, err := net.DialUDP("udp4", nil, baddr)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	laddr, err = net.ResolveUDPAddr("udp4", tempAddr.String())
	laddr.Port = 30000 //localListenPort

	//Creating local listening connections
	//localListenConn, err := net.ListenUDP("udp4", laddr)
	//if err != nil {
	//	return err
//	}

	//Creating listener on broadcast connection
	//broadcastListenConn, err := net.ListenUDP("udp", baddr)
	//if err != nil {
	//	localListenConn.Close()
	//	return err
	//}

	//go udp_receive_server(localListenConn, broadcastListenConn, 1024, receive_ch)
	//go udp_transmit_server(localListenConn, broadcastListenConn, send_ch)

	fmt.Printf("Generating local address: \t Network(): %s \t String(): %s \n", laddr.Network(), laddr.String())
	fmt.Printf("Generating broadcast address: \t Network(): %s \t String(): %s \n", baddr.Network(), baddr.String())
  return err
}

func main(){
  //send_ch := make (chan udp.Udp_message)
  //receive_ch := make (chan udp.Udp_message)
  udp_init()
}



/* Sender
// broadcastIP = #.#.#.255. First three bytes are from the local IP
addr = new InternetAddress(broadcastIP, port)
sendSock = new Socket(udp) // UDP, aka SOCK_DGRAM
sendSock.setOption(broadcast, true)
sendSock.sendTo(message, addr)
*/
/*
var laddr *net.UDPAddr //Local address

func receive(){
  var buffer;
  var fromWho;
  recvSock = new Socket(udp)
  for 1 > 0{
    buffer = 0;
    // fromWho will be modified by ref here. Or it's a return value. Depends.
    recvSock.receiveFrom(buffer, ref fromWho)
    if(fromWho.IP != localIP){      // check we are not receiving from ourselves
        // do stuff with buffer
    }
  }
}
*/
/*
//Receiver
byte[1024]          buffer
InternetAddress     fromWho
recvSock = new Socket(udp)
recvSock.bind(addr)         // same addr as sender
loop {
    buffer.clear

    // fromWho will be modified by ref here. Or it's a return value. Depends.
    recvSock.receiveFrom(buffer, ref fromWho)
    if(fromWho.IP != localIP){      // check we are not receiving from ourselves
        // do stuff with buffer
    }
}
*/
