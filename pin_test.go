package rotaryphone

import (
	"fmt"
	"testing"
	"time"

	"github.com/brian-armstrong/gpio"
)

func TestPins(t *testing.T) {
	watcher := gpio.NewWatcher()
	watcher.AddPin(dialActivePin)
	watcher.AddPin(dialIncrementPin)
	watcher.AddPin(latchActivePin)

	for n := range watcher.Notification {
		fmt.Println("Got notification", n)
	}
}

func TestDebouncedPins(t *testing.T) {
	d := 100 * time.Millisecond

	ac := debouncedPin(dialActivePin, d)
	ic := debouncedPin(dialIncrementPin, d)
	lc := debouncedPin(latchActivePin, d)

	for {
		select {
		case n := <-ac:
			fmt.Println("Got dial active", n)
		case n := <-ic:
			fmt.Println("Got dial increment", n)
		case n := <-lc:
			fmt.Println("Got latch", n)
		}
	}
}
