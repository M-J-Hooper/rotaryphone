package unit

import (
	"math/rand"
	"testing"
	"time"

	"github.com/M-J-Hooper/rotaryphone"
	"github.com/brian-armstrong/gpio"
)

func TestDebouncedWatcher(t *testing.T) {
	dbw := rotaryphone.DebouncedWatcher{
		Watcher:       NewTestWatcher(),
		BounceTime:    100 * time.Millisecond,
		CurrentValues: make(map[uint]uint),
	}
	// Eventually there will be 5 stable signals
	for i := 0; i < 5; i++ {
		dbw.Watch()
	}
}

type TestWatcher struct {
	notify chan gpio.WatcherNotification
}

func NewTestWatcher() *TestWatcher {
	rand.Seed(time.Now().UnixNano())
	c := make(chan gpio.WatcherNotification)
	go func(c chan gpio.WatcherNotification) {
		for {
			sleep := 40 * time.Millisecond
			time.Sleep(sleep)

			value := uint(rand.Intn(2))
			c <- gpio.WatcherNotification{Pin: 1, Value: value}
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
