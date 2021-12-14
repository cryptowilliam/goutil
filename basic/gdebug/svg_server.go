package gdebug

import (
	"bytes"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/grand"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"github.com/google/pprof/driver"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type (
	svgServer struct {
		listenPort int
		log        glog.Interface
		selfPath   string
		useGoTool  bool
	}
)

// Note: Graphviz required
func newSvgServer(listenPort int, log glog.Interface) (*svgServer, error) {
	selfPath, err := gproc.SelfPath()
	if err != nil {
		return nil, err
	}
	return &svgServer{listenPort: listenPort, log: log, selfPath: selfPath, useGoTool: false}, nil
}

func (s *svgServer) replyError(w http.ResponseWriter, err error, wrapMsg string) {
	if err == nil {
		return
	}
	err = gerrors.Wrap(err, wrapMsg)
	s.log.Erro(err)
	if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
		s.log.Erro(errWrite)
	}
}

func (s *svgServer) serve(w http.ResponseWriter, r *http.Request) {
	s.log.Infof("accept visual pprof request %s", r.URL.String())
	ss := strings.Split(r.URL.Path, "debug/visual-pprof/")
	if len(ss) == 0 {
		err := gerrors.New("invalid path %s", r.URL.Path)
		s.replyError(w, err, "")
		return
	}
	profile := strings.ToLower(ss[len(ss)-1])
	err := error(nil)

	var imgBuf []byte
	if s.useGoTool {
		imgPath := "profile-" + grand.RandomString(10) + ".svg"
		defer os.Remove(imgPath)
		profPath, err := CaptureToFile(profile, 10*time.Second)
		if err != nil {
			s.replyError(w, err, "capture profile error")
			return
		}
		s.log.Debgf("pprof capture file: %s", profPath)

		// convert .prof file to image
		result, err := exec.Command("go", "tool", "pprof", "-svg", "-output", imgPath, s.selfPath, profPath).CombinedOutput()
		if err != nil {
			s.replyError(w, err, fmt.Sprintf("execute shell returns %s, error", result))
			return
		}
		imgBuf, err = gfs.FileToBytes(imgPath)
		if err != nil {
			s.replyError(w, err, "read svg error")
			return
		}
	} else {
		fetchUrl := fmt.Sprintf("http://localhost:%d/debug/pprof/%s", s.listenPort, profile)
		if profile != "threadcreate" {
			fetchUrl += "?debug=1"
		}
		imgBuf, err = s.httpProfileToSVG(fetchUrl)
		if err != nil {
			s.replyError(w, err, "handle svg error")
			return
		}
	}

	if _, err = w.Write(imgBuf); err != nil {
		s.log.Erro(err)
	}
}

// Note: Graphviz required
// fetch text profile content from "profileHttpUrl" and convert it to SVG image.
func (s *svgServer) httpProfileToSVG(profileHttpUrl string) ([]byte, error) {
	tempFile := filepath.Join(os.TempDir(), grand.RandomString(20))
	f, err := os.Create(tempFile)
	if err != nil {
		return nil, err
	}
	f.Close()
	defer os.Remove(tempFile)

	// Note: Graphviz required
	// "-output=****" is only used to avoid reports error, don't delete it.
	result := bytes.Buffer{}
	err = driver.PProf(&driver.Options{
		Flagset: newFlagSet("-svg", "-output="+tempFile, profileHttpUrl),
		UI:      newFakeUI(),
		Writer:  newWriter(&result),
	})
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}