package gfs

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"time"
)

type (
	ChangedEvent struct {
		Path  string
		Event WatchEvent
		Err   error
	}

	FsWatcher struct {
		wr *watcher.Watcher
	}
)

// An WatchEvent is a type that is used to describe what type
// of event has occurred during the watching process.
type WatchEvent uint32

// Ops
const (
	Create = WatchEvent(watcher.Create)
	Write  = WatchEvent(watcher.Write)
	Remove = WatchEvent(watcher.Remove)
	Rename = WatchEvent(watcher.Rename)
	Chmod  = WatchEvent(watcher.Chmod)
	Move   = WatchEvent(watcher.Move)
)

func NewWatcher() *FsWatcher {
	r := new(FsWatcher)
	r.wr = watcher.New()
	return r
}

// this function is block
func (w *FsWatcher) Loop(path string, evts WatchEvent, interval time.Duration, notifyCh chan ChangedEvent) (err error) {
	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	//w.SetMaxEvents(1)

	// Only notify write events.
	w.wr.FilterOps(watcher.Op(evts))

	go func() {
		for {
			select {
			case event := <-w.wr.Event:
				fmt.Println(event) // Print the event's info.
				notifyCh <- ChangedEvent{Path: event.Path, Event: WatchEvent(event.Op)}
				fmt.Println("ok")
			case err := <-w.wr.Error:
				notifyCh <- ChangedEvent{Err: err}
				return
			case <-w.wr.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.wr.Add(path); err != nil {
		return err
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.wr.Start(interval); err != nil {
		return err
	}
	return nil
}
