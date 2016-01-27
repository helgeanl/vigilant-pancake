//
//  buzzer.c
//  
//
//  Created by Helge-André Langåker on 14.11.2015.
//
//

#include "buzzer.h"

/*
 Play music with your robot, version 2!
 Buzzer example for the Orangutan/ATmega168.
 This version is entirely interrupt-driven.
 
 This program uses Timer2 in "Clear Timer/Counter on Compare Match"
 mode to generate the buzzer frequency and to count down notes.
 The counter is incremented at 125 kHz, and resets when count = OCR2A.
 Upon reset, buzzer port pin is toggled.
 
 OCR2A=255 corresponds to 488 Hz, or buzzer frequency 244 Hz.
 Higher frequencies are set by reducing the value of OCR2A.
 
 The same technique is used to implement a general purpose 1 millisecond
 timer tick using timer 0 (unsigned integer ms_tick). Timer0 also plays the tune,
 counting down notes in 1 ms intervals.
 
 In this example, musical notes are defined according to the international tempered scale.
 
 The CPU clock frequency assumed to be 8 MHz. If different, TCCR0 and 2 dividers
 should be changed appropriately. Otherwise the timing will be incorrect.
 
 Produces about 834 bytes of loaded code with -O3
 sjames_remington at yahoo dot com
 */

#define F_CPU 8000000UL
#include <avr/io.h>
#include <avr/interrupt.h>
#include <util/delay.h>

/*
 
 global variables for music player
 
 Play_tune: set to 1 to play tune. Music routine resets to zero when done.
 
 *tune_array = pointer to tune array consisting of byte pairs [duration, note].
	duration = 8 for whole note, 4 for half note, etc. duration = 0 terminates tune.
	note = timer constant corresponding to predefined note value taken from table defined below.
 
 tempo = (time in milliseconds for whole note)/8
 
 In this version, the zeroth entry of the tune data array contains (tempo)/8
	the address of the first entry is passed to the music routine
 
 */

unsigned char play_tune, *tune_array;
unsigned int tempo;
unsigned int ms_tick=0;  //general purpose millisecond timer tick, incremented by TIMER2

/*
 TIMER 2 OCR2A interrupt service routine
 This routine is called when TIMER2 count (TCNT2) = OCR2A, which is set on the fly for each note.
 Operation is to simply toggle the buzzer line, producing frequency = 1/(2*timer period)
 */

ISR(TIMER2_COMPA_vect)
{
    PORTB ^= 1;					//toggle buzzer output line
}


/*
 TIMER 0 OCR0A interrupt service routine
 This routine is called every millisecond, when TCNT0 = OCR0A (125)
 
 Operation:
 -increment general purpose timer tick at 1 ms intervals
 -if tempo>0, continue playing current tune, pointed to by *tune_array
 -signals "done playing tune" by setting global variable play_tune->0
 */

ISR(TIMER0_COMPA_vect)
{
    
    static unsigned int duration=0;	//internal timer for current note
    volatile unsigned char dur,note;
    
    ms_tick++;								//increment timer tick at 1 ms intervals
    
    if(play_tune) {							//are we playing a tune?
        if(duration>0)	duration--;				//yes, count down this note's timer
        
        else {								//done with this note, get next [duration,note] pair
            dur=*tune_array++;
            duration = tempo*dur; 			//calculate time in ms to play this note
            if(dur==0) play_tune=0;		//duration entry=0, so signal done playing tune
            note=*tune_array++;
            if (note>0) { 				//check for note or silence to play
                OCR2A = note;			//set timer interrupt period and point to next entry in tune table
                DDRB |= 1; 				//turn on buzzer output line
            }
            else {
                OCR2A=255;				//silence, reduce interrupt load
                DDRB &=~(1); 			//and turn off buzzer
            }
        }
    }
}

