package main

import (
	"fmt"
	"github.com/cryptowilliam/goutil/container/gvolume"
)

func main() {
	vol, err := gvolume.ParseString("10 MB")
	fmt.Println(vol.String(), err)
}
