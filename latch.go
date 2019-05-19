package rotaryphone

import (
	"fmt"
	"time"
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
	c := debouncedPin(latchActivePin, 100*time.Millisecond)
	for n := range c {
		pin, value := n.Pin, n.Value

		fmt.Println("Latch got stable", pin, value)
		if value == 1 {
			l.Active <- true
		} else {
			l.Active <- false
		}
	}
}
