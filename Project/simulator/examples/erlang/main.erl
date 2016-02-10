-module(main).
-export([simtest/0]).

-define(comedi, 0).
-define(simulation, 1).

-define(LIGHT_COMMAND1, 16#300 + 13).

simtest() ->
    io:fwrite("hello from erlang\n"),
    ok = erlang:load_nif("./driver/elevator_driver", 0),
    io_init(?simulation),
    io_set_bit(?LIGHT_COMMAND1),
    timer:sleep(timer:seconds(1)),
    io:fwrite("goodbye from erlang\n").
    
io_init(_X) -> 
    exit(nif_library_not_loaded).
    
io_set_bit(_X) -> 
    exit(nif_library_not_loaded).

