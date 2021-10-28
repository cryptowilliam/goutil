package groutine

import (
	"fmt"
	"testing"
	"time"
)

func demo(args ...interface{}) {
	fmt.Println(args[0].(string))
}

func TestNewRoutine(t *testing.T) {
	tmo := time.Second * 10

	r := RunLoop(demo, time.Second, &tmo, "demo")
	r.Wait()
}
