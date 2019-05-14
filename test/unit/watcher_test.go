package unit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/M-J-Hooper/rotaryphone"
	"github.com/brian-armstrong/gpio"
)

func TestDebouncedWatcher(t *testing.T) {
	watcher := NewTestWatcher(4 * time.Millisecond)
	dbw := rotaryphone.NewDebouncedWatcher(watcher, 10*time.Millisecond)

	last := 2
	for i := 0; i < 5; i++ {
		// Eventually there will be 5 stable signals
		_, value := dbw.Watch()
		v := int(value)
		fmt.Println("Test got stable value", value)
		if last == v {
			t.Fatal("Successive signals with the same value")
		}
		last = v
	}
}

type TestWatcher struct {
	notification chan gpio.WatcherNotification
}

func NewTestWatcher(sleep time.Duration) *TestWatcher {
	rand.Seed(time.Now().UnixNano())
	c := make(chan gpio.WatcherNotification)
	go func(c chan gpio.WatcherNotification) {
		for {
			time.Sleep(sleep)
			value := uint(rand.Intn(2))
			c <- gpio.WatcherNotification{Pin: 1, Value: value}
		}
	}(c)
	return &TestWatcher{c}
}

func (t *TestWatcher) Watch() (uint, uint) {
	n := <-t.notification
	return n.Pin, n.Value
}
