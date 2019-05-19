package rotaryphone

import (
	"fmt"
	"time"

	"github.com/brian-armstrong/gpio"
)

//3.3v  => white wire
const latchActivePin = 18 // => dark green wire

type Latch struct {
	Active chan bool
}

func NewLatch() *Latch {
	l := &Latch{make(chan bool)}
	go l.run()
	return l
}

func (l Latch) run() {
	pinWatcher := gpio.NewWatcher()
	pinWatcher.AddPin(latchActivePin)
	defer pinWatcher.Close()

	c := make(chan interface{})
	go castGpioChannel(pinWatcher.Notification, c)

	watcher := NewDebouncedWatcher(c, 100*time.Millisecond)

	for {
		n := watcher.Watch().(gpio.WatcherNotification)
		pin, value := n.Pin, n.Value

		fmt.Println("Latch got stable", pin, value)
		if pin == latchActivePin {
			if value == 1 {
				l.Active <- true
			} else {
				l.Active <- false
			}
		}
	}
}
