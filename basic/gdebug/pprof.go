package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type (
	// Profile represents a pprof profile.
	Profile string
)

var (
	ProfileCPU          = Profile("cpu")
	ProfileHeap         = Profile("heap")
	ProfileBlock        = Profile("block") // Stack traces that led to blocking on synchronization primitives.
	ProfileMutex        = Profile("mutex") // Stack traces of holders of contended mutexes.
	ProfileGoRoutine    = Profile("goroutine")
	ProfileThreadCreate = Profile("threadcreate") // Stack traces that led to the creation of new OS threads.
)

func (p Profile) String() string {
	return string(p)
}

func newTemp() (*os.File, error) {
	f, err := ioutil.TempFile("", "profile-")
	if err != nil {
		return nil, gerrors.New("Cannot create new temp profile file: %v", err)
	}
	return f, nil
}

// `blockCapRate` is the fraction of goroutine blocking events that
// are reported in the blocking profile. The profiler aims to
// sample an average of one blocking event per rate nanoseconds spent blocked.
//
// If zero value is provided, it will include every blocking event
// in the profile.
func captureProfile(profile Profile, cpuCapDur time.Duration, blockCapRate int) (string, error) {
	if profile == ProfileBlock && blockCapRate > 0 {
		runtime.SetBlockProfileRate(blockCapRate)
	}

	switch profile {
	case ProfileCPU:
		if cpuCapDur <= 0 {
			cpuCapDur = 30 * time.Second
		}
		f, err := newTemp()
		if err != nil {
			return "", err
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			return "", err
		}
		time.Sleep(cpuCapDur)
		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			return "", err
		}
		return f.Name(), nil
	case ProfileHeap, ProfileBlock, ProfileMutex, ProfileGoRoutine, ProfileThreadCreate:
		f, err := newTemp()
		if err != nil {
			return "", err
		}
		if err := pprof.Lookup(profile.String()).WriteTo(f, 2); err != nil {
			return "", err
		}
		if err := f.Close(); err != nil {
			return "", err
		}
		return f.Name(), nil
	default:
		return "", gerrors.New("unsupported profile %s", profile.String())
	}
}
