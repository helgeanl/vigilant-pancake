package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

func thread_1() {
	for j := 0; j < 99999; j++ {
		i++
	}
}

func thread_2() {
	for k := 0; k < 99999; k++ {
		i--
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go thread_1()
	go thread_2()
	time.Sleep(100 * time.Millisecond)
	Println(i)
}
