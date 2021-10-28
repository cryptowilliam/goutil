package gprofile

import (
	"github.com/rakyll/autopprof"
	"time"
)

// call this function at main function
// run your program and send SIGQUIT to the process (or CTRL+\ on Mac).
// profile capturing will start. Pprof UI will be started once capture is completed.
func StartProfile(interval time.Duration) {
	autopprof.Capture(autopprof.CPUProfile{
		Duration: interval,
	})
}
