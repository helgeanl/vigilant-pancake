from ctypes import *
from time import sleep
from enum import *

from channels import *

class ElevatorType(IntEnum):
    comedi = 0,
    simulation = 1
    


io = cdll.LoadLibrary("driver/elevator_driver.so")

io.io_init(c_int(ElevatorType.simulation))
io.io_set_bit(c_int(Channel.LIGHT_COMMAND1))

sleep(2)