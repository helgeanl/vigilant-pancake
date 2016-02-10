What this is
============

An extension to `io.c` that lets you choose between a "hardware" elevator (the normal libcomedi interface) and a simulated one.

The source code is written in D. To use the simulator, you will need a D compiler (and runtime). I reccomend using the dmd compiler (since this is the only one I have tested it with), which you can get from [The D lang website](http://dlang.org/download.html#dmd).

Basic examples for using the simulator from various languages are also provided. Please add more examples if you have any.



Files
=====

 - `sim_frontend.d`: A program that communicates (over UDP localhost) with the simulated elevator, such that it can display the state of the elevator and take input (from keyboard) to simulate buttons and switches.
 - `sim_backend.d` (and `timer_event.d`): The simulated elevator, which must be compiled to either an object file or shared library.
 - `simulator.con`: A config file for the simulation.


Usage
=====

API changes
-----------
The only difference is:
 - `io_init(ElevatorType type)`
 - Takes either `ET_comedi` or `ET_simuation`

Compiling & linking the backend
-------------------------------
Compiling:
 - `dmd sim_backend.d timer_event.d -lib -ofsimelev`
 - You may also need to compile `io.c` into an object file: `gcc -c io.c`

Linking:
 - Give your linker of choice `simelev.a` and the D runtime `libphobos2.a` like any other object files.
 - `libphobos2.a` is (probably?) found in `/usr/lib/x86_64-linux-gnu/libphobos2.a`. This depends on your OS.
 - Example: `gcc [compile options] [c-files] simelev.a libphobos2.a -lpthread -lcomedi -lm` 
 
 
Running the frontend interface
------------------------------
The simulator interface is a standalone program, and is intended to run in its own window. It communicates over UDP localhost, so it does not need to be restarted even if the "Elevator" is restarted.  

Running:
 - `rdmd sim_frontend.d`

 
Keyboard controls
-----------------

<table>
    <thead>
        <tr>
            <th align="left" colspan="5">Controls</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td align="left"><strong>Button \ Floor</strong></td>
            <td align="center"><em>1</em></td>
            <td align="center"><em>2</em></td>
            <td align="center"><em>3</em></td>
            <td align="center"><em>4</em></td>
        </tr>
        <tr>
            <td align="right"><em>up</em></td>
            <td align="center">Q</td>
            <td align="center">W</td>
            <td align="center">E</td>
            <td align="center"></td> 
       </tr>
       <tr>
            <td align="right"><em>down</em></td>
            <td align="center"></td>
            <td align="center">S</td>
            <td align="center">D</td>
            <td align="center">F</td>
        </tr>
        <tr>
            <td align="right"><em>command</em></td>
            <td align="center">Z</td>
            <td align="center">X</td>
            <td align="center">C</td>
            <td align="center">V</td>
        </tr>
        <tr>
            <td align="left"colspan="5"><strong>Other</strong></td>
        </tr>
        <tr>
            <td align="right"><em>stop</em></td>
            <td align="left" colspan="4">T</td>
        </tr>
        <tr>
            <td align="right"><em>obstruction</em></td>
            <td align="left" colspan="4">G</td>
        </tr>
    </tbody>
</table>


A keypress must be followed by pressing Enter.  

The duration of a keypress is set in `simulator.con`.

Display
-------

```
+---------------+ +----+--------------+---------+
|   #           | |  up| 0* 1  2      | obstr:^ |
| 0 - 1*- 2 - 3 | |down|    1  2* 3*  | door:   |
|      <-       | | cmd| 0  1  2* 3   | stop:   |
+---------------+ +----+--------------+------103+
```

The ascii-art-style display is updated whenever the state of the simulated elevator is updated.

A print count (number of times a new state is printed) is shown in the lower right corner of the display. Try to avoid writing to the (simulated) hardware if nothing has happened, as writing to the screen is painfully slow. A jump of 20-50 in the printcount is fine (even expected), but if there are larger jumps or there is a continous upward count, it is time to re-evaluate some design choices.



 
 
 
 
 
