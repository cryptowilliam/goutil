package gheadless

import (
	"fmt"
	"testing"
)

func TestScreenShot(t *testing.T) {
	err := ScreenShot("https://coin999.cash/", "snap.png")
	fmt.Println(err)
}