/*
 Set up buzzer to play tones. Timer/Counter 2 is used in CTC mode, clear timer on compare/match
 Clocked at 125 kHz  = 8MHz/64
 
 OCR2A  overflow  buzzer
 value    freq    freq
 
 255      488     244 (Hz)
 128      976     488
 64      1953     976
 32      3906    1953
 
 For other CPU clock or timer ratios, multiply or divide accordingly
 
 */
void init_buzzer(void)
{
    DDRB |= 1;				//set data direction reg bit 0 to output
    PORTB &=~(1);			//set buzzer output low
    /*
     initialize timer 0
     */
    TCCR0A = 0; //tmr0 Clear Timer on Compare Match
    TCCR0B = _BV(WGM01)|_BV(CS01)|_BV(CS00);  //CTC mode with CPU clock/64 = 125 kHz
    OCR0A = 125;  //set for 1 ms interrupt
    TCNT0=0;
    TIMSK0 = (1<<OCIE0A); // enable interrupt on TCNT2=OCR0A
    /*
     initialize timer 2
     */
    TCCR2A= (1<<WGM21);	//mode 2, clear timer on compare match
    TCCR2B=(1<<CS22);				//CS22=1,CS21=0,CS20=0	=> counter increments at 8 MHz/64 = 125 kHz
    TCNT2=0;				//clear counter
    OCR2A=255;				//set OCR2A to minimum o/p frequency
    OCR2B=255;				//set OCR2B to max, irrelevant
    TIMSK2= (1<<OCIE2A);	//enable interrupt on compare match with OCR2A
}

void init_pwm(void)
{
    DDRB  |=  (1 << 2);		//set data direction reg B bit 2 to output
    PORTB &= ~(1 << 2);		//set PWM output low
    /*
     initialize timer 0
     */
    TCCR0A = 0; //tmr0 Clear Timer on Compare Match
    TCCR0B = _BV(WGM01)|_BV(CS01)|_BV(CS00);  //CTC mode with CPU clock/64 = 125 kHz
    OCR0A = 125;  //set for 1 ms interrupt
    TCNT0=0;
    TIMSK0 = (1<<OCIE0A); // enable interrupt on TCNT2=OCR0A
    /*
     initialize timer 2
     */
    TCCR2A= (1<<WGM21);	//mode 2, clear timer on compare match
    TCCR2B=(1<<CS22);				//CS22=1,CS21=0,CS20=0	=> counter increments at 8 MHz/64 = 125 kHz
    TCNT2=0;				//clear counter
    OCR2A=255;				//set OCR2A to minimum o/p frequency
    OCR2B=255;				//set OCR2B to max, irrelevant
    TIMSK2= (1<<OCIE2A);	//enable interrupt on compare match with OCR2A
}


/*
 
 Define three octaves of notes for coding convenience. x=sharp (#)
 This table is defined according to the International Equal-Tempered Scale,
 With A4 = 440.00 Hz frequency standard
 B3 here is 247.03 Hz, very close to the standard value of 246.96 Hz.
 
 The formula for the timer constant = 8E6/(128*frequency in Hz)
 
 */

#define p   0   //pause = silence
#define b3  253
#define c4  239
#define cx4 225
#define d4  213
#define dx4 201
#define e4  190
#define f4  179
#define fx4 169
#define g4  159
#define gx4 150
#define a4  142
#define ax4 134
#define b4  127
#define c5  119
#define cx5 113
#define d5  106
#define dx5 100
#define e5  95
#define f5  89
#define fx5 84
#define g5  80
#define gx5 75
#define a5  71
#define ax5 67
#define b5  63
#define c6  60
#define cx6 56
#define d6  53
#define dx6 50
#define e6  47
#define f6  45
#define fx6 42
#define g6  40
#define gx6 38
#define a6  36
#define b6  32
#define c7  30


