package gdebug

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

type (
	// Profile represents a pprof profile.
	Profile string

	// VisualizePprof uses official `go tool pprof` web UI to show visualized pprof data.
	VisualizePprof struct {
		listen string
		log glog.Interface
	}
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

func convertProfile(s string) (Profile, error) {
	switch Profile(s) {
	case ProfileCPU, ProfileHeap, ProfileBlock, ProfileMutex, ProfileGoRoutine, ProfileThreadCreate:
		return Profile(s), nil
	default:
		return ProfileCPU, gerrors.New("invalid Profile %s", s)
	}
}

func newVisualizePprof() *VisualizePprof {
	return &VisualizePprof{}
}

func (c *VisualizePprof) serveVisualPprof(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	if len(ss) == 0 {
		err := gerrors.New("invalid path %s", r.URL.Path)
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			c.log.Erro(err)
			return
		}
	}
	profile, err := convertProfile(strings.ToLower(ss[len(ss)-1]))
	if err != nil {
		c.log.Erro(err)
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			c.log.Erro(err)
			return
		}
	}
	filePath, err := captureProfile(profile, 10 * time.Second, 0)
	if err != nil {
		c.log.Erro(err)
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			c.log.Erro(err)
			return
		}
	}

	cmd := exec.Command("go", "tool", "pprof", "-http="+c.listen, filePath)
	if err := cmd.Run(); err != nil {
		c.log.Erro(err, "Cannot start pprof UI")
	} else {
		c.log.Infof("start pprof UI at %s", c.listen)
	}
}