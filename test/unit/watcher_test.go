package unit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/M-J-Hooper/rotaryphone"
)

func TestDebouncedWatcher(t *testing.T) {
	watcher := NewTestWatcher(4 * time.Millisecond)
	dbw := rotaryphone.NewDebouncedWatcher(watcher.Notification, 10*time.Millisecond)

	last := -1
	for i := 0; i < 5; i++ {
		// Eventually there will be 5 stable signals
		n := dbw.Watch().(int)
		fmt.Println("Test got stable value", n)
		if last == n {
			t.Fatal("Successive signals with the same value")
		}
		last = n
	}
}

type TestWatcher struct {
	Notification chan interface{}
}

func NewTestWatcher(sleep time.Duration) *TestWatcher {
	rand.Seed(time.Now().UnixNano())
	c := make(chan interface{})
	go func(c chan interface{}) {
		for {
			time.Sleep(sleep)
			c <- interface{}(rand.Intn(3))
		}
	}(c)
	return &TestWatcher{c}
}
