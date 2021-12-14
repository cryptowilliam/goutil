package gdebug

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/grand"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	// visualizePprof uses official `go tool pprof` web UI to show visualized pprof data.
	visualizePprof struct {
		log      glog.Interface
		selfPath string
		useGoTool bool
	}
)

func newTemp() (*os.File, error) {
	f, err := ioutil.TempFile("", "profile-")
	if err != nil {
		return nil, gerrors.New("Cannot create new temp profile file: %v", err)
	}
	return f, nil
}

func newVisualizePprof(log glog.Interface) (*visualizePprof, error) {
	selfPath, err := gproc.SelfPath()
	if err != nil {
		return nil, err
	}
	return &visualizePprof{log: log, selfPath: selfPath, useGoTool: false}, nil
}

func (c *visualizePprof) replyError(w http.ResponseWriter, err error, wrapMsg string) {
	if err == nil {
		return
	}
	err = gerrors.Wrap(err, wrapMsg)
	fmt.Println("before err log")
	c.log.Erro(err)
	fmt.Println("after err log")
	if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
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

	imgPath := "profile-" + grand.RandomString(10) + ".svg"
	var imgBuf []byte
	if c.useGoTool {
		profPath, err := CaptureToFile(profile, 10*time.Second)
		if err != nil {
			c.replyError(w, err, "capture profile error")
			return
		}
		c.log.Debgf("pprof capture file: %s", profPath)

		// convert .prof file to image
		result, err := exec.Command("go", "tool", "pprof", "-svg", "-output", imgPath, c.selfPath, profPath).CombinedOutput()
		if err != nil {
			c.replyError(w, err, fmt.Sprintf("execute shell returns %s, error", result))
			return
		}
		imgBuf, err = gfs.FileToBytes(imgPath)
		if err != nil {
			c.replyError(w, err, "read svg error")
			return
		}
	} else {
		prof, err := Capture(profile, 10*time.Second)
		if err != nil {
			c.replyError(w, err, "capture profile error")
			return
		}
		imgBuf, err = prof.ToSvg()
		if err != nil {
			c.replyError(w, err, "handle svg error")
			return
		}
	}

	fmt.Println("before w.Write")
	if _, err = w.Write(imgBuf); err != nil {
		c.log.Erro(err)
	}
	fmt.Println("after w.Write")
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