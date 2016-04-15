// Package hardware defines interactions with the lift hardware at the real time
// lab at The Department of Engineering Cybernetics at NTNU, Trondheim,
// Norway.
//
// This file is a golang port of elev.c from the hand out driver
// (https://github.com/TTK4145/Project)
package hardware

import (
	def "definitions"
	"log"
)

var lampChannelMatrix = [def.NumFloors][def.NumButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}
var btnChannelMatrix = [def.NumFloors][def.NumButtons]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

// Init initialises the lift hardware and moves the lift to a defined state.
// (Descending until it reaches a floor.)
func Init() int {
	// Init hardware
	if !ioInit() {
		return -1
	}
	// Zero all lamps
	SetDoorLamp(false)
	for f := 0; f < def.NumFloors; f++ {
		if f != 0 {
			SetBtnLamp(def.LightUpdate{Floor: f, Button: def.BtnHallDown, UpdateTo: false})
		}
		if f != def.NumFloors-1 {
			SetBtnLamp(def.LightUpdate{Floor: f, Button: def.BtnHallUp, UpdateTo: false})
		}
		SetBtnLamp(def.LightUpdate{Floor: f, Button: def.BtnCab, UpdateTo: false})
	}

	// Move to defined state
	SetMotorDir(def.DirDown)
	floor := GetFloor()
	for floor == -1 {
		floor = GetFloor()
	}
	SetMotorDir(def.DirStop)
	SetFloorLamp(floor)

	log.Println(def.ColG, "Hardware initialized.", def.ColN)
	return floor
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

func GetFloor() int {
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
	if floor < 0 || floor >= def.NumFloors {
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
	if ioReadBit(btnChannelMatrix[floor][btn]) {
		return true
	}
	return false
}

func SetBtnLamp(LightUpdate def.LightUpdate) {
	if LightUpdate.UpdateTo {
		ioSetBit(lampChannelMatrix[LightUpdate.Floor][LightUpdate.Button])
	} else {
		ioClearBit(lampChannelMatrix[LightUpdate.Floor][LightUpdate.Button])
	}
}
