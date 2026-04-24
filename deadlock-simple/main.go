package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock1 = sync.Mutex{}
	lock2 = sync.Mutex{}
)

func redLock() {
	for {

		fmt.Println("red required lock1")
		lock1.Lock()
		fmt.Println("red required lock2")
		lock2.Lock()
		fmt.Println("red has both lock")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("red lock released")
	}

}
func blueLock() {
	for {
		fmt.Println("blue required lock2")
		lock2.Lock()
		fmt.Println("blue required lock1")
		lock1.Lock()
		fmt.Println("blue has both lock")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("blue lock released")
	}

}
func main() {
	go redLock()
	go blueLock()
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}
