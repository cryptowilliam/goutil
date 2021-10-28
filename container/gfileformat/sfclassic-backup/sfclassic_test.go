package sfclassic

import (
	"os"
	"testing"
)

func TestIdentify(t *testing.T) {
	sf := New()
	f, _ := os.Open("classic/classic.sig")
	defer f.Close()
	ids, _ := sf.Identify(f, "classic.sig")
	if ids[0].String() != "fmt/883" {
		t.Fatal(ids)
	}
}
