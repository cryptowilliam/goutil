package gpanic

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/glog"
	"os"
)

// More panic features:
// https://github.com/AlexanderChen1989/ha

// Function:
// Take over all panic and handle them gracefully.
// Usage:
// Please call this function at the beginning of main and all go func() routines like "defer HandlePanic(true)".
func HandlePanic(exit bool) {
	if r := recover(); r != nil {
		// Let application recover from panicking state
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		glog.Fata(err)
		if exit {
			os.Exit(2)
		}
	}
}
