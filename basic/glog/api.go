package glog

import (
	"github.com/cryptowilliam/goutil/sys/gtime"
	"sync/atomic"
)

var DefaultLogger = &DefaultImpl{}
var initialized atomic.Value

func init() {
	initialized.Store(false)
	Init(false)
}

func Init(saveDisk bool) {
	conf, err := DefaultConfig()
	if err != nil {
		panic(err)
		return
	}
	conf.SaveDisk = saveDisk
	if err := DefaultLogger.Init(conf); err != nil {
		panic(err)
		return
	}
	initialized.Store(true)
}

// 返回符合logs.Logger接口的实例
// 这里不可以返回logger, 因为Init/WriteMsg/Destory/Flush函数都属于*logger而不是logger,
// 所以logger不符合logs.Logger接口要求，*logger才符合

func SetClock(c gtime.Clock) {
	DefaultLogger.SetClock(c)
}

func Debgf(format string, a ...interface{}) {
	DefaultLogger.Debgf(format, a...)
}

func Infof(format string, a ...interface{}) {
	DefaultLogger.Infof(format, a...)
}

func Warnf(format string, a ...interface{}) {
	DefaultLogger.Warnf(format, a...)
}

func Errof(format string, a ...interface{}) {
	DefaultLogger.Errof(format, a...)
}

func Fataf(format string, a ...interface{}) {
	DefaultLogger.Fataf(format, a...)
}

func Erro(err error, wrapMsg ...string) {
	DefaultLogger.Erro(err, wrapMsg...)
}

func Fata(err error, wrapMsg ...string) {
	DefaultLogger.Fata(err, wrapMsg...)
}

func AssertOk(err error, wrapMsg ...string) {
	DefaultLogger.AssertOk(err, wrapMsg...)
}

func AssertTrue(express bool, wrapMsg ...string) {
	DefaultLogger.AssertTrue(express, wrapMsg...)
}
