package gtime

import (
	"fmt"
	"testing"
)

func TestMonoClock_Now(t *testing.T) {
	fmt.Println(NewMonoClock().Now().String())
}
