// Exercise 6 - Phoenix

package main
import(
  //"log"
	//"net"
	"time"
	"os"
	//"bytes"
	"io/ioutil"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"bufio"
	"path/filepath"
	//"strings"
)

func checkErr(err error){
		if err != nil {
			fmt.Println("An unrecovarable error occured, ....", err.Error())
			os.Exit(0)
		}
}

func main(){

		_,filePath,_,_ :=runtime.Caller(0)
		directory, _ := filepath.Split(filePath)
  	const iAmAliveInterval = 1 * time.Second
		const iAmAliveTimeout = 3 * time.Second
		storageName := directory + "/Counter.dat"
		var counter = 0

  	print(string("----- Backup Phase -----\n"))
		// open output file
		storage, err := os.Open(storageName)
		if err != nil {
			newFile,_ := os.Create(storageName)
			newFile.Close()
			//checkErr(err)
		}

		for(true){
			fmt.Println("...Waiting")
			fileStat,err := os.Stat(storageName)
			checkErr(err)
				if time.Now().After(fileStat.ModTime().Add(iAmAliveTimeout))  {
					break;
				}
				time.Sleep(iAmAliveInterval) // No need to check more often than this
		}


  	print(string("----- Primary Phase -----\n"))

		// Run a new terminal window with backup
		//app := "osascript"
		//arg0 := "-e"
		//arg1 := "tell app \"Terminal\" to do script \"go run '"+filePath+"'\""
		
		app := "gnome-terminal"
		arg0 := "-x"
		arg1 := "go run "+filePath
		cmd := exec.Command( app,arg0,arg1)

		cmd.Start(); // Carefull...

		r4 := bufio.NewReader(storage)
    b4, _ := r4.Peek(8)

    counter,_ = strconv.Atoi(string(b4))
		storage.Close()
		storage,_=os.Create(storageName)

		defer func() {
				if err := storage.Close(); err != nil {
					panic(err)
				}
			}()

    for(true){
        counter++;
				counterString := strconv.Itoa(counter)
				ioutil.WriteFile(storageName, []byte(counterString), 0644)
				//w := bufio.NewWriter(storage)
		  	///_,err := w.WriteString(counterString)
				//checkErr(err)
        fmt.Println(counter)
        time.Sleep(iAmAliveInterval)
    }
}
