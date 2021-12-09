package gcmd

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	Cmd struct {
		wait bool
		waitTimeout time.Duration
		print bool
		//log glog.Interface
		cmdline string
		result []byte
		err error
		pid int
	}

	Runner Cmd
)

func New() *Cmd {
	return &Cmd{}
}

func (c *Cmd) Wait(wait bool, timeout time.Duration) *Cmd {
	c.wait = wait
	c.waitTimeout = timeout
	return c
}

func (c *Cmd) Print(print bool) *Cmd {
	c.print = print
	return c
}

func (c *Cmd) AsFunc(name string, arg ...string) *Cmd {
	c.cmdline = strings.Join(append([]string{name}, arg...), " ")
	return c
}

func (c *Cmd) AsShell(cmdline string) *Cmd {
	c.cmdline = cmdline
	return c
}

func (c *Cmd) Ready() *Runner {
	return (*Runner)(c)
}

func (c *Runner) Go() *Runner {
	cmder := exec.Command("sh", "-c", c.cmdline)
	c.pid = cmder.Process.Pid
	if c.print {
		cmder.Stdout = os.Stdout
		cmder.Stderr = os.Stderr
	}
	if c.wait {
		if c.waitTimeout > 0 {
			timer := time.NewTimer(c.waitTimeout)
			chDone := make(chan struct{})
			go func() {
				c.result, c.err = cmder.CombinedOutput()
				chDone <- struct{}{}
			}()
			select {
			case <-chDone:
				return (*Runner)(c)
			case <-timer.C:
				c.err = gerrors.ErrTimeout
				return (*Runner)(c)
			}
		} else {
			c.result, c.err = cmder.CombinedOutput()
		}
	} else {
		c.err = cmder.Start()
	}
	return (*Runner)(c)
}

func (c *Runner) Pid() int {
	return c.pid
}

func (c *Runner) Kill() error {
	return gproc.Terminate(gproc.ProcId(c.pid))
}

func (c *Runner) Result() ([]byte, error) {
	return c.result, c.err
}

