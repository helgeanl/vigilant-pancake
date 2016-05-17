package main

import (
	. "fmt"
	"runtime"
	"time"
	//"sync"
)

var i = 0
//var mutex = &sync.Mutex{}

func thread_1(done chan bool) {
	for j := 0; j < 10001; j++ {
		<-done
		i++
		done <- true
	}
}

func thread_2(done chan bool) {
	for k := 0; k < 10000; k++ {
		//for<-ch{}
		<-done
		i--
		done <- true
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//var a, b = 0, 0

	done := make(chan bool, 1);
	done<-true

	go thread_1(done)
	go thread_2(done)

	time.Sleep(100 * time.Millisecond)
	Println("Done, value: ", i)
	//Println("Thread one run: ", a, " times. Thread two run: ", b, " times.")
}
