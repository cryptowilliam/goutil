package ghdd

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"github.com/shirou/gopsutil/disk"
)

// https://github.com/cydev/du
// https://github.com/ricochet2200/go-disk-usage
// https://gist.github.com/lunny/9828326
// http://wendal.net/2012/1224.html
// https://github.com/lxn/win
// https://github.com/AllenDang/w32

type VolumeInfo struct {
	FilesystemType string
	Free           gvolume.Volume
	Total          gvolume.Volume
}

func GetVolumeInfo(volumePath string) (*VolumeInfo, error) {
	var vi VolumeInfo
	du, err := disk.Usage(volumePath)
	if err != nil {
		return nil, err
	}
	vi.Free, err = gvolume.FromByteSize(float64(du.Free))
	if err != nil {
		return nil, err
	}
	vi.Total, err = gvolume.FromByteSize(float64(du.Total))
	if err != nil {
		return nil, err
	}
	vi.FilesystemType = du.Fstype
	return &vi, nil
}

// ListVolumes returns partition path, same as mount point in Unix, or logical drive in Windows
func ListVolumes() (partitionPath []string, err error) {
	ps, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, gerrors.New("Get disk partitions error")
	}

	var result []string
	for _, item := range ps {
		result = append(result, item.Mountpoint)
	}
	return result, nil
}
