package gcmd

import (
	"fmt"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"os"
	"os/exec"
	"strings"
	"time"
)

// sudo 下执行ExecWait或者ExecNoWait接口时会有问题

type (
	Cmder struct {
		cmd *exec.Cmd
	}
)

func StartCmder(screenPrint bool, cmd string, arg ...string) *Cmder {
	var cmder Cmder
	cmder.cmd = exec.Command(cmd, arg...)

	fmt.Println("command line:")
	fmt.Println(cmd, strings.Join(append([]string{}, arg...), " "))
	if screenPrint {
		cmder.cmd.Stdout = os.Stdout
		cmder.cmd.Stderr = os.Stderr
		if err := cmder.cmd.Start(); err != nil {
			fmt.Println(err.Error())
		}
	}
	return &cmder
}

func (c *Cmder) GetPid() int32 {
	return int32(c.cmd.Process.Pid)
}

func (c *Cmder) Wait() error {
	return c.cmd.Wait()
}

func (c *Cmder) WaitWithTimeout(timeout time.Duration) (execTimeout bool, execError error) {
	tmr := time.NewTimer(timeout)
	err := error(nil)
	chDone := make(chan struct{})

	go func() {
		err = c.cmd.Wait()
		chDone <- struct{}{}
	}()

	select {
	case <-chDone:
		return false, err
	case <-tmr.C:
		return true, nil
	}
}

func (c *Cmder) Kill() error {
	return gproc.Terminate(gproc.ProcId(c.GetPid()))
}

func (c *Cmder) GetResultOutput() (string, error) {
	out, err := c.cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// will print screen
func ExecWaitPrintScreen(name string, arg ...string) error {
	var cmd *exec.Cmd

	cmd = exec.Command(name, arg...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// will not print screen
func ExecWaitReturn(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}

func ExecShell(shellCommand string) ([]byte, error) {
	return exec.Command("sh", "-c", shellCommand).CombinedOutput()
}