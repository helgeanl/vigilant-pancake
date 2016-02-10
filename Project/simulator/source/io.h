
#pragma once


typedef enum {
    ET_comedi,
    ET_simulation
} ElevatorType;


// Returns 0 on init failure
int io_init(ElevatorType type);

void io_set_bit(int channel);
void io_clear_bit(int channel);

int io_read_bit(int channel);

int io_read_analog(int channel);
void io_write_analog(int channel, int value);


