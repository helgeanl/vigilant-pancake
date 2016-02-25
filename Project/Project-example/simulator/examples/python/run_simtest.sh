mkdir driver;
cp ../../source/* driver;
mv driver/simulator.con .;
dmd -ofdriver/simelev.so -fPIC -lib driver/timer_event.d driver/sim_backend.d;
gcc -o driver/elevator_driver.so -std=c11 -fPIC -shared driver/io.c driver/simelev.so /usr/lib/x86_64-linux-gnu/libphobos2.so /usr/lib/libcomedi.so;
gnome-terminal -e "rdmd driver/sim_frontend.d";
sleep 2;
python3 main.py