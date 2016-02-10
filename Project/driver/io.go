// Wrapper for libComedi Elevator control.
// These functions provide an interface to the elevators in the real time lab

package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

//Dropping "Io" prefix due to the "io." calling mechanism in GO
func Init()int{    				
	return C.io_init()
}

func SetBit(channel int){
	return C.io_set_bit(channel)
}

func ClearBit(channel int){
	return C.io_clear_bit(channel)
}

func WriteAnalog(channel int, value int){
	return C.io_write_analog(channel, value)
}
func ReadBit(channel int)int{
	return C.io_read_bit(channel)
}
func ReadAnalog(channel int)int {
	return C.io_read_analog(channel)
}
