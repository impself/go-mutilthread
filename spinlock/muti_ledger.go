package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	totalAccounts  = 50000
	maxAmountMoved = 10
	initMoney      = 100
	thread         = 4
)

var rwLock sync.RWMutex

func perfmMovements(ledger *[totalAccounts]int32, totalTrans *int64) {
	for {
		accountA := rand.Intn(totalAccounts)
		accountB := rand.Intn(totalAccounts)
		for accountA == accountB {
			accountB = rand.Intn(totalAccounts)
		}
		amountToMove := rand.Int31n(maxAmountMoved)

		// lock in order to prevent deadlock
		a, b := accountA, accountB
		if a > b {
			a, b = b, a
		}
		rwLock.Lock()
		ledger[accountA] -= amountToMove
		ledger[accountB] += amountToMove
		rwLock.Unlock()

		atomic.AddInt64(totalTrans, 1)
	}
}

func main() {
	fmt.Printf("total thread: %d, total account: %d\n", thread, totalAccounts)
	var ledger [totalAccounts]int32
	var totalTrans int64
	for i := 0; i < totalAccounts; i++ {
		ledger[i] = initMoney
	}
	for i := 0; i < thread; i++ {
		go perfmMovements(&ledger, &totalTrans)
	}
	for {
		time.Sleep(2000 * time.Millisecond)
		rwLock.RLock()
		var sum int32
		for i := 0; i < totalAccounts; i++ {
			sum += ledger[i]
		}
		rwLock.RUnlock()
		fmt.Printf("transactions: %d, total money: %d\n", totalTrans, sum)
	}
}
