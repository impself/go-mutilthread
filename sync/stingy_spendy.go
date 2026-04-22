package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	money = 100
	lock  = sync.Mutex{}
)

func stingy() {
	for i := 0; i <= 2000; i++ {
		lock.Lock()
		money += 10
		lock.Unlock()
		time.Sleep(time.Millisecond)
	}
	fmt.Println("stingy done")
}

func spendy() {
	for i := 0; i <= 2000; i++ {
		lock.Lock()
		money -= 10
		lock.Unlock()
		time.Sleep(time.Microsecond)
	}
	fmt.Println("spendy done")
}
func main() {
	go stingy()
	go spendy()
	time.Sleep(5 * time.Second)
	print(money)
}
