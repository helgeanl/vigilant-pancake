package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
)

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

var laddr *net.UDPAddr //Local address
var baddr *net.UDPAddr //Broadcast address

func main() {
    /* Lets prepare a address at any address at port 10001*/
    //ServerAddr,err := net.ResolveUDPAddr("udp",":30303")
    //CheckError(err)

    var err error
    //Generating broadcast address
  	baddr, err = net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(30302))
    //var ip = ServerAddr.IP


    //Generating localaddress
  	tempConn, err := net.DialUDP("udp4", nil, baddr)
  	defer tempConn.Close()
  	tempAddr := tempConn.LocalAddr()
  	laddr, err = net.ResolveUDPAddr("udp4", tempAddr.String())
  	laddr.Port = 30000
    fmt.Printf("Server IP: %s\n", laddr.IP )

    /* Now listen at selected port */
    //ServerConn, err := net.ListenUDP("udp4", ServerAddr)
    ServerConn, err := net.ListenUDP("udp4", baddr)
    CheckError(err)
    defer ServerConn.Close()

    buf := make([]byte, 1024)

    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        fmt.Println("Received ",string(buf[0:n]), " from ",addr)

        if err != nil {
            fmt.Println("Error: ",err)
        }
    }
}
