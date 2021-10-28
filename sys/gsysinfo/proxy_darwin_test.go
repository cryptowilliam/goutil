package gsysinfo

import (
	"fmt"
	"testing"
)

func TestGetCurrentNetworkInterface(t *testing.T) {
	fmt.Println(GetCurrentNetworkInterface())
}
