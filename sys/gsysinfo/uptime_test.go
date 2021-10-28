package gsysinfo

import (
	"fmt"
	"testing"
)

func TestUpTime(t *testing.T) {
	fmt.Println(UpTime())
}

func TestUpDuration(t *testing.T) {
	fmt.Println(UpDuration())
}
