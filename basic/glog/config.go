package glog

import (
	"fmt"
	"github.com/cryptowilliam/goutil/sys/gmachineid"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"github.com/cryptowilliam/goutil/sys/gsysinfo"
	"path"
)

type (
	Config struct {
		SaveDisk       bool
		SaveDir        string
		FileNameFormat string
		PrintScreen    bool

		MachId  string
		AppName string
	}
)

func DefaultConfig() (*Config, error) {
	res := &Config{}

	// Disk save filename format...
	res.FileNameFormat = "2006-01-02.log" // YEAR-MONTH-DAY.log
	res.SaveDisk = true
	res.PrintScreen = true
	// Get default logs directory
	pi, err := gproc.GetProcInfo(gproc.GetPidOfMyself())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dir := gsysinfo.GetAppLogFolder(pi.Name)
	res.SaveDir = dir
	// Get machine Id.
	id, err := gmachineid.Get()
	if err != nil {
		return nil, err
	}
	res.MachId = id
	// Get app name.
	fn, err := gproc.SelfPath()
	if err != nil {
		return nil, err
	}
	res.AppName = path.Base(fn)

	return res, nil
}
