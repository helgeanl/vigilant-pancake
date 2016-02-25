
Getting started
---------------
(The instructions below are the "manual" version of the script `run_simtest.sh`)

 - Copy the contents of `../../source` except `simulator.con` into the `driver` folder.
 - Copy `../../source/simulator.con` into this folder.
 - Compile `simelev.so` in the `driver` folder:
   - `dmd -ofdriver/simelev.so -fPIC -lib driver/timer_event.d driver/sim_backend.d`
 - Compile `elevator_driver.so` in the `driver` folder:
   - `gcc -o driver/elevator_driver.so -std=c11 -fPIC -shared -I/usr/lib/erlang/erts-7.2/include/ io_nif.c driver/io.c driver/simelev.so /usr/lib/x86_64-linux-gnu/libphobos2.so /usr/lib/libcomedi.so`
 - Run `rdmd driver/sim_frontend.d` from this folder.
 - Run `erl` from this folder, and enter `c(main).`, `main:simtest().`
   - Or just run `echo "c(main), main:simtest()." | erl`
 - The simulator frontend should show a light on the command button on the bottom floor, and a print count of 1.
 
 
Installing Erlang
-----------------
https://www.erlang-solutions.com/resources/download.html has package installers.

If you are on the machines on the lab, you could run the following:

    wget https://packages.erlang-solutions.com/erlang/esl-erlang/FLAVOUR_1_general/esl-erlang_18.2-1~ubuntu~trusty_amd64.deb;
    sudo dpkg -i esl-erlang_18.2-1~ubuntu~trusty_amd64.deb;
    sudo apt-get install -f;
    rm esl-erlang_18.2-1~ubuntu~trusty_amd64.deb;
    
    