int main(void)
{
    /*
     Tune array entries are:
     first byte, tempo in (milliseconds for whole note) divided by 8
     following byte pairs: [duration , note]
     duration 8=whole note, 4=half note etc.
     note = 0 silence for this interval
     duration,note pair =[0,0] ends the tune.
     */
    /*
     unsigned char fuer_elise[]= //modified from AVR Butterfly code example. Thanks, ATMEL Norway!
     
     {		(120>>3),
     8,e5, 8,dx5, 8,e5, 8,dx5, 8,e5, 8,b4, 8,d5, 8,c5, 4,a4, 8,p,
     8,c4, 8,e4, 8,a4, 4,b4, 8,p, 8,e4, 8,gx4, 8,b4, 4,c5, 8,p, 8,e4,
     8,e5, 8,dx5, 8,e5, 8,dx5, 8,e5, 8,b4, 8,d5, 8,c5, 4,a4, 8,p, 8,c4,
     8,e4, 8,a4, 4,b4, 8,p, 8,e4, 8,c5, 8,b4, 4,a4,
     8,p,8,p,  //up one octave for test!
     8,e6, 8,dx6, 8,e6, 8,dx6, 8,e6, 8,b5, 8,d6, 8,c6, 4,a5, 8,p,
     8,c5, 8,e5, 8,a5, 4,b5, 8,p, 8,e5, 8,gx5, 8,b5, 4,c6, 8,p, 8,e5,
     8,e6, 8,dx6, 8,e6, 8,dx6, 8,e6, 8,b5, 8,d6, 8,c6, 4,a5, 8,p, 8,c5,
     8,e5, 8,a5, 4,b5, 8,p, 8,e5, 8,c6, 8,b5, 4,a5,
     0, 0
     };
     unsigned char chirp[]=  //bird chirp
     { (18>>3),1,100,1,95,1,90,1,85,1,80,1,75,1,70,1,65,1,60,1,55,1,50,1,45,1,40,
     1,35,2,32,2,30,2,28,2,26,2,24,2,22,2,20,2,18,2,16, 0, 0};
     */
    unsigned char sorcerer[]= //Sorcerer's Apprentice by Paul Dukas. 3/8 time, stored by bars #1-28
    {
        (210>>3),
        4,d5,4,p,4,p,
        4,a5,4,p,4,p,
        4,a4,4,b4,4,cx5,
        4,d5,4,p,4,f5,
        4,d5,4,p,4,f5,
        4,e5,4,d5,4,cx5,
        4,d5,4,p,4,f5,
        4,d5,4,p,4,f5,
        4,e5,4,d5,4,cx5,
        4,d5,4,p,4,f5,
        4,d5,4,f5,4,e5,
        4,d5,4,e5,4,f5,
        4,e5,4,g5,4,f5,
        4,e5,4,p,4,gx5,
        4,d5,4,p,4,f5,
        4,e5,4,g5,4,f5,
        4,e5,4,p,4,gx5,
        4,d5,4,p,4,f5,
        4,e5,4,f5,4,g5,
        3,a5,1,p,3,a5,1,p,3,a5,1,p,
        3,a5,1,p,3,a5,1,p,3,a5,1,p,
        4,a5,4,g5,4,f5,
        4,e5,4,g5,4,f5,
        4,e5,4,d5,4,c5,
        1,d5,3,c5,4,b4,4,a4,
        4,g5,4,p,4,p,
        10,g5,2,p,
        2,a5,0,0
    };
    
    init_buzzer();
    DDRB &= ~(1);			//turn off buzzer for now by setting PORTB.0 = input
    tempo=sorcerer[0];			//set tempo = ms per whole note, divided by 8 for note duration scaling
    tune_array=&sorcerer[1];	//point to current tune
    play_tune=1;			//tell timer interrupt routine to play tune
    sei();
    
    loop_until_bit_is_clear(play_tune,0);			//wait until done, or go do something else!
    
    DDRB &=~(1); 			//turn off buzzer and
    return 0;				//exit
}