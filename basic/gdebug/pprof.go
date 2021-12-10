package gdebug

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/sys/gcmd"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

type (
	// Profile represents a pprof profile.
	Profile string

	// VisualizePprof uses official `go tool pprof` web UI to show visualized pprof data.
	VisualizePprof struct {
		log glog.Interface
		historyPidList []int
		mu sync.RWMutex
		selfPath string
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

func newVisualizePprof(log glog.Interface) (*VisualizePprof, error) {
	selfPid := gproc.GetPidOfMyself()
	selfPath, err := gproc.GetExePathFromPid(int(selfPid))
	if err != nil {
		return nil, err
	}
	return &VisualizePprof{log: log, selfPath: selfPath}, nil
}

func (c *VisualizePprof) replyError(w http.ResponseWriter, err error, wrapMsg string) {
	if err  == nil {
		return
	}
	err = gerrors.Wrap(err, wrapMsg)
	if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
		c.log.Erro(err)
		c.log.Erro(errWrite)
	}
}

func (c *VisualizePprof) serveVisualPprof(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, v := range c.historyPidList {
		_ = gproc.Terminate(gproc.ProcId(v))
	}

	c.log.Infof("accept visual pprof request %s", r.URL.String())
	ss := strings.Split(r.URL.Path, "debug/visual-pprof/")
	if len(ss) == 0 {
		err := gerrors.New("invalid path %s", r.URL.Path)
		c.replyError(w, err, "")
		return
	}
	profile, err := convertProfile(strings.ToLower(ss[len(ss)-1]))
	if err != nil {
		c.replyError(w, err, "convert profile error")
		return
	}
	profPath, err := captureProfile(profile, 10 * time.Second, 0)
	if err != nil {
		c.replyError(w, err, "capture profile error")
		return
	}
	svgPath := profPath+".svg"

	cmdline := fmt.Sprintf("go tool pprof -svg '%s' '%s' > '%s'", c.selfPath, profPath, svgPath)
	_, err = gcmd.ExecShell(cmdline)
	if err != nil {
		c.replyError(w, err, "execute shell error")
		return
	}
	svgStr, err := gfs.FileToString(svgPath)
	if err != nil {
		c.replyError(w, err, "read svg error")
		return
	}
	svgHtml, err := gstring.SubstrBetween(svgStr, "<svg", "/svg>", true, false, true, true)
	if err != nil {
		c.replyError(w, err, "handle svg error")
		return
	}
	htmlTemplate := `<html>
<head>
    <meta charset="UTF-8">
    <title></title>
    <style>
        * {
            margin: 0;
            padding: 0;
        }
        html, body, .fullpage {
            width: 100%;
            height: 100%;
        }
        .fullpage {
            color: white;
            font-size: 35px;
            text-align: center;
        }
    </style>
</head>
<body>
	<div class="fullpage">%s</div>
</body>
</html>`
	htmlSrc := strings.Replace(htmlTemplate, "%s", svgHtml, 1)
	_, err = w.Write([]byte(htmlSrc))
	if err != nil {
		c.log.Erro(err)
	}

	// Use go tool inside http UI server, it is more powerful but hard to manage,
	// maybe it will be enabled later.
	/*cmd := exec.Command("go", "tool", "pprof", "-http="+c.listen, filePath)
	if err := cmd.Run(); err != nil {
		err = gerrors.Wrap(err, "start pprof UI error")
		c.log.Erro(err)
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			c.log.Erro(err)
		}
	} else {
		c.historyPidList = append(c.historyPidList, cmd.Process.Pid)
		info := fmt.Sprintf("start pprof UI at %s", c.listen)
		c.log.Infof(info)
		if _, errWrite := w.Write([]byte(info)); errWrite != nil {
			c.log.Erro(err)
		}
	}*/
}