import  std.stdio,
        std.file,
        std.datetime,
        std.conv,
        std.process,
        core.thread;



void main(){

    immutable iAmAliveInterval  = 1.seconds;
    immutable iAmAliveTimeout   = 3.seconds;
    immutable storage = "counter.dat";

    auto counter = 0;

    if(!storage.exists){
        std.file.write(storage, counter.to!string);
    }

    writeln(" --- Backup phase --- ");

    while(true){
        if(Clock.currTime > storage.timeLastModified + iAmAliveTimeout){
            break;
        }
        Thread.sleep(iAmAliveInterval); // No need to check more often than this
    }



    writeln(" --- Primary phase --- ");

    //spawnShell("gnome-terminal -x rdmd " ~ __FILE__);
    //spawnShell("osascript -e 'tell app "Terminal" to do script "ls"");
		spawnShell("osascript -e 'tell application \"Terminal\" to do script \"rdmd ~/Desktop/pheonix.d\"'" );
    // osascript -e 'tell app "Terminal" to do script rdmd'
    counter = std.file.readText(storage).to!(typeof(counter));

    while(true){
        counter++;
        std.file.write(storage, counter.to!string);
        counter.writeln;
        Thread.sleep(iAmAliveInterval);
    }
}
