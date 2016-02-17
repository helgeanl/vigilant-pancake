// Exercise 6 - Phoenix

package main
import(
  "log"
	"net"
	"time"
)



void main(){
  	const iAmAliveInterval = 1
		const iAmAliveTimeout = 3
		const storage = "Counter.dat"

		var counter = 0




    if !storage.exists {
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
    spawnShell("osascript -e 'tell app "Terminal" to do script "cat phoenix.d"");
    // osascript -e 'tell app "Terminal" to do script rdmd'
    counter = std.file.readText(storage).to!(typeof(counter));

    while(true){
        counter++;
        std.file.write(storage, counter.to!string);
        counter.writeln;
        Thread.sleep(iAmAliveInterval);
    }
}
