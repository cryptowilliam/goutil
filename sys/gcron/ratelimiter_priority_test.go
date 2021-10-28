package gcron

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterPriority_TakeWait(t *testing.T) {
	rlp := NewRateLimiterPriority(time.Second)

	var wg sync.WaitGroup
	var r string
	var rMu sync.RWMutex

	wg.Add(3)
	fmt.Println("please wait 30 seconds...")

	go func() {
		defer wg.Add(-1)
		for i := 0; i < 10; i++ {
			rlp.TakeWait(1)
			//fmt.Println("1", time.Now())

			rMu.Lock()
			r += "1"
			rMu.Unlock()
		}
	}()

	go func() {
		defer wg.Add(-1)
		for i := 0; i < 10; i++ {
			rlp.TakeWait(3)
			//fmt.Println("3", time.Now())

			rMu.Lock()
			r += "3"
			rMu.Unlock()
		}
	}()

	go func() {
		defer wg.Add(-1)
		time.Sleep(time.Second * 10)
		for i := 0; i < 10; i++ {
			rlp.TakeWait(0)
			//fmt.Println("0", time.Now())

			rMu.Lock()
			r += "0"
			rMu.Unlock()
		}
	}()

	wg.Wait()
	if r != "111111111100000000003333333333" {
		t.Errorf("RateLimiterPriority Test Failed")
		return
	}
}
