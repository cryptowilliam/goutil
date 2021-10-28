package gusers

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"os/exec"
	"runtime"
)

func Logout() error {
	var err error
	if runtime.GOOS == "darwin" {
		_, err = exec.Command("sudo", "pkill", "loginwindow").CombinedOutput()
	} else if runtime.GOOS == "linux" {
		_, err = exec.Command("exit").CombinedOutput()
	} else if runtime.GOOS == "windows" {
		_, err = exec.Command("logoff").CombinedOutput()
	} else {
		err = gerrors.New("Unsupported OS " + runtime.GOOS)
	}
	return err
}

// https://github.com/noonien/pm
func Sleep() error {
	var err error
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		_, err = exec.Command("sudo", "shutdown", "-s", "now").CombinedOutput()
	} else if runtime.GOOS == "windows" {
		_, err = exec.Command("shutdown", "-h", "-t", "0").CombinedOutput()
	} else {
		err = gerrors.New("Unsupported OS " + runtime.GOOS)
	}
	return err
}

func Shutdown() error {
	var err error
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		_, err = exec.Command("sudo", "shutdown", "-h", "now").CombinedOutput()
	} else if runtime.GOOS == "windows" {
		_, err = exec.Command("shutdown", "-s", "-t", "0").CombinedOutput()
	} else {
		err = gerrors.New("Unsupported OS " + runtime.GOOS)
	}
	return err
}

// 参考 https://github.com/jpg0/rebooter
// ubuntu/centos下直接reboot就可以了，是不是因为登录用户是root？
func Reboot() error {
	var err error
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		_, err = exec.Command("sudo", "shutdown", "-r", "now").CombinedOutput()
	} else if runtime.GOOS == "windows" {
		_, err = exec.Command("shutdown", "-r", "-t", "0").CombinedOutput()
	} else {
		err = gerrors.New("Unsupported OS " + runtime.GOOS)
	}
	return err
}
