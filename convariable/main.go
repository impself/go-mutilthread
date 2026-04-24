package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	money    = 100
	lock     = sync.Mutex{}
	variable = sync.NewCond(&lock)
)

// 存钱
func stingy() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		money += 10
		variable.Signal()
		fmt.Println("存钱后余额：", money)
		lock.Unlock()
		time.Sleep(time.Millisecond)

	}
	fmt.Println("stingy done")

}

// 取钱
func spendy() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		for money-20 < 0 {
			variable.Wait()
		}
		money -= 20
		lock.Unlock()
		fmt.Println("取钱后余额：", money)
		time.Sleep(time.Millisecond)
	}
	fmt.Println("spendy done")
}

func main() {
	go stingy()
	go spendy()
	time.Sleep(3000 * time.Millisecond)
	fmt.Println("余额：", money)
}
