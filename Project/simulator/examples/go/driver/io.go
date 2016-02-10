package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: ${SRCDIR}/simelev.a /usr/lib/x86_64-linux-gnu/libphobos2.a -lpthread -lcomedi -lm
#include "io.h"
*/
import "C"

func Init(elevatorType int) int {
	return int(C.io_init(C.ElevatorType(elevatorType)))
}

func Set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}

func Clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func Write_analog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

func Read_bit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}

func Read_analog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}
