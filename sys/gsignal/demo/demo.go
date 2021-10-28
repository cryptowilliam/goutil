package main

import (
	"fmt"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gsignal"
	"os"
	"time"
)

func exitcb(sig os.Signal, closemsg string) {
	s := fmt.Sprintf("exit signal %s, closemsg %s", sig.String(), closemsg)
	fmt.Println(s)
	gfs.AppendStringToFile(s, "demo.log")
}

func main() {
	gsignal.RegisterExitCallback(exitcb)
	/*go func() {
		type obj struct {
			name string
		}
		o := new(obj)
		o = nil
		o.name = ""
		time.Sleep(time.Second * 5)
		os.Exit(-1)
	}()*/

	for {
		time.Sleep(time.Second)
	}
}
