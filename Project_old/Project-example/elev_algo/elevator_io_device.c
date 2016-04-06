
#include "elevator_io_device.h"

#include <assert.h>
#include <stdlib.h>

#include "con_load.h"
#include "driver/channels.h"
#include "driver/io.h"


static int elev_read_floorSensor(void);
static int elev_read_requestButton(int floor, Button button);
static int elev_read_stopButton(void);
static int elev_read_obstruction(void);

static void elev_write_floorIndicator(int floor);
static void elev_write_requestButtonLight(int floor, Button button, int value);
static void elev_write_doorLight(int value);
static void elev_write_stopButtonLight(int value);
static void elev_write_motorDirection(Dirn dirn);



static void __attribute__((constructor)) elev_init(void){

    ElevatorType et;

    con_load("elevator.con",
        con_enum("elevatorType", et,
            con_match(ET_simulation)
            con_match(ET_comedi)
        )
    )

    int success = io_init(et);
    assert(success && "Elevator hardware initialization failed");

    for(int floor = 0; floor < N_FLOORS; floor++) {
        for(Button btn = 0; btn < N_BUTTONS; btn++){
            elev_write_requestButtonLight(floor, btn, 0);
        }
    }

    elev_write_stopButtonLight(0);
    elev_write_doorLight(0);
    elev_write_floorIndicator(0);
}


ElevInputDevice elevio_getInputDevice(void){
    return (ElevInputDevice){
        .floorSensor    = &elev_read_floorSensor,
        .requestButton  = &elev_read_requestButton,
        .stopButton     = &elev_read_stopButton,
        .obstruction    = &elev_read_obstruction
    };
}


ElevOutputDevice elevio_getOutputDevice(void){
    return (ElevOutputDevice){
        .floorIndicator     = &elev_write_floorIndicator,
        .requestButtonLight = &elev_write_requestButtonLight,
        .doorLight          = &elev_write_doorLight,
        .stopButtonLight    = &elev_write_stopButtonLight,
        .motorDirection     = &elev_write_motorDirection
    };
}


char* elevio_dirn_toString(Dirn d){
    return
        d == D_Up    ? "D_Up"         :
        d == D_Down  ? "D_Down"       :
        d == D_Stop  ? "D_Stop"       :
                       "D_UNDEFINED"  ;
}


char* elevio_button_toString(Button b){
    return
        b == B_HallUp       ? "B_HallUp"        :
        b == B_HallDown     ? "B_HallDown"      :
        b == B_Cab          ? "B_Cab"           :
                              "B_UNDEFINED"     ;
}





static const int floorSensorChannels[N_FLOORS] = {
    SENSOR_FLOOR1,
    SENSOR_FLOOR2,
    SENSOR_FLOOR3,
    SENSOR_FLOOR4,
};

static int elev_read_floorSensor(void){
    for(int f = 0; f < N_FLOORS; f++){
        if(io_read_bit(floorSensorChannels[f])){
            return f;
        }
    }
    return -1;
}


static const int buttonChannels[N_FLOORS][N_BUTTONS] = {
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
};

static int elev_read_requestButton(int floor, Button button){
    assert(floor >= 0);
    assert(floor < N_FLOORS);
    assert(button < N_BUTTONS);

    if(io_read_bit(buttonChannels[floor][button])){
        return 1;
    } else {
        return 0;
    }
}


static int elev_read_stopButton(void){
    return io_read_bit(STOP);
}


static int elev_read_obstruction(void){
    return io_read_bit(OBSTRUCTION);
}





static void elev_write_floorIndicator(int floor){
    assert(floor >= 0);
    assert(floor < N_FLOORS);

    if(floor & 0x02){
        io_set_bit(LIGHT_FLOOR_IND1);
    } else {
        io_clear_bit(LIGHT_FLOOR_IND1);
    }

    if(floor & 0x01){
        io_set_bit(LIGHT_FLOOR_IND2);
    } else {
        io_clear_bit(LIGHT_FLOOR_IND2);
    }
}


static const int buttonLightChannels[N_FLOORS][N_BUTTONS] = {
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
};

static void elev_write_requestButtonLight(int floor, Button button, int value){
    assert(floor >= 0);
    assert(floor < N_FLOORS);
    assert(button < N_BUTTONS);

    if(value){
        io_set_bit(buttonLightChannels[floor][button]);
    } else {
        io_clear_bit(buttonLightChannels[floor][button]);
    }
}


static void elev_write_doorLight(int value){
    if(value){
        io_set_bit(LIGHT_DOOR_OPEN);
    } else {
        io_clear_bit(LIGHT_DOOR_OPEN);
    }
}


static void elev_write_stopButtonLight(int value){
    if(value){
        io_set_bit(LIGHT_STOP);
    } else {
        io_clear_bit(LIGHT_STOP);
    }
}


static void elev_write_motorDirection(Dirn dirn){
    switch(dirn){
    case D_Up:
        io_clear_bit(MOTORDIR);
        io_write_analog(MOTOR, 2800);
        break;
    case D_Down:
        io_set_bit(MOTORDIR);
        io_write_analog(MOTOR, 2800);
        break;
    case D_Stop:
    default:
        io_write_analog(MOTOR, 0);
        break;
    }
}





