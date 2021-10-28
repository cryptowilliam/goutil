package gcounter

import (
	"sync"
	"testing"
	"time"
)

const (
	_testMin = 0
	_testMax = 1024
)

func BenchmarkMutex(b *testing.B) {
	c, _ := NewAccumulator(_testMin, _testMax)
	for i := 0; i < b.N; i++ {
		c.Incr()
	}
}

func TestMutex(t *testing.T) {
	c, _ := NewAccumulator(_testMin, _testMax)
	for i := 0; i < 100; i++ {
		c.Incr()
	}
	if c.Reset() != _testMin {
		t.Fail()
	}
}

func TestConcurrent(t *testing.T) {
	c, _ := NewAccumulator(_testMin, _testMax)

	num := 100
	threads := 30
	wg := sync.WaitGroup{}

	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Add(-1)
			for j := 0; j < num; j++ {
				c.Incr()
			}
		}()
	}
	wg.Wait()

	var correct int64
	if threads*num <= _testMax {
		correct = int64(num * threads)
	} else {
		correct = int64(num * threads % (_testMax - _testMin + 1))
	}
	if c.Get() != correct {
		t.Errorf("Correct result is %d, in fact returns %d", correct, c.Get())
	}
}

func BenchmarkMutexWithFlush(b *testing.B) {
	c, _ := NewAccumulator(_testMin, _testMax)
	for i := 0; i < b.N; i++ {
		c.Incr()
		if i%10000 == 0 {
			c.Reset()
		}
	}
}

func BenchmarkMutexWithConcurrentFlush(b *testing.B) {
	c, _ := NewAccumulator(_testMin, _testMax)
	sleep, _ := time.ParseDuration("10ms")
	go func() {
		for {
			time.Sleep(sleep)
			c.Reset()
		}
	}()
	for i := 0; i < b.N; i++ {
		c.Incr()
		if i%10000 == 0 {
			c.Reset()
		}
	}
}
