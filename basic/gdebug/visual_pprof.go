package gdebug

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/crypto/gbase"
	"github.com/cryptowilliam/goutil/sys/gcmd"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

type (
	// profile represents a pprof profile.
	profile string

	// visualizePprof uses official `go tool pprof` web UI to show visualized pprof data.
	visualizePprof struct {
		log      glog.Interface
		selfPath string
	}
)

var (
	profileCPU          = profile("cpu")
	profileHeap         = profile("heap")
	profileBlock        = profile("block") // Stack traces that led to blocking on synchronization primitives.
	profileMutex        = profile("mutex") // Stack traces of holders of contended mutexes.
	profileAllocs       = profile("allocs")
	profileGoroutine    = profile("goroutine")
	profileThreadCreate = profile("threadcreate") // Stack traces that led to the creation of new OS threads.
)

func (p profile) String() string {
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
func captureProfile(profile profile, cpuCapDur time.Duration, blockCapRate int) (string, error) {
	if profile == profileBlock && blockCapRate > 0 {
		runtime.SetBlockProfileRate(blockCapRate)
	}

	switch profile {
	case profileCPU:
		if cpuCapDur <= 0 {
			cpuCapDur = 30 * time.Second
		}
		f, err := newTemp()
		if err != nil {
			return "", err
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			return "", err
		}
		time.Sleep(cpuCapDur)
		pprof.StopCPUProfile()
		return f.Name(), nil

	case profileHeap, profileBlock, profileMutex, profileAllocs, profileGoroutine, profileThreadCreate:
		f, err := newTemp()
		if err != nil {
			return "", err
		}
		defer f.Close()
		if err := pprof.Lookup(profile.String()).WriteTo(f, 2); err != nil {
			return "", err
		}
		return f.Name(), nil

	default:
		return "", gerrors.New("unsupported profile %s", profile.String())
	}
}

func convertProfile(s string) (profile, error) {
	switch profile(s) {
	case profileCPU, profileHeap, profileBlock, profileMutex, profileAllocs, profileGoroutine, profileThreadCreate:
		return profile(s), nil
	default:
		return profileCPU, gerrors.New("invalid profile %s", s)
	}
}

func newVisualizePprof(log glog.Interface) (*visualizePprof, error) {
	selfPid := gproc.GetPidOfMyself()
	selfPath, err := gproc.GetExePathFromPid(int(selfPid))
	if err != nil {
		return nil, err
	}
	return &visualizePprof{log: log, selfPath: selfPath}, nil
}

func (c *visualizePprof) replyError(w http.ResponseWriter, err error, wrapMsg string) {
	if err == nil {
		return
	}
	err = gerrors.Wrap(err, wrapMsg)
	if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
		c.log.Erro(err)
		c.log.Erro(errWrite)
	}
}

func (c *visualizePprof) serveVisualPprof(w http.ResponseWriter, r *http.Request) {
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
	profPath, err := captureProfile(profile, 10*time.Second, 0)
	if err != nil {
		c.replyError(w, err, "capture profile error")
		return
	}
	c.log.Debgf("pprof capture file: %s", profPath)

	// convert .prof file to image
	imgType := "png"
	imgPath := profPath + "." + imgType
	cmdline := fmt.Sprintf("go tool pprof -%s '%s' '%s' > '%s'", imgType, c.selfPath, profPath, imgPath)
	result, err := gcmd.ExecShell(cmdline)
	if err != nil {
		c.replyError(w, err, fmt.Sprintf("execute shell returns %s, error", result))
		return
	}
	imgBuf, err := gfs.FileToBytes(imgPath)
	if err != nil {
		c.replyError(w, err, "read svg error")
		return
	}

	// convert image file to html source
	imgHtml := ""
	if imgType == "svg" {
		imgHtml, err = gstring.SubstrBetween(string(imgBuf), "<svg", "/svg>", true, false, true, true)
		if err != nil {
			c.replyError(w, err, "handle svg error")
			return
		}
	} else if imgType == "png" {
		imgHtml = `<img src="data:image/png;base64,` + gbase.Base64Encode(imgBuf) + `" />`
	} else {
		err = gerrors.New("unknown image type %s", imgType)
		if err != nil {
			c.replyError(w, err, "handle svg error")
			return
		}
	}

	// insert image into html
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
	htmlSrc := strings.Replace(htmlTemplate, "%s", imgHtml, 1)
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
