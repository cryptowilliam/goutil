package glog

import (
	"sync"
	"testing"
)

func Test_Infof(t *testing.T) {
	Init(false)
	wg := sync.WaitGroup{}

	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Add(-1)
			for j := 0; j < 100000; j++ {
				Infof("%d, %d", n, j)
			}
		}(i)
	}

	wg.Wait()
}
