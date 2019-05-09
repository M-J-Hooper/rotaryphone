package rotaryphone

import (
	"fmt"
	"time"

	"github.com/brian-armstrong/gpio"
)

type Watcher interface {
	Watch() (uint, uint)
	AddPin(uint)
	Close()
}

type DebouncedWatcher struct {
	Watcher
	BounceTime    time.Duration
	CurrentValues map[uint]uint
}

func NewDebouncedWatcher(BounceTime time.Duration) *DebouncedWatcher {
	return &DebouncedWatcher{
		gpio.NewWatcher(),
		BounceTime,
		make(map[uint]uint),
	}
}

func (dbw *DebouncedWatcher) Watch() (uint, uint) {
	var waiting bool
	stop := make(chan struct{})
	stable := make(chan gpio.WatcherNotification)
	notify := make(chan gpio.WatcherNotification)

	go func(notify chan gpio.WatcherNotification) {
		for {
			p, v := dbw.Watcher.Watch()
			notify <- gpio.WatcherNotification{p, v}
		}
	}(notify)

	for {
		select {
		case n := <-notify:
			fmt.Println("Notification", n.Pin, n.Value)
			if n.Value == dbw.CurrentValues[n.Pin] {
				if waiting {
					stop <- struct{}{}
					waiting = false
				}
			} else {
				if !waiting {
					go dbw.Wait(n, stable, stop)
					waiting = true
				}
			}
		case n := <-stable:
			fmt.Println("Got stable", n.Pin, n.Value)
			dbw.CurrentValues[n.Pin] = n.Value
			return n.Pin, n.Value
		}
	}
	return 0, 0
}

func (dbw *DebouncedWatcher) Wait(
	n gpio.WatcherNotification,
	stable chan gpio.WatcherNotification,
	stop chan struct{},
) {
	start := time.Now()
	for {
		select {
		case <-stop:
			fmt.Println("Got stop")
			return
		default:
			if time.Since(start) > dbw.BounceTime {
				fmt.Println("Sending stable", n.Pin, n.Value)
				stable <- n
				return
			}
		}
	}
}
