package gnum

import (
	"fmt"
	"testing"
)

func TestFibonacci(t *testing.T) {
	for i := 0; i < 20; i++ {
		fmt.Println(Fibonacci(i))
	}
}
