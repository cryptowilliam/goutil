package gsingle

// test succeed under macOS but !sometimes! failed under linux: https://github.com/marcsauter/single
// test succeed under linux: https://github.com/allan-simon/go-singleinstance

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/marcsauter/single"
)

type Instance struct {
	lock *single.Single
}

func New(appName string) *Instance {
	res := Instance{}
	res.lock = single.New(appName)
	return &res
}

func (l *Instance) IsSingle() (bool, error) {
	if err := l.lock.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		return false, nil
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		return false, gerrors.Errorf("failed to acquire exclusive app lock: %v", err)
	}

	return true, nil
}

func (l *Instance) UnLock() error {
	return l.lock.TryUnlock()
}
