package gdebug

import (
	"bytes"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/google/pprof/profile"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"time"
)

type (
	Profile profile.Profile
)

func newTemp() (*os.File, error) {
	f, err := ioutil.TempFile("", "profile-")
	if err != nil {
		return nil, gerrors.New("Cannot create new temp profile file: %v", err)
	}
	return f, nil
}

// Capture captures profile and returns content.
func Capture(profile string, cpuCapDur time.Duration) (*Profile, error) {
	switch profile {
	case "profile":
		if cpuCapDur <= 0 {
			cpuCapDur = 30 * time.Second
		}
		f := &bytes.Buffer{}
		if err := pprof.StartCPUProfile(f); err != nil {
			return nil, err
		}
		time.Sleep(cpuCapDur)
		pprof.StopCPUProfile()
		gfs.BytesToFile(f.Bytes(), fmt.Sprintf("cpu-profile-%s.txt", time.Now().String()))
		return ParseProfile(f.Bytes())

	case "heap", "block", "mutex", "allocs", "goroutine", "threadcreate":
		f := &bytes.Buffer{}
		if err := pprof.Lookup(profile).WriteTo(f, 2); err != nil {
			return nil, err
		}
		gfs.BytesToFile(f.Bytes(), fmt.Sprintf("%s-profile-%s.txt", profile, time.Now().String()))
		return ParseProfile(f.Bytes())

	default:
		return nil, gerrors.New("unsupported ProfileName %s", profile)
	}
}

// CaptureToFile captures profile and save it to file.
// Note: go tool pprof required.
func CaptureToFile(profile string, cpuCapDur time.Duration) (string, error) {
	switch profile {
	case "profile":
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

	case "heap", "block", "mutex", "allocs", "goroutine", "threadcreate":
		f, err := newTemp()
		if err != nil {
			return "", err
		}
		defer f.Close()
		if err := pprof.Lookup(profile).WriteTo(f, 2); err != nil {
			return "", err
		}
		return f.Name(), nil

	default:
		return "", gerrors.New("unsupported ProfileName %s", profile)
	}
}

// ParseProfile parses buffer into Profile
func ParseProfile(b []byte) (*Profile, error) {
	result, err := profile.ParseData(b)
	if err != nil {
		return nil, err
	}
	return (*Profile)(result), nil
}

// VerifyProfile verify if b is valid profile content.
func VerifyProfile(b []byte) error {
	_, err := profile.ParseData(b)
	return err
}

/*
// ToDotGraph convert profile to dot graph,
// which is an image format created by Graphviz.
func (p *Profile) ToDotGraph() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := (*profile.Profile)(p).Write(&buf); err != nil {
		return nil, err
	}

	result := bytes.Buffer{}
	err := driver.PProf(&driver.Options{
		Fetch:   &fetcher{b: buf.Bytes()},
		Flagset: newFlagSet("-dot"),
		UI:      &fakeUI{},
		Writer:  &writer{&result},
	})
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

// ToSvg convert profile to SVG image.
// FIXME: output image is totally different from "go tool pprof -svg -output imagePath binaryPath profilePath"
func (p *Profile) ToSvg() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := (*profile.Profile)(p).Write(&buf); err != nil {
		return nil, err
	}

	tempFile, err := newTemp()
	if err != nil {
		return nil, err
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	selfPath, err := gproc.SelfPath()
	if err != nil {
		return nil, err
	}

	result := bytes.Buffer{}
	err = driver.PProf(&driver.Options{
		Fetch:   newFetcher(buf.Bytes()),
//-output=*** is necessary, it no output argument, will report error
		Flagset: newFlagSet("-svg", "-output="+tempPath, "-source_path="+selfPath),
		UI:      newFakeUI(),
		Writer:  newWriter(&result),
	})
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
*/
