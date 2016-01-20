package main

import (
	. "fmt"
	"runtime"
	"time"
	"sync"
)

var i = 0
var mutex = &sync.Mutex{}

func thread_1(a *int) {
	for j := 0; j < 1000001; j++ {
		mutex.Lock()
		i++
		mutex.Unlock()
		*a = j
	}
}

func thread_2(a *int) {
	for k := 0; k < 1000000; k++ {
		mutex.Lock()
		i--
		mutex.Unlock()
		*a = k
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var a, b = 0, 0
	chan = make(chan int, 1);
	<- chan
	go thread_1(&a)

	go thread_2(&b)

	time.Sleep(10000 * time.Millisecond)
	Println("Done, value: ", i)
	Println("Thread one run: ", a, " times. Thread two run: ", b, " times.")
}
