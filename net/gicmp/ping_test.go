package gicmp

import (
	"fmt"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	sz := uint16(1200)
	fmt.Println(Ping("baidu.com", &sz, time.Second*5))
}
