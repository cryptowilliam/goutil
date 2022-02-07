package gfs

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type PathInfo struct {
	Exist        bool
	IsFolder     bool
	ModifiedTime time.Time
}

const (
	InvalidFilenameCharsWindows = "\"\\:/*?<>|“”"
)

func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, gerrors.Wrap(err, fmt.Sprintf("Stat(%s)", path))
	}
	if fi.IsDir() {
		return 0, gerrors.Errorf("path(%s) is directory", path)
	}
	return fi.Size(), nil
}

func GetPathInfo(path string) (*PathInfo, error) {
	var pi PathInfo
	fi, err := os.Stat(path)
	if err == nil {
		pi.Exist = true
		pi.IsFolder = fi.IsDir()
		pi.ModifiedTime = fi.ModTime()
		return &pi, nil
	} else if err != nil && os.IsNotExist(err) {
		pi.Exist = false
		return &pi, nil
	} else {
		return &pi, err
	}
}

func FileExits(filename string) bool {
	pi, err := GetPathInfo(filename)
	if err != nil {
		return false
	}
	return !pi.IsFolder && pi.Exist
}

// Combine absolute path and relative path to get a new absolute path
func PathJoin(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}

// "/root/home/abc.txt" -> "abc.txt"
// note: this function doesn't work if file name contains '/', like "mydir/a/b.txt" and real file name is "a/b.txt"
func PathBase(path string) string {
	return filepath.Base(path)
}

// Replace illegal chars for short filename / dir name, not multi-level directory
func RefactShortPathName(path string) string {
	var illegalChars = "/\\:*\"<>|"
	for _, c := range illegalChars {
		path = strings.Replace(path, string(c), "-", -1)
	}
	return path
}
