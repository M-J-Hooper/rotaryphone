package test

import (
	"testing"
	"time"

	"github.com/M-J-Hooper/rotaryphone"
	"github.com/brian-armstrong/gpio"
)

func TestDebouncedWatcher(t *testing.T) {
	dbw := rotaryphone.DebouncedWatcher{
		NewTestWatcher(),
		100 * time.Millisecond,
		make(map[uint]uint),
	}

	dbw.Watch()
}

type TestWatcher struct {
	notify chan gpio.WatcherNotification
}

func NewTestWatcher() *TestWatcher {
	c := make(chan gpio.WatcherNotification)
	go func(c chan gpio.WatcherNotification) {
		for i := 0; i < 20; i++ {
			sleep := time.Duration(i*10) * time.Millisecond
			time.Sleep(sleep)
			c <- gpio.WatcherNotification{1, uint(i % 2)}
		}
	}(c)
	return &TestWatcher{c}
}

func (t *TestWatcher) Watch() (uint, uint) {
	n := <-t.notify
	return n.Pin, n.Value
}

func (t *TestWatcher) AddPin(uint) {}
func (t *TestWatcher) Close()      {}
