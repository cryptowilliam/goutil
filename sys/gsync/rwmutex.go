package gsync

import "github.com/jonhoo/drwmutex"

// warning: test required!

/*
The default Go implementation of sync.RWMutex does not scale well to multiple cores,
as all readers contend on the same memory location when they all try to atomically increment it.
This repository provides an n-way RWMutex, also known as a "big reader" lock,
which gives each CPU core its own RWMutex. Readers take only a read lock local to their core,
whereas writers must take all locks in order.
*/

type RWMutex drwmutex.DRWMutex

func (rm RWMutex) RLock() {
	(drwmutex.DRWMutex)(rm).RLock()
}

func (rm RWMutex) Lock() {
	(drwmutex.DRWMutex)(rm).Lock()
}

func (rm RWMutex) Unlock() {
	(drwmutex.DRWMutex)(rm).Unlock()
}
