package main

import (
	"fmt"
	"time"
)

func waitBarrier(name string, sleep int, barrier *Barrier) {
	for {
		fmt.Println(name, "running")
		time.Sleep(time.Duration(sleep) * time.Second)
		fmt.Println(name, "is waiting on barrier")
		barrier.Wait()
	}
}

func main() {
	barrier := NewBarrier(2)
	go waitBarrier("red", 2, barrier)
	go waitBarrier("blue", 4, barrier)
	time.Sleep(time.Duration(20) * time.Second)
}
