package glog

type Interface interface {
	Debgf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errof(format string, a ...interface{})
	Fataf(format string, a ...interface{})
	Erro(err error, wrapMsg ...string)
	Fata(err error, wrapMsg ...string)
	AssertOk(err error, wrapMsg ...string)
	AssertTrue(express bool, wrapMsg ...string)
}
