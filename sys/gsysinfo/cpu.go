package gsysinfo

import (
	"flag"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"log"
	"os"
	"strconv"
	//"github.com/klauspost/cpuid" // x86/x64 is supported only for now
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"time"
)

// Get unique serial number of CPU
func GetSerialNumber() (string, error) {
	return "Unsupported for now", nil
}

// 获取所有CPU的使用百分比，以数组返回
func GetAllUsedPercent(duration time.Duration) ([]float64, error) {
	return cpu.Percent(duration, true)
}

// 获取所有CPU的使用百分比，组合成总百分比后返回
func GetCombinedUsedPercent(duration time.Duration) (float64, error) {
	p, err := cpu.Percent(duration, false)
	if err != nil {
		return 0, err
	}
	return p[0], err
}

func GetCpuCount() int {
	return runtime.NumCPU()
}

var (
	fExe     = flag.String("e", "", "the name of the executable to watch and limit")
	fLimit   = flag.Int("l", 50, "the percent (between 1 and 100) to limit the processes CPU usage to")
	fTimeout = flag.Int("t", 0, "timeout (seconds) to exit after if there is no suitable target process (lazy mode)")
	fPid     = flag.Int("p", 0, "pid of the process")
)

func getProcesses(processes []*os.Process, targets []string) []*os.Process {
	dh, err := os.Open("/proc")
	if err != nil {
		log.Fatalf("cannot open /proc: %s", err)
	}
	defer dh.Close()
	fis, err := dh.Readdir(-1)
	if err != nil {
		log.Fatalf("cannot read /proc: %s", err)
	}
	var dst string
	if processes == nil {
		processes = make([]*os.Process, 0, len(fis))
	}
	var ok bool
	for _, fi := range fis {
		if !fi.Mode().IsDir() {
			continue
		}
		if !isAllDigit(fi.Name()) {
			continue
		}
		pid, err := strconv.Atoi(fi.Name())
		if err != nil {
			continue
		}
		if len(targets) == 0 {
			ok = true
		} else {
			if dst, err = os.Readlink("/proc/" + fi.Name() + "/exe"); err != nil {
				continue
			}
			for _, exe := range targets {
				if exe == dst {
					ok = true
					break
					//log.Printf("dst=%q =?= exe=%q", dst, exe)
				}
			}
		}
		if !ok {
			continue
		}
		p, err := os.FindProcess(pid)
		if err != nil {
			log.Printf("cannot find process %d: %s", pid, err)
		}
		processes = append(processes, p)
	}
	return processes
}

func isAllDigit(name string) bool {
	for _, c := range name {
		if c < '0' || c >= '9' {
			return false
		}
	}
	return true
}

// 控制CPU使用率，动态调整sleep时间
type DyncSleep struct {
	cpuUsage      float64 // 允许的CPU百分比
	lastSleepTime time.Duration
}

func NewDyncSleep(cpuUsage float64) (*DyncSleep, error) {
	if cpuUsage <= 0 || cpuUsage >= 100 {
		return nil, gerrors.Errorf("Invalid cpuUsage %f", cpuUsage)
	}
	return &DyncSleep{cpuUsage: cpuUsage, lastSleepTime: time.Millisecond}, nil
}

func (s *DyncSleep) Sleep() {
	used, err := GetCombinedUsedPercent(time.Second)
	if err != nil {
		time.Sleep(s.lastSleepTime)
	} else {
		if used > s.cpuUsage {
			s.lastSleepTime += time.Millisecond
		}
		if used < s.cpuUsage {
			s.lastSleepTime -= time.Millisecond

		}
		if s.lastSleepTime <= 0 {
			s.lastSleepTime = time.Millisecond
		}
		time.Sleep(s.lastSleepTime)
	}
}
