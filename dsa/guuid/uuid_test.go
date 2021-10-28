package guuid

import (
	"fmt"
	"testing"
)

func TestNewBigInt(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(NewBigInt().String())
	}
}
