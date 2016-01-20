package main

import (
	. "fmt"
	"runtime"
	"time"
)


func thread_1(ch chan int) {
	var k = 0
	for j := 0; j < 10000; j++ {
		k = <-ch
		k++
		ch <- k
	}
}


func thread_2(ch chan int) {
	var j = 0
	for k := 0; k < 10000; k++ {
		j = <-ch
		j--
		ch <- j
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//var a, b = 0, 0
	var i = 0
	ch := make(chan int, 1);
	ch <- i

	go thread_1(ch)
	go thread_2(ch)

	time.Sleep(500 * time.Millisecond)
	Println("Done, value: ", <-ch)
	//Println("Thread one run: ", a, " times. Thread two run: ", b, " times.")
}
