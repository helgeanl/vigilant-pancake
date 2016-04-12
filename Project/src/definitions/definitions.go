package definitions

import (
	"log"
	"os/exec"
	"time"
)

// Global system constants
const NumButtons = 3
const NumFloors = 4
const DoorOpenTime = 3 * time.Second
const SpamInterval = 400 * time.Millisecond
const CostReplyTimeoutDuration = 10 * time.Second
const RequestTimeoutDuration = 30 * time.Second
const ElevTimeoutDuration = 3 * time.Second

const (
	BtnHallUp int = iota
	BtnHallDown
	BtnCab
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

type BtnPress struct {
	Button int
	Floor  int
}

// Message serves as a ...
type Message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string `json:"-"`
}

type Elevator struct {
	floor     int
	dir      int
	behaviour int
}

type LightUpdate struct{
	Floor int
	Button int
	UpdateTo bool
}

type MessageChan struct {
	// Network interaction
	Outgoing chan Message
	Incoming chan Message
}
type HardwareChan struct {
	// Hardware interaction
	MotorDir   		chan int
	FloorLamp	  	chan int
	DoorLamp  	 	chan bool
	BtnPressed		chan BtnPress
	BtnLightChan 	chan BtnPress
	// Door timer
	doorTimerReset 	chan bool
}
type EventChan struct {
	NewRequest     chan bool//?????????????????????????????????????????????
	FloorReached   chan int
	DoorTimeout    chan bool
	DeadElevator   chan int
	RequestTimeout chan BtnPress
}
// Network message category constants
const (
	Alive int = iota + 1
	NewRequest
	CompleteRequest
	Cost
)

//var SyncLightsChan = make(chan bool)
var CloseConnectionChan = make(chan bool)
var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "main")
// Restart program
//func Restart(err error) {
//	start := exec.Command("gnome-terminal", "-x", "sh", "-c", "main")
//	start.Run()
//	log.Fatal(err)
//}

// Colors for printing to console
const Col0 = "\x1b[30;1m" // Dark grey
const ColR = "\x1b[31;1m" // Red
const ColG = "\x1b[32;1m" // Green
const ColY = "\x1b[33;1m" // Yellow
const ColB = "\x1b[34;1m" // Blue
const ColM = "\x1b[35;1m" // Magenta
const ColC = "\x1b[36;1m" // Cyan
const ColW = "\x1b[37;1m" // White
const ColN = "\x1b[0m"    // Grey (neutral)
