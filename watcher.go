package rotaryphone

import (
	"fmt"
	"time"
)

type Watcher interface {
	Watch() (uint, uint)
}

type Event struct {
	Time  time.Time
	Key   uint
	Value uint
}

type DebouncedWatcher struct {
	Watcher
	BounceTime   time.Duration
	notification chan Event
	stop         chan Event
	notified     map[uint]Event
	reported     map[uint]Event
}

func NewDebouncedWatcher(watcher Watcher, bounceTime time.Duration) *DebouncedWatcher {
	dbw := DebouncedWatcher{
		Watcher:      watcher,
		BounceTime:   bounceTime,
		notified:     make(map[uint]Event),
		reported:     make(map[uint]Event),
		notification: make(chan Event),
		stop:         make(chan Event),
	}
	go dbw.notify()
	go dbw.wait()
	return &dbw
}

func (dbw *DebouncedWatcher) notify() {
	for {
		k, v := dbw.Watcher.Watch()
		fmt.Println("Run notification", k, v, dbw.notified)
		dbw.notification <- Event{
			Time:  time.Now(),
			Key:   k,
			Value: v,
		}
	}
}

func (dbw *DebouncedWatcher) Watch() (uint, uint) {
	e := <-dbw.stop
	return e.Key, e.Value
}

func (dbw *DebouncedWatcher) wait() {
	for {
		select {
		case e := <-dbw.notification:
			fmt.Println("Got notification with current", e)
			curr, ok := dbw.notified[e.Key]
			if !ok || e.Value != curr.Value {
				dbw.notified[e.Key] = e
			}
		default:
			for k, n := range dbw.notified {
				r, ok := dbw.reported[k]
				if !ok || r.Value != n.Value {
					if time.Since(n.Time) > dbw.BounceTime {
						fmt.Println("Sending stable at", time.Now())
						dbw.reported[k] = n
						dbw.stop <- n
					}
				}
			}
		}
	}
}
