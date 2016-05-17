package main

import(
	"log"
	"net"
	"time"
)


func main() {
	
	//set up send-socket
	remote_addr, _ := net.ResolveUDPAddr("udp", "129.241.187.255:20007")
	socket_send, err := net.DialUDP("udp", nil, remote_addr)

	if err != nil {
		log.Fatal(err)
	}
	
	//set up listen-socket
	port, _ := net.ResolveUDPAddr("udp", ":20007")
	socket_listen, err := net.ListenUDP("udp", port)

	if err != nil {
		log.Fatal(err)
	}

	//close sockets when done
	defer socket_listen.Close()
	defer socket_send.Close()

	for {
		//send message
		msg := "Hello terminal 7:"
		socket_send.Write([]byte(msg))

		//listen to message
		var buffer[64] byte
		length, addr, err := socket_listen.ReadFromUDP(buffer[:])
		log.Println(length)
		log.Println(addr)
		log.Println(err)
		log.Println(string(buffer[:]), "\n")
		time.Sleep(1*time.Second)
	}

}
