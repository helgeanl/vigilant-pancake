package hardware

/*
#cgo LDFLAGS: -lcomedi -lm
#include "channels.h"
*/
import "C"
import (
	def "definitions"
	"log"
)

var lampChannelMatrix = [def.NumFloors][def.NumButtons]int{
	{C.LIGHT_UP1, C.LIGHT_DOWN1, C.LIGHT_COMMAND1},
	{C.LIGHT_UP2, C.LIGHT_DOWN2, C.LIGHT_COMMAND2},
	{C.LIGHT_UP3, C.LIGHT_DOWN3, C.LIGHT_COMMAND3},
	{C.LIGHT_UP4, C.LIGHT_DOWN4, C.LIGHT_COMMAND4},
}
var btnChannelMatrix = [def.NumFloors][def.NumButtons]int{
	{C.BUTTON_UP1, C.BUTTON_DOWN1, C.BUTTON_COMMAND1},
	{C.BUTTON_UP2, C.BUTTON_DOWN2, C.BUTTON_COMMAND2},
	{C.BUTTON_UP3, C.BUTTON_DOWN3, C.BUTTON_COMMAND3},
	{C.BUTTON_UP4, C.BUTTON_DOWN4, C.BUTTON_COMMAND4},
}

// Init initialises the lift hardware and moves the lift to a defined state.
// (Descending until it reaches a floor.)
func Init() int {
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
		ioWriteAnalog(C.MOTOR, 0)
	} else if dirn > 0 {
		ioClearBit(C.MOTORDIR)
		ioWriteAnalog(C.MOTOR, 2800)
	} else if dirn < 0 {
		ioSetBit(C.MOTORDIR)
		ioWriteAnalog(C.MOTOR, 2800)
	}
}

func SetDoorLamp(value bool) {
	if value {
		ioSetBit(C.LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(C.LIGHT_DOOR_OPEN)
	}
}

func GetFloor() int {
	if ioReadBit(C.SENSOR_FLOOR1) {
		return 0
	} else if ioReadBit(C.SENSOR_FLOOR2) {
		return 1
	} else if ioReadBit(C.SENSOR_FLOOR3) {
		return 2
	} else if ioReadBit(C.SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func SetFloorLamp(floor int) {
	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		ioSetBit(C.LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(C.LIGHT_FLOOR_IND1)
	}
	if floor&0x01 > 0 {
		ioSetBit(C.LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(C.LIGHT_FLOOR_IND2)
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
