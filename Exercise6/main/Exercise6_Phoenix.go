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
	//"strings"
)

func checkErr(err error){
	if err != nil {
		fmt.Println("An unrecovarable error occured", err.Error())
		os.Exit(0)
	}
}

func main(){

	_,filePath,_,_ :=runtime.Caller(0)
	directory, _ := os.Getwd()
  	const iAmAliveInterval = 1 * time.Second
		const iAmAliveTimeout = 3 * time.Second
		 storageName := "Counter.dat"
		 fmt.Println("Directory: " + directory )
		 fmt.Println("Path: " + filePath )
		var counter = 0

  print(string("----- Backup Phase -----\n"))
	// open output file
	storage, err := os.Open(storageName)
	if err != nil {
		panic(err);
	}
	// close fo on exit and check for its returned error
//	defer func() {
//		if err := storage.Close(); err != nil {
//			panic(err)
//		}
//	}()


		for(true){
			fmt.Println("...Waiting")
			fileStat,_ := os.Stat(storageName)
				if time.Now().After(fileStat.ModTime().Add(iAmAliveTimeout))  {
					break;
				}
				time.Sleep(iAmAliveInterval) // No need to check more often than this
		}


  	print(string("----- Primary Phase -----\n"))


		app := "osascript"
		arg0 := "-e"
		arg1 := "tell application \"Terminal\" to do script \"go run '"+filePath+"'\""
		cmd := exec.Command( app,arg0,arg1)
		//if false {
			cmd.Start(); // Carefull...
		//}
		checkErr(err)

		/*
		// Write
		w := bufio.NewWriter(f)
    n4, err := w.WriteString("buffered\n")

		// Read
		r4 := bufio.NewReader(f)
    b4, err := r4.Peek(5)
    check(err)
		*/

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
        //storage.Write(counter)
				ioutil.WriteFile(storageName, []byte(counterString), 0644)
				//w := bufio.NewWriter(storage)
		  	///_,err := w.WriteString(counterString)
				//checkErr(err)
        fmt.Println(counter)
        time.Sleep(iAmAliveInterval)
    }

}
