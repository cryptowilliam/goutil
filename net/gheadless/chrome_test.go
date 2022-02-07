package gheadless

import (
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/sys/gsysinfo"
	"path/filepath"
	"testing"
	"time"
)

func TestChromeHeadless_Screenshot(t *testing.T) {
	target := "https://www.yahoo.com"
	homeDir, err := gsysinfo.GetHomeDir()
	gtest.Assert(t, err)

	image, err := NewChrome().Screenshot(target, "", glog.DefaultLogger, 20 * time.Second)
	gtest.Assert(t, err)
	err = bufToFile(image, filepath.Join(homeDir, "Downloads/yahoo1.png"))
	gtest.Assert(t, err)
}

func TestChromeHeadless_GetFullHtml(t *testing.T) {
	target := "https://www.yahoo.com"
	homeDir, err := gsysinfo.GetHomeDir()
	gtest.Assert(t, err)

	fullHtml, err := NewChrome().GetFullHtml(target, "", glog.DefaultLogger, 20 * time.Second)
	gtest.Assert(t, err)
	err = bufToFile(fullHtml, filepath.Join(homeDir, "Downloads/yahoo1.txt"))
	gtest.Assert(t, err)
}

func TestChromeHeadless_DoTask(t *testing.T) {
	target := "https://www.yahoo.com"
	homeDir, err := gsysinfo.GetHomeDir()
	gtest.Assert(t, err)

	result, err := NewChrome().DoTask(target, "", []string{TaskScreenshot, TaskFullHtml}, glog.DefaultLogger, 20 * time.Second)
	gtest.Assert(t, err)
	err = bufToFile(result[TaskScreenshot], filepath.Join(homeDir, "Downloads/yahoo2.png"))
	gtest.Assert(t, err)
	err = bufToFile(result[TaskFullHtml], filepath.Join(homeDir, "Downloads/yahoo2.txt"))
	gtest.Assert(t, err)
}