package config

import (
	"os/exec"
)

// Global system constants
const NumButtons = 3
const NumFloors = 4

const (
	BtnUp int = iota
	BtnDown
	BtnInside
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

// Local IP address
var Laddr string

type Keypress struct {
	Button int
	Floor  int
}

// Generic network message. No other messages are ever sent on the network.
type Message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string `json:"-"`
}

// Network message category constants
const (
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

var SyncLightsChan = make(chan bool)
var CloseConnectionChan = make(chan bool)

// Start a new terminal when restart.Run()
var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "main")

// Colours for printing to console
const Col0 = "\x1b[30;1m" // Dark grey
const ColR = "\x1b[31;1m" // Red
const ColG = "\x1b[32;1m" // Green
const ColY = "\x1b[33;1m" // Yellow
const ColB = "\x1b[34;1m" // Blue
const ColM = "\x1b[35;1m" // Magenta
const ColC = "\x1b[36;1m" // Cyan
const ColW = "\x1b[37;1m" // White
const ColN = "\x1b[0m"    // Grey (neutral)
