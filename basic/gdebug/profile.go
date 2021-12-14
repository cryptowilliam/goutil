package gdebug

import (
	"bytes"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/google/pprof/driver"
	"github.com/google/pprof/profile"
	"runtime/pprof"
	"time"
)


type (
	Profile profile.Profile

	// ProfileName represents a pprof profile.
	ProfileName string
)

var (
	profileCPU          = ProfileName("cpu")
	profileHeap         = ProfileName("heap")
	profileBlock        = ProfileName("block") // Stack traces that led to blocking on synchronization primitives.
	profileMutex        = ProfileName("mutex") // Stack traces of holders of contended mutexes.
	profileAllocs       = ProfileName("allocs")
	profileGoroutine    = ProfileName("goroutine")
	profileThreadCreate = ProfileName("threadcreate") // Stack traces that led to the creation of new OS threads.
)

func (p ProfileName) String() string {
	return string(p)
}

func convertProfile(s string) (ProfileName, error) {
	switch ProfileName(s) {
	case profileCPU, profileHeap, profileBlock, profileMutex, profileAllocs, profileGoroutine, profileThreadCreate:
		return ProfileName(s), nil
	default:
		return profileCPU, gerrors.New("invalid ProfileName %s", s)
	}
}

// Capture captures profile and returns content.
func Capture(profile ProfileName, cpuCapDur time.Duration) (*Profile, error) {
	switch profile {
	case profileCPU:
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

	case profileHeap, profileBlock, profileMutex, profileAllocs, profileGoroutine, profileThreadCreate:
		f := &bytes.Buffer{}
		if err := pprof.Lookup(profile.String()).WriteTo(f, 2); err != nil {
			return nil, err
		}
		gfs.BytesToFile(f.Bytes(), fmt.Sprintf("%s-profile-%s.txt", profile.String(), time.Now().String()))
		return ParseProfile(f.Bytes())

	default:
		return nil, gerrors.New("unsupported ProfileName %s", profile.String())
	}
}

// CaptureToFile captures profile and save it to file.
func CaptureToFile(profile ProfileName, cpuCapDur time.Duration) (string, error) {
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
		return "", gerrors.New("unsupported ProfileName %s", profile.String())
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

	result := bytes.Buffer{}
	err = driver.PProf(&driver.Options{
		Fetch:   newFetcher(buf.Bytes()),
		Flagset: newFlagSet("-svg", "-output="+tempPath),
		UI:      newFakeUI(),
		Writer:  newWriter(&result),
	})
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}