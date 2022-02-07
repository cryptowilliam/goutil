package gproc

// https://github.com/fastly/gopkg/blob/master/executable/executable.go
// https://github.com/crquan/coremem/blob/master/coremem.go
// https://github.com/janimo/memchart/blob/master/memchart.go

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/kardianos/osext"
	"github.com/shirou/gopsutil/process"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// on linux
// 32768 by default, you can read the value on your system in /proc/sys/kernel/pid_max,
// and you can set the value higher (up to 32768 for 32 bit systems or 4194304 for 64 bit) with:
// echo 4194303 > /proc/sys/kernel/pid_max
// on windows
// process id is a DWORD value, so max id value is int32 max value 4294967295
type ProcId int32

const (
	InvalidProcId ProcId = -1
)

type ProcInfo struct {
	Name         string // short filename
	Filename     string // full file path
	Cmdline      string // run command line
	Param        string // run param
	MemUsedBytes uint64
}

func GetAllPids() ([]ProcId, error) {
	pids, err := process.Pids()
	var result []ProcId
	for _, pid := range pids {
		result = append(result, ProcId(pid))
	}
	return result, err
}

func GetPidOfMyself() ProcId {
	pid := os.Getpid()
	return ProcId(pid)
}

func GetPidByProcFullFilename(filename string) ([]ProcId, error) {
	allIds, err := GetAllPids()
	if err != nil {
		return nil, err
	}
	r := []ProcId{}
	for _, id := range allIds {
		info, err := GetProcInfo(id)
		if err != nil {
			return nil, err
		}
		if strings.ToLower(info.Filename) == strings.ToLower(filename) {
			r = append(r, id)
		}
	}
	return r, nil
}

func GetPidByProcName(procName string) ([]ProcId, error) {
	if procName == "" {
		return nil, gerrors.New("empty process name")
	}
	allIds, err := GetAllPids()
	if err != nil {
		return nil, err
	}
	var r []ProcId
	for _, id := range allIds {
		info, err := GetProcInfo(id)
		if err != nil {
			return nil, err
		}
		if strings.ToLower(info.Name) == strings.ToLower(procName) {
			r = append(r, id)
		}
	}
	return r, nil
}

func GetProcCreateTime(procName string) (map[ProcId]time.Time, error) {
	r := make(map[ProcId]time.Time)

	ids, err := GetPidByProcName(procName)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		tm, err := GetPidCreateTime(id)
		if err != nil {
			return nil, err
		}
		r[id] = tm
	}

	return r, nil
}

func GetPidCreateTime(pid ProcId) (time.Time, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return time.Time{}, err
	}
	epochMillis, err := proc.CreateTime()
	if err != nil {
		return time.Time{}, err
	}

	// TODO: fix import cycle
	// time.Unix(0, epochMillis * 1000000) equals clock.EpochMillisToTime(epochMillis)
	return time.Unix(0, epochMillis*1000000), nil
}

// 替换自己的二进制文件
func ReplaceMySelfFile(newfilepath string) error {
	return nil
}

// 还可以参考：https://github.com/rcrowley/goagain/blob/master/goagain.go#L77
// https://stackoverflow.com/questions/68201595/how-to-restart-itself-in-go-daemon-process
// Restart current process, with same parameters.
func RestartMyself() error {
	argv0, err := SelfPath()
	if err != nil {
		return err
	}
	files := make([]*os.File, syscall.Stderr+1)
	files[syscall.Stdin] = os.Stdin
	files[syscall.Stdout] = os.Stdout
	files[syscall.Stderr] = os.Stderr
	wd, err := os.Getwd()
	if nil != err {
		return err
	}
	_, err = os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})
	os.Exit(0)
	return err
}

// Notice:
// _, b, _, _ := runtime.Caller(0)
// return filepath.Dir(b)
// this is wrong
//
// Get process file folder, not working folder
func SelfDir() (string, error) {
	p, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}

func SelfPath() (string, error) {
	return osext.Executable()
}

// returns last element of path: short filename
func SelfBase() (string, error) {
	return osext.Executable()
}

func Terminate(pid ProcId) error {
	if pid < 0 {
		return gerrors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}
	return proc.Terminate()
}

// 隐藏进程,通过ps等命令看不到进程信息, 一般是通过把pid改为0做到的,windows下应该有api
func Hide(pid ProcId) error {
	if pid < 0 {
		return gerrors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}
	return nil
}

func Show(pid ProcId) error {
	if pid < 0 {
		return gerrors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}
	return nil
}

func GetProcInfo(pid ProcId) (*ProcInfo, error) {
	if pid < 0 {
		return nil, gerrors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}

	var pi ProcInfo
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return &pi, err
	}
	pi.Name, err = proc.Name()
	if err != nil {
		pi.Name = ""
	}
	cmdline, err := proc.Cmdline()
	pi.Cmdline = cmdline

	// get path
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		path, err := GetExePathFromPid(int(pid))
		if err != nil {
			return nil, err
		}
		pi.Filename = path
	} else {
		pi.Filename, err = proc.Exe()
	}

	pi.Param = strings.Replace(cmdline, pi.Filename, "", 1)
	pi.Param = strings.TrimSpace(pi.Param)
	mi, err := proc.MemoryInfo()
	if err != nil {
		pi.MemUsedBytes = 0
	} else {
		pi.MemUsedBytes = mi.Swap
	}
	return &pi, nil
}
