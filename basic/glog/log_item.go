package glog

import (
	"strings"
	"time"
)

type (
	Level string

	LogItem struct {
		Time  time.Time
		Level Level
		Text  string
		Tags  map[string]string `json:"ExtTags,omitempty" bson:"ExtTags,omitempty"`
	}
)

const (
	LevelDebg Level = "DEBG"
	LevelInfo Level = "INFO"
	LevelWarn Level = "WARN"
	LevelErro Level = "ERRO"
	LevelFata Level = "FATA"
)

func NewLogItem(level Level, time time.Time, text string) LogItem {
	return LogItem{Level: level, Time: time, Text: text}
}

func (l *LogItem) SetExtTag(key, value string) *LogItem {
	if l.Tags == nil {
		l.Tags = make(map[string]string)
	}
	l.Tags[key] = value
	return l
}

func (l *LogItem) GetExtTagEx(key string) (string, bool) {
	if l.Tags == nil {
		return "", false
	}
	value, ok := l.Tags[key]
	if !ok {
		return "", false
	}
	return value, true
}

func (l *LogItem) GetExtTag(key string) string {
	value, ok := l.GetExtTagEx(key)
	if !ok {
		return ""
	}
	return value
}

// 清理日志中不需要出现的信息
func (l *LogItem) clear() {
	l.Text = strings.Trim(l.Text, "\n")
	lines := strings.Split(l.Text, "\n")
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
		/*if strings.Contains(ln, "runtime.main") {
			break
		}*/
		ln = strings.Replace(ln, "/usr/local/go/", "$ROOT/", -1)
		ln = strings.Replace(ln, "github.com", "gitlab.com", -1)
		ln = strings.Replace(ln, "bitbucket.org", "gitlab.com", -1)
		ln = strings.Replace(ln, "golang.org", "", -1)
		ln = strings.Replace(ln, "gopkg.in", "", -1)
		ln = strings.Replace(ln, "v2ray.com", "", -1)
		ln = strings.Replace(ln, "go", "c", -1)
		ln = strings.Replace(ln, "taci", "rafw", -1)
		ln = strings.Replace(ln, "cryptowilliam/goutil", "rafael", -1)
		ln = strings.Replace(ln, "kcp", "qtp", -1)
		result = append(result, ln)
	}

	l.Text = strings.Join(result, "\n")
}

