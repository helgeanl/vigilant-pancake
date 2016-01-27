package main

import (
    "fmt"
    "net"
    //"time"
    "strconv"
)


var laddr *net.UDPAddr //Local address
var baddr *net.UDPAddr //Broadcast address

func udp_init()(err error) {
  //Generating broadcast address

  baddr,err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(30000))
  //broadcastListenConn, err := net.ListenUDP("udp", baddr)

  fmt.Printf("Generating broadcast address: \t Network(): %s \t String(): %s \n", baddr.Network(), baddr.String())
  return err
}

func main(){
  udp_init()
}
