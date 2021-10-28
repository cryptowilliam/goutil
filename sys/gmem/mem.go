package gmem

import (
	"fmt"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"github.com/shirou/gopsutil/mem"
)

type MemUsage struct {
	TotalBytes     uint64 `json:"total"`
	UsedBytes      uint64 `json:"used"`
	AvailableBytes uint64 `json:"Available"`
	SelfBytes      uint64 `json:"self"`
}

func (mu MemUsage) String() string {
	total, err := gvolume.FromByteSize(float64(mu.TotalBytes))
	totalString := total.String()
	if err != nil {
		totalString = err.Error()
	}
	used, err := gvolume.FromByteSize(float64(mu.UsedBytes))
	usedString := used.String()
	if err != nil {
		usedString = err.Error()
	}
	available, err := gvolume.FromByteSize(float64(mu.AvailableBytes))
	availableString := available.String()
	if err != nil {
		availableString = err.Error()
	}
	return fmt.Sprintln(
		"Total:", totalString,
		"Used:", usedString,
		"Available:", availableString,
	)
}

func GetUsage() (MemUsage, error) {
	var mu MemUsage
	vms, err := mem.VirtualMemory()
	if err != nil {
		return mu, err
	}
	mu.TotalBytes = vms.Total
	mu.UsedBytes = vms.Used
	mu.AvailableBytes = vms.Available
	return mu, nil
}

/* Custom implement for Unix
func getUsage() MemUsage {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	mem := MemUsage{}
	mem.Self = memStat.Alloc // usage of current process
	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err == nil {
		mem.All = sysInfo.Totalram * uint32(syscall.Getpagesize())
		mem.Free = sysInfo.Freeram * uint32(syscall.Getpagesize())
		mem.Used = mem.All - mem.Free
	}
	return mem
}
*/
