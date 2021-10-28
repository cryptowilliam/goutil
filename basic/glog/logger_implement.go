package glog

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"github.com/sttts/color"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	DefaultImpl struct {
		clock           gtime.Clock
		conf            *Config
		currLogFilename string
		currLogFile     *os.File
		currLogFileMu   sync.Mutex
		printMu         sync.Mutex
	}
)

func NewInsideLogger(c gtime.Clock) *DefaultImpl {
	return &DefaultImpl{clock: c}
}

func (lgz *DefaultImpl) SetClock(c gtime.Clock) {
	lgz.clock = c
}

func (lgz *DefaultImpl) Logging(log LogItem) error {
	return nil
}

func (lgz *DefaultImpl) LoggingEx(level Level, text string, tags map[string]string) error {
	return nil
}

// Create or update log file.
func (lgz *DefaultImpl) getFile(tm time.Time) (*os.File, error) {
	if len(lgz.conf.SaveDir) == 0 {
		return nil, gerrors.New("Empty log dierctory")
	}

	var err error
	var newLogFilename string
	newLogFilename = lgz.conf.SaveDir + "/" + tm.Format(lgz.conf.FileNameFormat)

	// Update log filename and log file descriptor
	if lgz.currLogFilename != newLogFilename {
		if lgz.currLogFile != nil {
			lgz.currLogFile.Close()
			lgz.currLogFile = nil
		}

		for {
			pi, err := gfs.GetPathInfo(newLogFilename)
			if err != nil {
				return nil, err
			}
			if !pi.Exist {
				lgz.currLogFile, err = os.Create(newLogFilename)
				if err != nil {
					return nil, err
				}
			}
			if pi.Exist && pi.IsFolder {
				err := os.RemoveAll(newLogFilename)
				if err != nil {
					return nil, err
				} else {
					continue
				}
			}
			break
		}

		lgz.currLogFilename = newLogFilename
	}

	// O log file
	if lgz.currLogFile == nil {
		lgz.currLogFile, err = os.OpenFile(lgz.currLogFilename, os.O_RDWR|os.O_CREATE, 0755)
		lgz.currLogFile.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, err
		}
	}
	return lgz.currLogFile, nil
}

// 注意，这里的receiver必须用*DefaultImpl，不可以用logger，否则conf将无法保存进l里面去
func (lgz *DefaultImpl) Init(config *Config) error {
	err := error(nil)

	if config == nil {
		config, err = DefaultConfig()
		if err != nil {
			return err
		}
	}

	lgz.conf = config
	lgz.clock = gtime.GetSysClock()

	if lgz.conf.SaveDisk {
		fmt.Println(fmt.Sprintf("items logging into %s", lgz.conf.SaveDir))

		// 检查日志输出文件夹是否正常
		if err := os.MkdirAll(lgz.conf.SaveDir, os.ModePerm); err != nil {
			return err
		}
		pi, err := gfs.GetPathInfo(lgz.conf.SaveDir)
		if err != nil {
			return err
		}
		if !pi.Exist {
			return gerrors.New(lgz.conf.SaveDir + " create failed")
		}
		if !pi.IsFolder {
			return gerrors.New(lgz.conf.SaveDir + " create failed because it is not a folder")
		}
	}
	return nil
}

// Output to disk and screen if user want it
func (lgz *DefaultImpl) WriteMsg(when time.Time, msg string, level Level) error {
	msg = lgz.clock.Now().Format("2006-01-02 15:04:05.000 -07 [") + string(level) + "] " + msg

	if lgz.conf.SaveDisk {
		f, err := DefaultLogger.getFile(when)
		if err != nil {
			return gerrors.Wrap(err, "WriteMsg")
		}

		// Output log file
		DefaultLogger.currLogFileMu.Lock()
		if f != nil {
			_, err := f.Write([]byte(msg + "\n"))
			if err != nil {
				DefaultLogger.currLogFileMu.Unlock()
				return err
			}
		}
		DefaultLogger.currLogFileMu.Unlock()

	}

	// Print screen
	if lgz.conf.PrintScreen {
		lgz.printMu.Lock()
		switch level {
		case LevelDebg:
			color.Println(color.Green(msg))
		case LevelInfo:
			color.Println(color.Cyan(msg))
		case LevelWarn:
			color.Println(color.Yellow(msg))
		case LevelErro:
			color.Println(color.Red(msg))
		case LevelFata:
			color.Println(color.Magenta(msg))
		}
		lgz.printMu.Unlock()
	}

	return nil
}

