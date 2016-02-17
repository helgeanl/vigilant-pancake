package UDP

import(
	"log"
	"net"
	"time"
)


func main() {
	addr, err := net.ResolveUDPAddr("udp", ":30000")

	if err != nil {
		log.Fatal(err)
	}

	socket, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatal(err)
	}

	defer socket.Close()

	for {
			var buffer[64] byte
			length, addr, err := socket.ReadFromUDP(buffer[:])
			log.Println(length)
			log.Println(addr)
			log.Println(err)
			log.Println(string(buffer[:]), "\n")
			time.Sleep(1*time.Second)

	}

}
