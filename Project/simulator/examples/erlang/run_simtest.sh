mkdir driver;
cp ../../source/* driver;
mv driver/simulator.con .;
dmd -ofdriver/simelev.so -fPIC -lib driver/timer_event.d driver/sim_backend.d;
gcc -o driver/elevator_driver.so -std=c11 -fPIC -shared -I/usr/lib/erlang/erts-7.2/include/ io_nif.c driver/io.c driver/simelev.so /usr/lib/x86_64-linux-gnu/libphobos2.so /usr/lib/libcomedi.so;
gnome-terminal -e "rdmd driver/sim_frontend.d";
sleep 2;
echo "c(main), main:simtest()." | erl;