// Package hw defines interactions with the lift hardware at the real time
// lab at The Department of Engineering Cybernetics at NTNU, Trondheim,
// Norway.
//
// This file is a golang port of elev.c from the hand out driver
// (https://github.com/TTK4145/Project)
package hardware

import (
	"config"
	"errors"
	"log"
)

var lampChannelMatrix = [config.NumFloors][config.NumButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}
var btnChannelMatrix = [config.NumFloors][config.NumButtons]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

// Init initialises the lift hardware and moves the lift to a defined state.
// (Descending until it reaches a floor.)
func Init() (int, error) {
	// Init hardware
	if !ioInit() {
		return -1, errors.New("Hardware driver: ioInit() failed!")
	}

	// Zero all floor button lamps
	for f := 0; f < config.NumFloors; f++ {
		if f != 0 {
			SetBtnLamp(f, config.BtnDown, false)
		}
		if f != config.NumFloors-1 {
			SetBtnLamp(f, config.BtnUp, false)
		}
		SetBtnLamp(f, config.BtnInside, false)
	}

	SetStopLamp(false)
	SetDoorLamp(false)

	// Move to defined state
	SetMotorDir(config.DirDown)
	floor := Floor()
	for floor == -1 {
		floor = Floor()
	}
	SetMotorDir(config.DirStop)
	SetFloorLamp(floor)

	log.Println(config.ColG, "Hardware initialised.", config.ColN)
	return floor, nil
}

func SetMotorDir(dirn int) {
	if dirn == 0 {
		ioWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		ioClearBit(MOTORDIR)
		ioWriteAnalog(MOTOR, 2800)
	} else if dirn < 0 {
		ioSetBit(MOTORDIR)
		ioWriteAnalog(MOTOR, 2800)
	}
}

func SetDoorLamp(value bool) {
	if value {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

func Floor() int {
	if ioReadBit(SENSOR_FLOOR1) {
		return 0
	} else if ioReadBit(SENSOR_FLOOR2) {
		return 1
	} else if ioReadBit(SENSOR_FLOOR3) {
		return 2
	} else if ioReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func SetFloorLamp(floor int) {
	if floor < 0 || floor >= config.NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		log.Println("No floor indicator will be set.")
		return
	}

	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
}

func ReadBtn(floor int, btn int) bool {
	if floor < 0 || floor >= config.NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if btn < 0 || btn >= config.NumButtons {
		log.Printf("Error: Button %d out of range!\n", btn)
		return false
	}
	if btn == config.BtnUp && floor == config.NumFloors-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if btn == config.BtnDown && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	if ioReadBit(btnChannelMatrix[floor][btn]) {
		return true
	} else {
		return false
	}
}

func SetBtnLamp(floor int, btn int, value bool) {
	if floor < 0 || floor >= config.NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return
	}
	if btn == config.BtnUp && floor == config.NumFloors-1 {
		log.Println("Button up from top floor does not exist!")
		return
	}
	if btn == config.BtnDown && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return
	}
	if btn != config.BtnUp &&
		btn != config.BtnDown &&
		btn != config.BtnInside {
		log.Printf("Invalid button %d\n", btn)
		return
	}

	if value {
		ioSetBit(lampChannelMatrix[floor][btn])
	} else {
		ioClearBit(lampChannelMatrix[floor][btn])
	}
}

func SetStopLamp(value bool) {
	if value {
		ioSetBit(LIGHT_STOP)
	} else {
		ioClearBit(LIGHT_STOP)
	}
}

func GetObstructionSignal() bool {
	return ioReadBit(OBSTRUCTION)
}

func GetStopSignal() bool {
	return ioReadBit(STOP)
}
