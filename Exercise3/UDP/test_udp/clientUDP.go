package main

import (
    "fmt"
    "net"
    "time"
    "strconv"
)

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
var laddr *net.UDPAddr //Local address
var laddr2 *net.UDPAddr //Local address
var baddr *net.UDPAddr //Broadcast address
var baddr2 *net.UDPAddr //Broadcast address

func main() {
    //ServerAddr,err := net.ResolveUDPAddr("udp","10.20.76.1:30303")
    //CheckError(err)
    var err error
    var err2 error
    //Generating broadcast address
  baddr, err = net.ResolveUDPAddr("udp4", "129.241.187.255:"+strconv.Itoa(30000))
  baddr2, err2 = net.ResolveUDPAddr("udp4", "129.241.187.255:"+strconv.Itoa(20007))
  
    //Generating localaddress
  	tempConn, err := net.DialUDP("udp4", nil, baddr)
  	defer tempConn.Close()
  	tempAddr := tempConn.LocalAddr()
  	laddr, err = net.ResolveUDPAddr("udp4", tempAddr.String())
  	laddr.Port = 30007
    fmt.Printf("Client IP: %s\n", laddr.IP )
	CheckError(err)
	
	
	
	//Generating localaddress
  	tempConn2, err2 := net.DialUDP("udp4", nil, baddr2)
  	defer tempConn2.Close()
  	tempAddr2 := tempConn2.LocalAddr()
  	laddr2, err2 = net.ResolveUDPAddr("udp4", tempAddr2.String())
  	laddr2.Port = 20007
    fmt.Printf("Client IP: %s\n", laddr2.IP )
	CheckError(err2)

	
    //LocalAddr, err := net.ResolveUDPAddr("udp", "10.20.76.1:0")
    //CheckError(err)
    //127.0.0.1:0
	go send()
	go receive()
	for{
	
	}
    

}

func send(){
	Conn, err := net.DialUDP("udp", nil, baddr2)
    CheckError(err)
    if err != nil {
		log.Fatal(err)
	}
	
    defer Conn.Close()
	
    i := 0
    for {
        msg := "Hello back again:  " + strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := Conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second * 1)
    }

}

func receive(){
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
	time.Sleep(time.Second * 1)
   }

}

