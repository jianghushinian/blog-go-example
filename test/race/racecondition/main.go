package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Account struct {
	mu      sync.Mutex
	balance int
}

func (a *Account) Withdraw(amount int) bool {
	a.mu.Lock()
	ok := a.balance >= amount
	a.mu.Unlock()
	if !ok {
		return false
	}

	time.Sleep(10 * time.Millisecond)

	a.mu.Lock()
	a.balance -= amount
	a.mu.Unlock()
	return true
}

func main() {
	a := &Account{balance: 100}

	var (
		success atomic.Int64
		wg      sync.WaitGroup
	)

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			if a.Withdraw(100) {
				success.Add(1)
			}
		})
	}
	wg.Wait()
	fmt.Println("balance:", a.balance, "success:", success.Load())
}