func (lgz *DefaultImpl) Destroy() {
	lgz.currLogFileMu.Lock()
	if lgz.currLogFile != nil {
		lgz.currLogFile.Sync()
		lgz.currLogFile.Close()
	}
	lgz.currLogFileMu.Unlock()
}

func (lgz *DefaultImpl) Flush() {
	f, err := DefaultLogger.getFile(lgz.clock.Now())
	if err == nil {
		DefaultLogger.currLogFileMu.Lock()
		if f != nil {
			f.Sync()
		}
		DefaultLogger.currLogFileMu.Unlock()
	}
}

func clear(Text string) string {
	Text = strings.Trim(Text, "\n")
	lines := strings.Split(Text, "\n")
	var result []string
	removeTags := []string{
		"runtime/asm_amd64",
		"runtime.goexit",
		"support/xerror",
		"runtime.main",
	}
	for _, ln := range lines {
		jump := false
		for _, tag := range removeTags {
			if strings.Contains(ln, tag) {
				jump = true
				break
			}
		}
		if jump {
			continue
		}
		result = append(result, ln)
	}

	return strings.Join(result, "\n")
}

func (lgz *DefaultImpl) logging(message string, level Level) {
	message = clear(message)
	if initialized.Load().(bool) {
		now := lgz.clock.Now()

		// Output logs to disk and screen if user want it.
		if err := lgz.WriteMsg(now, message, level); err != nil {
			fmt.Println(err)
		}

		// Output to channel.
		/*if DefaultLogger.out != nil && len(DefaultLogger.out) < cap(DefaultLogger.out) {
			item := Log{}
			item.T = now
			item.Level = level.String()
			item.Message = message
			DefaultLogger.out <- item
		}*/
	} else {
		fmt.Println("xlog inited flag false, please call xlog.Init first")
	}
}

func (lgz *DefaultImpl) Debgf(format string, a ...interface{}) {
	lgz.logging(fmt.Sprintf(format, a...), LevelDebg)
}

func (lgz *DefaultImpl) Infof(format string, a ...interface{}) {
	lgz.logging(fmt.Sprintf(format, a...), LevelInfo)
}

func (lgz *DefaultImpl) Warnf(format string, a ...interface{}) {
	lgz.logging(fmt.Sprintf(format, a...), LevelWarn)
}

func (lgz *DefaultImpl) Errof(format string, a ...interface{}) {
	lgz.logging(fmt.Sprintf(format, a...), LevelErro)
}

func (lgz *DefaultImpl) Fataf(format string, a ...interface{}) {
	lgz.logging(fmt.Sprintf(format, a...), LevelFata)
}

func (lgz *DefaultImpl) Erro(err error, wrapMsg ...string) {
	stack := gerrors.GetStack(err)
	if len(wrapMsg) > 0 {
		lgz.logging(strings.Join(append([]string{}, wrapMsg...), ",")+": "+err.Error()+"\nstack: "+stack, LevelErro)
	} else {
		lgz.logging(err.Error()+"\nstack: "+stack, LevelErro)
	}
}

func (lgz *DefaultImpl) Fata(err error, wrapMsg ...string) {
	if len(wrapMsg) > 0 {
		lgz.logging(strings.Join(append([]string{}, wrapMsg...), ",")+": "+err.Error(), LevelFata)
	} else {
		lgz.logging(err.Error(), LevelFata)
	}
}

func (lgz *DefaultImpl) AssertOk(err error, wrapMsg ...string) {
	if err != nil {
		lgz.Erro(err, wrapMsg...)
		os.Exit(-1)
	}
}

func (lgz *DefaultImpl) AssertTrue(express bool, wrapMsg ...string) {
	if !express {
		lgz.Erro(gerrors.Errorf("express MUST be true"), wrapMsg...)
		os.Exit(-1)
	}
}
