package gsignal

import (
	// import for a
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// syscall.SIGTNT (os.Interrupt): ctrl +c
// syscall.SIGTERM: kill, killall
// syscall.SIGHUP: hangup (terminal closed)

var reged = func() atomic.Value { s := atomic.Value{}; s.Store(false); return s }()
var exitSigCh = make(chan os.Signal, 1)
var closeCh = make(chan string)
var onExitCallbacks = make([]ExitHandler, 0)
var onExitMu sync.Mutex

type ExitHandler func(sig os.Signal, closemsg string)

// adds exit callback, callbacks will be executed before exit
// Notice:
// Can catch kill / killall / Hangup(close terminal)
// Can't catch os.Exit, main() return and exception panic crash
func RegisterExitCallback(f ExitHandler) {
	onExitMu.Lock()
	onExitCallbacks = append(onExitCallbacks, f)
	onExitMu.Unlock()

	// catch signals
	if !reged.Load().(bool) {
		signal.Notify(exitSigCh, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		reged.Store(true)
		go waitExitSignals()
	}
}

// Wait the signals for exit
// Return caught signal
func waitExitSignals() {
	var sig os.Signal
	var closemsg string

	// wait
	select {
	case sig = <-exitSigCh:
	case closemsg = <-closeCh:
	}

	// execute callback functions
	onExitMu.Lock()
	cbs := onExitCallbacks
	onExitMu.Unlock()
	for _, f := range cbs {
		f(sig, closemsg)
	}

	// quit application
	os.Exit(0)
}

// Exit send close signal with given message
func CloseWaiter(message string) {
	closeCh <- message
}
