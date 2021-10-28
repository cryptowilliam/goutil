package gstring

import (
	"fmt"
	"strings"
	"testing"
)

func TestRemoveDuplicate(t *testing.T) {
	src := []string{"", "", "a", "ab", "ab"}
	expect := []string{"", "a", "ab"}
	dst := RemoveDuplicate(src)
	if len(dst) != len(expect) {
		t.Errorf("RemoveDuplicate error")
		return
	}
	fmt.Println("output:", strings.Join(dst, ", "), ".")
}
